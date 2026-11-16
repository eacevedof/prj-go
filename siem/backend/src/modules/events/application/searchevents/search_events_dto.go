package searchevents

import (
    "siem-backend/src/modules/shared/application"

    "siem-backend/src/modules/events/domain/enums"
)

type SearchEventsDto struct {
    application.AbstractInputDto
    isDashboardAuthenticated bool
    // los parametros de filtro de la query tal cual (from, to, search, <campo>, not_<campo>): solo primitivos
    filters  map[string]any
    page     int
    pageSize int
}

func SearchEventsDtoFromPrimitives(primitives map[string]any) *SearchEventsDto {
    filters, _ := primitives["filters"].(map[string]any)
    if filters == nil {
        filters = map[string]any{}
    }
    return &SearchEventsDto{
        AbstractInputDto:         application.NewAbstractInputDto(primitives),
        isDashboardAuthenticated: application.GetInt(primitives, "isDashboardAuthenticated") == 1,
        filters:                  filters,
        page:                     getIntOrDefaultIfEmpty(primitives, "page", enums.EventsLimitEnum.PAGE_DEFAULT),
        pageSize:                 getIntOrDefaultIfEmpty(primitives, "pageSize", enums.EventsLimitEnum.PAGE_SIZE_DEFAULT),
    }
}

func (d *SearchEventsDto) IsDashboardAuthenticated() bool {
    return d.isDashboardAuthenticated
}

func (d *SearchEventsDto) GetFilters() map[string]any {
    return d.filters
}

func (d *SearchEventsDto) GetPage() int {
    return d.page
}

func (d *SearchEventsDto) GetPageSize() int {
    return d.pageSize
}

// getIntOrDefaultIfEmpty un parametro de query que no viene llega como "": cuenta como ausente.
func getIntOrDefaultIfEmpty(primitives map[string]any, key string, defaultValue int) int {
    if application.GetTrimmedString(primitives, key) == "" {
        return defaultValue
    }
    return application.GetInt(primitives, key)
}
