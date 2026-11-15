package geteventsfacets

import "siem-backend/src/modules/shared/application"

// GetEventsFacetsDto GET /api/v1/events/facets: los mismos filtros que la lista, sin paginacion.
type GetEventsFacetsDto struct {
    application.AbstractInputDto
    isDashboardAuthenticated bool
    // los parametros de filtro de la query tal cual: solo primitivos
    filters map[string]any
}

func GetEventsFacetsDtoFromPrimitives(primitives map[string]any) *GetEventsFacetsDto {
    filters, _ := primitives["filters"].(map[string]any)
    if filters == nil {
        filters = map[string]any{}
    }
    return &GetEventsFacetsDto{
        AbstractInputDto:         application.NewAbstractInputDto(primitives),
        isDashboardAuthenticated: application.GetInt(primitives, "isDashboardAuthenticated") == 1,
        filters:                  filters,
    }
}

func (d *GetEventsFacetsDto) IsDashboardAuthenticated() bool {
    return d.isDashboardAuthenticated
}

func (d *GetEventsFacetsDto) GetFilters() map[string]any {
    return d.filters
}
