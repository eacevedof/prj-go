package enums

// ShipperMessageEnum mensajes de error del envio ({0}... los rellena config.Tr).
var ShipperMessageEnum = struct {
    BACKEND_REJECTED_API_TOKEN     string
    BACKEND_UNAVAILABLE_X          string
    BACKEND_UNEXPECTED_STATUS_X_X  string
}{
    BACKEND_REJECTED_API_TOKEN:    "Backend rejected the api token (401): check APP_SIEM_API_TOKEN",
    BACKEND_UNAVAILABLE_X:         "Backend unavailable: {0}",
    BACKEND_UNEXPECTED_STATUS_X_X: "Backend answered {0}: {1}",
}
