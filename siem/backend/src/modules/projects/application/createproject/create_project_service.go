package createproject

import (
    "siem-backend/src/core/config"

    sharedenums "siem-backend/src/modules/shared/domain/enums"
    "siem-backend/src/modules/shared/infrastructure/components/datetimer"
    "siem-backend/src/modules/shared/infrastructure/components/uuid"
    "siem-backend/src/modules/shared/infrastructure/traits"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"

    "siem-backend/src/modules/projects/domain/enums"
    "siem-backend/src/modules/projects/domain/entities"
    "siem-backend/src/modules/projects/domain/services"
    "siem-backend/src/modules/projects/infrastructure/repositories"
    "siem-backend/src/modules/projects/domain/exceptions"
)

// CreateProjectService alta de un proyecto (lo que se monitoriza): desde el panel (POST) y por consola (project-create).
type CreateProjectService struct {
    traits.ErrorStackTrait

    dateTimer                      *datetimer.DateTimer
    uuidGenerator                  *uuid.UuidGenerator
    projectsReaderSqliteRepository *repositories.ProjectsReaderSqliteRepository
    projectsWriterSqliteRepository *repositories.ProjectsWriterSqliteRepository

    createProjectDto *CreateProjectDto
}

func CreateProjectServiceGetInstance() *CreateProjectService {
    return &CreateProjectService{
        dateTimer:                      datetimer.GetInstance(),
        uuidGenerator:                  uuid.GetInstance(),
        projectsReaderSqliteRepository: repositories.ProjectsReaderSqliteRepositoryGetInstance(),
        projectsWriterSqliteRepository: repositories.ProjectsWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *FieldsException (400), *ProjectsException (409) o uno inesperado.
func (s *CreateProjectService) Invoke(createProjectDto *CreateProjectDto) (*CreateProjectResultDto, error) {
    s.createProjectDto = createProjectDto

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }
    if err := s.failIfProjectNameExists(); err != nil {
        return nil, err
    }

    projectId, err := s.projectsWriterSqliteRepository.InsertProject(
        entities.ProjectEntityFromPrimitives(map[string]any{
            enums.ProjectsKeyEnum.PROJECT_UUID:        s.uuidGenerator.GetUniqueIdWithPrefix(entities.PROJECT_UUID_PREFIX),
            enums.ProjectsKeyEnum.PROJECT_NAME:        createProjectDto.GetProjectName(),
            enums.ProjectsKeyEnum.PROJECT_DESCRIPTION: createProjectDto.GetProjectDescription(),
            enums.ProjectsKeyEnum.CREATED_AT:          s.dateTimer.GetNowDtIntoTargetTz(),
        }),
    )
    if err != nil {
        return nil, err
    }

    projectEntity, err := s.projectsReaderSqliteRepository.GetProjectById(projectId)
    if err != nil {
        return nil, err
    }
    return CreateProjectResultDtoFromPrimitives(map[string]any{
        "project": entities.ToPrimitivesOrEmpty(projectEntity),
    }), nil
}

func (s *CreateProjectService) failIfWrongInput() error {
    errors := services.ProjectFieldsValidatorGetInstance().GetErrors(
        s.createProjectDto.GetProjectName(), s.createProjectDto.GetProjectDescription(),
    )
    if len(errors) == 0 {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
    )
}

func (s *CreateProjectService) failIfProjectNameExists() error {
    projectName := s.createProjectDto.GetProjectName()
    doesProjectNameExist, err := s.projectsReaderSqliteRepository.DoesProjectNameExist(projectName)
    if err != nil {
        return err
    }
    if doesProjectNameExist {
        return exceptions.ConflictCustom(config.Tr(enums.ProjectsMessageEnum.PROJECT_NAME_X_ALREADY_EXISTS, projectName))
    }
    return nil
}
