package collectcontainerlogs

import (
    "os"
    "strings"
    "time"

    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/docker/domain/entities"
    "siem-client/src/modules/docker/domain/enums"
    "siem-client/src/modules/docker/domain/services"
    "siem-client/src/modules/docker/infrastructure/repositories"
    "siem-client/src/modules/docker/infrastructure/traits"
)

// CollectContainerLogsService las lineas nuevas del log de cada contenedor (los que casan con logs_include y no con
// logs_exclude) como docker_log. El cursor de cada uno es la marca de tiempo de su ultima linea enviada
// (state/docker-logs.json). Reglas:
//   - primera pasada del agente: los contenedores que ya existen arrancan en "ahora" (no se reenvia el historico).
//   - un contenedor que aparece despues se lee desde su principio (hasta logs_max_lines_per_container).
//   - el contenedor del propio agente nunca: cada linea que imprime generaria otro evento (bucle).
//   - uno parado solo se relee si se paro despues de su cursor (sus ultimas lineas antes de caer).
type CollectContainerLogsService struct {
    traits.DockerCollectorTrait

    containersReaderSocketRepository *repositories.ContainersReaderSocketRepository
    dockerLogStreamParser            *services.DockerLogStreamParser
    dockerEventBuilder               *services.DockerEventBuilder
}

func CollectContainerLogsServiceGetInstance() *CollectContainerLogsService {
    return &CollectContainerLogsService{
        DockerCollectorTrait:             traits.NewDockerCollectorTrait(),
        containersReaderSocketRepository: repositories.ContainersReaderSocketRepositoryGetInstance(),
        dockerLogStreamParser:            services.DockerLogStreamParserGetInstance(),
        dockerEventBuilder:               services.DockerEventBuilderGetInstance(),
    }
}

// Invoke error: *AgentConfigException (400, 404), *DockerException (503) o uno inesperado (disco). Un contenedor cuyo
// log no se puede leer se salta en esta pasada.
func (s *CollectContainerLogsService) Invoke(collectContainerLogsDto *CollectContainerLogsDto) (*CollectContainerLogsResultDto, error) {
    agentConfigEntity, dockerConfigEntity, err := s.GetValidConfigOrFail()
    if err != nil {
        return nil, err
    }
    if !dockerConfigEntity.IsLogsEnabled() {
        return CollectContainerLogsResultDtoFromPrimitives(map[string]any{"isDisabled": 1}), nil
    }

    logsState, err := s.StateReaderJsonRepository.GetState(s.GetStatePath(), enums.DockerStateFileEnum.LOGS)
    if err != nil {
        return nil, err
    }
    socketPath := dockerConfigEntity.GetSocketPath()
    containerEntities, err := s.containersReaderSocketRepository.GetAllContainerEntities(socketPath)
    if err != nil {
        return nil, s.GetDockerUnavailableException(socketPath, err)
    }

    isInitialized := logsState[enums.DockerStateKeyEnum.IS_INITIALIZED] == true
    lastTsById, _ := logsState[enums.DockerStateKeyEnum.LAST_TS_BY_CONTAINER_ID].(map[string]any)
    nextLastTsById := map[string]any{}
    selfContainerIdPrefix := s.getSelfContainerIdPrefix()
    now := time.Now().UTC()
    events := []map[string]any{}
    containersRead := 0
    for _, containerEntity := range containerEntities {
        if !dockerConfigEntity.IsContainerLogged(containerEntity.GetName()) || s.isSelf(containerEntity, selfContainerIdPrefix) {
            continue
        }
        lastTs, hasCursor := s.getCursor(lastTsById, containerEntity.GetId())
        if !hasCursor && !isInitialized {
            nextLastTsById[containerEntity.GetId()] = now.Format(time.RFC3339Nano)
            continue
        }
        nextLastTsById[containerEntity.GetId()] = lastTs.Format(time.RFC3339Nano)
        if hasCursor && !s.mayHaveNewLines(containerEntity, lastTs) {
            continue
        }

        containerEvents, newestTs := s.getContainerLogEvents(socketPath, containerEntity, lastTs, dockerConfigEntity)
        containersRead++
        events = append(events, containerEvents...)
        if !newestTs.IsZero() {
            nextLastTsById[containerEntity.GetId()] = newestTs.Format(time.RFC3339Nano)
        }
    }

    spooledBatches, err := s.SaveEventsInSpool(events, agentConfigEntity.GetMaxEventsPerBatch())
    if err != nil {
        return nil, err
    }
    err = s.StateWriterJsonRepository.SaveState(s.GetStatePath(), enums.DockerStateFileEnum.LOGS, map[string]any{
        enums.DockerStateKeyEnum.IS_INITIALIZED:          true,
        enums.DockerStateKeyEnum.LAST_TS_BY_CONTAINER_ID: nextLastTsById,
    })
    if err != nil {
        return nil, err
    }
    return CollectContainerLogsResultDtoFromPrimitives(map[string]any{
        "containersRead": containersRead,
        "logLines":       len(events),
        "spooledBatches": spooledBatches,
    }), nil
}

