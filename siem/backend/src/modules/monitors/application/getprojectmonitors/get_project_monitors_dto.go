package getprojectmonitors

import "siem-backend/src/modules/shared/application"

// GetProjectMonitorsDto GET /api/v1/projects/:project_uuid/monitors
type GetProjectMonitorsDto struct {
    application.AbstractInputDto
    projectUuid string
}

func GetProjectMonitorsDtoFromPrimitives(primitives map[string]any) *GetProjectMonitorsDto {
    return &GetProjectMonitorsDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        projectUuid:      application.GetTrimmedString(primitives, "projectUuid"),
    }
}

func (d *GetProjectMonitorsDto) GetProjectUuid() string {
    return d.projectUuid
}
