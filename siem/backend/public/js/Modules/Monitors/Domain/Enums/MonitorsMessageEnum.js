/**
 * Errores (en ingles, igual que la API: no se traducen).
 */
export const MonitorsMessageEnum = Object.freeze({
    PROJECT_UUID_IS_REQUIRED: "project_uuid: field is required",
    MONITOR_UUID_IS_REQUIRED: "monitor_uuid: field is required",
    MONITOR_NAME_IS_NOT_VALID: "monitor_name: required, 1-100 chars of letters, digits and . _ -",
    MONITOR_TYPE_IS_NOT_VALID: "monitor_type: must be one of cpu, memory, disk, gpu, network, folder, log_pattern, connections",
    INTERVAL_MUST_BE_10_TO_86400: "interval_seconds: must be an integer between 10 and 86400",
    WARNING_THRESHOLD_IS_NOT_VALID: "warning_threshold: must be a number greater than or equal to 0 (empty = no threshold)",
    CRITICAL_THRESHOLD_IS_NOT_VALID: "critical_threshold: must be a number greater than or equal to 0 (empty = no threshold)",
    THRESHOLD_MAX_IS_100: "thresholds: a percentage cannot be greater than 100",
    WARNING_MUST_BE_LOWER_THAN_CRITICAL: "warning_threshold: must be lower than critical_threshold",
    CONFIG_PATH_IS_REQUIRED: "config.path: field is required",
    CONFIG_MATCH_IS_REQUIRED: "config.match: field is required (the regular expression the lines must match)",
    CONFIG_MATCH_IS_NOT_VALID: "config.match: is not a valid regular expression",
    CONFIG_SEVERITY_IS_NOT_VALID: "config.severity: must be empty or one of debug, info, warning, error, critical",
    CONFIG_IGNORE_PROCESSES_IS_NOT_VALID: "config.ignore_processes: up to 50 process names of 1-64 chars, without spaces, separated by commas",
    CONFIG_IGNORE_PORTS_IS_NOT_VALID: "config.ignore_ports: up to 50 ports between 1 and 65535, separated by commas",
    UNEXPECTED_ERROR: "An unexpected error has occurred (is the backend running?)",
});
