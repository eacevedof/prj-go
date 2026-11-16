package createagent

import (
    "siem-backend/src/core/config"

    sharedenums "siem-backend/src/modules/shared/domain/enums"
    "siem-backend/src/modules/shared/infrastructure/components/datetimer"
    "siem-backend/src/modules/shared/infrastructure/components/token"
    "siem-backend/src/modules/shared/infrastructure/components/uuid"
    "siem-backend/src/modules/shared/infrastructure/traits"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"

    projectsentities "siem-backend/src/modules/projects/domain/entities"
    projectsrepositories "siem-backend/src/modules/projects/infrastructure/repositories"

    "siem-backend/src/modules/agents/domain/enums"
    "siem-backend/src/modules/agents/domain/entities"
    "siem-backend/src/modules/agents/domain/services"
    "siem-backend/src/modules/agents/infrastructure/repositories"
    "siem-backend/src/modules/agents/domain/exceptions"
)

// CreateAgentService alta de un agente (una instalacion del siem-client) dentro de un proyecto. Devuelve su
// api_token en claro: es la unica vez que se ve (en BD solo queda el sha256). Solo por consola: ./console agent-create
type CreateAgentService struct {
    traits.ErrorStackTrait

    dateTimer                      *datetimer.DateTimer
    tokenGenerator                 *token.TokenGenerator
    uuidGenerator                  *uuid.UuidGenerator
    projectsReaderSqliteRepository *projectsrepositories.ProjectsReaderSqliteRepository
    agentsReaderSqliteRepository   *repositories.AgentsReaderSqliteRepository
    agentsWriterSqliteRepository   *repositories.AgentsWriterSqliteRepository

    createAgentDto *CreateAgentDto
    projectEntity  *projectsentities.ProjectEntity
}

func CreateAgentServiceGetInstance() *CreateAgentService {
    return &CreateAgentService{
        dateTimer:                      datetimer.GetInstance(),
        tokenGenerator:                 token.GetInstance(),
        uuidGenerator:                  uuid.GetInstance(),
        projectsReaderSqliteRepository: projectsrepositories.ProjectsReaderSqliteRepositoryGetInstance(),
        agentsReaderSqliteRepository:   repositories.AgentsReaderSqliteRepositoryGetInstance(),
        agentsWriterSqliteRepository:   repositories.AgentsWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *FieldsException (400), *AgentsException (404, 409) o uno inesperado.
func (s *CreateAgentService) Invoke(createAgentDto *CreateAgentDto) (*CreateAgentResultDto, error) {
    s.createAgentDto = createAgentDto

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }
    if err := s.loadProjectOrFail(); err != nil {
        return nil, err
    }
    if err := s.failIfAgentNameExistsInProject(); err != nil {
        return nil, err
    }

    apiToken := s.tokenGenerator.GetNewToken()
    agentId, err := s.agentsWriterSqliteRepository.InsertAgent(
        entities.AgentEntityFromPrimitives(map[string]any{
            enums.AgentsKeyEnum.PROJECT_ID:                 s.projectEntity.GetId(),
            enums.AgentsKeyEnum.AGENT_UUID:                 s.uuidGenerator.GetUniqueIdWithPrefix(entities.AGENT_UUID_PREFIX),
            enums.AgentsKeyEnum.AGENT_NAME:                 createAgentDto.GetAgentName(),
            enums.AgentsKeyEnum.API_TOKEN_HASH:             s.tokenGenerator.GetTokenHash(apiToken),
            enums.AgentsKeyEnum.HEARTBEAT_INTERVAL_SECONDS: createAgentDto.GetHeartbeatIntervalSeconds(),
            enums.AgentsKeyEnum.CREATED_AT:                 s.dateTimer.GetNowDtIntoTargetTz(),
        }),
    )
    if err != nil {
        return nil, err
    }

    agentEntity, err := s.agentsReaderSqliteRepository.GetAgentById(agentId)
    if err != nil {
        return nil, err
    }
    return CreateAgentResultDtoFromPrimitives(map[string]any{
        "agent":    entities.ToPrimitivesOrEmpty(agentEntity),
        "apiToken": apiToken,
    }), nil
}

func (s *CreateAgentService) failIfWrongInput() error {
    errors := services.AgentFieldsValidatorGetInstance().GetErrors(
        s.createAgentDto.GetProjectUuid(),
        s.createAgentDto.GetAgentName(),
        s.createAgentDto.GetHeartbeatIntervalSeconds(),
    )
    if len(errors) == 0 {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
    )
}

func (s *CreateAgentService) loadProjectOrFail() error {
    projectUuid := s.createAgentDto.GetProjectUuid()
    projectEntity, err := s.projectsReaderSqliteRepository.GetProjectByProjectUuid(projectUuid)
    if err != nil {
        return err
    }
    if projectEntity == nil {
        return exceptions.NotFoundCustom(config.Tr(enums.AgentsMessageEnum.PROJECT_X_NOT_FOUND, projectUuid))
    }
    s.projectEntity = projectEntity
    return nil
}

func (s *CreateAgentService) failIfAgentNameExistsInProject() error {
    agentName := s.createAgentDto.GetAgentName()
    doesAgentNameExist, err := s.agentsReaderSqliteRepository.DoesAgentNameExistInProject(s.projectEntity.GetId(), agentName)
    if err != nil {
        return err
    }
    if doesAgentNameExist {
        return exceptions.ConflictCustom(config.Tr(enums.AgentsMessageEnum.AGENT_NAME_X_ALREADY_EXISTS_IN_PROJECT, agentName))
    }
    return nil
}
