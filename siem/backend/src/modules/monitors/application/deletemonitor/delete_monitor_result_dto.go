package deletemonitor

import "siem-backend/src/modules/shared/application"

// DeleteMonitorPrimitives lo que devuelve la API al borrar.
type DeleteMonitorPrimitives struct {
    MonitorUuid string `json:"monitor_uuid"`
}

type DeleteMonitorResultDto struct {
    monitorUuid string
}

func DeleteMonitorResultDtoFromPrimitives(primitives map[string]any) *DeleteMonitorResultDto {
    return &DeleteMonitorResultDto{monitorUuid: application.GetTrimmedString(primitives, "monitorUuid")}
}

func (d *DeleteMonitorResultDto) ToArray() DeleteMonitorPrimitives {
    return DeleteMonitorPrimitives{MonitorUuid: d.monitorUuid}
}