// getContainerLogEvents las lineas posteriores al cursor (Docker repite la del propio since) y la hora de la mas nueva.
func (s *CollectContainerLogsService) getContainerLogEvents(socketPath string, containerEntity *entities.ContainerEntity, lastTs time.Time, dockerConfigEntity *agentconfigentities.DockerConfigEntity) ([]map[string]any, time.Time) {
    body, err := s.containersReaderSocketRepository.GetContainerLogsBody(socketPath, containerEntity.GetId(), lastTs, dockerConfigEntity.GetLogsMaxLinesPerContainer())
    if err != nil {
        return nil, time.Time{}
    }
    events := []map[string]any{}
    var newestTs time.Time
    for _, containerLogLineEntity := range s.dockerLogStreamParser.GetContainerLogLineEntities(body, containerEntity.IsTty()) {
        if !containerLogLineEntity.GetOccurredAt().After(lastTs) {
            continue
        }
        events = append(events, s.dockerEventBuilder.GetContainerLogEvent(containerEntity, containerLogLineEntity))
        if containerLogLineEntity.GetOccurredAt().After(newestTs) {
            newestTs = containerLogLineEntity.GetOccurredAt()
        }
    }
    return events, newestTs
}

// mayHaveNewLines en marcha, siempre; parado, solo si se paro despues del cursor.
func (s *CollectContainerLogsService) mayHaveNewLines(containerEntity *entities.ContainerEntity, lastTs time.Time) bool {
    if containerEntity.IsRunning() {
        return true
    }
    finishedAt, err := time.Parse(time.RFC3339Nano, containerEntity.GetFinishedAt())
    return err == nil && finishedAt.After(lastTs)
}

// getCursor sin cursor: tiempo cero (= desde el principio del log).
func (s *CollectContainerLogsService) getCursor(lastTsById map[string]any, containerId string) (time.Time, bool) {
    rawLastTs, hasCursor := lastTsById[containerId].(string)
    if !hasCursor {
        return time.Time{}, false
    }
    lastTs, err := time.Parse(time.RFC3339Nano, rawLastTs)
    if err != nil {
        return time.Time{}, false
    }
    return lastTs, true
}

// getSelfContainerIdPrefix dentro de un contenedor el hostname es su id corto (salvo que se cambie con hostname:); en
// un host normal no casara con ningun id.
func (s *CollectContainerLogsService) getSelfContainerIdPrefix() string {
    hostname, err := os.Hostname()
    if err != nil || len(hostname) < entities.SHORT_ID_LENGTH {
        return ""
    }
    return hostname
}

func (s *CollectContainerLogsService) isSelf(containerEntity *entities.ContainerEntity, selfContainerIdPrefix string) bool {
    return selfContainerIdPrefix != "" && strings.HasPrefix(containerEntity.GetId(), selfContainerIdPrefix)
}
