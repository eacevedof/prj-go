/**
 * Tipos de evento: una linea de un fichero vigilado, el resultado de un comando, de Docker (estado de un contenedor,
 * muestra de sus recursos, evento del daemon, linea de su log) o una metrica del host (monitor del panel).
 */
export const EventTypeEnum = Object.freeze({
    LOG_LINE: "log_line",
    COMMAND_RESULT: "command_result",
    DOCKER_STATE: "docker_state",
    DOCKER_STATS: "docker_stats",
    DOCKER_EVENT: "docker_event",
    DOCKER_LOG: "docker_log",
    HOST_METRIC: "host_metric",
});
