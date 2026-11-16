package updateproject

import "siem-backend/src/modules/shared/application"

// UpdateProjectDto PUT /api/v1/projects/:project_uuid {project_name, project_description}
type UpdateProjectDto struct {
    application.AbstractInputDto
    projectUuid        string
    projectName        string
    projectDescription string
}

func UpdateProjectDtoFromPrimitives(primitives map[string]any) *UpdateProjectDto {
    return &UpdateProjectDto{
        AbstractInputDto:   application.NewAbstractInputDto(primitives),
        projectUuid:        application.GetTrimmedString(primitives, "projectUuid"),
        projectName:        application.GetTrimmedString(primitives, "projectName"),
        projectDescription: application.GetTrimmedString(primitives, "projectDescription"),
    }
}

func (d *UpdateProjectDto) GetProjectUuid() string {
    return d.projectUuid
}

func (d *UpdateProjectDto) GetProjectName() string {
    return d.projectName
}

func (d *UpdateProjectDto) GetProjectDescription() string {
    return d.projectDescription
}
