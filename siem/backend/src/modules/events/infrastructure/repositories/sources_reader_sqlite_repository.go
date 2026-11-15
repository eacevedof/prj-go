package repositories

import (
    "fmt"

    "siem-backend/src/modules/shared/application"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

type SourcesReaderSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func SourcesReaderSqliteRepositoryGetInstance() *SourcesReaderSqliteRepository {
    return &SourcesReaderSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// GetSourceIdByAgentIdAndTypeAndName 0 si el agente aun no tiene esa source.
func (r *SourcesReaderSqliteRepository) GetSourceIdByAgentIdAndTypeAndName(agentId int, sourceType string, sourceName string) (int, error) {
    sql := fmt.Sprintf(`
        -- getSourceIdByAgentIdAndTypeAndName
        SELECT
            id
        FROM `+"`sources`"+`
        WHERE 1=1
        AND agent_id = %d
        AND source_type = ?
        AND source_name = ?
        `, agentId)
    rows, err := r.Query(sql, sourceType, sourceName)
    if err != nil || len(rows) == 0 {
        return 0, err
    }
    return application.ToInt(rows[0]["id"]), nil
}

// GetSourcesByProjectId las sources de todos los agentes de un proyecto, con su ultimo estado.
func (r *SourcesReaderSqliteRepository) GetSourcesByProjectId(projectId int) ([]map[string]any, error) {
    sql := fmt.Sprintf(`
        -- getSourcesByProjectId
        SELECT
            s.agent_id, s.source_type, s.source_name, s.events_count, s.last_event_at, s.last_severity,
            s.last_exit_code, s.last_message
        FROM `+"`sources`"+` s
        INNER JOIN `+"`agents`"+` a ON a.id = s.agent_id
        WHERE 1=1
        AND a.project_id = %d
        ORDER BY s.agent_id ASC, s.source_type ASC, s.source_name ASC
        `, projectId)
    return r.Query(sql)
}
