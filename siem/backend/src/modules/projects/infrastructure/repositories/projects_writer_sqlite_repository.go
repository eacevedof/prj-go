package repositories

import (
    "fmt"

    "siem-backend/src/modules/projects/domain/entities"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
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
        INSERT INTO ` + "`projects`" + `
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

// DeleteProjectById se lleva por delante (ON DELETE CASCADE) sus agentes, monitores, sources y eventos.
func (r *ProjectsWriterSqliteRepository) DeleteProjectById(projectId int) (int, error) {
    sql := fmt.Sprintf(`
        -- deleteProjectById
        DELETE FROM `+"`projects`"+`
        WHERE 1=1
        AND id = %d
        `, projectId)
    return r.Execute(sql)
}
