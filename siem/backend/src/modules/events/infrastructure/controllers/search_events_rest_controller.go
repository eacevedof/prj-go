package controllers

import (
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/components/httprequester"
    "siem-backend/src/modules/shared/infrastructure/components/httpresponder"
    sharedcontrollers "siem-backend/src/modules/shared/infrastructure/controllers"

    "siem-backend/src/modules/events/domain/enums"
    "siem-backend/src/modules/events/application/geteventsfacets"
    "siem-backend/src/modules/events/application/geteventsfilteroptions"
    "siem-backend/src/modules/events/application/geteventshistogram"
    "siem-backend/src/modules/events/application/searchevents"
)

// Los tres controllers REST de busqueda del visor de eventos (Bearer del dashboard) comparten los mismos parametros:
//   from, to (rango UTC sobre occurred_at) · search (texto en el mensaje)
//   project_uuid, agent_uuid, monitor_uuid, severity, event_type, source_type, source_name -> incluir (con coma)
//   not_<campo>                                                                            -> excluir
// El cuarto, filter-options, da el catalogo de los desplegables Proyecto / Agente / Monitor.

// getFilterPrimitives los parametros de filtro de la query tal cual (clave de la API -> valor).
func getFilterPrimitives(httpRequest *httprequester.HttpRequester) map[string]any {
    filters := map[string]any{
        enums.EventsKeyEnum.FROM:   httpRequest.GetGetString(enums.EventsKeyEnum.FROM),
        enums.EventsKeyEnum.TO:     httpRequest.GetGetString(enums.EventsKeyEnum.TO),
        enums.EventsKeyEnum.SEARCH: httpRequest.GetGetString(enums.EventsKeyEnum.SEARCH),
    }
    for _, field := range enums.EventsFilterFieldEnum.VALUES {
        filters[field] = httpRequest.GetGetString(field)
        filters[enums.EventsFilterFieldEnum.NOT_PREFIX+field] = httpRequest.GetGetString(enums.EventsFilterFieldEnum.NOT_PREFIX + field)
    }
    return filters
}

// SearchEventsRestController GET /api/v1/events?<filtros>&page=&page_size=  -> {events, page, page_size, has_more}
type SearchEventsRestController struct {
    sharedcontrollers.AbstractRestController
}

func SearchEventsRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *SearchEventsRestController {
    return &SearchEventsRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *SearchEventsRestController) Invoke() *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "SearchEventsRestController.Invoke")
    }
    httpRequest := c.GetHttpRequest()
    searchEventsResultDto, err := searchevents.SearchEventsServiceGetInstance().Invoke(
        searchevents.SearchEventsDtoFromPrimitives(application.MergePrimitives(
            httpRequest.GetMetadata(),
            dashboardAuthPrimitives,
            map[string]any{
                "filters":  getFilterPrimitives(httpRequest),
                "page":     httpRequest.GetGetString(enums.EventsKeyEnum.PAGE),
                "pageSize": httpRequest.GetGetString(enums.EventsKeyEnum.PAGE_SIZE),
            },
        )),
    )
    if err != nil {
        return c.RespondError(err, "SearchEventsRestController.Invoke")
    }
    return c.RespondOk(searchEventsResultDto.ToArray())
}

// GetEventsHistogramRestController GET /api/v1/events/histogram?<filtros>
// -> {from, to, bucket_seconds, buckets: [{bucket_at, total, debug, info, warning, error, critical}]}
type GetEventsHistogramRestController struct {
    sharedcontrollers.AbstractRestController
}

func GetEventsHistogramRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *GetEventsHistogramRestController {
    return &GetEventsHistogramRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *GetEventsHistogramRestController) Invoke() *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "GetEventsHistogramRestController.Invoke")
    }
    httpRequest := c.GetHttpRequest()
    getEventsHistogramResultDto, err := geteventshistogram.GetEventsHistogramServiceGetInstance().Invoke(
        geteventshistogram.GetEventsHistogramDtoFromPrimitives(application.MergePrimitives(
            httpRequest.GetMetadata(), dashboardAuthPrimitives, map[string]any{"filters": getFilterPrimitives(httpRequest)},
        )),
    )
    if err != nil {
        return c.RespondError(err, "GetEventsHistogramRestController.Invoke")
    }
    return c.RespondOk(getEventsHistogramResultDto.ToArray())
}

// GetEventsFacetsRestController GET /api/v1/events/facets?<filtros>
// -> {facets: {<campo>: [{value, filter_field, filter_value, count}]}}
type GetEventsFacetsRestController struct {
    sharedcontrollers.AbstractRestController
}

func GetEventsFacetsRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *GetEventsFacetsRestController {
    return &GetEventsFacetsRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *GetEventsFacetsRestController) Invoke() *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "GetEventsFacetsRestController.Invoke")
    }
    httpRequest := c.GetHttpRequest()
    getEventsFacetsResultDto, err := geteventsfacets.GetEventsFacetsServiceGetInstance().Invoke(
        geteventsfacets.GetEventsFacetsDtoFromPrimitives(application.MergePrimitives(
            httpRequest.GetMetadata(), dashboardAuthPrimitives, map[string]any{"filters": getFilterPrimitives(httpRequest)},
        )),
    )
    if err != nil {
        return c.RespondError(err, "GetEventsFacetsRestController.Invoke")
    }
    return c.RespondOk(getEventsFacetsResultDto.ToArray())
}

// GetEventsFilterOptionsRestController GET /api/v1/events/filter-options
// -> {projects: [{project_uuid, project_name}], agents: [{agent_uuid, agent_name, project_uuid}],
//     monitors: [{monitor_uuid, monitor_name, monitor_type, project_uuid, agent_uuid}]}
type GetEventsFilterOptionsRestController struct {
    sharedcontrollers.AbstractRestController
}

func GetEventsFilterOptionsRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *GetEventsFilterOptionsRestController {
    return &GetEventsFilterOptionsRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *GetEventsFilterOptionsRestController) Invoke() *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "GetEventsFilterOptionsRestController.Invoke")
    }
    getEventsFilterOptionsResultDto, err := geteventsfilteroptions.GetEventsFilterOptionsServiceGetInstance().Invoke(
        geteventsfilteroptions.GetEventsFilterOptionsDtoFromPrimitives(application.MergePrimitives(
            c.GetHttpRequest().GetMetadata(), dashboardAuthPrimitives,
        )),
    )
    if err != nil {
        return c.RespondError(err, "GetEventsFilterOptionsRestController.Invoke")
    }
    return c.RespondOk(getEventsFilterOptionsResultDto.ToArray())
}
