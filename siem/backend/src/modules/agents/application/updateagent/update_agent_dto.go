package updateagent

import (
    "siem-backend/src/modules/agents/domain/entities"
    "siem-backend/src/modules/shared/application"
)

// UpdateAgentDto PUT /api/v1/agents/:agent_uuid {agent_name, heartbeat_interval_seconds, is_enabled}. is_enabled
// ausente = activo; heartbeat ausente = el de por defecto.
type UpdateAgentDto struct {
    application.AbstractInputDto
    agentUuid                string
    agentName                string
    heartbeatIntervalSeconds int
    isEnabled                bool
}

func UpdateAgentDtoFromPrimitives(primitives map[string]any) *UpdateAgentDto {
    heartbeatIntervalSeconds := entities.DEFAULT_HEARTBEAT_INTERVAL_SECONDS
    if application.GetTrimmedString(primitives, "heartbeatIntervalSeconds") != "" {
        heartbeatIntervalSeconds = application.GetInt(primitives, "heartbeatIntervalSeconds")
    }
    isEnabled := true
    if rawIsEnabled, isSet := primitives["isEnabled"]; isSet && rawIsEnabled != nil {
        isEnabled = application.ToInt(rawIsEnabled) == 1
    }
    return &UpdateAgentDto{
        AbstractInputDto:         application.NewAbstractInputDto(primitives),
        agentUuid:                application.GetTrimmedString(primitives, "agentUuid"),
        agentName:                application.GetTrimmedString(primitives, "agentName"),
        heartbeatIntervalSeconds: heartbeatIntervalSeconds,
        isEnabled:                isEnabled,
    }
}

func (d *UpdateAgentDto) GetAgentUuid() string {
    return d.agentUuid
}

func (d *UpdateAgentDto) GetAgentName() string {
    return d.agentName
}

func (d *UpdateAgentDto) GetHeartbeatIntervalSeconds() int {
    return d.heartbeatIntervalSeconds
}

func (d *UpdateAgentDto) IsEnabled() bool {
    return d.isEnabled
}
