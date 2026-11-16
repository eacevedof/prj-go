package enums

// ShipperKeyEnum claves del cuerpo de POST /api/v1/ingest-batches (contrato con el backend) y del estado del envio.
var ShipperKeyEnum = struct {
    BATCH_UUID     string
    AGENT          string
    HOST_NAME      string
    HOST_OS        string
    CLIENT_VERSION string
    CLIENT_PATH    string
    EVENTS         string

    // state/shipper.json
    LAST_SENT_AT_UNIX string
}{
    BATCH_UUID:     "batch_uuid",
    AGENT:          "agent",
    HOST_NAME:      "host_name",
    HOST_OS:        "host_os",
    CLIENT_VERSION: "client_version",
    CLIENT_PATH:    "client_path",
    EVENTS:         "events",

    LAST_SENT_AT_UNIX: "last_sent_at_unix",
}

// ShipperStateFileEnum fichero de estado del envio (dentro de APP_STATE_PATH).
var ShipperStateFileEnum = struct {
    SHIPPER string
}{
    SHIPPER: "shipper.json",
}
