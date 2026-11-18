package enums

// ProtocolEnum el protocolo de un socket. El fichero de /proc/net del que sale dice ademas si es IPv4 o IPv6.
var ProtocolEnum = struct {
    TCP string
    UDP string
}{
    TCP: "tcp",
    UDP: "udp",
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

// ConnectionKeyEnum claves de los primitivos de las entidades del modulo: los de una conexion (ConnectionEntity.
// ToPrimitives: los attributes del net_connection y lo que se guarda por conexion en state/connections.json), los de
// un socket y un proceso (SocketEntityFromPrimitives, ProcessEntityFromPrimitives) y los de un cambio.
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
    // socket
    UID   string
    INODE string
    // cambio (ConnectionChangeEntityFromPrimitives): la conexion y, en NetConnectionAttributeEnum, accion y duracion
    CONNECTION string
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
    UID:         "uid",
    INODE:       "inode",
    CONNECTION:  "connection",
}

// UNKNOWN_PROCESS nombre de proceso cuando no se puede saber de quien es el socket (sin permiso para leer sus fd).
const UNKNOWN_PROCESS = "?"

// ConnectionFormatEnum como se escriben la clave de una conexion entre dos fotos y sus extremos.
//   - KEY_SEPARATOR: protocolo|local|remoto|inodo, y |pid en la de la conexion.
//   - ANY_ENDPOINT: el extremo remoto de un socket en escucha.
//   - PROCESS_LABEL: "firefox(4321)".
var ConnectionFormatEnum = struct {
    KEY_SEPARATOR string
    ANY_ENDPOINT  string
    PROCESS_LABEL string
}{
    KEY_SEPARATOR: "|",
    ANY_ENDPOINT:  "*",
    PROCESS_LABEL: "{0}({1})",
}

// ConnectionEventMessageEnum el mensaje de un net_connection ({0}... los rellena config.Tr; las partes se unen con
// MESSAGE_PART_SEPARATOR porque Tr recorta los espacios de los extremos):
//
//   opened tcp outbound firefox(4321) 192.168.1.5:52780 -> 142.250.1.1:443
//   closed tcp outbound firefox(4321) 192.168.1.5:52780 -> 142.250.1.1:443 after 12 s
//   listen_opened tcp sshd(812) 0.0.0.0:22 (already open)
var ConnectionEventMessageEnum = struct {
    // {0} accion · {1} protocolo · {2} direccion · {3} proceso · {4} local · {5} remoto
    CONNECTION_X string
    // {0} accion · {1} protocolo · {2} proceso · {3} local
    LISTEN_X string
    // {0} segundos
    AFTER_X_SECONDS string
    ALREADY_OPEN    string
    // {0} tipo de monitor · {1} motivo
    X_NOT_AVAILABLE_X      string
    MESSAGE_PART_SEPARATOR string
}{
    CONNECTION_X:      "{0} {1} {2} {3} {4} -> {5}",
    LISTEN_X:          "{0} {1} {2} {3}",
    AFTER_X_SECONDS:   "after {0} s",
    ALREADY_OPEN:      "(already open)",
    X_NOT_AVAILABLE_X:      "{0} not available: {1}",
    MESSAGE_PART_SEPARATOR: " ",
}

// NET_CONNECTION_EVENT_TYPE el evento de este modulo (EventTypeEnum.NET_CONNECTION del backend: mismo valor).
const NET_CONNECTION_EVENT_TYPE = "net_connection"

// EVENT_DATE_TIME_FORMAT occurred_at del evento: UTC "YYYY-MM-DD HH:MM:SS" (el formato que espera el backend).
const EVENT_DATE_TIME_FORMAT = "2006-01-02 15:04:05"

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
