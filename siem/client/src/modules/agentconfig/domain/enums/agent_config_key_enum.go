package enums

// AgentConfigKeyEnum claves de siem-client.json (y de la salida de config-check).
var AgentConfigKeyEnum = struct {
    COLLECT_INTERVAL_SECONDS   string
    SEND_INTERVAL_SECONDS      string
    HEARTBEAT_INTERVAL_SECONDS string
    MAX_EVENTS_PER_BATCH       string
    MAX_SPOOL_BATCHES          string
    WATCHERS                   string
    COMMANDS                   string
    DOCKER                     string

    // watcher
    NAME         string
    PATH         string
    PATTERN      string
    IS_RECURSIVE string
    START_AT     string

    // command
    COMMAND          string
    INTERVAL_SECONDS string
    TIMEOUT_SECONDS  string

    // config-check
    CONFIG_PATH   string
    ENDPOINT_URL  string
    HAS_API_TOKEN string
}{
    COLLECT_INTERVAL_SECONDS:   "collect_interval_seconds",
    SEND_INTERVAL_SECONDS:      "send_interval_seconds",
    HEARTBEAT_INTERVAL_SECONDS: "heartbeat_interval_seconds",
    MAX_EVENTS_PER_BATCH:       "max_events_per_batch",
    MAX_SPOOL_BATCHES:          "max_spool_batches",
    WATCHERS:                   "watchers",
    COMMANDS:                   "commands",
    DOCKER:                     "docker",

    NAME:         "name",
    PATH:         "path",
    PATTERN:      "pattern",
    IS_RECURSIVE: "is_recursive",
    START_AT:     "start_at",

    COMMAND:          "command",
    INTERVAL_SECONDS: "interval_seconds",
    TIMEOUT_SECONDS:  "timeout_seconds",

    CONFIG_PATH:   "config_path",
    ENDPOINT_URL:  "endpoint_url",
    HAS_API_TOKEN: "has_api_token",
}
