package enums

// EventsLimitEnum limites de events: longitudes de lo que se guarda, paginacion y filtros del visor, histograma y
// retencion. Si cambian, cambiar tambien los mensajes de EventsMessageEnum y los limites del front.
var EventsLimitEnum = struct {
    // EventNormalizer: lo que no cabe se recorta (mensaje, ruta) o se descarta (source_name, attributes, monitor_uuid)
    MESSAGE_MAX_LENGTH     int
    SOURCE_NAME_MAX_LENGTH int
    FILE_PATH_MAX_LENGTH   int
    // attributes son unos pocos datos estructurados (estado, exit_code, % de cpu...): mas que esto no es un attribute
    ATTRIBUTES_MAX_BYTES int
    // events.monitor_uuid VARCHAR(50)
    MONITOR_UUID_MAX_LENGTH int
    // una hora del cliente que va por delante de la de llegada mas que esto es un reloj descolocado
    OCCURRED_AT_MAX_AHEAD_MINUTES int
    // sources.last_message
    LAST_MESSAGE_MAX_LENGTH int

    // GET /api/v1/events
    PAGE_DEFAULT      int
    PAGE_MIN          int
    PAGE_SIZE_DEFAULT int
    PAGE_SIZE_MIN     int
    PAGE_SIZE_MAX     int
    SEARCH_MAX_LENGTH int

    // histograma: ~60 barras; sin rango, las ultimas 24 h
    HISTOGRAM_TARGET_BUCKETS      int
    HISTOGRAM_DEFAULT_RANGE_HOURS int

    // ./console events-purge
    RETENTION_DAYS_MIN int
    RETENTION_DAYS_MAX int
    SECONDS_PER_DAY    int
    // un lote ya procesado se guarda 24 h solo para reconocer un reenvio del cliente
    DONE_BATCHES_RETENTION_SECONDS int
}{
    MESSAGE_MAX_LENGTH:            8000,
    SOURCE_NAME_MAX_LENGTH:        255,
    FILE_PATH_MAX_LENGTH:          500,
    ATTRIBUTES_MAX_BYTES:          4000,
    MONITOR_UUID_MAX_LENGTH:       50,
    OCCURRED_AT_MAX_AHEAD_MINUTES: 5,
    LAST_MESSAGE_MAX_LENGTH:       500,

    PAGE_DEFAULT:      1,
    PAGE_MIN:          1,
    PAGE_SIZE_DEFAULT: 50,
    PAGE_SIZE_MIN:     1,
    PAGE_SIZE_MAX:     200,
    SEARCH_MAX_LENGTH: 200,

    HISTOGRAM_TARGET_BUCKETS:      60,
    HISTOGRAM_DEFAULT_RANGE_HOURS: 24,

    RETENTION_DAYS_MIN:             1,
    RETENTION_DAYS_MAX:             3650,
    SECONDS_PER_DAY:                86400,
    DONE_BATCHES_RETENTION_SECONDS: 86400,
}
