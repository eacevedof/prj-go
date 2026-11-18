package checkagentconfig

import "siem-client/src/modules/agentconfig/infrastructure/traits"

// CheckAgentConfigService valida la configuracion completa del agente y la devuelve con los defaults aplicados. Lo
// usan ./console config-check y el arranque de agent-run (si no vale, el agente no arranca).
type CheckAgentConfigService struct {
    traits.AgentConfigLoaderTrait

    checkAgentConfigDto *CheckAgentConfigDto
}

func CheckAgentConfigServiceGetInstance() *CheckAgentConfigService {
    return &CheckAgentConfigService{AgentConfigLoaderTrait: traits.NewAgentConfigLoaderTrait()}
}

// Invoke error: *AgentConfigException (400, 404).
func (s *CheckAgentConfigService) Invoke(checkAgentConfigDto *CheckAgentConfigDto) (*CheckAgentConfigResultDto, error) {
    s.checkAgentConfigDto = checkAgentConfigDto

    agentConfigEntity, err := s.GetValidAgentConfigEntityOrFail()
    if err != nil {
        return nil, err
    }

    metricMonitors := []map[string]any{}
    for _, metricMonitorEntity := range agentConfigEntity.GetMetricMonitorEntities() {
        metricMonitors = append(metricMonitors, metricMonitorEntity.ToPrimitives())
    }
    // nil en el JSON de config-check = sin monitor connections
    var connectionsMonitor map[string]any
    if connectionsMonitorEntity := agentConfigEntity.GetConnectionsMonitorEntityOrNil(); connectionsMonitorEntity != nil {
        connectionsMonitor = connectionsMonitorEntity.ToPrimitives()
    }
    isRemoteConfigEnabled := 0
    if s.IsRemoteConfigEnabled() {
        isRemoteConfigEnabled = 1
    }
    hasApiToken := 0
    if s.GetSiemApiToken() != "" {
        hasApiToken = 1
    }
    return CheckAgentConfigResultDtoFromPrimitives(map[string]any{
        "configPath":  s.GetConfigPath(),
        "endpointUrl": s.GetSiemEndpointUrl(),
        "hasApiToken": hasApiToken,
        "agentConfig": agentConfigEntity.ToPrimitives(),
        "isRemoteConfigEnabled": isRemoteConfigEnabled,
        "metricMonitors":        metricMonitors,
        "connectionsMonitor":    connectionsMonitor,
        "remoteConfigWarnings":  agentConfigEntity.GetRemoteConfigWarnings(),
    }), nil
}
