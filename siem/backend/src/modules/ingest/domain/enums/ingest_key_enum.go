package enums

// IngestKeyEnum claves del contrato de POST /api/v1/ingest-batches (lo que manda el siem-client) y de las
// respuestas de la cola.
var IngestKeyEnum = struct {
    BATCH_UUID     string
    AGENT          string
    HOST_NAME      string
    HOST_OS        string
    CLIENT_VERSION string
    EVENTS         string

    QUEUED_EVENTS string
    IS_DUPLICATE  string

    // ./console ingest-queue-process
    MAX_BATCHES       string
    PROCESSED_BATCHES string
    FAILED_BATCHES    string
    INSERTED_EVENTS   string
    SKIPPED_EVENTS    string

    // columnas internas de ingest_queue
    ID          string
    AGENT_ID    string
    PROJECT_ID  string
    PAYLOAD     string
    ATTEMPTS    string
    RECEIVED_AT string
}{
    BATCH_UUID:     "batch_uuid",
    AGENT:          "agent",
    HOST_NAME:      "host_name",
    HOST_OS:        "host_os",
    CLIENT_VERSION: "client_version",
    EVENTS:         "events",

    QUEUED_EVENTS: "queued_events",
    IS_DUPLICATE:  "is_duplicate",

    MAX_BATCHES:       "max_batches",
    PROCESSED_BATCHES: "processed_batches",
    FAILED_BATCHES:    "failed_batches",
    INSERTED_EVENTS:   "inserted_events",
    SKIPPED_EVENTS:    "skipped_events",

    ID:          "id",
    AGENT_ID:    "agent_id",
    PROJECT_ID:  "project_id",
    PAYLOAD:     "payload",
    ATTEMPTS:    "attempts",
    RECEIVED_AT: "received_at",
}
