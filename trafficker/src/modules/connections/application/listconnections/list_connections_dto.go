package listconnections

import "trafficker/src/modules/shared/application"

// ListConnectionsDto filter: terminos separados por comas, todos se cumplen ("established,remote"). Vacio = todas.
type ListConnectionsDto struct {
    filter string
}

func ListConnectionsDtoFromPrimitives(primitives map[string]any) *ListConnectionsDto {
    return &ListConnectionsDto{filter: application.GetTrimmedString(primitives, "filter")}
}

func (d *ListConnectionsDto) GetFilter() string {
    return d.filter
}
