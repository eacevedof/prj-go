package enums

// CollectorsKeyEnum claves de un evento (el contrato con POST /api/v1/ingest-batches del backend) y de los ficheros
// de estado de los colectores.
var CollectorsKeyEnum = struct {
    EVENT_TYPE  string
    SOURCE_NAME string
    FILE_PATH   string
    MESSAGE     string
    EXIT_CODE   string
    DURATION_MS string
    OCCURRED_AT string
    SEVERITY    string
    ATTRIBUTES  string

    // attributes de las lineas que salen de un monitor del panel
    MONITOR_UUID string
    MONITOR_TYPE string

    // state/file-offsets.json
    OFFSETS              string
    INITIALIZED_WATCHERS string
}{
    EVENT_TYPE:  "event_type",
    SOURCE_NAME: "source_name",
    FILE_PATH:   "file_path",
    MESSAGE:     "message",
    EXIT_CODE:   "exit_code",
    DURATION_MS: "duration_ms",
    OCCURRED_AT: "occurred_at",
    SEVERITY:    "severity",
    ATTRIBUTES:  "attributes",

    MONITOR_UUID: "monitor_uuid",
    MONITOR_TYPE: "monitor_type",

    OFFSETS:              "offsets",
    INITIALIZED_WATCHERS: "initialized_watchers",
}

// CollectorsStateFileEnum ficheros de estado (dentro de APP_STATE_PATH).
var CollectorsStateFileEnum = struct {
    FILE_OFFSETS      string
    COMMANDS_LAST_RUN string
}{
    FILE_OFFSETS:      "file-offsets.json",
    COMMANDS_LAST_RUN: "commands-last-run.json",
}
