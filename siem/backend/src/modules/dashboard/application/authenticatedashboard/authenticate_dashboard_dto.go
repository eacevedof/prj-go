package authenticatedashboard

import "siem-backend/src/modules/shared/application"

type AuthenticateDashboardDto struct {
    application.AbstractInputDto
    dashboardToken string
}

func AuthenticateDashboardDtoFromPrimitives(primitives map[string]any) *AuthenticateDashboardDto {
    return &AuthenticateDashboardDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        dashboardToken:   application.GetTrimmedString(primitives, "dashboardToken"),
    }
}

func (d *AuthenticateDashboardDto) GetDashboardToken() string {
    return d.dashboardToken
}
