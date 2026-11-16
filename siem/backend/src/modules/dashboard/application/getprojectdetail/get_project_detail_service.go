package getprojectdetail

import (
    "siem-backend/src/core/config"
    agentsentities "siem-backend/src/modules/agents/domain/entities"
    agentsenums "siem-backend/src/modules/agents/domain/enums"
    agentsservices "siem-backend/src/modules/agents/domain/services"
    agentsrepositories "siem-backend/src/modules/agents/infrastructure/repositories"
    "siem-backend/src/modules/dashboard/domain/enums"
    "siem-backend/src/modules/dashboard/domain/exceptions"
    "siem-backend/src/modules/dashboard/domain/services"
    "siem-backend/src/modules/dashboard/infrastructure/repositories"
    eventsenums "siem-backend/src/modules/events/domain/enums"
    eventsrepositories "siem-backend/src/modules/events/infrastructure/repositories"
    projectsentities "siem-backend/src/modules/projects/domain/entities"
    projectsrepositories "siem-backend/src/modules/projects/infrastructure/repositories"
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/components/datetimer"
)

const last24hSeconds = 86400

// GetProjectDetailService el detalle de un proyecto para el front: sus agentes (estado, host, contadores de 24 h) y,
// por agente, sus sources (cada carpeta vigilada y cada comando) con su ultimo estado.
type GetProjectDetailService struct {
    dateTimer                       *datetimer.DateTimer
    agentStatusResolver             *agentsservices.AgentStatusResolver
    eventsStatsCounter              *services.EventsStatsCounter
    projectsReaderSqliteRepository  *projectsrepositories.ProjectsReaderSqliteRepository
    agentsReaderSqliteRepository    *agentsrepositories.AgentsReaderSqliteRepository
    sourcesReaderSqliteRepository   *eventsrepositories.SourcesReaderSqliteRepository
    dashboardReaderSqliteRepository *repositories.DashboardReaderSqliteRepository

    getProjectDetailDto *GetProjectDetailDto
    projectEntity       *projectsentities.ProjectEntity
    nowDt               string
}

