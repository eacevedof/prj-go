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

// FilterTermEnum terminos especiales del filtro de connections-list (el resto se busca como texto).
var FilterTermEnum = struct {
    // solo conexiones con IP remota publica (fuera loopback, redes privadas, link-local, CGNAT)
    REMOTE string
}{
    REMOTE: "remote",
}

// UNKNOWN_PROCESS nombre de proceso cuando no se puede saber de quien es el socket (sin permiso para leer sus fd).
const UNKNOWN_PROCESS = "?"
