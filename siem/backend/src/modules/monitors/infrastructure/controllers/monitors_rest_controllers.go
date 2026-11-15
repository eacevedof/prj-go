package controllers

import (
    "siem-backend/src/modules/monitors/application/createmonitor"
    "siem-backend/src/modules/monitors/application/deletemonitor"
    "siem-backend/src/modules/monitors/application/getagentconfig"
    "siem-backend/src/modules/monitors/application/getprojectmonitors"
    "siem-backend/src/modules/monitors/application/updatemonitor"
    "siem-backend/src/modules/monitors/domain/enums"
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/components/httprequester"
    "siem-backend/src/modules/shared/infrastructure/components/httpresponder"
    sharedcontrollers "siem-backend/src/modules/shared/infrastructure/controllers"
)

// Los cinco controllers REST de monitors van en este fichero: son el mismo esqueleto (autenticar -> servicio ->
// responder) y solo cambian la ruta y el caso de uso. Los del panel piden el Bearer del dashboard; el de
// agent-config, el del agente. Los errores controlados (400, 401, 404, 409) viajan tal cual (RespondError).

// getMonitorInputPrimitives el body JSON (snake_case, contrato de la API) a los primitivos del dto.
func getMonitorInputPrimitives(httpRequest *httprequester.HttpRequester) map[string]any {
    body := httpRequest.GetPost()
    return map[string]any{
        "agentUuid":         body[enums.MonitorsKeyEnum.AGENT_UUID],
        "monitorType":       body[enums.MonitorsKeyEnum.MONITOR_TYPE],
        "monitorName":       body[enums.MonitorsKeyEnum.MONITOR_NAME],
        "isEnabled":         body[enums.MonitorsKeyEnum.IS_ENABLED],
        "intervalSeconds":   body[enums.MonitorsKeyEnum.INTERVAL_SECONDS],
        "config":            body[enums.MonitorsKeyEnum.CONFIG],
        "warningThreshold":  body[enums.MonitorsKeyEnum.WARNING_THRESHOLD],
        "criticalThreshold": body[enums.MonitorsKeyEnum.CRITICAL_THRESHOLD],
    }
}

// GetProjectMonitorsRestController GET /api/v1/projects/:project_uuid/monitors -> {monitors: [monitor + states]}
type GetProjectMonitorsRestController struct {
    sharedcontrollers.AbstractRestController
}

func GetProjectMonitorsRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *GetProjectMonitorsRestController {
    return &GetProjectMonitorsRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *GetProjectMonitorsRestController) Invoke(projectUuid string) *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "GetProjectMonitorsRestController.Invoke")
    }
    getProjectMonitorsResultDto, err := getprojectmonitors.GetProjectMonitorsServiceGetInstance().Invoke(
        getprojectmonitors.GetProjectMonitorsDtoFromPrimitives(application.MergePrimitives(
            c.GetHttpRequest().GetMetadata(), dashboardAuthPrimitives, map[string]any{"projectUuid": projectUuid},
        )),
    )
    if err != nil {
        return c.RespondError(err, "GetProjectMonitorsRestController.Invoke")
    }
    return c.RespondOk(getProjectMonitorsResultDto.ToArray())
}

// CreateMonitorRestController POST /api/v1/projects/:project_uuid/monitors -> 201 el monitor
type CreateMonitorRestController struct {
    sharedcontrollers.AbstractRestController
}

func CreateMonitorRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *CreateMonitorRestController {
    return &CreateMonitorRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *CreateMonitorRestController) Invoke(projectUuid string) *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "CreateMonitorRestController.Invoke")
    }
    createMonitorResultDto, err := createmonitor.CreateMonitorServiceGetInstance().Invoke(
        createmonitor.CreateMonitorDtoFromPrimitives(application.MergePrimitives(
            c.GetHttpRequest().GetMetadata(), dashboardAuthPrimitives, getMonitorInputPrimitives(c.GetHttpRequest()),
            map[string]any{"projectUuid": projectUuid},
        )),
    )
    if err != nil {
        return c.RespondError(err, "CreateMonitorRestController.Invoke")
    }
    return c.RespondCreated(createMonitorResultDto.ToArray())
}

// UpdateMonitorRestController PUT /api/v1/monitors/:monitor_uuid -> el monitor
type UpdateMonitorRestController struct {
    sharedcontrollers.AbstractRestController
}

func UpdateMonitorRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *UpdateMonitorRestController {
    return &UpdateMonitorRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *UpdateMonitorRestController) Invoke(monitorUuid string) *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "UpdateMonitorRestController.Invoke")
    }
    updateMonitorResultDto, err := updatemonitor.UpdateMonitorServiceGetInstance().Invoke(
        updatemonitor.UpdateMonitorDtoFromPrimitives(application.MergePrimitives(
            c.GetHttpRequest().GetMetadata(), dashboardAuthPrimitives, getMonitorInputPrimitives(c.GetHttpRequest()),
            map[string]any{"monitorUuid": monitorUuid},
        )),
    )
    if err != nil {
        return c.RespondError(err, "UpdateMonitorRestController.Invoke")
    }
    return c.RespondOk(updateMonitorResultDto.ToArray())
}

// DeleteMonitorRestController DELETE /api/v1/monitors/:monitor_uuid -> {monitor_uuid}
type DeleteMonitorRestController struct {
    sharedcontrollers.AbstractRestController
}

func DeleteMonitorRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *DeleteMonitorRestController {
    return &DeleteMonitorRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *DeleteMonitorRestController) Invoke(monitorUuid string) *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "DeleteMonitorRestController.Invoke")
    }
    deleteMonitorResultDto, err := deletemonitor.DeleteMonitorServiceGetInstance().Invoke(
        deletemonitor.DeleteMonitorDtoFromPrimitives(application.MergePrimitives(
            c.GetHttpRequest().GetMetadata(), dashboardAuthPrimitives, map[string]any{"monitorUuid": monitorUuid},
        )),
    )
    if err != nil {
        return c.RespondError(err, "DeleteMonitorRestController.Invoke")
    }
    return c.RespondOk(deleteMonitorResultDto.ToArray())
}

// GetAgentConfigRestController GET /api/v1/agent-config (Bearer del agente) -> {config_hash, monitors}
type GetAgentConfigRestController struct {
    sharedcontrollers.AbstractRestController
}

func GetAgentConfigRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *GetAgentConfigRestController {
    return &GetAgentConfigRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *GetAgentConfigRestController) Invoke() *httpresponder.HttpResponder {
    authAgentPrimitives, err := c.GetAuthAgentOrFail()
    if err != nil {
        return c.RespondError(err, "GetAgentConfigRestController.Invoke")
    }
    getAgentConfigResultDto, err := getagentconfig.GetAgentConfigServiceGetInstance().Invoke(
        getagentconfig.GetAgentConfigDtoFromPrimitives(application.MergePrimitives(
            c.GetHttpRequest().GetMetadata(), authAgentPrimitives,
        )),
    )
    if err != nil {
        return c.RespondError(err, "GetAgentConfigRestController.Invoke")
    }
    return c.RespondOk(getAgentConfigResultDto.ToArray())
}
