package enums

// EventTypeEnum que es un evento: una linea de un fichero vigilado, el resultado de un comando programado o, de Docker,
// un cambio de estado de un contenedor, una muestra de sus recursos, un evento del daemon o una linea de su log.
// Mismos valores que el CHECK de events.event_type.
var EventTypeEnum = struct {
    LOG_LINE       string
    COMMAND_RESULT string
    DOCKER_STATE   string
    DOCKER_STATS   string
    DOCKER_EVENT   string
    DOCKER_LOG     string
    HOST_METRIC    string
    VALUES         []string
}{
    LOG_LINE:       "log_line",
    COMMAND_RESULT: "command_result",
    DOCKER_STATE:   "docker_state",
    DOCKER_STATS:   "docker_stats",
    DOCKER_EVENT:   "docker_event",
    DOCKER_LOG:     "docker_log",
    HOST_METRIC:    "host_metric",
    VALUES:         []string{"log_line", "command_result", "docker_state", "docker_stats", "docker_event", "docker_log", "host_metric"},
}

// SourceTypeByEventType a que tipo de source va cada tipo de evento: por contenedor, el estado y los eventos del daemon
// comparten source (docker); recursos y logs van aparte para que su ultimo valor no tape el estado.
var SourceTypeByEventType = map[string]string{
    EventTypeEnum.LOG_LINE:       SourceTypeEnum.FILE,
    EventTypeEnum.COMMAND_RESULT: SourceTypeEnum.COMMAND,
    EventTypeEnum.DOCKER_STATE:   SourceTypeEnum.DOCKER,
    EventTypeEnum.DOCKER_EVENT:   SourceTypeEnum.DOCKER,
    EventTypeEnum.DOCKER_STATS:   SourceTypeEnum.DOCKER_STATS,
    EventTypeEnum.DOCKER_LOG:     SourceTypeEnum.DOCKER_LOG,
    EventTypeEnum.HOST_METRIC:    SourceTypeEnum.METRIC,
}
