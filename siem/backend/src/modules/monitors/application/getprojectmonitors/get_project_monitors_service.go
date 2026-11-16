package getprojectmonitors

import (
    "strings"

    "siem-backend/src/core/config"

    "siem-backend/src/modules/shared/application"

    projectsrepositories "siem-backend/src/modules/projects/infrastructure/repositories"

    "siem-backend/src/modules/monitors/domain/enums"
    "siem-backend/src/modules/monitors/infrastructure/repositories"
    "siem-backend/src/modules/monitors/domain/exceptions"
)

// GetProjectMonitorsService los monitores de un proyecto (activos o no) para la tabla del panel, cada uno con su
// ultimo valor en cada agente que ya haya mandado algo.
type GetProjectMonitorsService struct {
    projectsReaderSqliteRepository *projectsrepositories.ProjectsReaderSqliteRepository
    monitorsReaderSqliteRepository *repositories.MonitorsReaderSqliteRepository
}

func GetProjectMonitorsServiceGetInstance() *GetProjectMonitorsService {
    return &GetProjectMonitorsService{
        projectsReaderSqliteRepository: projectsrepositories.ProjectsReaderSqliteRepositoryGetInstance(),
        monitorsReaderSqliteRepository: repositories.MonitorsReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *MonitorsException (404) o uno inesperado.
func (s *GetProjectMonitorsService) Invoke(getProjectMonitorsDto *GetProjectMonitorsDto) (*GetProjectMonitorsResultDto, error) {
    projectEntity, err := s.projectsReaderSqliteRepository.GetProjectByProjectUuid(getProjectMonitorsDto.GetProjectUuid())
    if err != nil {
        return nil, err
    }
    if projectEntity == nil {
        return nil, exceptions.NotFoundCustom(config.Tr(enums.MonitorsMessageEnum.PROJECT_X_NOT_FOUND, getProjectMonitorsDto.GetProjectUuid()))
    }

    monitorEntities, err := s.monitorsReaderSqliteRepository.GetMonitorsByProjectId(projectEntity.GetId())
    if err != nil {
        return nil, err
    }
    statesByMonitorName, err := s.getStatesByMonitorName(projectEntity.GetId())
    if err != nil {
        return nil, err
    }

    monitors := make([]MonitorWithStatesPrimitives, 0, len(monitorEntities))
    for _, monitorEntity := range monitorEntities {
        states := statesByMonitorName[strings.ToLower(monitorEntity.GetMonitorName())]
        if states == nil {
            states = []MonitorStatePrimitives{}
        }
        monitors = append(monitors, MonitorWithStatesPrimitives{MonitorPrimitives: monitorEntity.ToPrimitives(), States: states})
    }
    return GetProjectMonitorsResultDtoFromPrimitives(map[string]any{"monitors": monitors}), nil
}

// getStatesByMonitorName por nombre en minusculas: monitor_name es unico en el proyecto sin distinguir mayusculas.
func (s *GetProjectMonitorsService) getStatesByMonitorName(projectId int) (map[string][]MonitorStatePrimitives, error) {
    rows, err := s.monitorsReaderSqliteRepository.GetMonitorStatesByProjectId(projectId)
    if err != nil {
        return nil, err
    }
    statesByMonitorName := map[string][]MonitorStatePrimitives{}
    for _, row := range rows {
        monitorName := strings.ToLower(application.GetTrimmedString(row, enums.MonitorsKeyEnum.SOURCE_NAME))
        statesByMonitorName[monitorName] = append(statesByMonitorName[monitorName], MonitorStatePrimitives{
            AgentUuid:    application.GetTrimmedString(row, enums.MonitorsKeyEnum.AGENT_UUID),
            AgentName:    application.GetTrimmedString(row, enums.MonitorsKeyEnum.AGENT_NAME),
            LastEventAt:  application.GetTrimmedString(row, enums.MonitorsKeyEnum.LAST_EVENT_AT),
            LastSeverity: application.GetTrimmedString(row, enums.MonitorsKeyEnum.LAST_SEVERITY),
            LastMessage:  application.GetTrimmedString(row, enums.MonitorsKeyEnum.LAST_MESSAGE),
        })
    }
    return statesByMonitorName, nil
}
