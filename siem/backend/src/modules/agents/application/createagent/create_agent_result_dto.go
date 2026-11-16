package createagent

import (
    "siem-backend/src/modules/shared/application"

    "siem-backend/src/modules/agents/domain/entities"
)

type CreateAgentResultDto struct {
    // el agente creado (AgentEntity.ToPrimitives) + su api_token en claro, solo aqui
    agent    entities.AgentPrimitives
    apiToken string
}

func CreateAgentResultDtoFromPrimitives(primitives map[string]any) *CreateAgentResultDto {
    agent, _ := primitives["agent"].(entities.AgentPrimitives)
    return &CreateAgentResultDto{
        agent:    agent,
        apiToken: application.GetTrimmedString(primitives, "apiToken"),
    }
}

func (d *CreateAgentResultDto) GetAgent() entities.AgentPrimitives {
    return d.agent
}

func (d *CreateAgentResultDto) GetApiToken() string {
    return d.apiToken
}

func (d *CreateAgentResultDto) ToArray() entities.AgentWithApiTokenPrimitives {
    return entities.AgentWithApiTokenPrimitives{AgentPrimitives: d.agent, ApiToken: d.apiToken}
}
