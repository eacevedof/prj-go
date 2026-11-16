package updatemonitor

import (
    "siem-backend/src/modules/shared/application"

    "siem-backend/src/modules/monitors/application/monitorinput"
)

// UpdateMonitorDto PUT /api/v1/monitors/:monitor_uuid  (todos los campos editables; monitor_type no cambia)
type UpdateMonitorDto struct {
    application.AbstractInputDto
    monitorinput.MonitorInput
    monitorUuid string
}

func UpdateMonitorDtoFromPrimitives(primitives map[string]any) *UpdateMonitorDto {
    return &UpdateMonitorDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        MonitorInput:     monitorinput.MonitorInputFromPrimitives(primitives),
        monitorUuid:      application.GetTrimmedString(primitives, "monitorUuid"),
    }
}

func (d *UpdateMonitorDto) GetMonitorUuid() string {
    return d.monitorUuid
}
