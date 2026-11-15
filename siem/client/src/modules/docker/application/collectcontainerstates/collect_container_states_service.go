package collectcontainerstates

import (
    "time"

    "siem-client/src/modules/docker/domain/entities"
    "siem-client/src/modules/docker/domain/enums"
    "siem-client/src/modules/docker/domain/services"
    "siem-client/src/modules/docker/infrastructure/repositories"
    "siem-client/src/modules/docker/infrastructure/traits"
)

// CollectContainerStatesService una pasada por todos los contenedores: los que han cambiado desde la anterior (estado,
// reinicios, health, borrados) pasan al spool como eventos docker_state. Lo visto queda en
// state/docker-containers.json. En la primera pasada, un evento por contenedor con su estado de partida.
type CollectContainerStatesService struct {
    traits.DockerCollectorTrait

    containersReaderSocketRepository *repositories.ContainersReaderSocketRepository
    containerChangesDetector         *services.ContainerChangesDetector
    dockerEventBuilder               *services.DockerEventBuilder
}

func CollectContainerStatesServiceGetInstance() *CollectContainerStatesService {
    return &CollectContainerStatesService{
        DockerCollectorTrait:             traits.NewDockerCollectorTrait(),
        containersReaderSocketRepository: repositories.ContainersReaderSocketRepositoryGetInstance(),
        containerChangesDetector:         services.ContainerChangesDetectorGetInstance(),
        dockerEventBuilder:               services.DockerEventBuilderGetInstance(),
    }
}

// Invoke error: *AgentConfigException (400, 404), *DockerException (503: no responde) o uno inesperado (disco).
// Con Docker apagado en la config no hace nada (isDisabled).
func (s *CollectContainerStatesService) Invoke(collectContainerStatesDto *CollectContainerStatesDto) (*CollectContainerStatesResultDto, error) {
    agentConfigEntity, dockerConfigEntity, err := s.GetValidConfigOrFail()
    if err != nil {
        return nil, err
    }
    if !dockerConfigEntity.IsEnabled() {
        return CollectContainerStatesResultDtoFromPrimitives(map[string]any{"isDisabled": 1}), nil
    }

    containersState, err := s.StateReaderJsonRepository.GetState(s.GetStatePath(), enums.DockerStateFileEnum.CONTAINERS)
    if err != nil {
        return nil, err
    }
    containerEntities, err := s.containersReaderSocketRepository.GetAllContainerEntities(dockerConfigEntity.GetSocketPath())
    if err != nil {
        return nil, s.GetDockerUnavailableException(dockerConfigEntity.GetSocketPath(), err)
    }

    isFirstScan := containersState[enums.DockerStateKeyEnum.IS_INITIALIZED] != true
    containerChangeEntities := s.containerChangesDetector.GetContainerChangeEntities(
        s.getPreviousContainersById(containersState), containerEntities, isFirstScan,
    )
    now := time.Now().UTC()
    events := make([]map[string]any, 0, len(containerChangeEntities))
    for _, containerChangeEntity := range containerChangeEntities {
        events = append(events, s.dockerEventBuilder.GetContainerStateEvent(containerChangeEntity, now))
    }

    spooledBatches, err := s.SaveEventsInSpool(events, agentConfigEntity.GetMaxEventsPerBatch())
    if err != nil {
        return nil, err
    }
    if err := s.saveContainersState(containerEntities); err != nil {
        return nil, err
    }

    return CollectContainerStatesResultDtoFromPrimitives(map[string]any{
        "containers":        len(containerEntities),
        "runningContainers": s.getRunningCount(containerEntities),
        "changes":           len(events),
        "spooledBatches":    spooledBatches,
        "results":           events,
    }), nil
}

func (s *CollectContainerStatesService) getPreviousContainersById(containersState map[string]any) map[string]map[string]any {
    previousById := map[string]map[string]any{}
    rawContainers, _ := containersState[enums.DockerStateKeyEnum.CONTAINERS].(map[string]any)
    for containerId, rawContainer := range rawContainers {
        if previous, isMap := rawContainer.(map[string]any); isMap {
            previousById[containerId] = previous
        }
    }
    return previousById
}

func (s *CollectContainerStatesService) saveContainersState(containerEntities []*entities.ContainerEntity) error {
    containers := map[string]any{}
    for _, containerEntity := range containerEntities {
        containers[containerEntity.GetId()] = containerEntity.ToStatePrimitives()
    }
    return s.StateWriterJsonRepository.SaveState(s.GetStatePath(), enums.DockerStateFileEnum.CONTAINERS, map[string]any{
        enums.DockerStateKeyEnum.IS_INITIALIZED: true,
        enums.DockerStateKeyEnum.CONTAINERS:     containers,
    })
}

func (s *CollectContainerStatesService) getRunningCount(containerEntities []*entities.ContainerEntity) int {
    runningCount := 0
    for _, containerEntity := range containerEntities {
        if containerEntity.IsRunning() {
            runningCount++
        }
    }
    return runningCount
}
