package sendspooledevents

import (
    "runtime"
    "strconv"
    "time"

    "siem-client/src/core/config"
    agentconfigtraits "siem-client/src/modules/agentconfig/infrastructure/traits"
    "siem-client/src/modules/shared/application"
    "siem-client/src/modules/shared/infrastructure/components/server"
    "siem-client/src/modules/shared/infrastructure/components/uuid"
    "siem-client/src/modules/shared/infrastructure/repositories/configuration"
    "siem-client/src/modules/shared/infrastructure/repositories/spool"
    "siem-client/src/modules/shared/infrastructure/repositories/state"
    "siem-client/src/modules/shipper/domain/enums"
    "siem-client/src/modules/shipper/domain/exceptions"
    "siem-client/src/modules/shipper/infrastructure/repositories"
)

const responseBodyMaxLengthInError = 300

// SendSpooledEventsService vacia el spool hacia el backend, del lote mas antiguo al mas nuevo:
//   - 202: el backend lo tiene; se borra del spool.
//   - 400: el backend no lo admitira nunca; se aparta como <batch_uuid>.rejected.
//   - 401 / sin respuesta / otro codigo: se para la pasada y el lote sigue en el spool (reintento en la siguiente).
//
// El batch_uuid es el nombre del fichero: un reintento de un lote que si llego no duplica (el backend lo reconoce).
// Sin nada que enviar, cada heartbeat_interval_seconds manda un lote vacio: el latido que mantiene el agente online.
// Si el spool pasa de max_spool_batches (backend caido mucho tiempo), se descartan los lotes mas antiguos.
type SendSpooledEventsService struct {
    agentconfigtraits.AgentConfigLoaderTrait

    serverIdentifier                  *server.ServerIdentifier
    uuidGenerator                     *uuid.UuidGenerator
    environmentReaderRawRepository    *configuration.EnvironmentReaderRawRepository
    spoolReaderDiskRepository         *spool.SpoolReaderDiskRepository
    spoolWriterDiskRepository         *spool.SpoolWriterDiskRepository
    stateReaderJsonRepository         *state.StateReaderJsonRepository
    stateWriterJsonRepository         *state.StateWriterJsonRepository
    ingestBatchesWriterHttpRepository *repositories.IngestBatchesWriterHttpRepository

    sendSpooledEventsDto *SendSpooledEventsDto
    sentBatches          int
    sentEvents           int
    rejectedBatches      int
}

func SendSpooledEventsServiceGetInstance() *SendSpooledEventsService {
    return &SendSpooledEventsService{
        AgentConfigLoaderTrait:            agentconfigtraits.NewAgentConfigLoaderTrait(),
        serverIdentifier:                  server.GetInstance(),
        uuidGenerator:                     uuid.GetInstance(),
        environmentReaderRawRepository:    configuration.GetInstance(),
        spoolReaderDiskRepository:         spool.SpoolReaderDiskRepositoryGetInstance(),
        spoolWriterDiskRepository:         spool.SpoolWriterDiskRepositoryGetInstance(),
        stateReaderJsonRepository:         state.StateReaderJsonRepositoryGetInstance(),
        stateWriterJsonRepository:         state.StateWriterJsonRepositoryGetInstance(),
        ingestBatchesWriterHttpRepository: repositories.IngestBatchesWriterHttpRepositoryGetInstance(),
    }
}

// Invoke error: *AgentConfigException (400, 404), *ShipperException (401, 503) o uno inesperado (disco).
func (s *SendSpooledEventsService) Invoke(sendSpooledEventsDto *SendSpooledEventsDto) (*SendSpooledEventsResultDto, error) {
    s.sendSpooledEventsDto = sendSpooledEventsDto

    agentConfigEntity, err := s.GetValidAgentConfigEntityOrFail()
    if err != nil {
        return nil, err
    }

    batchUuids, err := s.spoolReaderDiskRepository.GetAllSpoolBatchUuids(s.GetSpoolPath())
    if err != nil {
        return nil, err
    }
    droppedBatches, err := s.dropOldestBatchesOverLimit(batchUuids, agentConfigEntity.GetMaxSpoolBatches())
    if err != nil {
        return nil, err
    }
    batchUuids = batchUuids[droppedBatches:]

    for _, batchUuid := range batchUuids[:min(len(batchUuids), sendSpooledEventsDto.GetMaxBatches())] {
        if err := s.sendSpoolBatch(batchUuid); err != nil {
            return nil, err
        }
    }

    isHeartbeatSent := 0
    if len(batchUuids) == 0 {
        isHeartbeatDue, err := s.isHeartbeatDue(agentConfigEntity.GetHeartbeatIntervalSeconds())
        if err != nil {
            return nil, err
        }
        if isHeartbeatDue {
            if err := s.sendIngestBatch(s.uuidGenerator.GetUniqueIdWithPrefix("hb-"), []map[string]any{}); err != nil {
                return nil, err
            }
            isHeartbeatSent = 1
        }
    }

    return SendSpooledEventsResultDtoFromPrimitives(map[string]any{
        "sentBatches":     s.sentBatches,
        "sentEvents":      s.sentEvents,
        "rejectedBatches": s.rejectedBatches,
        "droppedBatches":  droppedBatches,
        "pendingBatches":  len(batchUuids) - s.sentBatches - s.rejectedBatches,
        "isHeartbeatSent": isHeartbeatSent,
    }), nil
}

