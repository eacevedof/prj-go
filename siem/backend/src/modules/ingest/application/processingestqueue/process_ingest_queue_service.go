package processingestqueue

import (
    "encoding/json"
    "errors"
    "slices"

    "siem-backend/src/core/config"

    sharedenums "siem-backend/src/modules/shared/domain/enums"
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/components/datetimer"
    "siem-backend/src/modules/shared/infrastructure/traits"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"

    monitorsentities "siem-backend/src/modules/monitors/domain/entities"
    monitorsrepositories "siem-backend/src/modules/monitors/infrastructure/repositories"

    eventsenums "siem-backend/src/modules/events/domain/enums"
    eventsentities "siem-backend/src/modules/events/domain/entities"
    eventsservices "siem-backend/src/modules/events/domain/services"
    eventsrepositories "siem-backend/src/modules/events/infrastructure/repositories"

    "siem-backend/src/modules/ingest/domain/enums"
    "siem-backend/src/modules/ingest/infrastructure/repositories"
)

// ProcessIngestQueueService vacia ingest_queue: cada lote pendiente (FIFO) se normaliza y pasa a events, con su
// source al dia. Lo llaman el worker de la API (cada APP_QUEUE_WORKER_INTERVAL_SECONDS) y
// ./console ingest-queue-process.
// Un lote = una transaccion: o entran todos sus eventos y queda done, o no entra ninguno y se reintenta.
type ProcessIngestQueueService struct {
    traits.ErrorStackTrait

    dateTimer                         *datetimer.DateTimer
    eventNormalizer                   *eventsservices.EventNormalizer
    ingestQueueReaderSqliteRepository *repositories.IngestQueueReaderSqliteRepository
    ingestQueueWriterSqliteRepository *repositories.IngestQueueWriterSqliteRepository
    eventsWriterSqliteRepository      *eventsrepositories.EventsWriterSqliteRepository
    sourcesReaderSqliteRepository     *eventsrepositories.SourcesReaderSqliteRepository
    sourcesWriterSqliteRepository     *eventsrepositories.SourcesWriterSqliteRepository
    monitorsReaderSqliteRepository    *monitorsrepositories.MonitorsReaderSqliteRepository

    processIngestQueueDto *ProcessIngestQueueDto
    processedBatches      int
    failedBatches         int
    insertedEvents        int
    skippedEvents         int
}

func ProcessIngestQueueServiceGetInstance() *ProcessIngestQueueService {
    return &ProcessIngestQueueService{
        dateTimer:                         datetimer.GetInstance(),
        eventNormalizer:                   eventsservices.EventNormalizerGetInstance(),
        ingestQueueReaderSqliteRepository: repositories.IngestQueueReaderSqliteRepositoryGetInstance(),
        ingestQueueWriterSqliteRepository: repositories.IngestQueueWriterSqliteRepositoryGetInstance(),
        eventsWriterSqliteRepository:      eventsrepositories.EventsWriterSqliteRepositoryGetInstance(),
        sourcesReaderSqliteRepository:     eventsrepositories.SourcesReaderSqliteRepositoryGetInstance(),
        sourcesWriterSqliteRepository:     eventsrepositories.SourcesWriterSqliteRepositoryGetInstance(),
        monitorsReaderSqliteRepository:    monitorsrepositories.MonitorsReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *FieldsException (400) o uno inesperado (la BD). Un lote que falla no es un error del servicio: se
// cuenta en failedBatches y se sigue con el siguiente.
func (s *ProcessIngestQueueService) Invoke(processIngestQueueDto *ProcessIngestQueueDto) (*ProcessIngestQueueResultDto, error) {
    s.processIngestQueueDto = processIngestQueueDto

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }

    pendingBatches, err := s.ingestQueueReaderSqliteRepository.GetPendingBatches(processIngestQueueDto.GetMaxBatches())
    if err != nil {
        return nil, err
    }
    for _, pendingBatch := range pendingBatches {
        if err := s.moveBatchToEventsOrMarkAsFailed(pendingBatch); err != nil {
            return nil, err
        }
    }

    return ProcessIngestQueueResultDtoFromPrimitives(map[string]any{
        "processedBatches": s.processedBatches,
        "failedBatches":    s.failedBatches,
        "insertedEvents":   s.insertedEvents,
        "skippedEvents":    s.skippedEvents,
    }), nil
}

func (s *ProcessIngestQueueService) failIfWrongInput() error {
    maxBatches := s.processIngestQueueDto.GetMaxBatches()
    if maxBatches >= enums.IngestLimitEnum.MAX_BATCHES_MIN && maxBatches <= enums.IngestLimitEnum.MAX_BATCHES_MAX {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(
            config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS),
            []string{enums.IngestMessageEnum.MAX_BATCHES_MUST_BE_1_TO_1000},
        ),
    )
}

