package getdashboard

import (
    agentsentities "siem-backend/src/modules/agents/domain/entities"
    agentsenums "siem-backend/src/modules/agents/domain/enums"
    agentsservices "siem-backend/src/modules/agents/domain/services"
    agentsrepositories "siem-backend/src/modules/agents/infrastructure/repositories"
    "siem-backend/src/modules/dashboard/domain/enums"
    "siem-backend/src/modules/dashboard/domain/exceptions"
    "siem-backend/src/modules/dashboard/domain/services"
    "siem-backend/src/modules/dashboard/infrastructure/repositories"
    eventsenums "siem-backend/src/modules/events/domain/enums"
    ingestenums "siem-backend/src/modules/ingest/domain/enums"
    projectsentities "siem-backend/src/modules/projects/domain/entities"
    projectsrepositories "siem-backend/src/modules/projects/infrastructure/repositories"
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/components/datetimer"
    "siem-backend/src/modules/shared/infrastructure/repositories/configuration"
)

const last24hSeconds = 86400

// GetDashboardService la portada del front: el estado del servidor (cola, BD, uptime) y un semaforo por proyecto
// con sus agentes online y sus eventos de las ultimas 24 h.
type GetDashboardService struct {
    dateTimer                       *datetimer.DateTimer
    environmentReaderRawRepository  *configuration.EnvironmentReaderRawRepository
    agentStatusResolver             *agentsservices.AgentStatusResolver
    projectHealthResolver           *services.ProjectHealthResolver
    eventsStatsCounter              *services.EventsStatsCounter
    projectsReaderSqliteRepository  *projectsrepositories.ProjectsReaderSqliteRepository
    agentsReaderSqliteRepository    *agentsrepositories.AgentsReaderSqliteRepository
    dashboardReaderSqliteRepository *repositories.DashboardReaderSqliteRepository

    getDashboardDto *GetDashboardDto
    nowDt           string
}

