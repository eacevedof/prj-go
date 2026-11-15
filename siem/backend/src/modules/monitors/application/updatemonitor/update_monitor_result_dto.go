package updatemonitor

import "siem-backend/src/modules/monitors/domain/entities"

type UpdateMonitorResultDto struct {
    // el monitor ya editado (MonitorEntity.ToPrimitives)
    monitor entities.MonitorPrimitives
}

func UpdateMonitorResultDtoFromPrimitives(primitives map[string]any) *UpdateMonitorResultDto {
    monitor, _ := primitives["monitor"].(entities.MonitorPrimitives)
    return &UpdateMonitorResultDto{monitor: monitor}
}

func (d *UpdateMonitorResultDto) ToArray() entities.MonitorPrimitives {
    return d.monitor
}
