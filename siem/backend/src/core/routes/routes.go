package routes

import (
    "siem-backend/src/core/boot"
    "siem-backend/src/modules/shared/domain/enums"
    "siem-backend/src/modules/shared/infrastructure/components/httprequester"
    "siem-backend/src/modules/shared/infrastructure/components/httpresponder"
    sharedcontrollers "siem-backend/src/modules/shared/infrastructure/controllers"
    agentscontrollers "siem-backend/src/modules/agents/infrastructure/controllers"
    dashboardcontrollers "siem-backend/src/modules/dashboard/infrastructure/controllers"
    downloadscontrollers "siem-backend/src/modules/downloads/infrastructure/controllers"
    eventscontrollers "siem-backend/src/modules/events/infrastructure/controllers"
    ingestcontrollers "siem-backend/src/modules/ingest/infrastructure/controllers"
    monitorscontrollers "siem-backend/src/modules/monitors/infrastructure/controllers"
    projectscontrollers "siem-backend/src/modules/projects/infrastructure/controllers"
)

// GetRoutes rutas de la API (ApiRouter). Url con tags ":nombre" -> args con ese nombre para el controller.
// Allowed: metodos HTTP en minusculas (HEAD va implicito en GET). Una misma url puede repetirse con otro metodo.
// Auth: no lo decide el router. Cada controller autentica con GetAuthAgentOrFail() (Bearer del agente: la ingesta y
// su configuracion) o GetDashboardAuthOrFail() (Bearer del dashboard: todo lo que lee y escribe el panel).
// El front (public/) y los binarios (public/downloads/) los sirve PublicFilesServer, no pasan por aqui.
func GetRoutes() []boot.Route {
    return []boot.Route{
        {Url: "/health", Allowed: []string{enums.HttpEnum.METHOD_GET}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return sharedcontrollers.HealthControllerGetInstance(httpRequester).Invoke()
        }},

        //@agente (siem-client -> backend)
        {Url: "/api/v1/ingest-batches", Allowed: []string{enums.HttpEnum.METHOD_POST}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return ingestcontrollers.EnqueueIngestBatchRestControllerGetInstance(httpRequester).Invoke()
        }},
        {Url: "/api/v1/agent-config", Allowed: []string{enums.HttpEnum.METHOD_GET}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return monitorscontrollers.GetAgentConfigRestControllerGetInstance(httpRequester).Invoke()
        }},

        //@panel: portada y proyectos
        {Url: "/api/v1/dashboard", Allowed: []string{enums.HttpEnum.METHOD_GET}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return dashboardcontrollers.GetDashboardRestControllerGetInstance(httpRequester).Invoke()
        }},
        {Url: "/api/v1/projects", Allowed: []string{enums.HttpEnum.METHOD_POST}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return projectscontrollers.CreateProjectRestControllerGetInstance(httpRequester).Invoke()
        }},
        {Url: "/api/v1/projects/:project_uuid", Allowed: []string{enums.HttpEnum.METHOD_GET}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return dashboardcontrollers.GetProjectDetailRestControllerGetInstance(httpRequester).Invoke(args["project_uuid"])
        }},
        {Url: "/api/v1/projects/:project_uuid", Allowed: []string{enums.HttpEnum.METHOD_PUT}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return projectscontrollers.UpdateProjectRestControllerGetInstance(httpRequester).Invoke(args["project_uuid"])
        }},
        {Url: "/api/v1/projects/:project_uuid", Allowed: []string{enums.HttpEnum.METHOD_DELETE}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return projectscontrollers.DeleteProjectRestControllerGetInstance(httpRequester).Invoke(args["project_uuid"])
        }},

        //@panel: agentes
        {Url: "/api/v1/projects/:project_uuid/agents", Allowed: []string{enums.HttpEnum.METHOD_POST}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return agentscontrollers.CreateAgentRestControllerGetInstance(httpRequester).Invoke(args["project_uuid"])
        }},
        {Url: "/api/v1/agents/:agent_uuid", Allowed: []string{enums.HttpEnum.METHOD_PUT}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return agentscontrollers.UpdateAgentRestControllerGetInstance(httpRequester).Invoke(args["agent_uuid"])
        }},
        {Url: "/api/v1/agents/:agent_uuid", Allowed: []string{enums.HttpEnum.METHOD_DELETE}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return agentscontrollers.DeleteAgentRestControllerGetInstance(httpRequester).Invoke(args["agent_uuid"])
        }},

        //@panel: monitores (recursos a vigilar)
        {Url: "/api/v1/projects/:project_uuid/monitors", Allowed: []string{enums.HttpEnum.METHOD_GET}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return monitorscontrollers.GetProjectMonitorsRestControllerGetInstance(httpRequester).Invoke(args["project_uuid"])
        }},
        {Url: "/api/v1/projects/:project_uuid/monitors", Allowed: []string{enums.HttpEnum.METHOD_POST}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return monitorscontrollers.CreateMonitorRestControllerGetInstance(httpRequester).Invoke(args["project_uuid"])
        }},
        {Url: "/api/v1/monitors/:monitor_uuid", Allowed: []string{enums.HttpEnum.METHOD_PUT}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return monitorscontrollers.UpdateMonitorRestControllerGetInstance(httpRequester).Invoke(args["monitor_uuid"])
        }},
        {Url: "/api/v1/monitors/:monitor_uuid", Allowed: []string{enums.HttpEnum.METHOD_DELETE}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return monitorscontrollers.DeleteMonitorRestControllerGetInstance(httpRequester).Invoke(args["monitor_uuid"])
        }},

        //@panel: visor de eventos
        {Url: "/api/v1/events", Allowed: []string{enums.HttpEnum.METHOD_GET}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return eventscontrollers.SearchEventsRestControllerGetInstance(httpRequester).Invoke()
        }},
        {Url: "/api/v1/events/histogram", Allowed: []string{enums.HttpEnum.METHOD_GET}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return eventscontrollers.GetEventsHistogramRestControllerGetInstance(httpRequester).Invoke()
        }},
        {Url: "/api/v1/events/facets", Allowed: []string{enums.HttpEnum.METHOD_GET}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return eventscontrollers.GetEventsFacetsRestControllerGetInstance(httpRequester).Invoke()
        }},
        {Url: "/api/v1/events/filter-options", Allowed: []string{enums.HttpEnum.METHOD_GET}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return eventscontrollers.GetEventsFilterOptionsRestControllerGetInstance(httpRequester).Invoke()
        }},

        //@panel: binarios del siem-client
        {Url: "/api/v1/downloads", Allowed: []string{enums.HttpEnum.METHOD_GET}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return downloadscontrollers.GetDownloadsRestControllerGetInstance(httpRequester).Invoke()
        }},
    }
}
