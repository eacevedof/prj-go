/**
 * Claves de un agente tal cual las da GET /api/v1/projects/:project_uuid (snake_case).
 */
export const AgentsKeyEnum = Object.freeze({
    AGENT_ID: "agent_id",
    AGENT_UUID: "agent_uuid",
    AGENT_NAME: "agent_name",
    AGENT_STATUS: "agent_status",
    HEARTBEAT_INTERVAL_SECONDS: "heartbeat_interval_seconds",
    IS_ENABLED: "is_enabled",
    HOST_NAME: "host_name",
    HOST_OS: "host_os",
    HOST_IP: "host_ip",
    CLIENT_VERSION: "client_version",
    LAST_SEEN_AT: "last_seen_at",
    EVENTS_24H: "events_24h",
    WARNINGS_24H: "warnings_24h",
    ERRORS_24H: "errors_24h",
    CRITICALS_24H: "criticals_24h",
    SOURCES: "sources",
});