// moveBatchToEventsOrMarkAsFailed solo devuelve error si ni siquiera se puede anotar el fallo (la BD no responde).
func (s *ProcessIngestQueueService) moveBatchToEventsOrMarkAsFailed(pendingBatch map[string]any) error {
    batchError := s.moveBatchToEvents(pendingBatch)
    if batchError == nil {
        return nil
    }

    s.failedBatches++
    queueStatus := enums.QueueStatusEnum.PENDING
    if application.GetInt(pendingBatch, enums.IngestKeyEnum.ATTEMPTS)+1 >= enums.IngestLimitEnum.MAX_ATTEMPTS {
        queueStatus = enums.QueueStatusEnum.FAILED
    }
    _, err := s.ingestQueueWriterSqliteRepository.UpdateBatchAfterFailureById(
        application.GetInt(pendingBatch, enums.IngestKeyEnum.ID),
        queueStatus,
        batchError.Error(),
        s.dateTimer.GetNowDtIntoTargetTz(),
    )
    return err
}

func (s *ProcessIngestQueueService) moveBatchToEvents(pendingBatch map[string]any) error {
    rawEvents := []map[string]any{}
    payload := application.ToString(pendingBatch[enums.IngestKeyEnum.PAYLOAD])
    if err := json.Unmarshal([]byte(payload), &rawEvents); err != nil {
        return errors.New(enums.IngestMessageEnum.BATCH_PAYLOAD_IS_NOT_A_JSON_LIST)
    }

    eventEntities := s.getEventEntities(rawEvents, pendingBatch)
    skippedEvents := len(rawEvents) - len(eventEntities)

    isBatchClaimed := false
    err := s.ingestQueueWriterSqliteRepository.RunInSharedTransaction(
        func() error {
            claimedBatches, err := s.ingestQueueWriterSqliteRepository.UpdatePendingBatchAsDoneById(
                application.GetInt(pendingBatch, enums.IngestKeyEnum.ID), s.dateTimer.GetNowDtIntoTargetTz(),
            )
            if err != nil || claimedBatches == 0 {
                return err
            }
            isBatchClaimed = true
            return s.insertEventsAndUpdateSources(eventEntities)
        },
        &s.eventsWriterSqliteRepository.AbstractSqliteRepository,
        &s.sourcesReaderSqliteRepository.AbstractSqliteRepository,
        &s.sourcesWriterSqliteRepository.AbstractSqliteRepository,
    )
    if err != nil {
        return err
    }
    if !isBatchClaimed {
        return nil
    }

    s.processedBatches++
    s.insertedEvents += len(eventEntities)
    s.skippedEvents += skippedEvents
    return nil
}

func (s *ProcessIngestQueueService) getEventEntities(rawEvents []map[string]any, pendingBatch map[string]any) []*eventsentities.EventEntity {
    projectId := application.GetInt(pendingBatch, enums.IngestKeyEnum.PROJECT_ID)
    agentId := application.GetInt(pendingBatch, enums.IngestKeyEnum.AGENT_ID)
    receivedAt := application.GetTrimmedString(pendingBatch, enums.IngestKeyEnum.RECEIVED_AT)

    monitorEntitiesByUuid := s.getMonitorEntitiesByUuid(rawEvents, projectId)
    eventEntities := make([]*eventsentities.EventEntity, 0, len(rawEvents))
    for _, rawEvent := range rawEvents {
        s.loadMonitorThresholds(rawEvent, monitorEntitiesByUuid)
        if eventEntity := s.eventNormalizer.GetEventEntityOrNil(rawEvent, projectId, agentId, receivedAt); eventEntity != nil {
            eventEntities = append(eventEntities, eventEntity)
        }
    }
    return eventEntities
}

