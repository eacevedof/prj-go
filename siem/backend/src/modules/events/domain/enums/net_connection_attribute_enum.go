package enums

// NetConnectionAttributeEnum claves de attributes de un evento net_connection (monitor connections del siem-client):
// una conexion o un puerto en escucha que abre o cierra. is_unavailable = 1 es el aviso del agente de que en su
// sistema no puede muestrear (solo Linux); lo manda una vez, no en cada pasada.
var NetConnectionAttributeEnum = struct {
    ACTION         string
    PROTOCOL       string
    DIRECTION      string
    LOCAL_IP       string
    LOCAL_PORT     string
    REMOTE_IP      string
    REMOTE_PORT    string
    PID            string
    PROCESS        string
    IS_UNAVAILABLE string
}{
    ACTION:         "action",
    PROTOCOL:       "protocol",
    DIRECTION:      "direction",
    LOCAL_IP:       "local_ip",
    LOCAL_PORT:     "local_port",
    REMOTE_IP:      "remote_ip",
    REMOTE_PORT:    "remote_port",
    PID:            "pid",
    PROCESS:        "process",
    IS_UNAVAILABLE: "is_unavailable",
}

// NetConnectionAttributeValueEnum valores de action y direction de un net_connection.
var NetConnectionAttributeValueEnum = struct {
    ACTION_OPENED        string
    ACTION_CLOSED        string
    ACTION_LISTEN_OPENED string
    ACTION_LISTEN_CLOSED string
    DIRECTION_INBOUND    string
    DIRECTION_OUTBOUND   string
}{
    ACTION_OPENED:        "opened",
    ACTION_CLOSED:        "closed",
    ACTION_LISTEN_OPENED: "listen_opened",
    ACTION_LISTEN_CLOSED: "listen_closed",
    DIRECTION_INBOUND:    "inbound",
    DIRECTION_OUTBOUND:   "outbound",
}
