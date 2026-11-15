package deleteproject

import (
    "siem-backend/src/core/config"
    "siem-backend/src/modules/projects/domain/enums"
    "siem-backend/src/modules/projects/domain/exceptions"
    "siem-backend/src/modules/projects/infrastructure/repositories"
)

// DeleteProjectService baja de un proyecto desde el panel. La BD se lleva en cascada sus agentes, monitores, fuentes
// y eventos: no hay vuelta atras (el panel pide confirmacion antes de llamar).
type DeleteProjectService struct {
    projectsReaderSqliteRepository *repositories.ProjectsReaderSqliteRepository
    projectsWriterSqliteRepository *repositories.ProjectsWriterSqliteRepository
}

func DeleteProjectServiceGetInstance() *DeleteProjectService {
    return &DeleteProjectService{
        projectsReaderSqliteRepository: repositories.ProjectsReaderSqliteRepositoryGetInstance(),
        projectsWriterSqliteRepository: repositories.ProjectsWriterSqliteRepositoryGetInstance(),
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
    if _, err := s.projectsWriterSqliteRepository.DeleteProjectById(projectEntity.GetId()); err != nil {
        return nil, err
    }
    return DeleteProjectResultDtoFromPrimitives(map[string]any{"projectUuid": projectEntity.GetProjectUuid()}), nil
}
