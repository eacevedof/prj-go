package controllers

import (
    "errors"

    "siem-backend/src/modules/dashboard/application/getdashboard"
    "siem-backend/src/modules/dashboard/domain/exceptions"
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/components/httprequester"
    "siem-backend/src/modules/shared/infrastructure/components/httpresponder"
    sharedcontrollers "siem-backend/src/modules/shared/infrastructure/controllers"
)

// GetDashboardRestController GET /api/v1/dashboard  (Bearer del dashboard) -> {server, projects}
type GetDashboardRestController struct {
    sharedcontrollers.AbstractRestController
}

func GetDashboardRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *GetDashboardRestController {
    return &GetDashboardRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *GetDashboardRestController) Invoke() *httpresponder.HttpResponder {
    getDashboardResultDto, err := c.getDashboardResultDto()
    if err == nil {
        return c.RespondOk(getDashboardResultDto.ToArray())
    }

    var dashboardException *exceptions.DashboardException
    if errors.As(err, &dashboardException) {
        return c.RespondDomainError(dashboardException)
    }
    return c.RespondUnexpectedError(err, "GetDashboardRestController.Invoke")
}

func (c *GetDashboardRestController) getDashboardResultDto() (*getdashboard.GetDashboardResultDto, error) {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return nil, err
    }
    return getdashboard.GetDashboardServiceGetInstance().Invoke(
        getdashboard.GetDashboardDtoFromPrimitives(application.MergePrimitives(
            c.GetHttpRequest().GetMetadata(),
            dashboardAuthPrimitives,
        )),
    )
}
