package geteventshistogram

import (
    "time"

    "siem-backend/src/core/config"
    dashboardexceptions "siem-backend/src/modules/dashboard/domain/exceptions"
    "siem-backend/src/modules/events/domain/entities"
    "siem-backend/src/modules/events/domain/enums"
    "siem-backend/src/modules/events/domain/services"
    "siem-backend/src/modules/events/infrastructure/repositories"
    "siem-backend/src/modules/shared/application"
    sharedenums "siem-backend/src/modules/shared/domain/enums"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"
    "siem-backend/src/modules/shared/infrastructure/traits"
)

// GetEventsHistogramService cuantos eventos hay por tramo de tiempo y gravedad, con los mismos filtros que la lista:
// las barras apiladas del visor. Sin rango, las ultimas 24 h. El ancho del tramo se ajusta a ~60 barras.
type GetEventsHistogramService struct {
    traits.ErrorStackTrait

    eventsSearchValidator        *services.EventsSearchValidator
    histogramBucketCalculator    *services.HistogramBucketCalculator
    eventsReaderSqliteRepository *repositories.EventsReaderSqliteRepository
}

func GetEventsHistogramServiceGetInstance() *GetEventsHistogramService {
    return &GetEventsHistogramService{
        eventsSearchValidator:        services.EventsSearchValidatorGetInstance(),
        histogramBucketCalculator:    services.HistogramBucketCalculatorGetInstance(),
        eventsReaderSqliteRepository: repositories.EventsReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *DashboardException (401), *FieldsException (400) o uno inesperado.
func (s *GetEventsHistogramService) Invoke(getEventsHistogramDto *GetEventsHistogramDto) (*GetEventsHistogramResultDto, error) {
    if !getEventsHistogramDto.IsDashboardAuthenticated() {
        return nil, dashboardexceptions.UnauthorizedInvalidToken()
    }
    eventsFilterEntity := entities.EventsFilterEntityFromPrimitives(getEventsHistogramDto.GetFilters())
    if errors := s.eventsSearchValidator.GetFilterErrors(eventsFilterEntity); len(errors) > 0 {
        return nil, sharedexceptions.FieldsExceptionBadRequestCustom(
            s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
        )
    }

    fromTime, toTime := s.histogramBucketCalculator.GetRange(eventsFilterEntity.GetFrom(), eventsFilterEntity.GetTo(), time.Now())
    from := s.histogramBucketCalculator.GetFormatted(fromTime.Unix())
    to := s.histogramBucketCalculator.GetFormatted(toTime.Unix())
    bucketSeconds := s.histogramBucketCalculator.GetBucketSeconds(fromTime, toTime)

    rows, err := s.eventsReaderSqliteRepository.GetEventsCountByBucketAndSeverity(eventsFilterEntity.GetWithRange(from, to), bucketSeconds)
    if err != nil {
        return nil, err
    }
    return GetEventsHistogramResultDtoFromPrimitives(map[string]any{
        "from":          from,
        "to":            to,
        "bucketSeconds": bucketSeconds,
        "buckets":       s.getBuckets(rows, s.histogramBucketCalculator.GetBucketStarts(fromTime, toTime, bucketSeconds)),
    }), nil
}

// getBuckets todos los tramos del rango en orden; los que no traen filas quedan a cero.
func (s *GetEventsHistogramService) getBuckets(rows []map[string]any, bucketStarts []int64) []HistogramBucketPrimitives {
    bucketsByStart := map[int64]*HistogramBucketPrimitives{}
    buckets := make([]HistogramBucketPrimitives, len(bucketStarts))
    for i, bucketStart := range bucketStarts {
        buckets[i].BucketAt = s.histogramBucketCalculator.GetFormatted(bucketStart)
        bucketsByStart[bucketStart] = &buckets[i]
    }

    for _, row := range rows {
        bucket := bucketsByStart[int64(application.GetInt(row, "bucket_unix"))]
        if bucket == nil {
            continue
        }
        eventsCount := application.GetInt(row, "n_events")
        bucket.Total += eventsCount
        switch application.GetTrimmedString(row, enums.EventsKeyEnum.SEVERITY) {
        case enums.EventSeverityEnum.DEBUG:
            bucket.Debug += eventsCount
        case enums.EventSeverityEnum.INFO:
            bucket.Info += eventsCount
        case enums.EventSeverityEnum.WARNING:
            bucket.Warning += eventsCount
        case enums.EventSeverityEnum.ERROR:
            bucket.Error += eventsCount
        case enums.EventSeverityEnum.CRITICAL:
            bucket.Critical += eventsCount
        }
    }
    return buckets
}
