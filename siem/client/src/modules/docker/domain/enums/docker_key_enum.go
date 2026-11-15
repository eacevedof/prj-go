package enums

// DockerEventTypeEnum los tipos de evento de Docker que entiende el backend (EventTypeEnum de siem-backend).
var DockerEventTypeEnum = struct {
    DOCKER_STATE string
    DOCKER_STATS string
    DOCKER_EVENT string
    DOCKER_LOG   string
}{
    DOCKER_STATE: "docker_state",
    DOCKER_STATS: "docker_stats",
    DOCKER_EVENT: "docker_event",
    DOCKER_LOG:   "docker_log",
}

// DockerAttributeEnum claves de attributes de los eventos de Docker. status, health, exit_code, is_oom_killed,
// change_kind, cpu_percent, memory_percent y action las usa el backend para decidir la gravedad.
var DockerAttributeEnum = struct {
    CONTAINER_ID    string
    IMAGE           string
    CHANGE_KIND     string
    STATUS          string
    PREVIOUS_STATUS string
    HEALTH          string
    PREVIOUS_HEALTH string
    EXIT_CODE       string
    IS_OOM_KILLED   string
    RESTART_COUNT   string

    CPU_PERCENT        string
    MEMORY_USAGE_BYTES string
    MEMORY_LIMIT_BYTES string
    MEMORY_PERCENT     string
    NET_RX_BYTES       string
    NET_TX_BYTES       string

    TYPE     string
    ACTION   string
    ACTOR_ID string
    NAME     string
    SIGNAL   string
}{
    CONTAINER_ID:    "container_id",
    IMAGE:           "image",
    CHANGE_KIND:     "change_kind",
    STATUS:          "status",
    PREVIOUS_STATUS: "previous_status",
    HEALTH:          "health",
    PREVIOUS_HEALTH: "previous_health",
    EXIT_CODE:       "exit_code",
    IS_OOM_KILLED:   "is_oom_killed",
    RESTART_COUNT:   "restart_count",

    CPU_PERCENT:        "cpu_percent",
    MEMORY_USAGE_BYTES: "memory_usage_bytes",
    MEMORY_LIMIT_BYTES: "memory_limit_bytes",
    MEMORY_PERCENT:     "memory_percent",
    NET_RX_BYTES:       "net_rx_bytes",
    NET_TX_BYTES:       "net_tx_bytes",

    TYPE:     "type",
    ACTION:   "action",
    ACTOR_ID: "actor_id",
    NAME:     "name",
    SIGNAL:   "signal",
}

// ContainerChangeKindEnum que le ha pasado a un contenedor entre dos pasadas.
var ContainerChangeKindEnum = struct {
    // primera pasada del agente: el estado de partida de cada contenedor
    FIRST_SEEN     string
    APPEARED       string
    STATUS_CHANGED string
    // mismo estado pero con mas reinicios: se cayo y la restart policy lo levanto entre dos pasadas
    RESTARTED      string
    HEALTH_CHANGED string
    REMOVED        string
}{
    FIRST_SEEN:     "first_seen",
    APPEARED:       "appeared",
    STATUS_CHANGED: "status_changed",
    RESTARTED:      "restarted",
    HEALTH_CHANGED: "health_changed",
    REMOVED:        "removed",
}

// ContainerStatusEnum los State.Status de Docker que importan aqui (+ removed, que pone el agente).
var ContainerStatusEnum = struct {
    RUNNING string
    EXITED  string
    REMOVED string
}{
    RUNNING: "running",
    EXITED:  "exited",
    REMOVED: "removed",
}

// DockerStateKeyEnum claves de los ficheros de estado del modulo docker.
var DockerStateKeyEnum = struct {
    CONTAINERS          string
    IS_INITIALIZED      string
    SINCE               string
    LAST_RUN_UNIX       string
    CPU_BY_CONTAINER_ID string
    CPU_TOTAL           string
    SYSTEM_TOTAL        string
    LAST_TS_BY_CONTAINER_ID string
}{
    CONTAINERS:              "containers",
    IS_INITIALIZED:          "is_initialized",
    SINCE:                   "since",
    LAST_RUN_UNIX:           "last_run_unix",
    CPU_BY_CONTAINER_ID:     "cpu_by_container_id",
    CPU_TOTAL:               "cpu_total",
    SYSTEM_TOTAL:            "system_total",
    LAST_TS_BY_CONTAINER_ID: "last_ts_by_container_id",
}

// DockerStateFileEnum ficheros de estado (dentro de APP_STATE_PATH), uno por caso de uso.
var DockerStateFileEnum = struct {
    CONTAINERS string
    EVENTS     string
    STATS      string
    LOGS       string
}{
    CONTAINERS: "docker-containers.json",
    EVENTS:     "docker-events.json",
    STATS:      "docker-stats.json",
    LOGS:       "docker-logs.json",
}

// DockerMessageEnum mensajes de error ({0}... los rellena config.Tr).
var DockerMessageEnum = struct {
    DOCKER_UNAVAILABLE_X_X string
}{
    DOCKER_UNAVAILABLE_X_X: "Docker not reachable at {0}: {1} (mount the socket into the agent)",
}
