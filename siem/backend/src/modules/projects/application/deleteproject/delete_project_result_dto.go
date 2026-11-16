package deleteproject

import "siem-backend/src/modules/shared/application"

// DeletedAgentPrimitives un agente del proyecto borrado: lo justo para saber donde quedo instalado su siem-client
// (el borrado no lo desinstala; deja de ser aceptado, 401, y hay que quitarlo a mano de su host).
type DeletedAgentPrimitives struct {
    AgentName  string `json:"agent_name"`
    HostName   string `json:"host_name"`
    HostOs     string `json:"host_os"`
    HostIp     string `json:"host_ip"`
    ClientPath string `json:"client_path"`
}

// DeleteProjectPrimitives lo que devuelve la API al borrar.
type DeleteProjectPrimitives struct {
    ProjectUuid string                   `json:"project_uuid"`
    ProjectName string                   `json:"project_name"`
    Agents      []DeletedAgentPrimitives `json:"agents"`
}

type DeleteProjectResultDto struct {
    projectUuid string
    projectName string
    agents      []DeletedAgentPrimitives
}

func DeleteProjectResultDtoFromPrimitives(primitives map[string]any) *DeleteProjectResultDto {
    agents, _ := primitives["agents"].([]DeletedAgentPrimitives)
    if agents == nil {
        agents = []DeletedAgentPrimitives{}
    }
    return &DeleteProjectResultDto{
        projectUuid: application.GetTrimmedString(primitives, "projectUuid"),
        projectName: application.GetTrimmedString(primitives, "projectName"),
        agents:      agents,
    }
}

func (d *DeleteProjectResultDto) ToArray() DeleteProjectPrimitives {
    return DeleteProjectPrimitives{ProjectUuid: d.projectUuid, ProjectName: d.projectName, Agents: d.agents}
}
