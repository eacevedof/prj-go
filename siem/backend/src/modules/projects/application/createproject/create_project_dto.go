package createproject

import "siem-backend/src/modules/shared/application"

type CreateProjectDto struct {
    application.AbstractInputDto
    projectName        string
    projectDescription string
}

func CreateProjectDtoFromPrimitives(primitives map[string]any) *CreateProjectDto {
    return &CreateProjectDto{
        AbstractInputDto:   application.NewAbstractInputDto(primitives),
        projectName:        application.GetTrimmedString(primitives, "projectName"),
        projectDescription: application.GetTrimmedString(primitives, "projectDescription"),
    }
}

func (d *CreateProjectDto) GetProjectName() string {
    return d.projectName
}

func (d *CreateProjectDto) GetProjectDescription() string {
    return d.projectDescription
}
