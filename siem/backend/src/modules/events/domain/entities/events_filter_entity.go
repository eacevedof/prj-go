package entities

import (
    "strings"

    "siem-backend/src/modules/events/domain/enums"
    "siem-backend/src/modules/shared/application"
)

// EventsFilterEntity los filtros de la busqueda de eventos, los mismos para la lista, el histograma y las facetas:
// rango de fechas sobre occurred_at, texto en el mensaje y, por cada campo filtrable (EventsFilterFieldEnum), los
// valores que se incluyen y los que se excluyen.
type EventsFilterEntity struct {
    from                  string
    to                    string
    search                string
    includedValuesByField map[string][]string
    excludedValuesByField map[string][]string
}

// EventsFilterEntityFromPrimitives primitives: los parametros de la query tal cual (claves de la API). Un campo admite
// varios valores separados por coma (severity=error,critical); "not_<campo>" los excluye.
func EventsFilterEntityFromPrimitives(primitives map[string]any) *EventsFilterEntity {
    includedValuesByField := map[string][]string{}
    excludedValuesByField := map[string][]string{}
    for _, field := range enums.EventsFilterFieldEnum.VALUES {
        if values := getValues(application.GetTrimmedString(primitives, field)); len(values) > 0 {
            includedValuesByField[field] = values
        }
        if values := getValues(application.GetTrimmedString(primitives, enums.EventsFilterFieldEnum.NOT_PREFIX+field)); len(values) > 0 {
            excludedValuesByField[field] = values
        }
    }
    return &EventsFilterEntity{
        from:                  application.GetTrimmedString(primitives, enums.EventsKeyEnum.FROM),
        to:                    application.GetTrimmedString(primitives, enums.EventsKeyEnum.TO),
        search:                application.GetTrimmedString(primitives, enums.EventsKeyEnum.SEARCH),
        includedValuesByField: includedValuesByField,
        excludedValuesByField: excludedValuesByField,
    }
}

// GetFrom "" = sin limite inferior.
func (e *EventsFilterEntity) GetFrom() string {
    return e.from
}

// GetTo "" = sin limite superior.
func (e *EventsFilterEntity) GetTo() string {
    return e.to
}

func (e *EventsFilterEntity) GetSearch() string {
    return e.search
}

func (e *EventsFilterEntity) GetIncludedValues(field string) []string {
    return e.includedValuesByField[field]
}

func (e *EventsFilterEntity) GetExcludedValues(field string) []string {
    return e.excludedValuesByField[field]
}

// GetWithRange una copia con otro rango (el histograma y las facetas, sin rango, miran las ultimas 24 h).
func (e *EventsFilterEntity) GetWithRange(from string, to string) *EventsFilterEntity {
    eventsFilterEntity := *e
    eventsFilterEntity.from = from
    eventsFilterEntity.to = to
    return &eventsFilterEntity
}

// getValues "error, critical,," -> ["error", "critical"] (sin vacios ni repetidos).
func getValues(rawValues string) []string {
    values := []string{}
    for _, rawValue := range strings.Split(rawValues, ",") {
        value := strings.TrimSpace(rawValue)
        if value == "" || containsValue(values, value) {
            continue
        }
        values = append(values, value)
    }
    return values
}

func containsValue(values []string, value string) bool {
    for _, existingValue := range values {
        if existingValue == value {
            return true
        }
    }
    return false
}
