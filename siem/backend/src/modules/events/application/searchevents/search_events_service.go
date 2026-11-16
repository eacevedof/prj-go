package searchevents

import (
    "siem-backend/src/core/config"

    sharedenums "siem-backend/src/modules/shared/domain/enums"
    "siem-backend/src/modules/shared/infrastructure/traits"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"

    dashboardexceptions "siem-backend/src/modules/dashboard/domain/exceptions"

    "siem-backend/src/modules/events/domain/entities"
    "siem-backend/src/modules/events/domain/services"
    "siem-backend/src/modules/events/infrastructure/repositories"
)

// SearchEventsService una pagina de eventos (los mas recientes primero) con filtros opcionales: rango de fechas,
// texto, e incluir / excluir valores por campo. Para el visor del front.
type SearchEventsService struct {
    traits.ErrorStackTrait

    eventsSearchValidator        *services.EventsSearchValidator
    eventsReaderSqliteRepository *repositories.EventsReaderSqliteRepository

    searchEventsDto *SearchEventsDto
}

func SearchEventsServiceGetInstance() *SearchEventsService {
    return &SearchEventsService{
        eventsSearchValidator:        services.EventsSearchValidatorGetInstance(),
        eventsReaderSqliteRepository: repositories.EventsReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *DashboardException (401), *FieldsException (400) o uno inesperado.
func (s *SearchEventsService) Invoke(searchEventsDto *SearchEventsDto) (*SearchEventsResultDto, error) {
    s.searchEventsDto = searchEventsDto

    eventsFilterEntity := entities.EventsFilterEntityFromPrimitives(searchEventsDto.GetFilters())
    if err := s.failIfWrongInput(eventsFilterEntity); err != nil {
        return nil, err
    }

    // una fila de mas: si llega, hay pagina siguiente (sin un COUNT(*) sobre toda la tabla)
    pageSize := searchEventsDto.GetPageSize()
    events, err := s.eventsReaderSqliteRepository.GetEventsBySearch(eventsFilterEntity, pageSize+1, (searchEventsDto.GetPage()-1)*pageSize)
    if err != nil {
        return nil, err
    }

    hasMore := 0
    if len(events) > pageSize {
        hasMore = 1
        events = events[:pageSize]
    }
    return SearchEventsResultDtoFromPrimitives(map[string]any{
        "page":     searchEventsDto.GetPage(),
        "pageSize": pageSize,
        "hasMore":  hasMore,
        "events":   events,
    }), nil
}

func (s *SearchEventsService) failIfWrongInput(eventsFilterEntity *entities.EventsFilterEntity) error {
    if !s.searchEventsDto.IsDashboardAuthenticated() {
        return dashboardexceptions.UnauthorizedInvalidToken()
    }

    errors := append(
        s.eventsSearchValidator.GetFilterErrors(eventsFilterEntity),
        s.eventsSearchValidator.GetPaginationErrors(s.searchEventsDto.GetPage(), s.searchEventsDto.GetPageSize())...,
    )
    if len(errors) == 0 {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
    )
}
