package getprojectmonitors

import "siem-backend/src/modules/monitors/domain/entities"

// MonitorStatePrimitives el ultimo valor de un monitor en un agente.
type MonitorStatePrimitives struct {
    AgentUuid    string `json:"agent_uuid"`
    AgentName    string `json:"agent_name"`
    LastEventAt  string `json:"last_event_at"`
    LastSeverity string `json:"last_severity"`
    LastMessage  string `json:"last_message"`
}

// MonitorWithStatesPrimitives un monitor + su ultimo valor por agente ([] si aun no ha mandado nada).
type MonitorWithStatesPrimitives struct {
    entities.MonitorPrimitives
    States []MonitorStatePrimitives `json:"states"`
}

// GetProjectMonitorsPrimitives lo que devuelve la API.
type GetProjectMonitorsPrimitives struct {
    Monitors []MonitorWithStatesPrimitives `json:"monitors"`
}

type GetProjectMonitorsResultDto struct {
    monitors []MonitorWithStatesPrimitives
}

func GetProjectMonitorsResultDtoFromPrimitives(primitives map[string]any) *GetProjectMonitorsResultDto {
    monitors, _ := primitives["monitors"].([]MonitorWithStatesPrimitives)
    if monitors == nil {
        monitors = []MonitorWithStatesPrimitives{}
    }
    return &GetProjectMonitorsResultDto{monitors: monitors}
}

func (d *GetProjectMonitorsResultDto) ToArray() GetProjectMonitorsPrimitives {
    return GetProjectMonitorsPrimitives{Monitors: d.monitors}
}
