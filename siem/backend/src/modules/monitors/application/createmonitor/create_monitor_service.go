package createmonitor

import (
    "siem-backend/src/core/config"
    "siem-backend/src/modules/monitors/domain/entities"
    "siem-backend/src/modules/monitors/domain/enums"
    "siem-backend/src/modules/monitors/domain/exceptions"
    "siem-backend/src/modules/monitors/infrastructure/traits"
    projectsrepositories "siem-backend/src/modules/projects/infrastructure/repositories"
    "siem-backend/src/modules/shared/infrastructure/components/datetimer"
    "siem-backend/src/modules/shared/infrastructure/components/uuid"
)

// CreateMonitorService alta de un recurso a vigilar en un proyecto (desde el panel). El agente lo recibe en su
// siguiente GET /api/v1/agent-config.
type CreateMonitorService struct {
    traits.MonitorSaverTrait

    dateTimer                      *datetimer.DateTimer
    uuidGenerator                  *uuid.UuidGenerator
    projectsReaderSqliteRepository *projectsrepositories.ProjectsReaderSqliteRepository
}

func CreateMonitorServiceGetInstance() *CreateMonitorService {
    return &CreateMonitorService{
        MonitorSaverTrait:              traits.NewMonitorSaverTrait(),
        dateTimer:                      datetimer.GetInstance(),
        uuidGenerator:                  uuid.GetInstance(),
        projectsReaderSqliteRepository: projectsrepositories.ProjectsReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *MonitorsException (404 proyecto, 409 nombre), *FieldsException (400) o uno inesperado.
func (s *CreateMonitorService) Invoke(createMonitorDto *CreateMonitorDto) (*CreateMonitorResultDto, error) {
    projectEntity, err := s.projectsReaderSqliteRepository.GetProjectByProjectUuid(createMonitorDto.GetProjectUuid())
    if err != nil {
        return nil, err
    }
    if projectEntity == nil {
        return nil, exceptions.NotFoundCustom(config.Tr(enums.MonitorsMessageEnum.PROJECT_X_NOT_FOUND, createMonitorDto.GetProjectUuid()))
    }

    monitorType := createMonitorDto.GetMonitorType()
    if err := s.FailIfWrongInput(createMonitorDto.MonitorInput, monitorType); err != nil {
        return nil, err
    }
    agentId, err := s.GetAgentIdOrFail(createMonitorDto.GetAgentUuid(), projectEntity.GetId())
    if err != nil {
        return nil, err
    }
    if err := s.FailIfMonitorNameExists(projectEntity.GetId(), createMonitorDto.GetMonitorName(), 0); err != nil {
        return nil, err
    }

    warningThreshold, criticalThreshold := s.GetThresholdsOrNil(createMonitorDto.MonitorInput, monitorType)
    monitorId, err := s.MonitorsWriterSqliteRepository.InsertMonitor(entities.MonitorEntityFromPrimitives(map[string]any{
        enums.MonitorsKeyEnum.MONITOR_UUID:       s.uuidGenerator.GetUniqueIdWithPrefix(entities.MONITOR_UUID_PREFIX),
        enums.MonitorsKeyEnum.PROJECT_ID:         projectEntity.GetId(),
        enums.MonitorsKeyEnum.AGENT_ID:           agentId,
        enums.MonitorsKeyEnum.MONITOR_TYPE:       monitorType,
        enums.MonitorsKeyEnum.MONITOR_NAME:       createMonitorDto.GetMonitorName(),
        enums.MonitorsKeyEnum.IS_ENABLED:         createMonitorDto.IsEnabled(),
        enums.MonitorsKeyEnum.INTERVAL_SECONDS:   createMonitorDto.GetIntervalSeconds(),
        enums.MonitorsKeyEnum.CONFIG:             s.MonitorFieldsValidator.GetNormalizedConfig(monitorType, createMonitorDto.GetMonitorConfig()),
        enums.MonitorsKeyEnum.WARNING_THRESHOLD:  warningThreshold,
        enums.MonitorsKeyEnum.CRITICAL_THRESHOLD: criticalThreshold,
        enums.MonitorsKeyEnum.CREATED_AT:         s.dateTimer.GetNowDtIntoTargetTz(),
    }))
    if err != nil {
        return nil, err
    }

    monitorEntity, err := s.MonitorsReaderSqliteRepository.GetMonitorById(monitorId)
    if err != nil {
        return nil, err
    }
    return CreateMonitorResultDtoFromPrimitives(map[string]any{
        "monitor": entities.ToPrimitivesOrEmpty(monitorEntity),
    }), nil
}
