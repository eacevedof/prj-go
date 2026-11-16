package enums

// DockerAttributeEnum claves de attributes que manda el siem-client en los eventos de Docker y que usa
// EventSeverityResolver (el resto de attributes se guarda tal cual).
var DockerAttributeEnum = struct {
    // docker_state
    CHANGE_KIND   string
    STATUS        string
    HEALTH        string
    EXIT_CODE     string
    IS_OOM_KILLED string

    // docker_stats
    CPU_PERCENT    string
    MEMORY_PERCENT string

    // docker_event
    ACTION string
}{
    CHANGE_KIND:   "change_kind",
    STATUS:        "status",
    HEALTH:        "health",
    EXIT_CODE:     "exit_code",
    IS_OOM_KILLED: "is_oom_killed",

    CPU_PERCENT:    "cpu_percent",
    MEMORY_PERCENT: "memory_percent",

    ACTION: "action",
}

// DockerAttributeValueEnum valores de attributes de Docker que cambian la gravedad (EventSeverityResolver): estado y
// health de un contenedor, tipo de cambio y acciones del daemon.
var DockerAttributeValueEnum = struct {
    STATUS_EXITED     string
    STATUS_DEAD       string
    STATUS_RESTARTING string
    HEALTH_UNHEALTHY  string
    CHANGE_RESTARTED  string

    ACTION_OOM                  string
    ACTION_DIE                  string
    ACTION_KILL                 string
    ACTION_HEALTH_STATUS_PREFIX string
}{
    STATUS_EXITED:     "exited",
    STATUS_DEAD:       "dead",
    STATUS_RESTARTING: "restarting",
    HEALTH_UNHEALTHY:  "unhealthy",
    CHANGE_RESTARTED:  "restarted",

    ACTION_OOM:                  "oom",
    ACTION_DIE:                  "die",
    ACTION_KILL:                 "kill",
    ACTION_HEALTH_STATUS_PREFIX: "health_status",
}

// DockerThresholdEnum umbrales de recursos (en %): por debajo, una muestra es debug (no ensucia los contadores).
var DockerThresholdEnum = struct {
    CPU_WARNING_PERCENT    float64
    MEMORY_WARNING_PERCENT float64
    MEMORY_ERROR_PERCENT   float64
}{
    CPU_WARNING_PERCENT:    90,
    MEMORY_WARNING_PERCENT: 85,
    MEMORY_ERROR_PERCENT:   95,
}

// MonitorAttributeEnum claves de attributes de un evento que viene de un monitor del panel. Los umbrales los escribe
// el backend al procesar el lote (desde la tabla monitors, por monitor_uuid): lo que mande el agente se pisa.
var MonitorAttributeEnum = struct {
    MONITOR_UUID       string
    VALUE              string
    IS_UNAVAILABLE     string
    WARNING_THRESHOLD  string
    CRITICAL_THRESHOLD string
}{
    MONITOR_UUID:       "monitor_uuid",
    VALUE:              "value",
    IS_UNAVAILABLE:     "is_unavailable",
    WARNING_THRESHOLD:  "warning_threshold",
    CRITICAL_THRESHOLD: "critical_threshold",
}
