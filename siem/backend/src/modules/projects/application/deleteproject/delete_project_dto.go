package deleteproject

import "siem-backend/src/modules/shared/application"

// DeleteProjectDto DELETE /api/v1/projects/:project_uuid
type DeleteProjectDto struct {
    application.AbstractInputDto
    projectUuid string
}

func DeleteProjectDtoFromPrimitives(primitives map[string]any) *DeleteProjectDto {
    return &DeleteProjectDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        projectUuid:      application.GetTrimmedString(primitives, "projectUuid"),
    }
}

func (d *DeleteProjectDto) GetProjectUuid() string {
    return d.projectUuid
}
