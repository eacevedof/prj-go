package createmonitor

import (
    "siem-backend/src/modules/shared/application"

    "siem-backend/src/modules/monitors/application/monitorinput"
)

// CreateMonitorDto POST /api/v1/projects/:project_uuid/monitors
type CreateMonitorDto struct {
    application.AbstractInputDto
    monitorinput.MonitorInput
    projectUuid string
}

func CreateMonitorDtoFromPrimitives(primitives map[string]any) *CreateMonitorDto {
    return &CreateMonitorDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        MonitorInput:     monitorinput.MonitorInputFromPrimitives(primitives),
        projectUuid:      application.GetTrimmedString(primitives, "projectUuid"),
    }
}

func (d *CreateMonitorDto) GetProjectUuid() string {
    return d.projectUuid
}
