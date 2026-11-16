package updateproject

import (
    "siem-backend/src/core/config"
    "siem-backend/src/modules/projects/domain/entities"
    "siem-backend/src/modules/projects/domain/enums"
    "siem-backend/src/modules/projects/domain/exceptions"
    "siem-backend/src/modules/projects/domain/services"
    "siem-backend/src/modules/projects/infrastructure/repositories"
    sharedenums "siem-backend/src/modules/shared/domain/enums"
    "siem-backend/src/modules/shared/infrastructure/components/datetimer"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"
    "siem-backend/src/modules/shared/infrastructure/traits"
)

// UpdateProjectService cambia nombre y descripcion de un proyecto desde el panel. El uuid no cambia: agentes,
// monitores y eventos siguen colgando de el.
type UpdateProjectService struct {
    traits.ErrorStackTrait

    dateTimer                      *datetimer.DateTimer
    projectsReaderSqliteRepository *repositories.ProjectsReaderSqliteRepository
    projectsWriterSqliteRepository *repositories.ProjectsWriterSqliteRepository

    updateProjectDto *UpdateProjectDto
    projectEntity    *entities.ProjectEntity
}

func UpdateProjectServiceGetInstance() *UpdateProjectService {
    return &UpdateProjectService{
        dateTimer:                      datetimer.GetInstance(),
        projectsReaderSqliteRepository: repositories.ProjectsReaderSqliteRepositoryGetInstance(),
        projectsWriterSqliteRepository: repositories.ProjectsWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *FieldsException (400), *ProjectsException (404, 409) o uno inesperado.
func (s *UpdateProjectService) Invoke(updateProjectDto *UpdateProjectDto) (*UpdateProjectResultDto, error) {
    s.updateProjectDto = updateProjectDto

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }
    if err := s.failIfProjectNotFound(); err != nil {
        return nil, err
    }
    if err := s.failIfProjectNameExists(); err != nil {
        return nil, err
    }

    if _, err := s.projectsWriterSqliteRepository.UpdateProjectById(
        s.projectEntity.GetId(),
        updateProjectDto.GetProjectName(),
        updateProjectDto.GetProjectDescription(),
        s.dateTimer.GetNowDtIntoTargetTz(),
    ); err != nil {
        return nil, err
    }

    projectEntity, err := s.projectsReaderSqliteRepository.GetProjectById(s.projectEntity.GetId())
    if err != nil {
        return nil, err
    }
    return UpdateProjectResultDtoFromPrimitives(map[string]any{
        "project": entities.ToPrimitivesOrEmpty(projectEntity),
    }), nil
}

func (s *UpdateProjectService) failIfWrongInput() error {
    errors := services.ProjectFieldsValidatorGetInstance().GetErrors(
        s.updateProjectDto.GetProjectName(), s.updateProjectDto.GetProjectDescription(),
    )
    if len(errors) == 0 {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
    )
}

func (s *UpdateProjectService) failIfProjectNotFound() error {
    projectUuid := s.updateProjectDto.GetProjectUuid()
    projectEntity, err := s.projectsReaderSqliteRepository.GetProjectByProjectUuid(projectUuid)
    if err != nil {
        return err
    }
    if projectEntity == nil {
        return exceptions.NotFoundCustom(config.Tr(enums.ProjectsMessageEnum.PROJECT_X_NOT_FOUND, projectUuid))
    }
    s.projectEntity = projectEntity
    return nil
}

// failIfProjectNameExists el nombre es unico sin distinguir mayusculas; el propio proyecto no cuenta (se puede
// guardar sin cambiar el nombre, o cambiar solo sus mayusculas).
func (s *UpdateProjectService) failIfProjectNameExists() error {
    projectName := s.updateProjectDto.GetProjectName()
    doesProjectNameExist, err := s.projectsReaderSqliteRepository.DoesProjectNameExistInOtherProject(projectName, s.projectEntity.GetId())
    if err != nil {
        return err
    }
    if doesProjectNameExist {
        return exceptions.ConflictCustom(config.Tr(enums.ProjectsMessageEnum.PROJECT_NAME_X_ALREADY_EXISTS, projectName))
    }
    return nil
}
