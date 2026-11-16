package entities

import (
    "siem-backend/src/modules/shared/application"

    "siem-backend/src/modules/projects/domain/enums"
)

const PROJECT_UUID_PREFIX = "project-"

// ProjectPrimitives representacion publica: el id interno no sale, el identificador es project_uuid.
type ProjectPrimitives struct {
    ProjectUuid        string `json:"project_uuid"`
    ProjectName        string `json:"project_name"`
    ProjectDescription string `json:"project_description"`
    CreatedAt          string `json:"created_at"`
    UpdatedAt          string `json:"updated_at"`
}

type ProjectEntity struct {
    id                 int
    projectUuid        string
    projectName        string
    projectDescription string
    createdAt          string
    updatedAt          string
}

// ProjectEntityFromPrimitives desde una fila de la tabla projects (columnas en snake_case).
func ProjectEntityFromPrimitives(primitives map[string]any) *ProjectEntity {
    return &ProjectEntity{
        id:                 application.GetInt(primitives, enums.ProjectsKeyEnum.ID),
        projectUuid:        application.GetTrimmedString(primitives, enums.ProjectsKeyEnum.PROJECT_UUID),
        projectName:        application.GetTrimmedString(primitives, enums.ProjectsKeyEnum.PROJECT_NAME),
        projectDescription: application.GetTrimmedString(primitives, enums.ProjectsKeyEnum.PROJECT_DESCRIPTION),
        createdAt:          application.GetTrimmedString(primitives, enums.ProjectsKeyEnum.CREATED_AT),
        updatedAt:          application.GetTrimmedString(primitives, enums.ProjectsKeyEnum.UPDATED_AT),
    }
}

func (e *ProjectEntity) ToPrimitives() ProjectPrimitives {
    return ProjectPrimitives{
        ProjectUuid:        e.projectUuid,
        ProjectName:        e.projectName,
        ProjectDescription: e.projectDescription,
        CreatedAt:          e.createdAt,
        UpdatedAt:          e.updatedAt,
    }
}

// ToPrimitivesOrEmpty una lectura que no encuentra nada da el proyecto vacio.
func ToPrimitivesOrEmpty(projectEntity *ProjectEntity) ProjectPrimitives {
    if projectEntity == nil {
        return ProjectPrimitives{}
    }
    return projectEntity.ToPrimitives()
}

func (e *ProjectEntity) GetId() int {
    return e.id
}

func (e *ProjectEntity) GetProjectUuid() string {
    return e.projectUuid
}

func (e *ProjectEntity) GetProjectName() string {
    return e.projectName
}

func (e *ProjectEntity) GetProjectDescription() string {
    return e.projectDescription
}

func (e *ProjectEntity) GetCreatedAt() string {
    return e.createdAt
}

func (e *ProjectEntity) GetUpdatedAt() string {
    return e.updatedAt
}
