package repositories

import (
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

// EventsFilterOptionsReaderSqliteRepository el catalogo de los desplegables del visor: lee projects, agents y monitors
// solo para listar uuid + nombre (modelo de lectura del visor, sin entidades de otros modulos).
type EventsFilterOptionsReaderSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func EventsFilterOptionsReaderSqliteRepositoryGetInstance() *EventsFilterOptionsReaderSqliteRepository {
    return &EventsFilterOptionsReaderSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

func (r *EventsFilterOptionsReaderSqliteRepository) GetProjectOptions() ([]map[string]any, error) {
    return r.Query(`
        -- getProjectOptions
        SELECT p.project_uuid, p.project_name
        FROM ` + "`app_projects`" + ` p
        ORDER BY p.project_name COLLATE NOCASE ASC
        `)
}

func (r *EventsFilterOptionsReaderSqliteRepository) GetAgentOptions() ([]map[string]any, error) {
    return r.Query(`
        -- getAgentOptions
        SELECT a.agent_uuid, a.agent_name, p.project_uuid
        FROM ` + "`app_agents`" + ` a
        INNER JOIN ` + "`app_projects`" + ` p ON p.id = a.project_id
        ORDER BY a.agent_name COLLATE NOCASE ASC
        `)
}

// GetMonitorOptions agent_uuid NULL = monitor de todos los agentes del proyecto.
func (r *EventsFilterOptionsReaderSqliteRepository) GetMonitorOptions() ([]map[string]any, error) {
    return r.Query(`
        -- getMonitorOptions
        SELECT m.monitor_uuid, m.monitor_name, m.monitor_type, p.project_uuid, a.agent_uuid
        FROM ` + "`app_monitors`" + ` m
        INNER JOIN ` + "`app_projects`" + ` p ON p.id = m.project_id
        LEFT JOIN ` + "`app_agents`" + ` a ON a.id = m.agent_id
        ORDER BY m.monitor_name COLLATE NOCASE ASC
        `)
}
