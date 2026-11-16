package repositories

import (
    "fmt"

    "siem-backend/src/modules/monitors/domain/entities"
    "siem-backend/src/modules/shared/application"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

const monitorColumnsSql = `
            m.id, m.monitor_uuid, m.project_id, m.agent_id, m.monitor_type, m.monitor_name, m.is_enabled,
            m.interval_seconds, m.config, m.warning_threshold, m.critical_threshold, m.created_at, m.updated_at,
            p.project_uuid, COALESCE(a.agent_uuid, '') AS agent_uuid, COALESCE(a.agent_name, '') AS agent_name
        FROM ` + "`app_monitors`" + ` m
        INNER JOIN ` + "`app_projects`" + ` p ON p.id = m.project_id
        LEFT JOIN ` + "`app_agents`" + ` a ON a.id = m.agent_id`

type MonitorsReaderSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func MonitorsReaderSqliteRepositoryGetInstance() *MonitorsReaderSqliteRepository {
    return &MonitorsReaderSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

func (r *MonitorsReaderSqliteRepository) GetMonitorById(monitorId int) (*entities.MonitorEntity, error) {
    sql := fmt.Sprintf(`
        -- getMonitorById
        SELECT %s
        WHERE 1=1
        AND m.id = %d
        `, monitorColumnsSql, monitorId)
    return r.getFirstMonitorEntity(r.Query(sql))
}

func (r *MonitorsReaderSqliteRepository) GetMonitorByMonitorUuid(monitorUuid string) (*entities.MonitorEntity, error) {
    sql := fmt.Sprintf(`
        -- getMonitorByMonitorUuid
        SELECT %s
        WHERE 1=1
        AND m.monitor_uuid = ?
        `, monitorColumnsSql)
    return r.getFirstMonitorEntity(r.Query(sql, monitorUuid))
}

// GetMonitorsByProjectId todos los del proyecto (activos o no), por tipo y nombre: la tabla del panel.
func (r *MonitorsReaderSqliteRepository) GetMonitorsByProjectId(projectId int) ([]*entities.MonitorEntity, error) {
    sql := fmt.Sprintf(`
        -- getMonitorsByProjectId
        SELECT %s
        WHERE 1=1
        AND m.project_id = %d
        ORDER BY m.monitor_type ASC, m.monitor_name COLLATE NOCASE ASC
        `, monitorColumnsSql, projectId)
    return r.getMonitorEntities(r.Query(sql))
}

// GetEnabledMonitorsForAgent lo que tiene que vigilar un agente: los activos de su proyecto que son de todos
// (agent_id NULL) o suyos. Orden estable: de el sale el config_hash.
func (r *MonitorsReaderSqliteRepository) GetEnabledMonitorsForAgent(projectId int, agentId int) ([]*entities.MonitorEntity, error) {
    sql := fmt.Sprintf(`
        -- getEnabledMonitorsForAgent
        SELECT %s
        WHERE 1=1
        AND m.project_id = %d
        AND m.is_enabled = 1
        AND (m.agent_id IS NULL OR m.agent_id = %d)
        ORDER BY m.id ASC
        `, monitorColumnsSql, projectId, agentId)
    return r.getMonitorEntities(r.Query(sql))
}

// DoesMonitorNameExistInProject sin distinguir mayusculas (indice uk_monitors_project_id_monitor_name).
// excludedMonitorId: al editar, el propio monitor no cuenta (0 en el alta).
func (r *MonitorsReaderSqliteRepository) DoesMonitorNameExistInProject(projectId int, monitorName string, excludedMonitorId int) (bool, error) {
    sql := fmt.Sprintf(`
        -- doesMonitorNameExistInProject
        SELECT
            COUNT(*) AS n_rows
        FROM `+"`app_monitors`"+`
        WHERE 1=1
        AND project_id = %d
        AND id <> %d
        AND monitor_name = ? COLLATE NOCASE
        `, projectId, excludedMonitorId)
    rows, err := r.Query(sql, monitorName)
    if err != nil || len(rows) == 0 {
        return false, err
    }
    return application.ToInt(rows[0]["n_rows"]) > 0, nil
}

// GetMonitorStatesByProjectId el ultimo estado de cada monitor en cada agente. Un monitor no tiene tabla de estados:
// sus eventos van a la source del agente con source_name = monitor_name (metric o file).
func (r *MonitorsReaderSqliteRepository) GetMonitorStatesByProjectId(projectId int) ([]map[string]any, error) {
    sql := fmt.Sprintf(`
        -- getMonitorStatesByProjectId
        SELECT
            s.source_name, a.agent_uuid, a.agent_name, s.last_event_at, s.last_severity, s.last_message
        FROM `+"`app_sources`"+` s
        INNER JOIN `+"`app_agents`"+` a ON a.id = s.agent_id
        INNER JOIN `+"`app_monitors`"+` m ON m.project_id = a.project_id AND m.monitor_name = s.source_name COLLATE NOCASE
        WHERE 1=1
        AND a.project_id = %d
        AND s.source_type IN ('metric', 'file')
        ORDER BY a.agent_name COLLATE NOCASE ASC
        `, projectId)
    return r.Query(sql)
}

func (r *MonitorsReaderSqliteRepository) getFirstMonitorEntity(rows []map[string]any, err error) (*entities.MonitorEntity, error) {
    if err != nil || len(rows) == 0 {
        return nil, err
    }
    return entities.MonitorEntityFromPrimitives(rows[0]), nil
}

func (r *MonitorsReaderSqliteRepository) getMonitorEntities(rows []map[string]any, err error) ([]*entities.MonitorEntity, error) {
    if err != nil {
        return nil, err
    }
    monitorEntities := make([]*entities.MonitorEntity, 0, len(rows))
    for _, row := range rows {
        monitorEntities = append(monitorEntities, entities.MonitorEntityFromPrimitives(row))
    }
    return monitorEntities, nil
}
