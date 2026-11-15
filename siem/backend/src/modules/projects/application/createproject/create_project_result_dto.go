package createproject

import "siem-backend/src/modules/projects/domain/entities"

type CreateProjectResultDto struct {
    // el proyecto creado (ProjectEntity.ToPrimitives)
    project entities.ProjectPrimitives
}

func CreateProjectResultDtoFromPrimitives(primitives map[string]any) *CreateProjectResultDto {
    project, _ := primitives["project"].(entities.ProjectPrimitives)
    return &CreateProjectResultDto{project: project}
}

func (d *CreateProjectResultDto) GetProject() entities.ProjectPrimitives {
    return d.project
}

func (d *CreateProjectResultDto) ToArray() entities.ProjectPrimitives {
    return d.project
}
