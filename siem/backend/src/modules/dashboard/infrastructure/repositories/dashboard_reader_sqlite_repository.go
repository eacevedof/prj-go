package repositories

import (
    "fmt"

    "siem-backend/src/modules/shared/application"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

// DashboardReaderSqliteRepository las lecturas agregadas que solo necesita el dashboard (contadores sobre events,
// ingest_queue y sources, tamaño de la BD). Las entidades las leen los repositorios de su modulo.
type DashboardReaderSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func DashboardReaderSqliteRepositoryGetInstance() *DashboardReaderSqliteRepository {
    return &DashboardReaderSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// GetQueueStatsByStatus por cada estado vivo de la cola (pending, failed): lotes, eventos y el lote mas antiguo.
func (r *DashboardReaderSqliteRepository) GetQueueStatsByStatus() ([]map[string]any, error) {
    sql := `
        -- getQueueStatsByStatus
        SELECT
            queue_status, COUNT(*) AS n_batches, COALESCE(SUM(events_count), 0) AS n_events,
            MIN(received_at) AS oldest_received_at
        FROM ` + "`app_ingest_queue`" + `
        WHERE 1=1
        AND queue_status IN ('pending', 'failed')
        GROUP BY queue_status
        `
    return r.Query(sql)
}

func (r *DashboardReaderSqliteRepository) GetEventsTotal() (int, error) {
    sql := `
        -- getEventsTotal
        SELECT
            COUNT(*) AS n_rows
        FROM ` + "`app_events`" + `
        `
    return r.getFirstInt(r.Query(sql))
}

// GetDbSizeBytes paginas de la BD x su tamaño (sin contar el fichero -wal).
func (r *DashboardReaderSqliteRepository) GetDbSizeBytes() (int, error) {
    sql := `
        -- getDbSizeBytes
        SELECT
            page_count * page_size AS n_rows
        FROM pragma_page_count(), pragma_page_size()
        `
    return r.getFirstInt(r.Query(sql))
}

// GetEventsStatsByProjectIdSince eventos del proyecto desde occurredSince agrupados por gravedad (usa el indice
// idx_events_project_id_occurred_at).
func (r *DashboardReaderSqliteRepository) GetEventsStatsByProjectIdSince(projectId int, occurredSince string) ([]map[string]any, error) {
    sql := fmt.Sprintf(`
        -- getEventsStatsByProjectIdSince
        SELECT
            severity, COUNT(*) AS n_events
        FROM `+"`app_events`"+`
        WHERE 1=1
        AND project_id = %d
        AND occurred_at >= ?
        GROUP BY severity
        `, projectId)
    return r.Query(sql, occurredSince)
}

// GetEventsStatsByAgentIdSince lo mismo por agente (usa idx_events_agent_id_occurred_at).
func (r *DashboardReaderSqliteRepository) GetEventsStatsByAgentIdSince(agentId int, occurredSince string) ([]map[string]any, error) {
    sql := fmt.Sprintf(`
        -- getEventsStatsByAgentIdSince
        SELECT
            severity, COUNT(*) AS n_events
        FROM `+"`app_events`"+`
        WHERE 1=1
        AND agent_id = %d
        AND occurred_at >= ?
        GROUP BY severity
        `, agentId)
    return r.Query(sql, occurredSince)
}

// GetLastEventAtByProjectId el ultimo evento de cada proyecto, sin tocar events: lo guardan sus sources.
func (r *DashboardReaderSqliteRepository) GetLastEventAtByProjectId() (map[int]string, error) {
    sql := `
        -- getLastEventAtByProjectId
        SELECT
            a.project_id, MAX(s.last_event_at) AS last_event_at
        FROM ` + "`app_sources`" + ` s
        INNER JOIN ` + "`app_agents`" + ` a ON a.id = s.agent_id
        GROUP BY a.project_id
        `
    rows, err := r.Query(sql)
    if err != nil {
        return nil, err
    }
    lastEventAtByProjectId := map[int]string{}
    for _, row := range rows {
        lastEventAtByProjectId[application.ToInt(row["project_id"])] = application.ToString(row["last_event_at"])
    }
    return lastEventAtByProjectId, nil
}

func (r *DashboardReaderSqliteRepository) getFirstInt(rows []map[string]any, err error) (int, error) {
    if err != nil || len(rows) == 0 {
        return 0, err
    }
    return application.ToInt(rows[0]["n_rows"]), nil
}
