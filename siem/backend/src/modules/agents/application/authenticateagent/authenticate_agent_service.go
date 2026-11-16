package authenticateagent

import (
    "siem-backend/src/core/config"

    "siem-backend/src/modules/shared/infrastructure/components/token"

    "siem-backend/src/modules/agents/domain/enums"
    "siem-backend/src/modules/agents/infrastructure/repositories"
    "siem-backend/src/modules/agents/domain/exceptions"
)

// AuthenticateAgentService el Bearer de un siem-client -> su agente. Mismo 401 para todo fallo (sin token, token
// desconocido, agente deshabilitado): no se da pista de cual.
type AuthenticateAgentService struct {
    tokenGenerator               *token.TokenGenerator
    agentsReaderSqliteRepository *repositories.AgentsReaderSqliteRepository

    authenticateAgentDto *AuthenticateAgentDto
}

func AuthenticateAgentServiceGetInstance() *AuthenticateAgentService {
    return &AuthenticateAgentService{
        tokenGenerator:               token.GetInstance(),
        agentsReaderSqliteRepository: repositories.AgentsReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *AgentsException (401) o uno inesperado.
func (s *AuthenticateAgentService) Invoke(authenticateAgentDto *AuthenticateAgentDto) (*AuthenticateAgentResultDto, error) {
    s.authenticateAgentDto = authenticateAgentDto

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }

    agentEntity, err := s.agentsReaderSqliteRepository.GetAgentByApiTokenHash(
        s.tokenGenerator.GetTokenHash(authenticateAgentDto.GetApiToken()),
    )
    if err != nil {
        return nil, err
    }
    if agentEntity == nil || !agentEntity.IsEnabled() {
        return nil, exceptions.UnauthorizedCustom(config.Tr(enums.AgentsMessageEnum.INVALID_OR_MISSING_API_TOKEN))
    }

    return AuthenticateAgentResultDtoFromPrimitives(map[string]any{
        "authAgentId":   agentEntity.GetId(),
        "authAgentUuid": agentEntity.GetAgentUuid(),
        "authProjectId": agentEntity.GetProjectId(),
    }), nil
}

func (s *AuthenticateAgentService) failIfWrongInput() error {
    if s.authenticateAgentDto.GetApiToken() == "" {
        return exceptions.UnauthorizedCustom(config.Tr(enums.AgentsMessageEnum.INVALID_OR_MISSING_API_TOKEN))
    }
    return nil
}
