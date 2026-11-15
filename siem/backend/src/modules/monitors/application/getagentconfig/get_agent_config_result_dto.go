package getagentconfig

import "siem-backend/src/modules/shared/application"

// AgentMonitorPrimitives lo que necesita el agente de un monitor. Los umbrales no viajan: la gravedad la decide el
// backend.
type AgentMonitorPrimitives struct {
    MonitorUuid     string         `json:"monitor_uuid"`
    MonitorType     string         `json:"monitor_type"`
    MonitorName     string         `json:"monitor_name"`
    IntervalSeconds int            `json:"interval_seconds"`
    Config          map[string]any `json:"config"`
}

// GetAgentConfigPrimitives lo que devuelve la API (claves en este orden).
type GetAgentConfigPrimitives struct {
    ConfigHash string                   `json:"config_hash"`
    Monitors   []AgentMonitorPrimitives `json:"monitors"`
}

type GetAgentConfigResultDto struct {
    // sha256 de monitors: el agente sabe si algo ha cambiado sin comparar la lista
    configHash string
    monitors   []AgentMonitorPrimitives
}

func GetAgentConfigResultDtoFromPrimitives(primitives map[string]any) *GetAgentConfigResultDto {
    monitors, _ := primitives["monitors"].([]AgentMonitorPrimitives)
    if monitors == nil {
        monitors = []AgentMonitorPrimitives{}
    }
    return &GetAgentConfigResultDto{
        configHash: application.GetTrimmedString(primitives, "configHash"),
        monitors:   monitors,
    }
}

func (d *GetAgentConfigResultDto) ToArray() GetAgentConfigPrimitives {
    return GetAgentConfigPrimitives{ConfigHash: d.configHash, Monitors: d.monitors}
}
