package searchevents

import "siem-backend/src/modules/shared/application"

// SearchEventsPrimitives la pagina de eventos tal cual sale por la API (claves en este orden).
type SearchEventsPrimitives struct {
    Page     int              `json:"page"`
    PageSize int              `json:"page_size"`
    HasMore  bool             `json:"has_more"`
    Events   []map[string]any `json:"events"`
}

type SearchEventsResultDto struct {
    page     int
    pageSize int
    hasMore  bool
    // filas de EventsReaderSqliteRepository.GetEventsBySearch (solo primitivos)
    events []map[string]any
}

func SearchEventsResultDtoFromPrimitives(primitives map[string]any) *SearchEventsResultDto {
    events, _ := primitives["events"].([]map[string]any)
    if events == nil {
        events = []map[string]any{}
    }
    return &SearchEventsResultDto{
        page:     application.GetInt(primitives, "page"),
        pageSize: application.GetInt(primitives, "pageSize"),
        hasMore:  application.GetInt(primitives, "hasMore") == 1,
        events:   events,
    }
}

func (d *SearchEventsResultDto) GetPage() int {
    return d.page
}

func (d *SearchEventsResultDto) GetPageSize() int {
    return d.pageSize
}

func (d *SearchEventsResultDto) HasMore() bool {
    return d.hasMore
}

func (d *SearchEventsResultDto) GetEvents() []map[string]any {
    return d.events
}

func (d *SearchEventsResultDto) ToArray() SearchEventsPrimitives {
    return SearchEventsPrimitives{Page: d.page, PageSize: d.pageSize, HasMore: d.hasMore, Events: d.events}
}