func GetProjectDetailServiceGetInstance() *GetProjectDetailService {
    return &GetProjectDetailService{
        dateTimer:                       datetimer.GetInstance(),
        agentStatusResolver:             agentsservices.AgentStatusResolverGetInstance(),
        eventsStatsCounter:              services.EventsStatsCounterGetInstance(),
        projectsReaderSqliteRepository:  projectsrepositories.ProjectsReaderSqliteRepositoryGetInstance(),
        agentsReaderSqliteRepository:    agentsrepositories.AgentsReaderSqliteRepositoryGetInstance(),
        sourcesReaderSqliteRepository:   eventsrepositories.SourcesReaderSqliteRepositoryGetInstance(),
        dashboardReaderSqliteRepository: repositories.DashboardReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *DashboardException (401, 404) o uno inesperado.
func (s *GetProjectDetailService) Invoke(getProjectDetailDto *GetProjectDetailDto) (*GetProjectDetailResultDto, error) {
    s.getProjectDetailDto = getProjectDetailDto
    s.nowDt = s.dateTimer.GetNowDtIntoTargetTz()

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }
    if err := s.loadProjectOrFail(); err != nil {
        return nil, err
    }

    agentEntities, err := s.agentsReaderSqliteRepository.GetAgentsByProjectId(s.projectEntity.GetId())
    if err != nil {
        return nil, err
    }
    sourcesByAgentId, err := s.getSourcesByAgentId()
    if err != nil {
        return nil, err
    }

    agents := make([]map[string]any, 0, len(agentEntities))
    for _, agentEntity := range agentEntities {
        agent, err := s.getAgent(agentEntity, sourcesByAgentId[agentEntity.GetId()])
        if err != nil {
            return nil, err
        }
        agents = append(agents, agent)
    }
    return GetProjectDetailResultDtoFromPrimitives(map[string]any{
        "project": s.projectEntity.ToPrimitives(),
        "agents":  agents,
    }), nil
}

func (s *GetProjectDetailService) failIfWrongInput() error {
    if !s.getProjectDetailDto.IsDashboardAuthenticated() {
        return exceptions.UnauthorizedInvalidToken()
    }
    return nil
}

func (s *GetProjectDetailService) loadProjectOrFail() error {
    projectUuid := s.getProjectDetailDto.GetProjectUuid()
    projectEntity, err := s.projectsReaderSqliteRepository.GetProjectByProjectUuid(projectUuid)
    if err != nil {
        return err
    }
    if projectEntity == nil {
        return exceptions.NotFoundCustom(config.Tr(enums.DashboardMessageEnum.PROJECT_X_NOT_FOUND, projectUuid))
    }
    s.projectEntity = projectEntity
    return nil
}

// getSourcesByAgentId las sources del proyecto repartidas por agente, sin la columna interna agent_id.
func (s *GetProjectDetailService) getSourcesByAgentId() (map[int][]map[string]any, error) {
    sources, err := s.sourcesReaderSqliteRepository.GetSourcesByProjectId(s.projectEntity.GetId())
    if err != nil {
        return nil, err
    }
    sourcesByAgentId := map[int][]map[string]any{}
    for _, source := range sources {
        agentId := application.ToInt(source[eventsenums.EventsKeyEnum.AGENT_ID])
        delete(source, eventsenums.EventsKeyEnum.AGENT_ID)
        sourcesByAgentId[agentId] = append(sourcesByAgentId[agentId], source)
    }
    return sourcesByAgentId, nil
}

func (s *GetProjectDetailService) getAgent(agentEntity *agentsentities.AgentEntity, sources []map[string]any) (map[string]any, error) {
    eventsStats, err := s.dashboardReaderSqliteRepository.GetEventsStatsByAgentIdSince(
        agentEntity.GetId(), s.dateTimer.GetNowDtMinusSeconds(last24hSeconds),
    )
    if err != nil {
        return nil, err
    }
    eventsCountBySeverity := map[string]int{}
    for _, severityStats := range eventsStats {
        eventsCountBySeverity[application.ToString(severityStats["severity"])] = application.ToInt(severityStats["n_events"])
    }
    if sources == nil {
        sources = []map[string]any{}
    }

    return map[string]any{
        enums.DashboardKeyEnum.AGENT_ID:      agentEntity.GetId(),
        agentsenums.AgentsKeyEnum.AGENT_UUID: agentEntity.GetAgentUuid(),
        agentsenums.AgentsKeyEnum.AGENT_NAME: agentEntity.GetAgentName(),
        agentsenums.AgentsKeyEnum.AGENT_STATUS: s.agentStatusResolver.GetAgentStatus(
            agentEntity.IsEnabled(), agentEntity.GetLastSeenAt(), agentEntity.GetHeartbeatIntervalSeconds(), s.nowDt,
        ),
        agentsenums.AgentsKeyEnum.HEARTBEAT_INTERVAL_SECONDS: agentEntity.GetHeartbeatIntervalSeconds(),
        agentsenums.AgentsKeyEnum.IS_ENABLED:                 agentEntity.IsEnabled(),
        agentsenums.AgentsKeyEnum.HOST_NAME:                  agentEntity.GetHostName(),
        agentsenums.AgentsKeyEnum.HOST_OS:                    agentEntity.GetHostOs(),
        agentsenums.AgentsKeyEnum.HOST_IP:                    agentEntity.GetHostIp(),
        agentsenums.AgentsKeyEnum.CLIENT_VERSION:             agentEntity.GetClientVersion(),
        agentsenums.AgentsKeyEnum.CLIENT_PATH:                agentEntity.GetClientPath(),
        agentsenums.AgentsKeyEnum.LAST_SEEN_AT:               agentEntity.GetLastSeenAt(),
        enums.DashboardKeyEnum.EVENTS_24H:                    s.eventsStatsCounter.GetTotal(eventsCountBySeverity),
        enums.DashboardKeyEnum.WARNINGS_24H:                  eventsCountBySeverity[eventsenums.EventSeverityEnum.WARNING],
        enums.DashboardKeyEnum.ERRORS_24H:                    eventsCountBySeverity[eventsenums.EventSeverityEnum.ERROR],
        enums.DashboardKeyEnum.CRITICALS_24H:                 eventsCountBySeverity[eventsenums.EventSeverityEnum.CRITICAL],
        enums.DashboardKeyEnum.SOURCES:                       sources,
    }, nil
}