// getMonitorEntitiesByUuid los monitores del proyecto, solo si el lote trae algun evento de monitor (la mayoria de los
// lotes de logs no: se ahorran la consulta). Si la consulta falla, el lote entra igual, sin umbrales.
func (s *ProcessIngestQueueService) getMonitorEntitiesByUuid(rawEvents []map[string]any, projectId int) map[string]*monitorsentities.MonitorEntity {
    monitorEntitiesByUuid := map[string]*monitorsentities.MonitorEntity{}
    hasMonitorEvents := slices.ContainsFunc(rawEvents, func(rawEvent map[string]any) bool {
        return s.getMonitorUuid(rawEvent) != ""
    })
    if !hasMonitorEvents {
        return monitorEntitiesByUuid
    }
    monitorEntities, err := s.monitorsReaderSqliteRepository.GetMonitorsByProjectId(projectId)
    if err != nil {
        return monitorEntitiesByUuid
    }
    for _, monitorEntity := range monitorEntities {
        monitorEntitiesByUuid[monitorEntity.GetMonitorUuid()] = monitorEntity
    }
    return monitorEntitiesByUuid
}

// loadMonitorThresholds los umbrales de un evento de monitor salen SIEMPRE de la tabla monitors (los del panel): lo
// que mande el agente en esas dos claves se descarta. Con ellos decide la gravedad EventSeverityResolver.
func (s *ProcessIngestQueueService) loadMonitorThresholds(rawEvent map[string]any, monitorEntitiesByUuid map[string]*monitorsentities.MonitorEntity) {
    attributes, hasAttributes := rawEvent[eventsenums.EventsKeyEnum.ATTRIBUTES].(map[string]any)
    if !hasAttributes {
        return
    }
    delete(attributes, eventsenums.MonitorAttributeEnum.WARNING_THRESHOLD)
    delete(attributes, eventsenums.MonitorAttributeEnum.CRITICAL_THRESHOLD)

    monitorEntity := monitorEntitiesByUuid[s.getMonitorUuid(rawEvent)]
    if monitorEntity == nil {
        return
    }
    if warningThreshold := monitorEntity.GetWarningThreshold(); warningThreshold != nil {
        attributes[eventsenums.MonitorAttributeEnum.WARNING_THRESHOLD] = *warningThreshold
    }
    if criticalThreshold := monitorEntity.GetCriticalThreshold(); criticalThreshold != nil {
        attributes[eventsenums.MonitorAttributeEnum.CRITICAL_THRESHOLD] = *criticalThreshold
    }
}

func (s *ProcessIngestQueueService) getMonitorUuid(rawEvent map[string]any) string {
    attributes, _ := rawEvent[eventsenums.EventsKeyEnum.ATTRIBUTES].(map[string]any)
    return application.GetTrimmedString(attributes, eventsenums.MonitorAttributeEnum.MONITOR_UUID)
}

// insertEventsAndUpdateSources dentro de la transaccion del lote. Cada source se resuelve una vez (se crea si es su
// primer evento) y al final se queda con su ultimo evento del lote como ultimo estado conocido.
func (s *ProcessIngestQueueService) insertEventsAndUpdateSources(eventEntities []*eventsentities.EventEntity) error {
    sourceIdsBySourceKey := map[string]int{}
    lastEventEntitiesBySourceId := map[int]*eventsentities.EventEntity{}
    eventsCountBySourceId := map[int]int{}

    for _, eventEntity := range eventEntities {
        sourceKey := eventEntity.GetSourceType() + "|" + eventEntity.GetSourceName()
        sourceId, isResolved := sourceIdsBySourceKey[sourceKey]
        if !isResolved {
            var err error
            if sourceId, err = s.getSourceIdOrCreate(eventEntity); err != nil {
                return err
            }
            sourceIdsBySourceKey[sourceKey] = sourceId
        }

        if _, err := s.eventsWriterSqliteRepository.InsertEvent(eventEntity, sourceId); err != nil {
            return err
        }
        lastEventEntitiesBySourceId[sourceId] = eventEntity
        eventsCountBySourceId[sourceId]++
    }

    for sourceId, lastEventEntity := range lastEventEntitiesBySourceId {
        if _, err := s.sourcesWriterSqliteRepository.UpdateLastEventById(sourceId, lastEventEntity, eventsCountBySourceId[sourceId]); err != nil {
            return err
        }
    }
    return nil
}

func (s *ProcessIngestQueueService) getSourceIdOrCreate(eventEntity *eventsentities.EventEntity) (int, error) {
    sourceId, err := s.sourcesReaderSqliteRepository.GetSourceIdByAgentIdAndTypeAndName(
        eventEntity.GetAgentId(), eventEntity.GetSourceType(), eventEntity.GetSourceName(),
    )
    if err != nil || sourceId > 0 {
        return sourceId, err
    }
    return s.sourcesWriterSqliteRepository.InsertSource(
        eventEntity.GetAgentId(), eventEntity.GetSourceType(), eventEntity.GetSourceName(), eventEntity.GetReceivedAt(),
    )
}
