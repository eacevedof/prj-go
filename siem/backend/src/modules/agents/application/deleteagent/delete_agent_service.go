package deleteagent

import (
    "siem-backend/src/core/config"
    "siem-backend/src/modules/agents/domain/enums"
    "siem-backend/src/modules/agents/domain/exceptions"
    "siem-backend/src/modules/agents/infrastructure/repositories"
)

// DeleteAgentService baja de un agente desde el panel: su token deja de valer en el acto (el siem-client de ese host
// recibira 401) y la BD se lleva en cascada sus fuentes, eventos, cola y monitores propios.
type DeleteAgentService struct {
    agentsReaderSqliteRepository *repositories.AgentsReaderSqliteRepository
    agentsWriterSqliteRepository *repositories.AgentsWriterSqliteRepository
}

func DeleteAgentServiceGetInstance() *DeleteAgentService {
    return &DeleteAgentService{
        agentsReaderSqliteRepository: repositories.AgentsReaderSqliteRepositoryGetInstance(),
        agentsWriterSqliteRepository: repositories.AgentsWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *AgentsException (404) o uno inesperado.
func (s *DeleteAgentService) Invoke(deleteAgentDto *DeleteAgentDto) (*DeleteAgentResultDto, error) {
    agentEntity, err := s.agentsReaderSqliteRepository.GetAgentByAgentUuid(deleteAgentDto.GetAgentUuid())
    if err != nil {
        return nil, err
    }
    if agentEntity == nil {
        return nil, exceptions.NotFoundCustom(config.Tr(enums.AgentsMessageEnum.AGENT_X_NOT_FOUND, deleteAgentDto.GetAgentUuid()))
    }
    if _, err := s.agentsWriterSqliteRepository.DeleteAgentById(agentEntity.GetId()); err != nil {
        return nil, err
    }
    return DeleteAgentResultDtoFromPrimitives(map[string]any{"agentUuid": agentEntity.GetAgentUuid()}), nil
}
