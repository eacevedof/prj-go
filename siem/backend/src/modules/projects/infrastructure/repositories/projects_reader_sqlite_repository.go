package repositories

import (
    "fmt"

    "siem-backend/src/modules/projects/domain/entities"
    "siem-backend/src/modules/shared/application"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

type ProjectsReaderSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func ProjectsReaderSqliteRepositoryGetInstance() *ProjectsReaderSqliteRepository {
    return &ProjectsReaderSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

func (r *ProjectsReaderSqliteRepository) GetProjectById(projectId int) (*entities.ProjectEntity, error) {
    sql := fmt.Sprintf(`
        -- getProjectById
        SELECT
            id, project_uuid, project_name, project_description, created_at, updated_at
        FROM `+"`app_projects`"+`
        WHERE 1=1
        AND id = %d
        `, projectId)
    return r.getFirstProjectEntity(r.Query(sql))
}

func (r *ProjectsReaderSqliteRepository) GetProjectByProjectUuid(projectUuid string) (*entities.ProjectEntity, error) {
    sql := `
        -- getProjectByProjectUuid
        SELECT
            id, project_uuid, project_name, project_description, created_at, updated_at
        FROM ` + "`app_projects`" + `
        WHERE 1=1
        AND project_uuid = ?
        `
    return r.getFirstProjectEntity(r.Query(sql, projectUuid))
}

// DoesProjectNameExist project_name es unico sin distinguir mayusculas (indice uk_projects_project_name).
func (r *ProjectsReaderSqliteRepository) DoesProjectNameExist(projectName string) (bool, error) {
    sql := `
        -- doesProjectNameExist
        SELECT
            COUNT(*) AS n_rows
        FROM ` + "`app_projects`" + `
        WHERE 1=1
        AND project_name = ? COLLATE NOCASE
        `
    rows, err := r.Query(sql, projectName)
    if err != nil || len(rows) == 0 {
        return false, err
    }
    return application.ToInt(rows[0]["n_rows"]) > 0, nil
}

// DoesProjectNameExistInOtherProject al editar: el nombre lo tiene otro proyecto que no es excludedProjectId.
func (r *ProjectsReaderSqliteRepository) DoesProjectNameExistInOtherProject(projectName string, excludedProjectId int) (bool, error) {
    sql := `
        -- doesProjectNameExistInOtherProject
        SELECT
            COUNT(*) AS n_rows
        FROM ` + "`app_projects`" + `
        WHERE 1=1
        AND project_name = ? COLLATE NOCASE
        AND id <> ?
        `
    rows, err := r.Query(sql, projectName, excludedProjectId)
    if err != nil || len(rows) == 0 {
        return false, err
    }
    return application.ToInt(rows[0]["n_rows"]) > 0, nil
}

func (r *ProjectsReaderSqliteRepository) getFirstProjectEntity(rows []map[string]any, err error) (*entities.ProjectEntity, error) {
    if err != nil || len(rows) == 0 {
        return nil, err
    }
    return entities.ProjectEntityFromPrimitives(rows[0]), nil
}

// GetAllProjects todos los proyectos, por nombre (el dashboard).
func (r *ProjectsReaderSqliteRepository) GetAllProjects() ([]*entities.ProjectEntity, error) {
    sql := `
        -- getAllProjects
        SELECT
            id, project_uuid, project_name, project_description, created_at, updated_at
        FROM ` + "`app_projects`" + `
        ORDER BY project_name COLLATE NOCASE ASC
        `
    rows, err := r.Query(sql)
    if err != nil {
        return nil, err
    }
    projectEntities := make([]*entities.ProjectEntity, 0, len(rows))
    for _, row := range rows {
        projectEntities = append(projectEntities, entities.ProjectEntityFromPrimitives(row))
    }
    return projectEntities, nil
}
