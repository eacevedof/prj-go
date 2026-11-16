package services

import (
    "slices"
    "time"
    "unicode/utf8"

    "siem-backend/src/core/config"

    "siem-backend/src/modules/events/domain/enums"
    "siem-backend/src/modules/events/domain/entities"
)

// EventsSearchValidator reglas de los filtros de eventos (lista, histograma y facetas) y de la paginacion de la lista.
// Devuelve la pila de errores (vacia = todo bien).
type EventsSearchValidator struct{}

func EventsSearchValidatorGetInstance() *EventsSearchValidator {
    return &EventsSearchValidator{}
}

// GetFilterErrors un campo sin valores = sin ese filtro. Los valores de severity, event_type y source_type son listas
// cerradas tanto al incluir como al excluir; los uuid y source_name son texto libre (van como parametro del SQL).
func (v *EventsSearchValidator) GetFilterErrors(eventsFilterEntity *entities.EventsFilterEntity) []string {
    errors := []string{}
    for _, closedField := range []struct {
        field         string
        allowedValues []string
        message       string
    }{
        {enums.EventsFilterFieldEnum.SEVERITY, enums.EventSeverityEnum.VALUES, enums.EventsMessageEnum.SEVERITY_IS_NOT_VALID},
        {enums.EventsFilterFieldEnum.EVENT_TYPE, enums.EventTypeEnum.VALUES, enums.EventsMessageEnum.EVENT_TYPE_IS_NOT_VALID},
        {enums.EventsFilterFieldEnum.SOURCE_TYPE, enums.SourceTypeEnum.VALUES, enums.EventsMessageEnum.SOURCE_TYPE_IS_NOT_VALID},
    } {
        values := append(slices.Clone(eventsFilterEntity.GetIncludedValues(closedField.field)), eventsFilterEntity.GetExcludedValues(closedField.field)...)
        if slices.ContainsFunc(values, func(value string) bool { return !slices.Contains(closedField.allowedValues, value) }) {
            errors = append(errors, closedField.message)
        }
    }
    for _, field := range enums.EventsFilterFieldEnum.VALUES {
        if len(eventsFilterEntity.GetIncludedValues(field)) > enums.EventsFilterFieldEnum.MAX_VALUES_PER_FIELD {
            errors = append(errors, config.Tr(enums.EventsMessageEnum.X_MAX_VALUES_IS_20, field))
        }
        if len(eventsFilterEntity.GetExcludedValues(field)) > enums.EventsFilterFieldEnum.MAX_VALUES_PER_FIELD {
            errors = append(errors, config.Tr(enums.EventsMessageEnum.X_MAX_VALUES_IS_20, enums.EventsFilterFieldEnum.NOT_PREFIX+field))
        }
    }
    if utf8.RuneCountInString(eventsFilterEntity.GetSearch()) > enums.EventsLimitEnum.SEARCH_MAX_LENGTH {
        errors = append(errors, enums.EventsMessageEnum.SEARCH_MAX_LENGTH_IS_200)
    }
    return append(errors, v.getRangeErrors(eventsFilterEntity.GetFrom(), eventsFilterEntity.GetTo())...)
}

func (v *EventsSearchValidator) GetPaginationErrors(page int, pageSize int) []string {
    errors := []string{}
    if page < enums.EventsLimitEnum.PAGE_MIN {
        errors = append(errors, enums.EventsMessageEnum.PAGE_MUST_BE_1_OR_GREATER)
    }
    if pageSize < enums.EventsLimitEnum.PAGE_SIZE_MIN || pageSize > enums.EventsLimitEnum.PAGE_SIZE_MAX {
        errors = append(errors, enums.EventsMessageEnum.PAGE_SIZE_MUST_BE_1_TO_200)
    }
    return errors
}

// getRangeErrors from y to son opcionales; si vienen, con el formato de la BD (se comparan como texto) y en orden.
func (v *EventsSearchValidator) getRangeErrors(from string, to string) []string {
    errors := []string{}
    fromTime, fromErr := time.Parse(dateTimeFormat, from)
    toTime, toErr := time.Parse(dateTimeFormat, to)
    if from != "" && fromErr != nil {
        errors = append(errors, config.Tr(enums.EventsMessageEnum.X_MUST_BE_A_UTC_DATE_TIME, enums.EventsKeyEnum.FROM))
    }
    if to != "" && toErr != nil {
        errors = append(errors, config.Tr(enums.EventsMessageEnum.X_MUST_BE_A_UTC_DATE_TIME, enums.EventsKeyEnum.TO))
    }
    if from != "" && to != "" && fromErr == nil && toErr == nil && fromTime.After(toTime) {
        errors = append(errors, enums.EventsMessageEnum.FROM_MUST_BE_BEFORE_TO)
    }
    return errors
}
