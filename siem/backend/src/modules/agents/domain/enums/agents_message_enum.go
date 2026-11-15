package enums

// AgentsMessageEnum mensajes de error de agents ({0}, {1}... los rellena config.Tr).
var AgentsMessageEnum = struct {
    INVALID_OR_MISSING_API_TOKEN                    string
    PROJECT_X_NOT_FOUND                             string
    AGENT_NAME_X_ALREADY_EXISTS_IN_PROJECT          string
    AGENT_NAME_LENGTH_MUST_BE_3_TO_100              string
    AGENT_NAME_HAS_INVALID_CHARS                    string
    HEARTBEAT_INTERVAL_SECONDS_MUST_BE_10_TO_86400  string
    PROJECT_UUID_IS_REQUIRED                        string
    AGENT_X_NOT_FOUND                               string
}{
    INVALID_OR_MISSING_API_TOKEN:                   "Invalid or missing api token",
    PROJECT_X_NOT_FOUND:                            "Project {0} not found",
    AGENT_NAME_X_ALREADY_EXISTS_IN_PROJECT:         "Agent name {0} already exists in the project",
    AGENT_NAME_LENGTH_MUST_BE_3_TO_100:             "agent_name: length must be between 3 and 100",
    AGENT_NAME_HAS_INVALID_CHARS:                   "agent_name: only letters, digits and . _ - are allowed",
    HEARTBEAT_INTERVAL_SECONDS_MUST_BE_10_TO_86400: "heartbeat_interval_seconds: must be between 10 and 86400",
    PROJECT_UUID_IS_REQUIRED:                       "project_uuid: field is required",
    AGENT_X_NOT_FOUND:                              "Agent {0} not found",
}
