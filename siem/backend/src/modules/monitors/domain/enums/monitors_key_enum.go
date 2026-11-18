package enums

// MonitorsKeyEnum claves del contrato de monitors: el JSON de la API, las columnas de la tabla y lo que lee el agente.
var MonitorsKeyEnum = struct {
    ID                 string
    MONITOR_UUID       string
    PROJECT_ID         string
    PROJECT_UUID       string
    AGENT_ID           string
    AGENT_UUID         string
    AGENT_NAME         string
    MONITOR_TYPE       string
    MONITOR_NAME       string
    IS_ENABLED         string
    INTERVAL_SECONDS   string
    CONFIG             string
    WARNING_THRESHOLD  string
    CRITICAL_THRESHOLD string
    CREATED_AT         string
    UPDATED_AT         string

    // config de cada tipo
    PATH         string
    INTERFACE    string
    PATTERN      string
    IS_RECURSIVE string
    MATCH        string
    SEVERITY     string
    // connections
    IS_LOOPBACK_INCLUDED string
    IGNORE_PROCESSES     string
    IGNORE_PORTS         string

    // ultimo estado por agente (de sources)
    SOURCE_NAME   string
    LAST_EVENT_AT string
    LAST_SEVERITY string
    LAST_MESSAGE  string
}{
    ID:                 "id",
    MONITOR_UUID:       "monitor_uuid",
    PROJECT_ID:         "project_id",
    PROJECT_UUID:       "project_uuid",
    AGENT_ID:           "agent_id",
    AGENT_UUID:         "agent_uuid",
    AGENT_NAME:         "agent_name",
    MONITOR_TYPE:       "monitor_type",
    MONITOR_NAME:       "monitor_name",
    IS_ENABLED:         "is_enabled",
    INTERVAL_SECONDS:   "interval_seconds",
    CONFIG:             "config",
    WARNING_THRESHOLD:  "warning_threshold",
    CRITICAL_THRESHOLD: "critical_threshold",
    CREATED_AT:         "created_at",
    UPDATED_AT:         "updated_at",

    PATH:         "path",
    INTERFACE:    "interface",
    PATTERN:      "pattern",
    IS_RECURSIVE: "is_recursive",
    MATCH:        "match",
    SEVERITY:     "severity",

    IS_LOOPBACK_INCLUDED: "is_loopback_included",
    IGNORE_PROCESSES:     "ignore_processes",
    IGNORE_PORTS:         "ignore_ports",

    SOURCE_NAME:   "source_name",
    LAST_EVENT_AT: "last_event_at",
    LAST_SEVERITY: "last_severity",
    LAST_MESSAGE:  "last_message",
}

// MonitorTypeEnum que se vigila. Mismos valores que el CHECK de monitors.monitor_type. Los METRIC_VALUES mandan
// eventos host_metric con un valor numerico (admiten umbrales); folder y log_pattern mandan lineas de log; connections
// manda eventos net_connection cuando una conexion o un puerto en escucha abre o cierra (solo Linux). INTERVAL_VALUES:
// los que usan interval_seconds (las metricas y el muestreo de conexiones); los de lineas van en cada pasada.
var MonitorTypeEnum = struct {
    CPU         string
    MEMORY      string
    DISK        string
    GPU         string
    NETWORK     string
    FOLDER      string
    LOG_PATTERN string
    CONNECTIONS string

    VALUES          []string
    METRIC_VALUES   []string
    PERCENT_VALUES  []string
    INTERVAL_VALUES []string
}{
    CPU:         "cpu",
    MEMORY:      "memory",
    DISK:        "disk",
    GPU:         "gpu",
    NETWORK:     "network",
    FOLDER:      "folder",
    LOG_PATTERN: "log_pattern",
    CONNECTIONS: "connections",

    VALUES:          []string{"cpu", "memory", "disk", "gpu", "network", "folder", "log_pattern", "connections"},
    METRIC_VALUES:   []string{"cpu", "memory", "disk", "gpu", "network"},
    PERCENT_VALUES:  []string{"cpu", "memory", "disk", "gpu"},
    INTERVAL_VALUES: []string{"cpu", "memory", "disk", "gpu", "network", "connections"},
}

