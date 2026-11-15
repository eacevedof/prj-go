package enums

// AgentConfigMessageEnum mensajes de error de la configuracion ({0}, {1}... los rellena config.Tr).
var AgentConfigMessageEnum = struct {
    CONFIG_FILE_X_NOT_FOUND           string
    CONFIG_FILE_X_IS_NOT_VALID_X      string
    CONFIG_HAS_ERRORS                 string
    ENDPOINT_URL_IS_NOT_VALID         string
    API_TOKEN_IS_REQUIRED             string
    NOTHING_TO_WATCH                  string
    X_MUST_BE_BETWEEN_X_AND_X         string
    X_NAME_IS_NOT_VALID               string
    X_NAME_X_IS_DUPLICATED            string
    WATCHER_X_PATH_MUST_BE_ABSOLUTE   string
    WATCHER_X_PATTERN_IS_NOT_VALID    string
    WATCHER_X_START_AT_IS_NOT_VALID   string
    COMMAND_X_COMMAND_IS_REQUIRED     string
    DOCKER_SOCKET_PATH_MUST_BE_ABSOLUTE string
    DOCKER_X_PATTERN_X_IS_NOT_VALID   string
    WATCHER_X_MATCH_IS_NOT_VALID      string
    WATCHER_X_SEVERITY_IS_NOT_VALID   string

    // config remota (GET /api/v1/agent-config)
    MONITOR_X_IS_NOT_VALID_X          string
    MONITOR_X_NAME_IS_ALREADY_USED    string
    REMOTE_CONFIG_REJECTED_API_TOKEN  string
    REMOTE_CONFIG_UNAVAILABLE_X       string
    REMOTE_CONFIG_UNEXPECTED_STATUS_X string
    REMOTE_CONFIG_IS_NOT_VALID_X      string
}{
    CONFIG_FILE_X_NOT_FOUND:         "Config file {0} not found (copy config/siem-client.example.json)",
    CONFIG_FILE_X_IS_NOT_VALID_X:    "Config file {0} is not valid: {1}",
    CONFIG_HAS_ERRORS:               "Config errors",
    ENDPOINT_URL_IS_NOT_VALID:       "APP_SIEM_ENDPOINT_URL: must be an http(s) url",
    API_TOKEN_IS_REQUIRED:           "APP_SIEM_API_TOKEN: is required (./console agent-create in the backend)",
    NOTHING_TO_WATCH:                "watchers, commands or docker: at least one is required",
    X_MUST_BE_BETWEEN_X_AND_X:       "{0}: must be between {1} and {2}",
    X_NAME_IS_NOT_VALID:             "{0}.name: required, 1-100 chars of letters, digits and . _ -",
    X_NAME_X_IS_DUPLICATED:          "{0}.name: {1} is duplicated",
    WATCHER_X_PATH_MUST_BE_ABSOLUTE: "watcher {0}: path must be an absolute folder path",
    WATCHER_X_PATTERN_IS_NOT_VALID:  "watcher {0}: pattern is not a valid glob",
    WATCHER_X_START_AT_IS_NOT_VALID: "watcher {0}: start_at must be end or beginning",
    COMMAND_X_COMMAND_IS_REQUIRED:   "command {0}: command is required",
    DOCKER_SOCKET_PATH_MUST_BE_ABSOLUTE: "docker.socket_path: must be an absolute path (/var/run/docker.sock)",
    DOCKER_X_PATTERN_X_IS_NOT_VALID:     "docker.{0}: {1} is not a valid glob",
    WATCHER_X_MATCH_IS_NOT_VALID:        "watcher {0}: match is not a valid regex",
    WATCHER_X_SEVERITY_IS_NOT_VALID:     "watcher {0}: severity must be debug, info, warning, error or critical",

    MONITOR_X_IS_NOT_VALID_X:          "remote monitor {0} ignored: {1}",
    MONITOR_X_NAME_IS_ALREADY_USED:    "remote monitor {0} ignored: its name is already used by a local watcher or another monitor",
    REMOTE_CONFIG_REJECTED_API_TOKEN:  "Remote config: backend rejected the api token (401): check APP_SIEM_API_TOKEN",
    REMOTE_CONFIG_UNAVAILABLE_X:       "Remote config unavailable (using the last one saved): {0}",
    REMOTE_CONFIG_UNEXPECTED_STATUS_X: "Remote config: backend answered {0} (using the last one saved)",
    REMOTE_CONFIG_IS_NOT_VALID_X:      "Remote config: unexpected response (using the last one saved): {0}",
}
