package updateproject

import "siem-backend/src/modules/projects/domain/entities"

type UpdateProjectResultDto struct {
    // el proyecto ya cambiado (ProjectEntity.ToPrimitives)
    project entities.ProjectPrimitives
}

func UpdateProjectResultDtoFromPrimitives(primitives map[string]any) *UpdateProjectResultDto {
    project, _ := primitives["project"].(entities.ProjectPrimitives)
    return &UpdateProjectResultDto{project: project}
}

func (d *UpdateProjectResultDto) GetProject() entities.ProjectPrimitives {
    return d.project
}

func (d *UpdateProjectResultDto) ToArray() entities.ProjectPrimitives {
    return d.project
}
