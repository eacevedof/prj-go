package controllers

import (
    "siem-backend/src/modules/projects/application/createproject"
    "siem-backend/src/modules/projects/application/deleteproject"
    "siem-backend/src/modules/projects/domain/enums"
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/components/httprequester"
    "siem-backend/src/modules/shared/infrastructure/components/httpresponder"
    sharedcontrollers "siem-backend/src/modules/shared/infrastructure/controllers"
)

// CreateProjectRestController POST /api/v1/projects (Bearer del dashboard) {project_name, project_description}
// -> 201 el proyecto. El mismo caso de uso que ./console project-create.
type CreateProjectRestController struct {
    sharedcontrollers.AbstractRestController
}

func CreateProjectRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *CreateProjectRestController {
    return &CreateProjectRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *CreateProjectRestController) Invoke() *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "CreateProjectRestController.Invoke")
    }
    httpRequest := c.GetHttpRequest()
    createProjectResultDto, err := createproject.CreateProjectServiceGetInstance().Invoke(
        createproject.CreateProjectDtoFromPrimitives(application.MergePrimitives(
            httpRequest.GetMetadata(), dashboardAuthPrimitives,
            map[string]any{
                "projectName":        httpRequest.GetPostString(enums.ProjectsKeyEnum.PROJECT_NAME),
                "projectDescription": httpRequest.GetPostString(enums.ProjectsKeyEnum.PROJECT_DESCRIPTION),
            },
        )),
    )
    if err != nil {
        return c.RespondError(err, "CreateProjectRestController.Invoke")
    }
    return c.RespondCreated(createProjectResultDto.ToArray())
}

// DeleteProjectRestController DELETE /api/v1/projects/:project_uuid (Bearer del dashboard) -> {project_uuid}
type DeleteProjectRestController struct {
    sharedcontrollers.AbstractRestController
}

func DeleteProjectRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *DeleteProjectRestController {
    return &DeleteProjectRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *DeleteProjectRestController) Invoke(projectUuid string) *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "DeleteProjectRestController.Invoke")
    }
    deleteProjectResultDto, err := deleteproject.DeleteProjectServiceGetInstance().Invoke(
        deleteproject.DeleteProjectDtoFromPrimitives(application.MergePrimitives(
            c.GetHttpRequest().GetMetadata(), dashboardAuthPrimitives, map[string]any{"projectUuid": projectUuid},
        )),
    )
    if err != nil {
        return c.RespondError(err, "DeleteProjectRestController.Invoke")
    }
    return c.RespondOk(deleteProjectResultDto.ToArray())
}
