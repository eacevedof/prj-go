package purgeevents

import (
    "siem-backend/src/core/config"

    sharedenums "siem-backend/src/modules/shared/domain/enums"
    "siem-backend/src/modules/shared/infrastructure/components/datetimer"
    "siem-backend/src/modules/shared/infrastructure/traits"
    "siem-backend/src/modules/shared/infrastructure/repositories/configuration"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"

    ingestrepositories "siem-backend/src/modules/ingest/infrastructure/repositories"

    "siem-backend/src/modules/events/domain/enums"
    "siem-backend/src/modules/events/infrastructure/repositories"
)

// PurgeEventsService retencion: borra los eventos con mas de N dias y los lotes ya procesados de la cola. Es lo que
// mantiene la BD ligera: se lanza por cron (./console events-purge).
type PurgeEventsService struct {
    traits.ErrorStackTrait

    dateTimer                        *datetimer.DateTimer
    environmentReaderRawRepository   *configuration.EnvironmentReaderRawRepository
    eventsWriterSqliteRepository     *repositories.EventsWriterSqliteRepository
    ingestQueueWriterSqliteRepository *ingestrepositories.IngestQueueWriterSqliteRepository

    purgeEventsDto *PurgeEventsDto
    retentionDays  int
}

func PurgeEventsServiceGetInstance() *PurgeEventsService {
    return &PurgeEventsService{
        dateTimer:                         datetimer.GetInstance(),
        environmentReaderRawRepository:    configuration.GetInstance(),
        eventsWriterSqliteRepository:      repositories.EventsWriterSqliteRepositoryGetInstance(),
        ingestQueueWriterSqliteRepository: ingestrepositories.IngestQueueWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *FieldsException (400) o uno inesperado.
func (s *PurgeEventsService) Invoke(purgeEventsDto *PurgeEventsDto) (*PurgeEventsResultDto, error) {
    s.purgeEventsDto = purgeEventsDto
    s.loadRetentionDays()

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }

    deletedBefore := s.dateTimer.GetNowDtMinusSeconds(s.retentionDays * enums.EventsLimitEnum.SECONDS_PER_DAY)
    deletedEvents, err := s.eventsWriterSqliteRepository.DeleteEventsReceivedBefore(deletedBefore)
    if err != nil {
        return nil, err
    }
    deletedBatches, err := s.ingestQueueWriterSqliteRepository.DeleteDoneBatchesProcessedBefore(
        s.dateTimer.GetNowDtMinusSeconds(enums.EventsLimitEnum.DONE_BATCHES_RETENTION_SECONDS),
    )
    if err != nil {
        return nil, err
    }

    return PurgeEventsResultDtoFromPrimitives(map[string]any{
        "retentionDays":  s.retentionDays,
        "deletedBefore":  deletedBefore,
        "deletedEvents":  deletedEvents,
        "deletedBatches": deletedBatches,
    }), nil
}

func (s *PurgeEventsService) loadRetentionDays() {
    s.retentionDays = s.purgeEventsDto.GetRetentionDays()
    if s.retentionDays == 0 {
        s.retentionDays = s.environmentReaderRawRepository.GetEventsRetentionDays()
    }
}

func (s *PurgeEventsService) failIfWrongInput() error {
    if s.retentionDays >= enums.EventsLimitEnum.RETENTION_DAYS_MIN && s.retentionDays <= enums.EventsLimitEnum.RETENTION_DAYS_MAX {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(
            config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS),
            []string{enums.EventsMessageEnum.RETENTION_DAYS_MUST_BE_1_TO_3650},
        ),
    )
}
