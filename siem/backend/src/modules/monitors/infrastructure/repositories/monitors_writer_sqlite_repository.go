package repositories

import (
    "fmt"

    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"

    "siem-backend/src/modules/monitors/domain/entities"
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
        INSERT INTO `+"`app_monitors`"+`
            (monitor_uuid, project_id, agent_id, monitor_type, monitor_name, is_enabled, interval_seconds, config,
             warning_threshold, critical_threshold, created_at, updated_at)
        VALUES
            (?, %d, ?, ?, ?, %d, %d, ?, ?, ?, ?, ?)
        `, monitorEntity.GetProjectId(), r.getIsEnabledAsInt(monitorEntity.IsEnabled()), monitorEntity.GetIntervalSeconds())
    _, err := r.Execute(sql,
        monitorEntity.GetMonitorUuid(),
        r.getIdOrNil(monitorEntity.GetAgentId()),
        monitorEntity.GetMonitorType(),
        monitorEntity.GetMonitorName(),
        monitorEntity.GetConfigJson(),
        r.getThresholdOrNil(monitorEntity.GetWarningThreshold()),
        r.getThresholdOrNil(monitorEntity.GetCriticalThreshold()),
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
        UPDATE `+"`app_monitors`"+`
        SET agent_id = ?, monitor_name = ?, is_enabled = %d, interval_seconds = %d, config = ?,
            warning_threshold = ?, critical_threshold = ?, updated_at = ?
        WHERE 1=1
        AND id = %d
        `, r.getIsEnabledAsInt(monitorEntity.IsEnabled()), monitorEntity.GetIntervalSeconds(), monitorEntity.GetId())
    return r.Execute(sql,
        r.getIdOrNil(monitorEntity.GetAgentId()),
        monitorEntity.GetMonitorName(),
        monitorEntity.GetConfigJson(),
        r.getThresholdOrNil(monitorEntity.GetWarningThreshold()),
        r.getThresholdOrNil(monitorEntity.GetCriticalThreshold()),
        monitorEntity.GetUpdatedAt(),
    )
}

func (r *MonitorsWriterSqliteRepository) DeleteMonitorById(monitorId int) (int, error) {
    sql := fmt.Sprintf(`
        -- deleteMonitorById
        DELETE FROM `+"`app_monitors`"+`
        WHERE 1=1
        AND id = %d
        `, monitorId)
    return r.Execute(sql)
}

func (r *MonitorsWriterSqliteRepository) getIsEnabledAsInt(isEnabled bool) int {
    if isEnabled {
        return 1
    }
    return 0
}

func (r *MonitorsWriterSqliteRepository) getIdOrNil(id int) any {
    if id == 0 {
        return nil
    }
    return id
}

func (r *MonitorsWriterSqliteRepository) getThresholdOrNil(threshold *float64) any {
    if threshold == nil {
        return nil
    }
    return *threshold
}
