package controllers

import (
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/components/httprequester"
    "siem-backend/src/modules/shared/infrastructure/components/httpresponder"
    sharedcontrollers "siem-backend/src/modules/shared/infrastructure/controllers"

    "siem-backend/src/modules/downloads/application/getdownloads"
)

// GetDownloadsRestController GET /api/v1/downloads (Bearer del dashboard)
// -> {endpoint_url, downloads: [{file_name, os, arch, version, size_bytes, sha256, url}]}
type GetDownloadsRestController struct {
    sharedcontrollers.AbstractRestController
}

func GetDownloadsRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *GetDownloadsRestController {
    return &GetDownloadsRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *GetDownloadsRestController) Invoke() *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "GetDownloadsRestController.Invoke")
    }
    getDownloadsResultDto, err := getdownloads.GetDownloadsServiceGetInstance().Invoke(
        getdownloads.GetDownloadsDtoFromPrimitives(application.MergePrimitives(c.GetHttpRequest().GetMetadata(), dashboardAuthPrimitives)),
    )
    if err != nil {
        return c.RespondError(err, "GetDownloadsRestController.Invoke")
    }
    return c.RespondOk(getDownloadsResultDto.ToArray())
}
