package enqueueingestbatch

import (
    "siem-backend/src/core/config"
    agentsentities "siem-backend/src/modules/agents/domain/entities"
    agentsenums "siem-backend/src/modules/agents/domain/enums"
    agentsrepositories "siem-backend/src/modules/agents/infrastructure/repositories"
    "siem-backend/src/modules/ingest/domain/enums"
    "siem-backend/src/modules/ingest/domain/exceptions"
    "siem-backend/src/modules/ingest/domain/services"
    "siem-backend/src/modules/ingest/infrastructure/repositories"
    sharedenums "siem-backend/src/modules/shared/domain/enums"
    "siem-backend/src/modules/shared/infrastructure/components/datetimer"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"
    "siem-backend/src/modules/shared/infrastructure/traits"
)

// EnqueueIngestBatchService lo que hace el endpoint de ingesta, y nada mas: anota el latido del agente y deja el lote
// crudo en ingest_queue. Convertirlo en eventos es de ProcessIngestQueueService (el worker): asi el siem-client
// recibe su 202 enseguida aunque lleguen muchos lotes a la vez.
// Un lote sin eventos es solo un latido: no se encola.
type EnqueueIngestBatchService struct {
    traits.ErrorStackTrait

    dateTimer                         *datetimer.DateTimer
    agentsWriterSqliteRepository      *agentsrepositories.AgentsWriterSqliteRepository
    ingestQueueWriterSqliteRepository *repositories.IngestQueueWriterSqliteRepository

    enqueueIngestBatchDto *EnqueueIngestBatchDto
}

func EnqueueIngestBatchServiceGetInstance() *EnqueueIngestBatchService {
    return &EnqueueIngestBatchService{
        dateTimer:                         datetimer.GetInstance(),
        agentsWriterSqliteRepository:      agentsrepositories.AgentsWriterSqliteRepositoryGetInstance(),
        ingestQueueWriterSqliteRepository: repositories.IngestQueueWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *IngestException (401), *FieldsException (400) o uno inesperado.
func (s *EnqueueIngestBatchService) Invoke(enqueueIngestBatchDto *EnqueueIngestBatchDto) (*EnqueueIngestBatchResultDto, error) {
    s.enqueueIngestBatchDto = enqueueIngestBatchDto

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }

    receivedAt := s.dateTimer.GetNowDtIntoTargetTz()
    if err := s.updateAgentHeartbeat(receivedAt); err != nil {
        return nil, err
    }

    eventsCount := enqueueIngestBatchDto.GetEventsCount()
    if eventsCount == 0 {
        return s.getResultDto(0, 0), nil
    }

    insertedBatches, err := s.ingestQueueWriterSqliteRepository.InsertBatchIfNew(
        enqueueIngestBatchDto.GetAuthAgentId(),
        enqueueIngestBatchDto.GetBatchUuid(),
        enqueueIngestBatchDto.GetEventsJson(),
        eventsCount,
        receivedAt,
    )
    if err != nil {
        return nil, err
    }
    if insertedBatches == 0 {
        return s.getResultDto(0, 1), nil
    }
    return s.getResultDto(eventsCount, 0), nil
}

func (s *EnqueueIngestBatchService) failIfWrongInput() error {
    if s.enqueueIngestBatchDto.GetAuthAgentId() < 1 {
        return exceptions.UnauthorizedCustom(config.Tr(enums.IngestMessageEnum.INVALID_OR_MISSING_API_TOKEN))
    }

    errors := services.IngestBatchValidatorGetInstance().GetErrors(
        s.enqueueIngestBatchDto.GetBatchUuid(), s.enqueueIngestBatchDto.GetEventsCount(),
    )
    if len(errors) == 0 {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
    )
}

func (s *EnqueueIngestBatchService) updateAgentHeartbeat(lastSeenAt string) error {
    _, err := s.agentsWriterSqliteRepository.UpdateHostDataAndLastSeenById(
        agentsentities.AgentEntityFromPrimitives(map[string]any{
            agentsenums.AgentsKeyEnum.ID:             s.enqueueIngestBatchDto.GetAuthAgentId(),
            agentsenums.AgentsKeyEnum.HOST_NAME:      s.getTruncated(s.enqueueIngestBatchDto.GetHostName(), 150),
            agentsenums.AgentsKeyEnum.HOST_OS:        s.getTruncated(s.enqueueIngestBatchDto.GetHostOs(), 50),
            agentsenums.AgentsKeyEnum.HOST_IP:        s.getTruncated(s.enqueueIngestBatchDto.GetUserIpAddress(), 45),
            agentsenums.AgentsKeyEnum.CLIENT_VERSION: s.getTruncated(s.enqueueIngestBatchDto.GetClientVersion(), 20),
            agentsenums.AgentsKeyEnum.LAST_SEEN_AT:   lastSeenAt,
        }),
    )
    return err
}

func (s *EnqueueIngestBatchService) getResultDto(queuedEvents int, isDuplicate int) *EnqueueIngestBatchResultDto {
    return EnqueueIngestBatchResultDtoFromPrimitives(map[string]any{
        "batchUuid":    s.enqueueIngestBatchDto.GetBatchUuid(),
        "queuedEvents": queuedEvents,
        "isDuplicate":  isDuplicate,
    })
}

func (s *EnqueueIngestBatchService) getTruncated(text string, maxLength int) string {
    runes := []rune(text)
    if len(runes) <= maxLength {
        return text
    }
    return string(runes[:maxLength])
}
