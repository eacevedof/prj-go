package getprojectdetail

import "siem-backend/src/modules/shared/application"

type GetProjectDetailDto struct {
    application.AbstractInputDto
    isDashboardAuthenticated bool
    projectUuid              string
}

func GetProjectDetailDtoFromPrimitives(primitives map[string]any) *GetProjectDetailDto {
    return &GetProjectDetailDto{
        AbstractInputDto:         application.NewAbstractInputDto(primitives),
        isDashboardAuthenticated: application.GetInt(primitives, "isDashboardAuthenticated") == 1,
        projectUuid:              application.GetTrimmedString(primitives, "projectUuid"),
    }
}

func (d *GetProjectDetailDto) IsDashboardAuthenticated() bool {
    return d.isDashboardAuthenticated
}

func (d *GetProjectDetailDto) GetProjectUuid() string {
    return d.projectUuid
}