func (s *SendSpooledEventsService) dropOldestBatchesOverLimit(batchUuids []string, maxSpoolBatches int) (int, error) {
    batchesToDrop := len(batchUuids) - maxSpoolBatches
    for _, batchUuid := range batchUuids[:max(0, batchesToDrop)] {
        if err := s.spoolWriterDiskRepository.DeleteSpoolBatch(s.GetSpoolPath(), batchUuid); err != nil {
            return 0, err
        }
    }
    return max(0, batchesToDrop), nil
}

// sendSpoolBatch un lote que no se puede leer (fichero corrupto) se aparta como rejected: no tiene arreglo.
func (s *SendSpooledEventsService) sendSpoolBatch(batchUuid string) error {
    events, err := s.spoolReaderDiskRepository.GetSpoolBatchEvents(s.GetSpoolPath(), batchUuid)
    if err != nil {
        s.rejectedBatches++
        return s.spoolWriterDiskRepository.RenameSpoolBatchAsRejected(s.GetSpoolPath(), batchUuid)
    }

    err = s.sendIngestBatch(batchUuid, events)
    if err == errBatchRejected {
        s.rejectedBatches++
        return s.spoolWriterDiskRepository.RenameSpoolBatchAsRejected(s.GetSpoolPath(), batchUuid)
    }
    if err != nil {
        return err
    }

    s.sentBatches++
    s.sentEvents += len(events)
    return s.spoolWriterDiskRepository.DeleteSpoolBatch(s.GetSpoolPath(), batchUuid)
}

// errBatchRejected senal interna (no sale del servicio): el backend contesto 400 a este lote.
var errBatchRejected = exceptions.ServiceUnavailableCustom("batch rejected")

func (s *SendSpooledEventsService) sendIngestBatch(batchUuid string, events []map[string]any) error {
    statusCode, responseBody, err := s.ingestBatchesWriterHttpRepository.SendIngestBatch(
        s.GetSiemEndpointUrl(),
        s.GetSiemApiToken(),
        map[string]any{
            enums.ShipperKeyEnum.BATCH_UUID: batchUuid,
            enums.ShipperKeyEnum.AGENT: map[string]any{
                enums.ShipperKeyEnum.HOST_NAME:      s.serverIdentifier.GetServerName(),
                enums.ShipperKeyEnum.HOST_OS:        runtime.GOOS,
                enums.ShipperKeyEnum.CLIENT_VERSION: s.environmentReaderRawRepository.GetAppVersion(),
            },
            enums.ShipperKeyEnum.EVENTS: events,
        },
    )
    if err != nil {
        return exceptions.ServiceUnavailableCustom(config.Tr(enums.ShipperMessageEnum.BACKEND_UNAVAILABLE_X, err.Error()))
    }

    switch statusCode {
    case 202:
        return s.saveLastSentAt()
    case 400:
        return errBatchRejected
    case 401:
        return exceptions.UnauthorizedCustom(enums.ShipperMessageEnum.BACKEND_REJECTED_API_TOKEN)
    }
    return exceptions.ServiceUnavailableCustom(config.Tr(
        enums.ShipperMessageEnum.BACKEND_UNEXPECTED_STATUS_X_X, strconv.Itoa(statusCode), s.getTruncated(responseBody),
    ))
}

func (s *SendSpooledEventsService) isHeartbeatDue(heartbeatIntervalSeconds int) (bool, error) {
    shipperState, err := s.stateReaderJsonRepository.GetState(s.GetStatePath(), enums.ShipperStateFileEnum.SHIPPER)
    if err != nil {
        return false, err
    }
    lastSentAtUnix := int64(application.ToInt(shipperState[enums.ShipperKeyEnum.LAST_SENT_AT_UNIX]))
    return time.Now().Unix()-lastSentAtUnix >= int64(heartbeatIntervalSeconds), nil
}

// saveLastSentAt cualquier lote aceptado cuenta como latido.
func (s *SendSpooledEventsService) saveLastSentAt() error {
    return s.stateWriterJsonRepository.SaveState(s.GetStatePath(), enums.ShipperStateFileEnum.SHIPPER, map[string]any{
        enums.ShipperKeyEnum.LAST_SENT_AT_UNIX: time.Now().Unix(),
    })
}

func (s *SendSpooledEventsService) getTruncated(text string) string {
    runes := []rune(text)
    if len(runes) <= responseBodyMaxLengthInError {
        return text
    }
    return string(runes[:responseBodyMaxLengthInError])
}
