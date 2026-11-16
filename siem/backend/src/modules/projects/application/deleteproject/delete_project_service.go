package deleteproject

import (
    "siem-backend/src/core/config"
    agentsentities "siem-backend/src/modules/agents/domain/entities"
    agentsrepositories "siem-backend/src/modules/agents/infrastructure/repositories"
    "siem-backend/src/modules/projects/domain/enums"
    "siem-backend/src/modules/projects/domain/exceptions"
    "siem-backend/src/modules/projects/infrastructure/repositories"
)

// DeleteProjectService baja de un proyecto desde el panel. La BD se lleva en cascada sus agentes, monitores, fuentes
// y eventos: no hay vuelta atras (el panel pide el nombre antes de llamar). Devuelve los agentes que tenia (host y
// ruta del ejecutable) porque sus siem-client siguen instalados y hay que desinstalarlos a mano.
type DeleteProjectService struct {
    projectsReaderSqliteRepository *repositories.ProjectsReaderSqliteRepository
    projectsWriterSqliteRepository *repositories.ProjectsWriterSqliteRepository
    agentsReaderSqliteRepository   *agentsrepositories.AgentsReaderSqliteRepository
}

func DeleteProjectServiceGetInstance() *DeleteProjectService {
    return &DeleteProjectService{
        projectsReaderSqliteRepository: repositories.ProjectsReaderSqliteRepositoryGetInstance(),
        projectsWriterSqliteRepository: repositories.ProjectsWriterSqliteRepositoryGetInstance(),
        agentsReaderSqliteRepository:   agentsrepositories.AgentsReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *ProjectsException (404) o uno inesperado.
func (s *DeleteProjectService) Invoke(deleteProjectDto *DeleteProjectDto) (*DeleteProjectResultDto, error) {
    projectEntity, err := s.projectsReaderSqliteRepository.GetProjectByProjectUuid(deleteProjectDto.GetProjectUuid())
    if err != nil {
        return nil, err
    }
    if projectEntity == nil {
        return nil, exceptions.NotFoundCustom(config.Tr(enums.ProjectsMessageEnum.PROJECT_X_NOT_FOUND, deleteProjectDto.GetProjectUuid()))
    }

    // antes del borrado: la cascada se lleva los agentes
    agentEntities, err := s.agentsReaderSqliteRepository.GetAgentsByProjectId(projectEntity.GetId())
    if err != nil {
        return nil, err
    }
    if _, err := s.projectsWriterSqliteRepository.DeleteProjectById(projectEntity.GetId()); err != nil {
        return nil, err
    }
    return DeleteProjectResultDtoFromPrimitives(map[string]any{
        "projectUuid": projectEntity.GetProjectUuid(),
        "projectName": projectEntity.GetProjectName(),
        "agents":      s.getDeletedAgents(agentEntities),
    }), nil
}

func (s *DeleteProjectService) getDeletedAgents(agentEntities []*agentsentities.AgentEntity) []DeletedAgentPrimitives {
    deletedAgents := make([]DeletedAgentPrimitives, 0, len(agentEntities))
    for _, agentEntity := range agentEntities {
        deletedAgents = append(deletedAgents, DeletedAgentPrimitives{
            AgentName:  agentEntity.GetAgentName(),
            HostName:   agentEntity.GetHostName(),
            HostOs:     agentEntity.GetHostOs(),
            HostIp:     agentEntity.GetHostIp(),
            ClientPath: agentEntity.GetClientPath(),
        })
    }
    return deletedAgents
}
