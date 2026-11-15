package collectdaemonevents

import (
    "strings"
    "time"

    "siem-client/src/modules/docker/domain/enums"
    "siem-client/src/modules/docker/domain/services"
    "siem-client/src/modules/docker/infrastructure/repositories"
    "siem-client/src/modules/docker/infrastructure/traits"
    "siem-client/src/modules/shared/application"
)

const (
    // un healthcheck o un docker exec generan exec_create / exec_start / exec_die en cada ejecucion: ruido
    ignoredActionPrefix = "exec_"
    // una rafaga (docker compose down de 50 contenedores) no debe generar un lote gigante
    MAX_EVENTS_PER_ROUND = 1000
)

// CollectDaemonEventsService los eventos del daemon (create, start, die, oom, kill, pull...) desde la pasada anterior,
// tal cual, como docker_event: la auditoria de lo que pasa en Docker. Desde donde se leyo queda en
// state/docker-events.json (marca RFC3339 con nanosegundos); la primera vez arranca en "ahora".
type CollectDaemonEventsService struct {
    traits.DockerCollectorTrait

    daemonEventsReaderSocketRepository *repositories.DaemonEventsReaderSocketRepository
    dockerEventBuilder                 *services.DockerEventBuilder
}

func CollectDaemonEventsServiceGetInstance() *CollectDaemonEventsService {
    return &CollectDaemonEventsService{
        DockerCollectorTrait:               traits.NewDockerCollectorTrait(),
        daemonEventsReaderSocketRepository: repositories.DaemonEventsReaderSocketRepositoryGetInstance(),
        dockerEventBuilder:                 services.DockerEventBuilderGetInstance(),
    }
}

// Invoke error: *AgentConfigException (400, 404), *DockerException (503) o uno inesperado (disco).
func (s *CollectDaemonEventsService) Invoke(collectDaemonEventsDto *CollectDaemonEventsDto) (*CollectDaemonEventsResultDto, error) {
    agentConfigEntity, dockerConfigEntity, err := s.GetValidConfigOrFail()
    if err != nil {
        return nil, err
    }
    if !dockerConfigEntity.IsEventsEnabled() {
        return CollectDaemonEventsResultDtoFromPrimitives(map[string]any{"isDisabled": 1}), nil
    }

    eventsState, err := s.StateReaderJsonRepository.GetState(s.GetStatePath(), enums.DockerStateFileEnum.EVENTS)
    if err != nil {
        return nil, err
    }
    // ventanas contiguas con nanosegundos: cada pasada lee [since, now] y la siguiente empieza en ese now
    now := time.Now().UTC()
    since, err := time.Parse(time.RFC3339Nano, application.GetTrimmedString(eventsState, enums.DockerStateKeyEnum.SINCE))
    if err != nil {
        return CollectDaemonEventsResultDtoFromPrimitives(map[string]any{}), s.saveSince(now)
    }
    if !now.After(since) {
        return CollectDaemonEventsResultDtoFromPrimitives(map[string]any{}), nil
    }

    rawDaemonEvents, err := s.daemonEventsReaderSocketRepository.GetRawDaemonEvents(dockerConfigEntity.GetSocketPath(), since, now)
    if err != nil {
        return nil, s.GetDockerUnavailableException(dockerConfigEntity.GetSocketPath(), err)
    }

    events := []map[string]any{}
    droppedEvents := 0
    for _, rawDaemonEvent := range rawDaemonEvents {
        if strings.HasPrefix(application.GetTrimmedString(rawDaemonEvent, "Action"), ignoredActionPrefix) {
            continue
        }
        if len(events) >= MAX_EVENTS_PER_ROUND {
            droppedEvents++
            continue
        }
        events = append(events, s.dockerEventBuilder.GetDaemonEvent(rawDaemonEvent))
    }

    spooledBatches, err := s.SaveEventsInSpool(events, agentConfigEntity.GetMaxEventsPerBatch())
    if err != nil {
        return nil, err
    }
    if err := s.saveSince(now); err != nil {
        return nil, err
    }
    return CollectDaemonEventsResultDtoFromPrimitives(map[string]any{
        "daemonEvents":   len(events),
        "droppedEvents":  droppedEvents,
        "spooledBatches": spooledBatches,
    }), nil
}

// saveSince sin cursor (primera vez o estado corrupto) arranca en "ahora": no se reenvia el historico.
func (s *CollectDaemonEventsService) saveSince(since time.Time) error {
    return s.StateWriterJsonRepository.SaveState(s.GetStatePath(), enums.DockerStateFileEnum.EVENTS, map[string]any{
        enums.DockerStateKeyEnum.SINCE: since.Format(time.RFC3339Nano),
    })
}
