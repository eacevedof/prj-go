package updateagent

import (
    "siem-backend/src/core/config"
    "siem-backend/src/modules/agents/domain/entities"
    "siem-backend/src/modules/agents/domain/enums"
    "siem-backend/src/modules/agents/domain/exceptions"
    "siem-backend/src/modules/agents/domain/services"
    "siem-backend/src/modules/agents/infrastructure/repositories"
    sharedenums "siem-backend/src/modules/shared/domain/enums"
    "siem-backend/src/modules/shared/infrastructure/components/datetimer"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"
    "siem-backend/src/modules/shared/infrastructure/traits"
)

// UpdateAgentService cambia nombre, latido y si esta activo un agente desde el panel. El token no cambia. Desactivado,
// su token deja de valer (AuthenticateAgent lo rechaza con 401) hasta que se reactive: el agente no se borra.
type UpdateAgentService struct {
    traits.ErrorStackTrait

    dateTimer                    *datetimer.DateTimer
    agentsReaderSqliteRepository *repositories.AgentsReaderSqliteRepository
    agentsWriterSqliteRepository *repositories.AgentsWriterSqliteRepository

    updateAgentDto *UpdateAgentDto
    agentEntity    *entities.AgentEntity
}

func UpdateAgentServiceGetInstance() *UpdateAgentService {
    return &UpdateAgentService{
        dateTimer:                    datetimer.GetInstance(),
        agentsReaderSqliteRepository: repositories.AgentsReaderSqliteRepositoryGetInstance(),
        agentsWriterSqliteRepository: repositories.AgentsWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *FieldsException (400), *AgentsException (404, 409) o uno inesperado.
func (s *UpdateAgentService) Invoke(updateAgentDto *UpdateAgentDto) (*UpdateAgentResultDto, error) {
    s.updateAgentDto = updateAgentDto

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }
    if err := s.loadAgentOrFail(); err != nil {
        return nil, err
    }
    if err := s.failIfAgentNameExistsInProject(); err != nil {
        return nil, err
    }

    if _, err := s.agentsWriterSqliteRepository.UpdateAgentById(
        s.agentEntity.GetId(),
        updateAgentDto.GetAgentName(),
        updateAgentDto.GetHeartbeatIntervalSeconds(),
        updateAgentDto.IsEnabled(),
        s.dateTimer.GetNowDtIntoTargetTz(),
    ); err != nil {
        return nil, err
    }

    agentEntity, err := s.agentsReaderSqliteRepository.GetAgentById(s.agentEntity.GetId())
    if err != nil {
        return nil, err
    }
    return UpdateAgentResultDtoFromPrimitives(map[string]any{
        "agent": entities.ToPrimitivesOrEmpty(agentEntity),
    }), nil
}

func (s *UpdateAgentService) failIfWrongInput() error {
    errors := services.AgentFieldsValidatorGetInstance().GetAgentFieldsErrors(
        s.updateAgentDto.GetAgentName(),
        s.updateAgentDto.GetHeartbeatIntervalSeconds(),
    )
    if s.updateAgentDto.GetAgentUuid() == "" {
        errors = append(errors, enums.AgentsMessageEnum.AGENT_UUID_IS_REQUIRED)
    }
    if len(errors) == 0 {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
    )
}

func (s *UpdateAgentService) loadAgentOrFail() error {
    agentUuid := s.updateAgentDto.GetAgentUuid()
    agentEntity, err := s.agentsReaderSqliteRepository.GetAgentByAgentUuid(agentUuid)
    if err != nil {
        return err
    }
    if agentEntity == nil {
        return exceptions.NotFoundCustom(config.Tr(enums.AgentsMessageEnum.AGENT_X_NOT_FOUND, agentUuid))
    }
    s.agentEntity = agentEntity
    return nil
}

// failIfAgentNameExistsInProject el nombre es unico dentro del proyecto; el propio agente no cuenta.
func (s *UpdateAgentService) failIfAgentNameExistsInProject() error {
    agentName := s.updateAgentDto.GetAgentName()
    doesAgentNameExist, err := s.agentsReaderSqliteRepository.DoesAgentNameExistInOtherAgent(
        s.agentEntity.GetProjectId(), agentName, s.agentEntity.GetId(),
    )
    if err != nil {
        return err
    }
    if doesAgentNameExist {
        return exceptions.ConflictCustom(config.Tr(enums.AgentsMessageEnum.AGENT_NAME_X_ALREADY_EXISTS_IN_PROJECT, agentName))
    }
    return nil
}
