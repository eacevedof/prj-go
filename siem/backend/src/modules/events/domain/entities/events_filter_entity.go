package entities

import (
    "slices"
    "strings"

    "siem-backend/src/modules/shared/application"

    "siem-backend/src/modules/events/domain/enums"
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
    eventsFilterEntity := &EventsFilterEntity{
        from:                  application.GetTrimmedString(primitives, enums.EventsKeyEnum.FROM),
        to:                    application.GetTrimmedString(primitives, enums.EventsKeyEnum.TO),
        search:                application.GetTrimmedString(primitives, enums.EventsKeyEnum.SEARCH),
        includedValuesByField: map[string][]string{},
        excludedValuesByField: map[string][]string{},
    }
    eventsFilterEntity.loadValuesByField(primitives)
    return eventsFilterEntity
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

// loadValuesByField por cada campo filtrable, los valores a incluir (<campo>) y a excluir (not_<campo>).
func (e *EventsFilterEntity) loadValuesByField(primitives map[string]any) {
    for _, field := range enums.EventsFilterFieldEnum.VALUES {
        if values := e.getValues(application.GetTrimmedString(primitives, field)); len(values) > 0 {
            e.includedValuesByField[field] = values
        }
        if values := e.getValues(application.GetTrimmedString(primitives, enums.EventsFilterFieldEnum.NOT_PREFIX+field)); len(values) > 0 {
            e.excludedValuesByField[field] = values
        }
    }
}

// getValues "error, critical,," -> ["error", "critical"] (sin vacios ni repetidos).
func (e *EventsFilterEntity) getValues(rawValues string) []string {
    values := []string{}
    for _, rawValue := range strings.Split(rawValues, enums.EventsFilterFieldEnum.VALUES_SEPARATOR) {
        value := strings.TrimSpace(rawValue)
        if value == "" || slices.Contains(values, value) {
            continue
        }
        values = append(values, value)
    }
    return values
}
