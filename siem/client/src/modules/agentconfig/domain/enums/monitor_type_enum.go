package enums

// MonitorTypeEnum tipos de monitor (recurso a vigilar) que crea el panel del backend y baja el agente con
// GET /api/v1/agent-config. Las metricas las toma el modulo hostmetrics; folder y log_pattern se vigilan como watchers;
// connections (uno por agente) lo muestrea el modulo connections.
var MonitorTypeEnum = struct {
    CPU         string
    MEMORY      string
    DISK        string
    GPU         string
    NETWORK     string
    FOLDER      string
    LOG_PATTERN string
    CONNECTIONS string

    METRIC_VALUES []string
    LINES_VALUES  []string
}{
    CPU:         "cpu",
    MEMORY:      "memory",
    DISK:        "disk",
    GPU:         "gpu",
    NETWORK:     "network",
    FOLDER:      "folder",
    LOG_PATTERN: "log_pattern",
    CONNECTIONS: "connections",

    METRIC_VALUES: []string{"cpu", "memory", "disk", "gpu", "network"},
    LINES_VALUES:  []string{"folder", "log_pattern"},
}

// RemoteConfigKeyEnum claves de la respuesta de GET /api/v1/agent-config y de state/remote-config.json (la ultima
// config remota guardada: se usa si el backend no responde).
var RemoteConfigKeyEnum = struct {
    DATA        string
    CONFIG_HASH string
    MONITORS    string
    SYNCED_AT   string

    MONITOR_UUID     string
    MONITOR_TYPE     string
    MONITOR_NAME     string
    INTERVAL_SECONDS string
    CONFIG           string

    // config de cada tipo
    PATH         string
    PATTERN      string
    IS_RECURSIVE string
    MATCH        string
    SEVERITY     string
    INTERFACE    string
    // connections
    IS_LOOPBACK_INCLUDED string
    IGNORE_PROCESSES     string
    IGNORE_PORTS         string

    STATE_FILE string
}{
    DATA:        "data",
    CONFIG_HASH: "config_hash",
    MONITORS:    "monitors",
    SYNCED_AT:   "synced_at",

    MONITOR_UUID:     "monitor_uuid",
    MONITOR_TYPE:     "monitor_type",
    MONITOR_NAME:     "monitor_name",
    INTERVAL_SECONDS: "interval_seconds",
    CONFIG:           "config",

    PATH:         "path",
    PATTERN:      "pattern",
    IS_RECURSIVE: "is_recursive",
    MATCH:        "match",
    SEVERITY:     "severity",
    INTERFACE:    "interface",

    IS_LOOPBACK_INCLUDED: "is_loopback_included",
    IGNORE_PROCESSES:     "ignore_processes",
    IGNORE_PORTS:         "ignore_ports",

    STATE_FILE: "remote-config.json",
}

// LogSeverityEnum gravedades que puede forzar un watcher / monitor log_pattern (las del backend).
var LogSeverityEnum = struct {
    VALUES []string
}{
    VALUES: []string{"debug", "info", "warning", "error", "critical"},
}
