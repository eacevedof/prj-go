package checkagentconfig

import (
    "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/shared/application"
)

// CheckAgentConfigPrimitives lo que imprime config-check (el token nunca: solo si esta puesto). config = la local
// (siem-client.json con defaults) mas los watchers que salen de monitores folder / log_pattern del panel;
// metric_monitors = los monitores cpu, memory, disk, gpu y network; connections_monitor = el monitor connections (null
// si no hay); remote_config_warnings = monitores ignorados.
type CheckAgentConfigPrimitives struct {
    ConfigPath            string                         `json:"config_path"`
    EndpointUrl           string                         `json:"endpoint_url"`
    HasApiToken           bool                           `json:"has_api_token"`
    IsRemoteConfigEnabled bool                           `json:"is_remote_config_enabled"`
    Config                entities.AgentConfigPrimitives `json:"config"`
    MetricMonitors        []map[string]any               `json:"metric_monitors"`
    ConnectionsMonitor    map[string]any                 `json:"connections_monitor"`
    RemoteConfigWarnings  []string                       `json:"remote_config_warnings"`
}

type CheckAgentConfigResultDto struct {
    configPath            string
    endpointUrl           string
    hasApiToken           bool
    isRemoteConfigEnabled bool
    agentConfig           entities.AgentConfigPrimitives
    metricMonitors        []map[string]any
    connectionsMonitor    map[string]any
    remoteConfigWarnings  []string
}

func CheckAgentConfigResultDtoFromPrimitives(primitives map[string]any) *CheckAgentConfigResultDto {
    agentConfig, _ := primitives["agentConfig"].(entities.AgentConfigPrimitives)
    metricMonitors, _ := primitives["metricMonitors"].([]map[string]any)
    if metricMonitors == nil {
        metricMonitors = []map[string]any{}
    }
    connectionsMonitor, _ := primitives["connectionsMonitor"].(map[string]any)
    remoteConfigWarnings, _ := primitives["remoteConfigWarnings"].([]string)
    if remoteConfigWarnings == nil {
        remoteConfigWarnings = []string{}
    }
    return &CheckAgentConfigResultDto{
        configPath:            application.GetTrimmedString(primitives, "configPath"),
        endpointUrl:           application.GetTrimmedString(primitives, "endpointUrl"),
        hasApiToken:           application.GetInt(primitives, "hasApiToken") == 1,
        isRemoteConfigEnabled: application.GetInt(primitives, "isRemoteConfigEnabled") == 1,
        agentConfig:           agentConfig,
        metricMonitors:        metricMonitors,
        connectionsMonitor:    connectionsMonitor,
        remoteConfigWarnings:  remoteConfigWarnings,
    }
}

// IsConnectionsEnabled hay un monitor connections en la config remota.
func (d *CheckAgentConfigResultDto) IsConnectionsEnabled() bool {
    return d.connectionsMonitor != nil
}

func (d *CheckAgentConfigResultDto) GetCollectIntervalSeconds() int {
    return d.agentConfig.CollectIntervalSeconds
}

func (d *CheckAgentConfigResultDto) GetSendIntervalSeconds() int {
    return d.agentConfig.SendIntervalSeconds
}

func (d *CheckAgentConfigResultDto) GetWatchersCount() int {
    return len(d.agentConfig.Watchers)
}

func (d *CheckAgentConfigResultDto) GetCommandsCount() int {
    return len(d.agentConfig.Commands)
}

func (d *CheckAgentConfigResultDto) GetMetricMonitorsCount() int {
    return len(d.metricMonitors)
}

func (d *CheckAgentConfigResultDto) GetRemoteConfigWarnings() []string {
    return d.remoteConfigWarnings
}

func (d *CheckAgentConfigResultDto) GetEndpointUrl() string {
    return d.endpointUrl
}

func (d *CheckAgentConfigResultDto) IsDockerEnabled() bool {
    return d.agentConfig.Docker.IsEnabled
}

func (d *CheckAgentConfigResultDto) IsRemoteConfigEnabled() bool {
    return d.isRemoteConfigEnabled
}

func (d *CheckAgentConfigResultDto) ToArray() CheckAgentConfigPrimitives {
    return CheckAgentConfigPrimitives{
        ConfigPath:            d.configPath,
        EndpointUrl:           d.endpointUrl,
        HasApiToken:           d.hasApiToken,
        IsRemoteConfigEnabled: d.isRemoteConfigEnabled,
        Config:                d.agentConfig,
        MetricMonitors:        d.metricMonitors,
        ConnectionsMonitor:    d.connectionsMonitor,
        RemoteConfigWarnings:  d.remoteConfigWarnings,
    }
}
