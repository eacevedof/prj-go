package authenticatedashboard

import (
    "siem-backend/src/modules/shared/infrastructure/components/token"
    "siem-backend/src/modules/shared/infrastructure/repositories/configuration"

    "siem-backend/src/modules/dashboard/domain/exceptions"
)

// AuthenticateDashboardService el Bearer del front contra APP_DASHBOARD_TOKEN (un unico token de solo lectura: el
// dashboard no tiene usuarios). Sin APP_DASHBOARD_TOKEN configurado no entra nadie.
type AuthenticateDashboardService struct {
    tokenGenerator                 *token.TokenGenerator
    environmentReaderRawRepository *configuration.EnvironmentReaderRawRepository

    authenticateDashboardDto *AuthenticateDashboardDto
}

func AuthenticateDashboardServiceGetInstance() *AuthenticateDashboardService {
    return &AuthenticateDashboardService{
        tokenGenerator:                 token.GetInstance(),
        environmentReaderRawRepository: configuration.GetInstance(),
    }
}

// Invoke error: *DashboardException (401).
func (s *AuthenticateDashboardService) Invoke(authenticateDashboardDto *AuthenticateDashboardDto) (*AuthenticateDashboardResultDto, error) {
    s.authenticateDashboardDto = authenticateDashboardDto

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }
    return AuthenticateDashboardResultDtoFromPrimitives(map[string]any{"isDashboardAuthenticated": 1}), nil
}

func (s *AuthenticateDashboardService) failIfWrongInput() error {
    expectedDashboardToken := s.environmentReaderRawRepository.GetDashboardToken()
    dashboardToken := s.authenticateDashboardDto.GetDashboardToken()
    if expectedDashboardToken == "" || dashboardToken == "" {
        return exceptions.UnauthorizedInvalidToken()
    }
    if !s.tokenGenerator.AreTokensEqual(expectedDashboardToken, dashboardToken) {
        return exceptions.UnauthorizedInvalidToken()
    }
    return nil
}
