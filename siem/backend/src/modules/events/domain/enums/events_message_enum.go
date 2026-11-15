package enums

// EventsMessageEnum mensajes de error de events ({0}... los rellena config.Tr).
var EventsMessageEnum = struct {
    SEVERITY_IS_NOT_VALID            string
    EVENT_TYPE_IS_NOT_VALID          string
    SOURCE_TYPE_IS_NOT_VALID         string
    X_MAX_VALUES_IS_20               string
    X_MUST_BE_A_UTC_DATE_TIME        string
    FROM_MUST_BE_BEFORE_TO           string
    PAGE_MUST_BE_1_OR_GREATER        string
    PAGE_SIZE_MUST_BE_1_TO_200       string
    SEARCH_MAX_LENGTH_IS_200         string
    RETENTION_DAYS_MUST_BE_1_TO_3650 string
}{
    SEVERITY_IS_NOT_VALID:            "severity: must be one of debug, info, warning, error, critical",
    EVENT_TYPE_IS_NOT_VALID:          "event_type: must be one of log_line, command_result, docker_state, docker_stats, docker_event, docker_log, host_metric",
    SOURCE_TYPE_IS_NOT_VALID:         "source_type: must be one of file, command, docker, docker_stats, docker_log, metric",
    X_MAX_VALUES_IS_20:               "{0}: max 20 values",
    X_MUST_BE_A_UTC_DATE_TIME:        "{0}: must be a UTC date time (YYYY-MM-DD HH:MM:SS)",
    FROM_MUST_BE_BEFORE_TO:           "from: must be before to",
    PAGE_MUST_BE_1_OR_GREATER:        "page: must be 1 or greater",
    PAGE_SIZE_MUST_BE_1_TO_200:       "page_size: must be between 1 and 200",
    SEARCH_MAX_LENGTH_IS_200:         "search: max length is 200",
    RETENTION_DAYS_MUST_BE_1_TO_3650: "retention_days: must be between 1 and 3650",
}
