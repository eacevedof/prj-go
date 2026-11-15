package controllers

import (
    "strings"

    "siem-backend/src/modules/agents/application/createagent"
    "siem-backend/src/modules/agents/application/deleteagent"
    "siem-backend/src/modules/agents/domain/entities"
    "siem-backend/src/modules/agents/domain/enums"
    "siem-backend/src/modules/shared/infrastructure/repositories/configuration"
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/components/httprequester"
    "siem-backend/src/modules/shared/infrastructure/components/httpresponder"
    sharedcontrollers "siem-backend/src/modules/shared/infrastructure/controllers"
)

// CreateAgentRestController POST /api/v1/projects/:project_uuid/agents (Bearer del dashboard)
// {agent_name, heartbeat_interval_seconds} -> 201 el agente + su api_token en claro (solo aqui, una vez: en BD queda
// su hash). El mismo caso de uso que ./console agent-create.
type CreateAgentRestController struct {
    sharedcontrollers.AbstractRestController
}

func CreateAgentRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *CreateAgentRestController {
    return &CreateAgentRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *CreateAgentRestController) Invoke(projectUuid string) *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "CreateAgentRestController.Invoke")
    }
    httpRequest := c.GetHttpRequest()
    createAgentResultDto, err := createagent.CreateAgentServiceGetInstance().Invoke(
        createagent.CreateAgentDtoFromPrimitives(application.MergePrimitives(
            httpRequest.GetMetadata(), dashboardAuthPrimitives,
            map[string]any{
                "projectUuid":              projectUuid,
                "agentName":                httpRequest.GetPostString(enums.AgentsKeyEnum.AGENT_NAME),
                "heartbeatIntervalSeconds": httpRequest.GetPostString(enums.AgentsKeyEnum.HEARTBEAT_INTERVAL_SECONDS),
            },
        )),
    )
    if err != nil {
        return c.RespondError(err, "CreateAgentRestController.Invoke")
    }
    // endpoint_url: lo que el panel pone en el .env que se copia junto al token (la misma url que da /downloads)
    return c.RespondCreated(createdAgentPrimitives{
        AgentWithApiTokenPrimitives: createAgentResultDto.ToArray(),
        EndpointUrl:                 strings.TrimRight(configuration.GetInstance().GetBaseUrl(), "/") + ingestPath,
    })
}

const ingestPath = "/api/v1/ingest-batches"

// createdAgentPrimitives la respuesta del alta por REST: el agente, su api_token (una vez) y a donde tiene que enviar.
type createdAgentPrimitives struct {
    entities.AgentWithApiTokenPrimitives
    EndpointUrl string `json:"endpoint_url"`
}

// DeleteAgentRestController DELETE /api/v1/agents/:agent_uuid (Bearer del dashboard) -> {agent_uuid}
type DeleteAgentRestController struct {
    sharedcontrollers.AbstractRestController
}

func DeleteAgentRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *DeleteAgentRestController {
    return &DeleteAgentRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *DeleteAgentRestController) Invoke(agentUuid string) *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "DeleteAgentRestController.Invoke")
    }
    deleteAgentResultDto, err := deleteagent.DeleteAgentServiceGetInstance().Invoke(
        deleteagent.DeleteAgentDtoFromPrimitives(application.MergePrimitives(
            c.GetHttpRequest().GetMetadata(), dashboardAuthPrimitives, map[string]any{"agentUuid": agentUuid},
        )),
    )
    if err != nil {
        return c.RespondError(err, "DeleteAgentRestController.Invoke")
    }
    return c.RespondOk(deleteAgentResultDto.ToArray())
}
