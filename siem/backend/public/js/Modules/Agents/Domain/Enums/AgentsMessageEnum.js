/**
 * Errores (en ingles, igual que la API: no se traducen).
 */
export const AgentsMessageEnum = Object.freeze({
    PROJECT_UUID_IS_REQUIRED: "project_uuid: field is required",
    AGENT_UUID_IS_REQUIRED: "agent_uuid: field is required",
    AGENT_NAME_IS_NOT_VALID: "agent_name: required, 1-100 chars of letters, digits and . _ -",
    HEARTBEAT_INTERVAL_MUST_BE_5_TO_3600: "heartbeat_interval_seconds: must be between 5 and 3600",
    UNEXPECTED_ERROR: "An unexpected error has occurred (is the backend running?)",
});
