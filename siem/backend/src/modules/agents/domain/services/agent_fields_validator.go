package services

import (
    "regexp"
    "unicode/utf8"

    "siem-backend/src/modules/agents/domain/enums"
)

var agentNamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// AgentFieldsValidator reglas de los campos de un agente. Devuelve la pila de errores (vacia = todo bien).
type AgentFieldsValidator struct{}

func AgentFieldsValidatorGetInstance() *AgentFieldsValidator {
    return &AgentFieldsValidator{}
}

// GetErrors en el alta: el proyecto donde se crea mas los campos del agente.
func (v *AgentFieldsValidator) GetErrors(projectUuid string, agentName string, heartbeatIntervalSeconds int) []string {
    errors := []string{}
    if projectUuid == "" {
        errors = append(errors, enums.AgentsMessageEnum.PROJECT_UUID_IS_REQUIRED)
    }
    return append(errors, v.GetAgentFieldsErrors(agentName, heartbeatIntervalSeconds)...)
}

// GetAgentFieldsErrors los campos editables de un agente (alta y edicion).
func (v *AgentFieldsValidator) GetAgentFieldsErrors(agentName string, heartbeatIntervalSeconds int) []string {
    errors := []string{}
    agentNameLength := utf8.RuneCountInString(agentName)
    if agentNameLength < enums.AgentFieldsLimitEnum.AGENT_NAME_MIN_LENGTH || agentNameLength > enums.AgentFieldsLimitEnum.AGENT_NAME_MAX_LENGTH {
        errors = append(errors, enums.AgentsMessageEnum.AGENT_NAME_LENGTH_MUST_BE_3_TO_100)
    } else if !agentNamePattern.MatchString(agentName) {
        errors = append(errors, enums.AgentsMessageEnum.AGENT_NAME_HAS_INVALID_CHARS)
    }
    if heartbeatIntervalSeconds < enums.AgentFieldsLimitEnum.HEARTBEAT_INTERVAL_SECONDS_MIN ||
        heartbeatIntervalSeconds > enums.AgentFieldsLimitEnum.HEARTBEAT_INTERVAL_SECONDS_MAX {
        errors = append(errors, enums.AgentsMessageEnum.HEARTBEAT_INTERVAL_SECONDS_MUST_BE_10_TO_86400)
    }
    return errors
}
