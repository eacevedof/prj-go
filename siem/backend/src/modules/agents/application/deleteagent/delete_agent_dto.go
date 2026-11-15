package deleteagent

import "siem-backend/src/modules/shared/application"

// DeleteAgentDto DELETE /api/v1/agents/:agent_uuid
type DeleteAgentDto struct {
    application.AbstractInputDto
    agentUuid string
}

func DeleteAgentDtoFromPrimitives(primitives map[string]any) *DeleteAgentDto {
    return &DeleteAgentDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        agentUuid:        application.GetTrimmedString(primitives, "agentUuid"),
    }
}

func (d *DeleteAgentDto) GetAgentUuid() string {
    return d.agentUuid
}
