package controllers

import (
    "errors"

    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/components/httprequester"
    "siem-backend/src/modules/shared/infrastructure/components/httpresponder"
    sharedcontrollers "siem-backend/src/modules/shared/infrastructure/controllers"

    "siem-backend/src/modules/dashboard/application/getprojectdetail"
    "siem-backend/src/modules/dashboard/domain/exceptions"
)

// GetProjectDetailRestController GET /api/v1/projects/:project_uuid  (Bearer del dashboard) -> {project, agents}
type GetProjectDetailRestController struct {
    sharedcontrollers.AbstractRestController
}

func GetProjectDetailRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *GetProjectDetailRestController {
    return &GetProjectDetailRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *GetProjectDetailRestController) Invoke(projectUuid string) *httpresponder.HttpResponder {
    getProjectDetailResultDto, err := c.getProjectDetailResultDto(projectUuid)
    if err == nil {
        return c.RespondOk(getProjectDetailResultDto.ToArray())
    }

    var dashboardException *exceptions.DashboardException
    if errors.As(err, &dashboardException) {
        return c.RespondDomainError(dashboardException)
    }
    return c.RespondUnexpectedError(err, "GetProjectDetailRestController.Invoke")
}

func (c *GetProjectDetailRestController) getProjectDetailResultDto(projectUuid string) (*getprojectdetail.GetProjectDetailResultDto, error) {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return nil, err
    }
    return getprojectdetail.GetProjectDetailServiceGetInstance().Invoke(
        getprojectdetail.GetProjectDetailDtoFromPrimitives(application.MergePrimitives(
            c.GetHttpRequest().GetMetadata(),
            dashboardAuthPrimitives,
            map[string]any{"projectUuid": projectUuid},
        )),
    )
}