func GetDashboardServiceGetInstance() *GetDashboardService {
    return &GetDashboardService{
        dateTimer:                       datetimer.GetInstance(),
        environmentReaderRawRepository:  configuration.GetInstance(),
        agentStatusResolver:             agentsservices.AgentStatusResolverGetInstance(),
        projectHealthResolver:           services.ProjectHealthResolverGetInstance(),
        eventsStatsCounter:              services.EventsStatsCounterGetInstance(),
        projectsReaderSqliteRepository:  projectsrepositories.ProjectsReaderSqliteRepositoryGetInstance(),
        agentsReaderSqliteRepository:    agentsrepositories.AgentsReaderSqliteRepositoryGetInstance(),
        dashboardReaderSqliteRepository: repositories.DashboardReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *DashboardException (401) o uno inesperado.
func (s *GetDashboardService) Invoke(getDashboardDto *GetDashboardDto) (*GetDashboardResultDto, error) {
    s.getDashboardDto = getDashboardDto
    s.nowDt = s.dateTimer.GetNowDtIntoTargetTz()

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }

    server, err := s.getServer()
    if err != nil {
        return nil, err
    }
    projects, err := s.getProjects()
    if err != nil {
        return nil, err
    }
    return GetDashboardResultDtoFromPrimitives(map[string]any{"server": server, "projects": projects}), nil
}

func (s *GetDashboardService) failIfWrongInput() error {
    if !s.getDashboardDto.IsDashboardAuthenticated() {
        return exceptions.UnauthorizedInvalidToken()
    }
    return nil
}

func (s *GetDashboardService) getServer() (map[string]any, error) {
    eventsTotal, err := s.dashboardReaderSqliteRepository.GetEventsTotal()
    if err != nil {
        return nil, err
    }
    dbSizeBytes, err := s.dashboardReaderSqliteRepository.GetDbSizeBytes()
    if err != nil {
        return nil, err
    }
    queueStatsByStatus, err := s.dashboardReaderSqliteRepository.GetQueueStatsByStatus()
    if err != nil {
        return nil, err
    }

    startedAt := s.environmentReaderRawRepository.GetAppStartedAt()
    server := map[string]any{
        enums.DashboardKeyEnum.APP_NAME:                s.environmentReaderRawRepository.GetAppName(),
        enums.DashboardKeyEnum.APP_VERSION:             s.environmentReaderRawRepository.GetAppVersion(),
        enums.DashboardKeyEnum.APP_ENV:                 s.environmentReaderRawRepository.GetEnvironment(),
        enums.DashboardKeyEnum.NOW:                     s.nowDt,
        enums.DashboardKeyEnum.STARTED_AT:              startedAt,
        enums.DashboardKeyEnum.UPTIME_SECONDS:          s.dateTimer.GetSecondsBetween(startedAt, s.nowDt),
        enums.DashboardKeyEnum.DB_SIZE_BYTES:           dbSizeBytes,
        enums.DashboardKeyEnum.EVENTS_TOTAL:            eventsTotal,
        enums.DashboardKeyEnum.EVENTS_RETENTION_DAYS:   s.environmentReaderRawRepository.GetEventsRetentionDays(),
        enums.DashboardKeyEnum.QUEUE_WORKER_INTERVAL:   s.environmentReaderRawRepository.GetQueueWorkerIntervalSeconds(),
        enums.DashboardKeyEnum.QUEUE_PENDING_BATCHES:   0,
        enums.DashboardKeyEnum.QUEUE_PENDING_EVENTS:    0,
        enums.DashboardKeyEnum.QUEUE_FAILED_BATCHES:    0,
        enums.DashboardKeyEnum.QUEUE_OLDEST_PENDING_AT: "",
    }
    for _, queueStats := range queueStatsByStatus {
        switch application.ToString(queueStats["queue_status"]) {
        case ingestenums.QueueStatusEnum.PENDING:
            server[enums.DashboardKeyEnum.QUEUE_PENDING_BATCHES] = application.ToInt(queueStats["n_batches"])
            server[enums.DashboardKeyEnum.QUEUE_PENDING_EVENTS] = application.ToInt(queueStats["n_events"])
            server[enums.DashboardKeyEnum.QUEUE_OLDEST_PENDING_AT] = application.ToString(queueStats["oldest_received_at"])
        case ingestenums.QueueStatusEnum.FAILED:
            server[enums.DashboardKeyEnum.QUEUE_FAILED_BATCHES] = application.ToInt(queueStats["n_batches"])
        }
    }
    return server, nil
}

func (s *GetDashboardService) getProjects() ([]map[string]any, error) {
    projectEntities, err := s.projectsReaderSqliteRepository.GetAllProjects()
    if err != nil {
        return nil, err
    }
    agentEntities, err := s.agentsReaderSqliteRepository.GetAllAgents()
    if err != nil {
        return nil, err
    }
    lastEventAtByProjectId, err := s.dashboardReaderSqliteRepository.GetLastEventAtByProjectId()
    if err != nil {
        return nil, err
    }

    projects := make([]map[string]any, 0, len(projectEntities))
    for _, projectEntity := range projectEntities {
        project, err := s.getProject(projectEntity, agentEntities, lastEventAtByProjectId[projectEntity.GetId()])
        if err != nil {
            return nil, err
        }
        projects = append(projects, project)
    }
    return projects, nil
}

func (s *GetDashboardService) getProject(projectEntity *projectsentities.ProjectEntity, agentEntities []*agentsentities.AgentEntity, lastEventAt string) (map[string]any, error) {
    agentsTotal, enabledAgentsTotal, agentsOnline := 0, 0, 0
    for _, agentEntity := range agentEntities {
        if agentEntity.GetProjectId() != projectEntity.GetId() {
            continue
        }
        agentsTotal++
        if agentEntity.IsEnabled() {
            enabledAgentsTotal++
        }
        agentStatus := s.agentStatusResolver.GetAgentStatus(
            agentEntity.IsEnabled(), agentEntity.GetLastSeenAt(), agentEntity.GetHeartbeatIntervalSeconds(), s.nowDt,
        )
        if agentStatus == agentsenums.AgentStatusEnum.ONLINE {
            agentsOnline++
        }
    }

    eventsCountBySeverity, err := s.getEventsCountBySeverity(projectEntity.GetId())
    if err != nil {
        return nil, err
    }
    errors24h := eventsCountBySeverity[eventsenums.EventSeverityEnum.ERROR]
    criticals24h := eventsCountBySeverity[eventsenums.EventSeverityEnum.CRITICAL]

    return map[string]any{
        enums.DashboardKeyEnum.PROJECT_ID:          projectEntity.GetId(),
        enums.DashboardKeyEnum.PROJECT_UUID:        projectEntity.GetProjectUuid(),
        enums.DashboardKeyEnum.PROJECT_NAME:        projectEntity.GetProjectName(),
        enums.DashboardKeyEnum.PROJECT_DESCRIPTION: projectEntity.GetProjectDescription(),
        enums.DashboardKeyEnum.PROJECT_HEALTH:      s.projectHealthResolver.GetProjectHealth(enabledAgentsTotal, agentsOnline, errors24h, criticals24h),
        enums.DashboardKeyEnum.AGENTS_TOTAL:        agentsTotal,
        enums.DashboardKeyEnum.AGENTS_ONLINE:       agentsOnline,
        enums.DashboardKeyEnum.EVENTS_24H:          s.eventsStatsCounter.GetTotal(eventsCountBySeverity),
        enums.DashboardKeyEnum.WARNINGS_24H:        eventsCountBySeverity[eventsenums.EventSeverityEnum.WARNING],
        enums.DashboardKeyEnum.ERRORS_24H:          errors24h,
        enums.DashboardKeyEnum.CRITICALS_24H:       criticals24h,
        enums.DashboardKeyEnum.LAST_EVENT_AT:       lastEventAt,
    }, nil
}

func (s *GetDashboardService) getEventsCountBySeverity(projectId int) (map[string]int, error) {
    eventsStats, err := s.dashboardReaderSqliteRepository.GetEventsStatsByProjectIdSince(
        projectId, s.dateTimer.GetNowDtMinusSeconds(last24hSeconds),
    )
    if err != nil {
        return nil, err
    }
    eventsCountBySeverity := map[string]int{}
    for _, severityStats := range eventsStats {
        eventsCountBySeverity[application.ToString(severityStats["severity"])] = application.ToInt(severityStats["n_events"])
    }
    return eventsCountBySeverity, nil
}
