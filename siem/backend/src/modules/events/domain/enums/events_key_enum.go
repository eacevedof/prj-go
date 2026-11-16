package enums

// EventsKeyEnum claves del contrato de events: lo que manda el siem-client por evento, los filtros de la busqueda y
// las columnas de la tabla.
var EventsKeyEnum = struct {
    EVENT_ID     string
    EVENT_TYPE   string
    SEVERITY     string
    MESSAGE      string
    SOURCE_TYPE  string
    SOURCE_NAME  string
    FILE_PATH    string
    EXIT_CODE    string
    DURATION_MS  string
    ATTRIBUTES   string
    MONITOR_UUID string
    OCCURRED_AT  string
    RECEIVED_AT  string

    // filtros y paginacion de GET /api/v1/events
    PROJECT_UUID string
    AGENT_UUID   string
    SEARCH       string
    FROM         string
    TO           string
    PAGE         string
    PAGE_SIZE    string
    EVENTS       string
    HAS_MORE     string

    // filas del histograma y de las facetas (alias del SQL)
    BUCKET_UNIX  string
    FACET_VALUE  string
    FILTER_VALUE string
    EVENTS_COUNT string

    // ./console events-purge
    RETENTION_DAYS  string
    DELETED_EVENTS  string
    DELETED_BATCHES string
    DELETED_BEFORE  string

    // columnas internas de BD
    PROJECT_ID string
    AGENT_ID   string
    SOURCE_ID  string
}{
    EVENT_ID:     "event_id",
    EVENT_TYPE:   "event_type",
    SEVERITY:     "severity",
    MESSAGE:      "message",
    SOURCE_TYPE:  "source_type",
    SOURCE_NAME:  "source_name",
    FILE_PATH:    "file_path",
    EXIT_CODE:    "exit_code",
    DURATION_MS:  "duration_ms",
    ATTRIBUTES:   "attributes",
    MONITOR_UUID: "monitor_uuid",
    OCCURRED_AT:  "occurred_at",
    RECEIVED_AT:  "received_at",

    PROJECT_UUID: "project_uuid",
    AGENT_UUID:   "agent_uuid",
    SEARCH:       "search",
    FROM:         "from",
    TO:           "to",
    PAGE:         "page",
    PAGE_SIZE:    "page_size",
    EVENTS:       "events",
    HAS_MORE:     "has_more",

    BUCKET_UNIX:  "bucket_unix",
    FACET_VALUE:  "facet_value",
    FILTER_VALUE: "filter_value",
    EVENTS_COUNT: "n_events",

    RETENTION_DAYS:  "retention_days",
    DELETED_EVENTS:  "deleted_events",
    DELETED_BATCHES: "deleted_batches",
    DELETED_BEFORE:  "deleted_before",

    PROJECT_ID: "project_id",
    AGENT_ID:   "agent_id",
    SOURCE_ID:  "source_id",
}
