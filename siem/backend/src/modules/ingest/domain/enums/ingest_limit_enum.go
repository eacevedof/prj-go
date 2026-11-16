package enums

// IngestLimitEnum limites de ingest: tamano de un lote, pasadas de la cola y longitudes de los datos del host que se
// guardan en el agente (columnas de agents). Si cambian, cambiar tambien los mensajes de IngestMessageEnum.
var IngestLimitEnum = struct {
    MAX_EVENTS_PER_BATCH int
    // eventsCount de un lote cuyo "events" no era un array (lo rechaza IngestBatchValidator)
    EVENTS_COUNT_NOT_AN_ARRAY int

    // ./console ingest-queue-process --max_batches
    MAX_BATCHES_DEFAULT int
    MAX_BATCHES_MIN     int
    MAX_BATCHES_MAX     int
    // lotes por pasada del worker de la API
    MAX_BATCHES_PER_WORKER_ROUND int
    // un lote que falla se reintenta en las siguientes pasadas; al tercer fallo queda como failed (visible en el
    // dashboard) y deja de bloquear la cola
    MAX_ATTEMPTS int
    // ingest_queue.error_message
    ERROR_MESSAGE_MAX_LENGTH int

    // latido del agente: lo que no cabe en su columna se recorta
    HOST_NAME_MAX_LENGTH      int
    HOST_OS_MAX_LENGTH        int
    HOST_IP_MAX_LENGTH        int
    CLIENT_VERSION_MAX_LENGTH int
    CLIENT_PATH_MAX_LENGTH    int
}{
    MAX_EVENTS_PER_BATCH:      2000,
    EVENTS_COUNT_NOT_AN_ARRAY: -1,

    MAX_BATCHES_DEFAULT:          100,
    MAX_BATCHES_MIN:              1,
    MAX_BATCHES_MAX:              1000,
    MAX_BATCHES_PER_WORKER_ROUND: 50,
    MAX_ATTEMPTS:                 3,
    ERROR_MESSAGE_MAX_LENGTH:     500,

    HOST_NAME_MAX_LENGTH:      150,
    HOST_OS_MAX_LENGTH:        50,
    HOST_IP_MAX_LENGTH:        45,
    CLIENT_VERSION_MAX_LENGTH: 20,
    CLIENT_PATH_MAX_LENGTH:    500,
}
