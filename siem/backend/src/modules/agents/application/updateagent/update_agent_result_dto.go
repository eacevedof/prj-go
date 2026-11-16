package updateagent

import "siem-backend/src/modules/agents/domain/entities"

type UpdateAgentResultDto struct {
    // el agente ya cambiado (AgentEntity.ToPrimitives: sin ids ni hash del token)
    agent entities.AgentPrimitives
}

func UpdateAgentResultDtoFromPrimitives(primitives map[string]any) *UpdateAgentResultDto {
    agent, _ := primitives["agent"].(entities.AgentPrimitives)
    return &UpdateAgentResultDto{agent: agent}
}

func (d *UpdateAgentResultDto) GetAgent() entities.AgentPrimitives {
    return d.agent
}

func (d *UpdateAgentResultDto) ToArray() entities.AgentPrimitives {
    return d.agent
}
