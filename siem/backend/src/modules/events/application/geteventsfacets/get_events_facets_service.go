package geteventsfacets

import (
    "time"

    "siem-backend/src/core/config"

    sharedenums "siem-backend/src/modules/shared/domain/enums"
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/traits"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"

    dashboardexceptions "siem-backend/src/modules/dashboard/domain/exceptions"

    "siem-backend/src/modules/events/domain/enums"
    "siem-backend/src/modules/events/domain/entities"
    "siem-backend/src/modules/events/domain/services"
    "siem-backend/src/modules/events/infrastructure/repositories"
)

// GetEventsFacetsService por cada campo del panel lateral del visor (gravedad, tipo, fuente, agente, proyecto), sus
// valores mas repetidos con su recuento, con los mismos filtros que la lista. Sin rango, las ultimas 24 h (como el
// histograma: contar sobre toda la tabla en cada refresco no escala).
type GetEventsFacetsService struct {
    traits.ErrorStackTrait

    eventsSearchValidator        *services.EventsSearchValidator
    histogramBucketCalculator    *services.HistogramBucketCalculator
    eventsReaderSqliteRepository *repositories.EventsReaderSqliteRepository
}

func GetEventsFacetsServiceGetInstance() *GetEventsFacetsService {
    return &GetEventsFacetsService{
        eventsSearchValidator:        services.EventsSearchValidatorGetInstance(),
        histogramBucketCalculator:    services.HistogramBucketCalculatorGetInstance(),
        eventsReaderSqliteRepository: repositories.EventsReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *DashboardException (401), *FieldsException (400) o uno inesperado.
func (s *GetEventsFacetsService) Invoke(getEventsFacetsDto *GetEventsFacetsDto) (*GetEventsFacetsResultDto, error) {
    if !getEventsFacetsDto.IsDashboardAuthenticated() {
        return nil, dashboardexceptions.UnauthorizedInvalidToken()
    }
    eventsFilterEntity := entities.EventsFilterEntityFromPrimitives(getEventsFacetsDto.GetFilters())
    if errors := s.eventsSearchValidator.GetFilterErrors(eventsFilterEntity); len(errors) > 0 {
        return nil, sharedexceptions.FieldsExceptionBadRequestCustom(
            s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
        )
    }

    fromTime, toTime := s.histogramBucketCalculator.GetRange(eventsFilterEntity.GetFrom(), eventsFilterEntity.GetTo(), time.Now())
    rangedEventsFilterEntity := eventsFilterEntity.GetWithRange(
        s.histogramBucketCalculator.GetFormatted(fromTime.Unix()), s.histogramBucketCalculator.GetFormatted(toTime.Unix()),
    )

    facets := map[string][]FacetValuePrimitives{}
    for _, facetField := range enums.EventsFacetFieldEnum.VALUES {
        rows, err := s.eventsReaderSqliteRepository.GetTopValuesByFacetField(rangedEventsFilterEntity, facetField, enums.EventsFacetFieldEnum.TOP)
        if err != nil {
            return nil, err
        }
        facetValues := make([]FacetValuePrimitives, 0, len(rows))
        for _, row := range rows {
            facetValues = append(facetValues, FacetValuePrimitives{
                Value:       application.GetTrimmedString(row, enums.EventsKeyEnum.FACET_VALUE),
                FilterField: s.eventsReaderSqliteRepository.GetFilterFieldByFacetField(facetField),
                FilterValue: application.GetTrimmedString(row, enums.EventsKeyEnum.FILTER_VALUE),
                Count:       application.GetInt(row, enums.EventsKeyEnum.EVENTS_COUNT),
            })
        }
        facets[facetField] = facetValues
    }
    return GetEventsFacetsResultDtoFromPrimitives(map[string]any{"facets": facets}), nil
}