// MonitorsMessageEnum mensajes de error ({0}... los rellena config.Tr).
var MonitorsMessageEnum = struct {
    MONITOR_X_NOT_FOUND                 string
    PROJECT_X_NOT_FOUND                 string
    AGENT_X_NOT_FOUND_IN_PROJECT        string
    MONITOR_NAME_X_ALREADY_EXISTS       string
    MONITOR_X_ALREADY_WATCHES_IT        string
    ALL_AGENTS                          string
    MONITOR_NAME_IS_NOT_VALID           string
    MONITOR_TYPE_IS_NOT_VALID           string
    INTERVAL_SECONDS_MUST_BE_10_TO_86400 string
    X_MUST_BE_BETWEEN_0_AND_100         string
    X_MUST_BE_GREATER_THAN_0            string
    WARNING_MUST_BE_LOWER_THAN_CRITICAL string
    CONFIG_X_MUST_BE_AN_ABSOLUTE_PATH   string
    CONFIG_INTERFACE_IS_NOT_VALID       string
    CONFIG_PATTERN_IS_NOT_A_VALID_GLOB  string
    CONFIG_MATCH_IS_REQUIRED            string
    CONFIG_MATCH_IS_NOT_A_VALID_REGEX   string
    CONFIG_SEVERITY_IS_NOT_VALID        string
    CONFIG_IGNORE_PROCESSES_IS_NOT_VALID string
    CONFIG_IGNORE_PORTS_IS_NOT_VALID     string
}{
    MONITOR_X_NOT_FOUND:                  "Monitor {0} not found",
    PROJECT_X_NOT_FOUND:                  "Project {0} not found",
    AGENT_X_NOT_FOUND_IN_PROJECT:         "agent_uuid: agent {0} not found in this project",
    MONITOR_NAME_X_ALREADY_EXISTS:        "monitor_name: {0} already exists in this project",
    MONITOR_X_ALREADY_WATCHES_IT:         "config: monitor {0} already watches the same {1} on {2}",
    ALL_AGENTS:                           "all agents",
    MONITOR_NAME_IS_NOT_VALID:            "monitor_name: required, 1-100 chars of letters, digits and . _ -",
    MONITOR_TYPE_IS_NOT_VALID:            "monitor_type: must be one of cpu, memory, disk, gpu, network, folder, log_pattern, connections",
    INTERVAL_SECONDS_MUST_BE_10_TO_86400: "interval_seconds: must be between 10 and 86400",
    X_MUST_BE_BETWEEN_0_AND_100:          "{0}: must be between 0 and 100",
    X_MUST_BE_GREATER_THAN_0:             "{0}: must be greater than 0",
    WARNING_MUST_BE_LOWER_THAN_CRITICAL:  "warning_threshold: must be lower than critical_threshold",
    CONFIG_X_MUST_BE_AN_ABSOLUTE_PATH:    "config.{0}: must be an absolute path",
    CONFIG_INTERFACE_IS_NOT_VALID:        "config.interface: up to 50 chars of letters, digits and . _ : -",
    CONFIG_PATTERN_IS_NOT_A_VALID_GLOB:   "config.pattern: is not a valid glob",
    CONFIG_MATCH_IS_REQUIRED:             "config.match: is required (the regex a line must match)",
    CONFIG_MATCH_IS_NOT_A_VALID_REGEX:    "config.match: is not a valid regex",
    CONFIG_SEVERITY_IS_NOT_VALID:         "config.severity: must be empty or one of debug, info, warning, error, critical",
    CONFIG_IGNORE_PROCESSES_IS_NOT_VALID: "config.ignore_processes: up to 50 process names of 1-64 chars, without spaces or commas",
    CONFIG_IGNORE_PORTS_IS_NOT_VALID:     "config.ignore_ports: up to 50 ports between 1 and 65535",
}
