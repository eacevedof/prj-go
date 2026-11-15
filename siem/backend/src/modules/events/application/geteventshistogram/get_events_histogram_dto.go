package geteventshistogram

import "siem-backend/src/modules/shared/application"

// GetEventsHistogramDto GET /api/v1/events/histogram: los mismos filtros que la lista, sin paginacion.
type GetEventsHistogramDto struct {
    application.AbstractInputDto
    isDashboardAuthenticated bool
    // los parametros de filtro de la query tal cual: solo primitivos
    filters map[string]any
}

func GetEventsHistogramDtoFromPrimitives(primitives map[string]any) *GetEventsHistogramDto {
    filters, _ := primitives["filters"].(map[string]any)
    if filters == nil {
        filters = map[string]any{}
    }
    return &GetEventsHistogramDto{
        AbstractInputDto:         application.NewAbstractInputDto(primitives),
        isDashboardAuthenticated: application.GetInt(primitives, "isDashboardAuthenticated") == 1,
        filters:                  filters,
    }
}

func (d *GetEventsHistogramDto) IsDashboardAuthenticated() bool {
    return d.isDashboardAuthenticated
}

func (d *GetEventsHistogramDto) GetFilters() map[string]any {
    return d.filters
}
