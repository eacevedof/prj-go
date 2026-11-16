package repositories

import (
    "fmt"

    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"

    "siem-backend/src/modules/projects/domain/entities"
)

type ProjectsWriterSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func ProjectsWriterSqliteRepositoryGetInstance() *ProjectsWriterSqliteRepository {
    return &ProjectsWriterSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// InsertProject en el alta updated_at = created_at.
func (r *ProjectsWriterSqliteRepository) InsertProject(projectEntity *entities.ProjectEntity) (int, error) {
    sql := `
        -- insertProject
        INSERT INTO ` + "`app_projects`" + `
            (project_uuid, project_name, project_description, created_at, updated_at)
        VALUES
            (?, ?, ?, ?, ?)
        `
    _, err := r.Execute(sql,
        projectEntity.GetProjectUuid(),
        projectEntity.GetProjectName(),
        projectEntity.GetProjectDescription(),
        projectEntity.GetCreatedAt(),
        projectEntity.GetCreatedAt(),
    )
    if err != nil {
        return 0, err
    }
    return r.GetLastInsertId(), nil
}

// UpdateProjectById nombre y descripcion (el uuid no cambia).
func (r *ProjectsWriterSqliteRepository) UpdateProjectById(projectId int, projectName string, projectDescription string, updatedAt string) (int, error) {
    sql := `
        -- updateProjectById
        UPDATE ` + "`app_projects`" + `
        SET
            project_name = ?, project_description = ?, updated_at = ?
        WHERE 1=1
        AND id = ?
        `
    return r.Execute(sql, projectName, projectDescription, updatedAt, projectId)
}

// DeleteProjectById se lleva por delante (ON DELETE CASCADE) sus agentes, monitores, sources y eventos.
func (r *ProjectsWriterSqliteRepository) DeleteProjectById(projectId int) (int, error) {
    sql := fmt.Sprintf(`
        -- deleteProjectById
        DELETE FROM `+"`app_projects`"+`
        WHERE 1=1
        AND id = %d
        `, projectId)
    return r.Execute(sql)
}
