package authenticatedashboard

import "siem-backend/src/modules/shared/application"

type AuthenticateDashboardResultDto struct {
    isDashboardAuthenticated bool
}

func AuthenticateDashboardResultDtoFromPrimitives(primitives map[string]any) *AuthenticateDashboardResultDto {
    return &AuthenticateDashboardResultDto{
        isDashboardAuthenticated: application.GetInt(primitives, "isDashboardAuthenticated") == 1,
    }
}

func (d *AuthenticateDashboardResultDto) IsDashboardAuthenticated() bool {
    return d.isDashboardAuthenticated
}

// ToAuthPrimitives lo que el controller mezcla en los primitivos del dto del caso de uso que viene despues.
func (d *AuthenticateDashboardResultDto) ToAuthPrimitives() map[string]any {
    isDashboardAuthenticated := 0
    if d.isDashboardAuthenticated {
        isDashboardAuthenticated = 1
    }
    return map[string]any{"isDashboardAuthenticated": isDashboardAuthenticated}
}
