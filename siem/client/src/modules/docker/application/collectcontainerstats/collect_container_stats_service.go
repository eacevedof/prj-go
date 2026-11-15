package collectcontainerstats

import (
    "time"

    "siem-client/src/modules/docker/domain/enums"
    "siem-client/src/modules/docker/domain/services"
    "siem-client/src/modules/docker/infrastructure/repositories"
    "siem-client/src/modules/docker/infrastructure/traits"
    "siem-client/src/modules/shared/application"
)

// CollectContainerStatsService cada stats_interval_seconds, una muestra de cpu, memoria y red de cada contenedor en
// marcha como docker_stats. La medida de cpu de cada uno se guarda en state/docker-stats.json para calcular el % en
// la siguiente muestra. Un contenedor que se para entre la lista y su muestra se salta.
type CollectContainerStatsService struct {
    traits.DockerCollectorTrait

    containersReaderSocketRepository *repositories.ContainersReaderSocketRepository
    containerStatsCalculator         *services.ContainerStatsCalculator
    dockerEventBuilder               *services.DockerEventBuilder
}

func CollectContainerStatsServiceGetInstance() *CollectContainerStatsService {
    return &CollectContainerStatsService{
        DockerCollectorTrait:             traits.NewDockerCollectorTrait(),
        containersReaderSocketRepository: repositories.ContainersReaderSocketRepositoryGetInstance(),
        containerStatsCalculator:         services.ContainerStatsCalculatorGetInstance(),
        dockerEventBuilder:               services.DockerEventBuilderGetInstance(),
    }
}

// Invoke error: *AgentConfigException (400, 404), *DockerException (503) o uno inesperado (disco). Si no toca
// todavia, no hace nada (isNotDue).
func (s *CollectContainerStatsService) Invoke(collectContainerStatsDto *CollectContainerStatsDto) (*CollectContainerStatsResultDto, error) {
    agentConfigEntity, dockerConfigEntity, err := s.GetValidConfigOrFail()
    if err != nil {
        return nil, err
    }
    if !dockerConfigEntity.IsStatsEnabled() {
        return CollectContainerStatsResultDtoFromPrimitives(map[string]any{"isDisabled": 1}), nil
    }

    statsState, err := s.StateReaderJsonRepository.GetState(s.GetStatePath(), enums.DockerStateFileEnum.STATS)
    if err != nil {
        return nil, err
    }
    now := time.Now().UTC()
    lastRunUnix := int64(application.GetInt(statsState, enums.DockerStateKeyEnum.LAST_RUN_UNIX))
    isDue := collectContainerStatsDto.IsForced() || now.Unix()-lastRunUnix >= int64(dockerConfigEntity.GetStatsIntervalSeconds())
    if !isDue {
        return CollectContainerStatsResultDtoFromPrimitives(map[string]any{"isNotDue": 1}), nil
    }

    socketPath := dockerConfigEntity.GetSocketPath()
    containerEntities, err := s.containersReaderSocketRepository.GetAllContainerEntities(socketPath)
    if err != nil {
        return nil, s.GetDockerUnavailableException(socketPath, err)
    }

    previousCpuById, _ := statsState[enums.DockerStateKeyEnum.CPU_BY_CONTAINER_ID].(map[string]any)
    nextCpuById := map[string]any{}
    events := []map[string]any{}
    for _, containerEntity := range containerEntities {
        if !containerEntity.IsRunning() {
            continue
        }
        rawStats, err := s.containersReaderSocketRepository.GetContainerRawStats(socketPath, containerEntity.GetId())
        if err != nil {
            continue
        }
        previousCpu, _ := previousCpuById[containerEntity.GetId()].(map[string]any)
        containerStatsEntity := s.containerStatsCalculator.GetContainerStatsEntity(
            rawStats, getFloat(previousCpu, enums.DockerStateKeyEnum.CPU_TOTAL), getFloat(previousCpu, enums.DockerStateKeyEnum.SYSTEM_TOTAL),
        )
        nextCpuById[containerEntity.GetId()] = map[string]any{
            enums.DockerStateKeyEnum.CPU_TOTAL:    containerStatsEntity.GetCpuTotal(),
            enums.DockerStateKeyEnum.SYSTEM_TOTAL: containerStatsEntity.GetSystemTotal(),
        }
        events = append(events, s.dockerEventBuilder.GetContainerStatsEvent(containerEntity, containerStatsEntity, now))
    }

    spooledBatches, err := s.SaveEventsInSpool(events, agentConfigEntity.GetMaxEventsPerBatch())
    if err != nil {
        return nil, err
    }
    err = s.StateWriterJsonRepository.SaveState(s.GetStatePath(), enums.DockerStateFileEnum.STATS, map[string]any{
        enums.DockerStateKeyEnum.LAST_RUN_UNIX:       now.Unix(),
        enums.DockerStateKeyEnum.CPU_BY_CONTAINER_ID: nextCpuById,
    })
    if err != nil {
        return nil, err
    }
    return CollectContainerStatsResultDtoFromPrimitives(map[string]any{
        "sampledContainers": len(events),
        "spooledBatches":    spooledBatches,
        "results":           events,
    }), nil
}

// getFloat los numeros del JSON de estado vuelven como float64 (las medidas de cpu son nanosegundos: no caben en int32).
func getFloat(primitives map[string]any, key string) float64 {
    value, _ := primitives[key].(float64)
    return value
}
