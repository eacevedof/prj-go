package enums

// AgentsKeyEnum claves de la API y de consola de agents, mas las columnas con las que se monta AgentEntity.
var AgentsKeyEnum = struct {
    AGENT_UUID                 string
    AGENT_NAME                 string
    AGENT_STATUS               string
    PROJECT_UUID               string
    API_TOKEN                  string
    HEARTBEAT_INTERVAL_SECONDS string
    IS_ENABLED                 string
    HOST_NAME                  string
    HOST_OS                    string
    HOST_IP                    string
    CLIENT_VERSION             string
    LAST_SEEN_AT               string
    CREATED_AT                 string
    UPDATED_AT                 string

    // columnas internas de BD: nunca salen por la API
    ID             string
    PROJECT_ID     string
    API_TOKEN_HASH string
}{
    AGENT_UUID:                 "agent_uuid",
    AGENT_NAME:                 "agent_name",
    AGENT_STATUS:               "agent_status",
    PROJECT_UUID:               "project_uuid",
    API_TOKEN:                  "api_token",
    HEARTBEAT_INTERVAL_SECONDS: "heartbeat_interval_seconds",
    IS_ENABLED:                 "is_enabled",
    HOST_NAME:                  "host_name",
    HOST_OS:                    "host_os",
    HOST_IP:                    "host_ip",
    CLIENT_VERSION:             "client_version",
    LAST_SEEN_AT:               "last_seen_at",
    CREATED_AT:                 "created_at",
    UPDATED_AT:                 "updated_at",

    ID:             "id",
    PROJECT_ID:     "project_id",
    API_TOKEN_HASH: "api_token_hash",
}
