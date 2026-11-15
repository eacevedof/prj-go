package deletemonitor

import "siem-backend/src/modules/shared/application"

// DeleteMonitorDto DELETE /api/v1/monitors/:monitor_uuid
type DeleteMonitorDto struct {
    application.AbstractInputDto
    monitorUuid string
}

func DeleteMonitorDtoFromPrimitives(primitives map[string]any) *DeleteMonitorDto {
    return &DeleteMonitorDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        monitorUuid:      application.GetTrimmedString(primitives, "monitorUuid"),
    }
}

func (d *DeleteMonitorDto) GetMonitorUuid() string {
    return d.monitorUuid
}
