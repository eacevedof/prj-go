package controllers

import (
    "siem-backend/src/modules/projects/application/createproject"
    "siem-backend/src/modules/projects/application/deleteproject"
    "siem-backend/src/modules/projects/application/updateproject"
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

// UpdateProjectRestController PUT /api/v1/projects/:project_uuid (Bearer del dashboard)
// {project_name, project_description} -> 200 el proyecto.
type UpdateProjectRestController struct {
    sharedcontrollers.AbstractRestController
}

func UpdateProjectRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *UpdateProjectRestController {
    return &UpdateProjectRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *UpdateProjectRestController) Invoke(projectUuid string) *httpresponder.HttpResponder {
    dashboardAuthPrimitives, err := c.GetDashboardAuthOrFail()
    if err != nil {
        return c.RespondError(err, "UpdateProjectRestController.Invoke")
    }
    httpRequest := c.GetHttpRequest()
    updateProjectResultDto, err := updateproject.UpdateProjectServiceGetInstance().Invoke(
        updateproject.UpdateProjectDtoFromPrimitives(application.MergePrimitives(
            httpRequest.GetMetadata(), dashboardAuthPrimitives,
            map[string]any{
                "projectUuid":        projectUuid,
                "projectName":        httpRequest.GetPostString(enums.ProjectsKeyEnum.PROJECT_NAME),
                "projectDescription": httpRequest.GetPostString(enums.ProjectsKeyEnum.PROJECT_DESCRIPTION),
            },
        )),
    )
    if err != nil {
        return c.RespondError(err, "UpdateProjectRestController.Invoke")
    }
    return c.RespondOk(updateProjectResultDto.ToArray())
}

// DeleteProjectRestController DELETE /api/v1/projects/:project_uuid (Bearer del dashboard) ->
// {project_uuid, project_name, agents: [{agent_name, host_name, host_os, host_ip, client_path}]}
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
