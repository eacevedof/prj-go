package createmonitor

import "siem-backend/src/modules/monitors/domain/entities"

type CreateMonitorResultDto struct {
    // el monitor creado (MonitorEntity.ToPrimitives)
    monitor entities.MonitorPrimitives
}

func CreateMonitorResultDtoFromPrimitives(primitives map[string]any) *CreateMonitorResultDto {
    monitor, _ := primitives["monitor"].(entities.MonitorPrimitives)
    return &CreateMonitorResultDto{monitor: monitor}
}

func (d *CreateMonitorResultDto) ToArray() entities.MonitorPrimitives {
    return d.monitor
}
