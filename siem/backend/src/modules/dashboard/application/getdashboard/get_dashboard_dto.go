package getdashboard

import "siem-backend/src/modules/shared/application"

type GetDashboardDto struct {
    application.AbstractInputDto
    isDashboardAuthenticated bool
}

func GetDashboardDtoFromPrimitives(primitives map[string]any) *GetDashboardDto {
    return &GetDashboardDto{
        AbstractInputDto:         application.NewAbstractInputDto(primitives),
        isDashboardAuthenticated: application.GetInt(primitives, "isDashboardAuthenticated") == 1,
    }
}

func (d *GetDashboardDto) IsDashboardAuthenticated() bool {
    return d.isDashboardAuthenticated
}
