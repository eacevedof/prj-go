package geteventsfilteroptions

// ProjectOptionPrimitives un proyecto del desplegable "Proyecto" del visor.
type ProjectOptionPrimitives struct {
    ProjectUuid string `json:"project_uuid"`
    ProjectName string `json:"project_name"`
}

// AgentOptionPrimitives un agente del desplegable "Agente": su proyecto acota la lista al elegir uno.
type AgentOptionPrimitives struct {
    AgentUuid   string `json:"agent_uuid"`
    AgentName   string `json:"agent_name"`
    ProjectUuid string `json:"project_uuid"`
}

// MonitorOptionPrimitives un monitor del desplegable "Monitor". agent_uuid "" = lo aplican todos los agentes del
// proyecto.
type MonitorOptionPrimitives struct {
    MonitorUuid string `json:"monitor_uuid"`
    MonitorName string `json:"monitor_name"`
    MonitorType string `json:"monitor_type"`
    ProjectUuid string `json:"project_uuid"`
    AgentUuid   string `json:"agent_uuid"`
}

// GetEventsFilterOptionsPrimitives lo que devuelve la API: el catalogo entero (no solo lo que tiene eventos en el
// rango, como las facetas), ordenado por nombre.
type GetEventsFilterOptionsPrimitives struct {
    Projects []ProjectOptionPrimitives `json:"projects"`
    Agents   []AgentOptionPrimitives   `json:"agents"`
    Monitors []MonitorOptionPrimitives `json:"monitors"`
}

type GetEventsFilterOptionsResultDto struct {
    projects []ProjectOptionPrimitives
    agents   []AgentOptionPrimitives
    monitors []MonitorOptionPrimitives
}

func GetEventsFilterOptionsResultDtoFromPrimitives(primitives map[string]any) *GetEventsFilterOptionsResultDto {
    projects, _ := primitives["projects"].([]ProjectOptionPrimitives)
    agents, _ := primitives["agents"].([]AgentOptionPrimitives)
    monitors, _ := primitives["monitors"].([]MonitorOptionPrimitives)
    if projects == nil {
        projects = []ProjectOptionPrimitives{}
    }
    if agents == nil {
        agents = []AgentOptionPrimitives{}
    }
    if monitors == nil {
        monitors = []MonitorOptionPrimitives{}
    }
    return &GetEventsFilterOptionsResultDto{projects: projects, agents: agents, monitors: monitors}
}

func (d *GetEventsFilterOptionsResultDto) ToArray() GetEventsFilterOptionsPrimitives {
    return GetEventsFilterOptionsPrimitives{Projects: d.projects, Agents: d.agents, Monitors: d.monitors}
}
