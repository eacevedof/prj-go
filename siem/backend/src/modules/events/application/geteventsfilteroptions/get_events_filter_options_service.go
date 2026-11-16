package geteventsfilteroptions

import (
    "siem-backend/src/modules/shared/application"

    dashboardexceptions "siem-backend/src/modules/dashboard/domain/exceptions"

    "siem-backend/src/modules/events/domain/enums"
    "siem-backend/src/modules/events/infrastructure/repositories"
)

// GetEventsFilterOptionsService los proyectos, agentes y monitores que hay, para los desplegables del visor de eventos
// (Proyecto -> Agente -> Monitor). Las facetas solo traen lo que tiene eventos en el rango; aqui va todo.
type GetEventsFilterOptionsService struct {
    eventsFilterOptionsReaderSqliteRepository *repositories.EventsFilterOptionsReaderSqliteRepository
}

func GetEventsFilterOptionsServiceGetInstance() *GetEventsFilterOptionsService {
    return &GetEventsFilterOptionsService{
        eventsFilterOptionsReaderSqliteRepository: repositories.EventsFilterOptionsReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *DashboardException (401) o uno inesperado.
func (s *GetEventsFilterOptionsService) Invoke(getEventsFilterOptionsDto *GetEventsFilterOptionsDto) (*GetEventsFilterOptionsResultDto, error) {
    if !getEventsFilterOptionsDto.IsDashboardAuthenticated() {
        return nil, dashboardexceptions.UnauthorizedInvalidToken()
    }
    projectRows, err := s.eventsFilterOptionsReaderSqliteRepository.GetProjectOptions()
    if err != nil {
        return nil, err
    }
    agentRows, err := s.eventsFilterOptionsReaderSqliteRepository.GetAgentOptions()
    if err != nil {
        return nil, err
    }
    monitorRows, err := s.eventsFilterOptionsReaderSqliteRepository.GetMonitorOptions()
    if err != nil {
        return nil, err
    }
    return GetEventsFilterOptionsResultDtoFromPrimitives(map[string]any{
        "projects": s.getProjectOptions(projectRows),
        "agents":   s.getAgentOptions(agentRows),
        "monitors": s.getMonitorOptions(monitorRows),
    }), nil
}

func (s *GetEventsFilterOptionsService) getProjectOptions(rows []map[string]any) []ProjectOptionPrimitives {
    projectOptions := make([]ProjectOptionPrimitives, 0, len(rows))
    for _, row := range rows {
        projectOptions = append(projectOptions, ProjectOptionPrimitives{
            ProjectUuid: application.GetTrimmedString(row, enums.EventsFilterOptionKeyEnum.PROJECT_UUID),
            ProjectName: application.GetTrimmedString(row, enums.EventsFilterOptionKeyEnum.PROJECT_NAME),
        })
    }
    return projectOptions
}

func (s *GetEventsFilterOptionsService) getAgentOptions(rows []map[string]any) []AgentOptionPrimitives {
    agentOptions := make([]AgentOptionPrimitives, 0, len(rows))
    for _, row := range rows {
        agentOptions = append(agentOptions, AgentOptionPrimitives{
            AgentUuid:   application.GetTrimmedString(row, enums.EventsFilterOptionKeyEnum.AGENT_UUID),
            AgentName:   application.GetTrimmedString(row, enums.EventsFilterOptionKeyEnum.AGENT_NAME),
            ProjectUuid: application.GetTrimmedString(row, enums.EventsFilterOptionKeyEnum.PROJECT_UUID),
        })
    }
    return agentOptions
}

func (s *GetEventsFilterOptionsService) getMonitorOptions(rows []map[string]any) []MonitorOptionPrimitives {
    monitorOptions := make([]MonitorOptionPrimitives, 0, len(rows))
    for _, row := range rows {
        monitorOptions = append(monitorOptions, MonitorOptionPrimitives{
            MonitorUuid: application.GetTrimmedString(row, enums.EventsFilterOptionKeyEnum.MONITOR_UUID),
            MonitorName: application.GetTrimmedString(row, enums.EventsFilterOptionKeyEnum.MONITOR_NAME),
            MonitorType: application.GetTrimmedString(row, enums.EventsFilterOptionKeyEnum.MONITOR_TYPE),
            ProjectUuid: application.GetTrimmedString(row, enums.EventsFilterOptionKeyEnum.PROJECT_UUID),
            AgentUuid:   application.GetTrimmedString(row, enums.EventsFilterOptionKeyEnum.AGENT_UUID),
        })
    }
    return monitorOptions
}
