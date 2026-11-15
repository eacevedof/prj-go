package getprojectdetail

import "siem-backend/src/modules/projects/domain/entities"

// GetProjectDetailPrimitives la respuesta de GET /api/v1/projects/:project_uuid (claves en este orden).
type GetProjectDetailPrimitives struct {
    Project entities.ProjectPrimitives `json:"project"`
    Agents  []map[string]any           `json:"agents"`
}

type GetProjectDetailResultDto struct {
    project entities.ProjectPrimitives
    // un map por agente: sus datos, su estado, sus contadores de 24 h y "sources" (lista de maps de primitivos)
    agents []map[string]any
}

func GetProjectDetailResultDtoFromPrimitives(primitives map[string]any) *GetProjectDetailResultDto {
    project, _ := primitives["project"].(entities.ProjectPrimitives)
    agents, _ := primitives["agents"].([]map[string]any)
    if agents == nil {
        agents = []map[string]any{}
    }
    return &GetProjectDetailResultDto{project: project, agents: agents}
}

func (d *GetProjectDetailResultDto) GetProject() entities.ProjectPrimitives {
    return d.project
}

func (d *GetProjectDetailResultDto) GetAgents() []map[string]any {
    return d.agents
}

func (d *GetProjectDetailResultDto) ToArray() GetProjectDetailPrimitives {
    return GetProjectDetailPrimitives{Project: d.project, Agents: d.agents}
}
