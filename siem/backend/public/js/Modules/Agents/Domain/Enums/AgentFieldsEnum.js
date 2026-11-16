/**
 * Limites de los campos de un agente: los mismos que AgentFieldsValidator del backend (aqui solo para avisar antes).
 */
export const AgentFieldsEnum = Object.freeze({
    AGENT_NAME_PATTERN: /^[A-Za-z0-9._-]{3,100}$/,
    HEARTBEAT_MIN_SECONDS: 10,
    HEARTBEAT_MAX_SECONDS: 86400,
});
