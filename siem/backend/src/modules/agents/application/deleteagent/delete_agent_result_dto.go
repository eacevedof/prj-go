package deleteagent

import "siem-backend/src/modules/shared/application"

// DeleteAgentPrimitives lo que devuelve la API al borrar.
type DeleteAgentPrimitives struct {
    AgentUuid string `json:"agent_uuid"`
}

type DeleteAgentResultDto struct {
    agentUuid string
}

func DeleteAgentResultDtoFromPrimitives(primitives map[string]any) *DeleteAgentResultDto {
    return &DeleteAgentResultDto{agentUuid: application.GetTrimmedString(primitives, "agentUuid")}
}

func (d *DeleteAgentResultDto) ToArray() DeleteAgentPrimitives {
    return DeleteAgentPrimitives{AgentUuid: d.agentUuid}
}
