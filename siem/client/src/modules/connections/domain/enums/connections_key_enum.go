package enums

// ProtocolEnum el protocolo de un socket. El fichero de /proc/net del que sale dice ademas si es IPv4 o IPv6.
var ProtocolEnum = struct {
    TCP string
    UDP string
}{
    TCP: "tcp",
    UDP: "udp",
}

// ProcNetFileEnum los ficheros de /proc/net que se leen y el protocolo / version de IP de cada uno.
var ProcNetFileEnum = struct {
    TCP  string
    TCP6 string
    UDP  string
    UDP6 string
}{
    TCP:  "tcp",
    TCP6: "tcp6",
    UDP:  "udp",
    UDP6: "udp6",
}

// DirectionEnum hacia donde va una conexion vista desde este host.
//   - listen: socket en escucha (un servicio que acepta conexiones).
//   - inbound: alguien de fuera se ha conectado a un puerto que este host escucha.
//   - outbound: un proceso de este host se ha conectado a otro.
var DirectionEnum = struct {
    LISTEN   string
    INBOUND  string
    OUTBOUND string
}{
    LISTEN:   "listen",
    INBOUND:  "inbound",
    OUTBOUND: "outbound",
}

// ConnectionKeyEnum claves de los primitivos de una conexion (ConnectionEntity.ToPrimitives): los attributes del
// net_connection y lo que se guarda por conexion en state/connections.json.
var ConnectionKeyEnum = struct {
    PROTOCOL    string
    IP_VERSION  string
    DIRECTION   string
    STATE       string
    LOCAL_IP    string
    LOCAL_PORT  string
    REMOTE_IP   string
    REMOTE_PORT string
    PID         string
    PROCESS     string
    EXE         string
    USER        string
}{
    PROTOCOL:    "protocol",
    IP_VERSION:  "ip_version",
    DIRECTION:   "direction",
    STATE:       "state",
    LOCAL_IP:    "local_ip",
    LOCAL_PORT:  "local_port",
    REMOTE_IP:   "remote_ip",
    REMOTE_PORT: "remote_port",
    PID:         "pid",
    PROCESS:     "process",
    EXE:         "exe",
    USER:        "user",
}

// UNKNOWN_PROCESS nombre de proceso cuando no se puede saber de quien es el socket (sin permiso para leer sus fd).
const UNKNOWN_PROCESS = "?"

// NET_CONNECTION_EVENT_TYPE el evento de este modulo (EventTypeEnum.NET_CONNECTION del backend: mismo valor).
const NET_CONNECTION_EVENT_TYPE = "net_connection"

// ConnectionActionEnum que ha pasado entre dos pasadas: una conexion o un puerto en escucha que aparece o desaparece
// (NetConnectionAttributeValueEnum del backend: mismos valores; la gravedad la decide el).
var ConnectionActionEnum = struct {
    OPENED        string
    CLOSED        string
    LISTEN_OPENED string
    LISTEN_CLOSED string
}{
    OPENED:        "opened",
    CLOSED:        "closed",
    LISTEN_OPENED: "listen_opened",
    LISTEN_CLOSED: "listen_closed",
}

// NetConnectionAttributeEnum claves de attributes de un net_connection, ademas de las de ConnectionEntity.ToPrimitives.
//   - is_initial: la primera pasada manda lo que ya estaba abierto (estado de partida), marcado asi.
//   - is_unavailable: el aviso (una vez) de que en este sistema no se puede muestrear.
var NetConnectionAttributeEnum = struct {
    MONITOR_UUID     string
    MONITOR_TYPE     string
    ACTION           string
    DURATION_SECONDS string
    IS_INITIAL       string
    IS_UNAVAILABLE   string
}{
    MONITOR_UUID:     "monitor_uuid",
    MONITOR_TYPE:     "monitor_type",
    ACTION:           "action",
    DURATION_SECONDS: "duration_seconds",
    IS_INITIAL:       "is_initial",
    IS_UNAVAILABLE:   "is_unavailable",
}

// ConnectionsStateEnum state/connections.json: de que monitor es la foto, cuando se tomo, si ya se aviso de que no
// esta disponible y la foto (clave de la conexion -> sus primitivos + opened_at_unix, para la duracion del cierre).
var ConnectionsStateEnum = struct {
    STATE_FILE     string
    MONITOR_UUID   string
    LAST_RUN_UNIX  string
    IS_UNAVAILABLE string
    CONNECTIONS    string
    OPENED_AT_UNIX string
}{
    STATE_FILE:     "connections.json",
    MONITOR_UUID:   "monitor_uuid",
    LAST_RUN_UNIX:  "last_run_unix",
    IS_UNAVAILABLE: "is_unavailable",
    CONNECTIONS:    "connections",
    OPENED_AT_UNIX: "opened_at_unix",
}

// ConnectionsMessageEnum motivos de "no disponible" (en ingles, como el resto del agente). {0}... los rellena config.Tr.
var ConnectionsMessageEnum = struct {
    OS_X_IS_NOT_SUPPORTED string
    CANNOT_READ_SOCKETS_X string
}{
    OS_X_IS_NOT_SUPPORTED: "OS {0} is not supported (only Linux: /proc/net)",
    CANNOT_READ_SOCKETS_X: "cannot read the sockets in /proc: {0}",
}
