package geteventsfacets

// FacetValuePrimitives un valor de un campo y cuantos eventos lo tienen.
type FacetValuePrimitives struct {
    Value string `json:"value"`
    // con que parametro y valor se filtra por el (agente y proyecto se ven por nombre y se filtran por uuid)
    FilterField string `json:"filter_field"`
    FilterValue string `json:"filter_value"`
    Count int    `json:"count"`
}

// GetEventsFacetsPrimitives lo que devuelve la API: por campo, su top de valores de mayor a menor.
type GetEventsFacetsPrimitives struct {
    Facets map[string][]FacetValuePrimitives `json:"facets"`
}

type GetEventsFacetsResultDto struct {
    facets map[string][]FacetValuePrimitives
}

func GetEventsFacetsResultDtoFromPrimitives(primitives map[string]any) *GetEventsFacetsResultDto {
    facets, _ := primitives["facets"].(map[string][]FacetValuePrimitives)
    if facets == nil {
        facets = map[string][]FacetValuePrimitives{}
    }
    return &GetEventsFacetsResultDto{facets: facets}
}

func (d *GetEventsFacetsResultDto) ToArray() GetEventsFacetsPrimitives {
    return GetEventsFacetsPrimitives{Facets: d.facets}
}
