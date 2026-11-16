package enums

// AgentFieldsLimitEnum limites de los campos editables de un agente (AgentFieldsValidator). Si cambian, cambiar
// tambien los mensajes de AgentsMessageEnum y los limites del front.
var AgentFieldsLimitEnum = struct {
    AGENT_NAME_MIN_LENGTH          int
    AGENT_NAME_MAX_LENGTH          int
    HEARTBEAT_INTERVAL_SECONDS_MIN int
    HEARTBEAT_INTERVAL_SECONDS_MAX int
}{
    AGENT_NAME_MIN_LENGTH:          3,
    AGENT_NAME_MAX_LENGTH:          100,
    HEARTBEAT_INTERVAL_SECONDS_MIN: 10,
    HEARTBEAT_INTERVAL_SECONDS_MAX: 86400,
}
