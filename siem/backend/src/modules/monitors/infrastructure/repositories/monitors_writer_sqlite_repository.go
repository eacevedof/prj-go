package repositories

import (
    "fmt"

    "siem-backend/src/modules/monitors/domain/entities"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

type MonitorsWriterSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func MonitorsWriterSqliteRepositoryGetInstance() *MonitorsWriterSqliteRepository {
    return &MonitorsWriterSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// InsertMonitor en el alta updated_at = created_at. agent_id 0 se guarda NULL (todos los agentes del proyecto).
func (r *MonitorsWriterSqliteRepository) InsertMonitor(monitorEntity *entities.MonitorEntity) (int, error) {
    sql := fmt.Sprintf(`
        -- insertMonitor
        INSERT INTO `+"`monitors`"+`
            (monitor_uuid, project_id, agent_id, monitor_type, monitor_name, is_enabled, interval_seconds, config,
             warning_threshold, critical_threshold, created_at, updated_at)
        VALUES
            (?, %d, ?, ?, ?, %d, %d, ?, ?, ?, ?, ?)
        `, monitorEntity.GetProjectId(), getBoolAsInt(monitorEntity.IsEnabled()), monitorEntity.GetIntervalSeconds())
    _, err := r.Execute(sql,
        monitorEntity.GetMonitorUuid(),
        getIdOrNil(monitorEntity.GetAgentId()),
        monitorEntity.GetMonitorType(),
        monitorEntity.GetMonitorName(),
        monitorEntity.GetConfigJson(),
        getThresholdOrNil(monitorEntity.GetWarningThreshold()),
        getThresholdOrNil(monitorEntity.GetCriticalThreshold()),
        monitorEntity.GetCreatedAt(),
        monitorEntity.GetCreatedAt(),
    )
    if err != nil {
        return 0, err
    }
    return r.GetLastInsertId(), nil
}

// UpdateMonitorById todo lo editable; monitor_type no cambia (sus eventos ya guardados son de ese tipo).
func (r *MonitorsWriterSqliteRepository) UpdateMonitorById(monitorEntity *entities.MonitorEntity) (int, error) {
    sql := fmt.Sprintf(`
        -- updateMonitorById
        UPDATE `+"`monitors`"+`
        SET agent_id = ?, monitor_name = ?, is_enabled = %d, interval_seconds = %d, config = ?,
            warning_threshold = ?, critical_threshold = ?, updated_at = ?
        WHERE 1=1
        AND id = %d
        `, getBoolAsInt(monitorEntity.IsEnabled()), monitorEntity.GetIntervalSeconds(), monitorEntity.GetId())
    return r.Execute(sql,
        getIdOrNil(monitorEntity.GetAgentId()),
        monitorEntity.GetMonitorName(),
        monitorEntity.GetConfigJson(),
        getThresholdOrNil(monitorEntity.GetWarningThreshold()),
        getThresholdOrNil(monitorEntity.GetCriticalThreshold()),
        monitorEntity.GetUpdatedAt(),
    )
}

func (r *MonitorsWriterSqliteRepository) DeleteMonitorById(monitorId int) (int, error) {
    sql := fmt.Sprintf(`
        -- deleteMonitorById
        DELETE FROM `+"`monitors`"+`
        WHERE 1=1
        AND id = %d
        `, monitorId)
    return r.Execute(sql)
}

func getBoolAsInt(value bool) int {
    if value {
        return 1
    }
    return 0
}

func getIdOrNil(id int) any {
    if id == 0 {
        return nil
    }
    return id
}

func getThresholdOrNil(threshold *float64) any {
    if threshold == nil {
        return nil
    }
    return *threshold
}
