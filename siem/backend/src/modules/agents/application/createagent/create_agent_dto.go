package createagent

import (
    "siem-backend/src/modules/agents/domain/entities"
    "siem-backend/src/modules/shared/application"
)

type CreateAgentDto struct {
    application.AbstractInputDto
    projectUuid              string
    agentName                string
    heartbeatIntervalSeconds int
}

func CreateAgentDtoFromPrimitives(primitives map[string]any) *CreateAgentDto {
    heartbeatIntervalSeconds := entities.DEFAULT_HEARTBEAT_INTERVAL_SECONDS
    // "" (opcion de consola sin pasar) cuenta como ausente: se queda el valor por defecto
    if application.GetTrimmedString(primitives, "heartbeatIntervalSeconds") != "" {
        heartbeatIntervalSeconds = application.GetInt(primitives, "heartbeatIntervalSeconds")
    }
    return &CreateAgentDto{
        AbstractInputDto:         application.NewAbstractInputDto(primitives),
        projectUuid:              application.GetTrimmedString(primitives, "projectUuid"),
        agentName:                application.GetTrimmedString(primitives, "agentName"),
        heartbeatIntervalSeconds: heartbeatIntervalSeconds,
    }
}

func (d *CreateAgentDto) GetProjectUuid() string {
    return d.projectUuid
}

func (d *CreateAgentDto) GetAgentName() string {
    return d.agentName
}

func (d *CreateAgentDto) GetHeartbeatIntervalSeconds() int {
    return d.heartbeatIntervalSeconds
}
