/**
 * Errores (en ingles, igual que la API: no se traducen). Los limites son los del backend (AgentFieldsValidator).
 */
export const AgentsMessageEnum = Object.freeze({
    PROJECT_UUID_IS_REQUIRED: "project_uuid: field is required",
    AGENT_UUID_IS_REQUIRED: "agent_uuid: field is required",
    AGENT_NAME_IS_NOT_VALID: "agent_name: required, 3-100 chars of letters, digits and . _ -",
    HEARTBEAT_INTERVAL_MUST_BE_10_TO_86400: "heartbeat_interval_seconds: must be between 10 and 86400",
    UNEXPECTED_ERROR: "An unexpected error has occurred (is the backend running?)",
});
