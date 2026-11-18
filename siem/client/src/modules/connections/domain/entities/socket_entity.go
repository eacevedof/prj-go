package entities

import (
    "net/netip"
    "strconv"
    "strings"

    "siem-client/src/modules/shared/application"

    "siem-client/src/modules/connections/domain/enums"
)

// SocketEntity una linea de /proc/net/{tcp,tcp6,udp,udp6}: los dos extremos, el estado, el uid del duenyo y el inodo
// (lo que une el socket con el proceso que lo tiene abierto).
type SocketEntity struct {
    protocol   string
    ipVersion  int
    localAddr  netip.Addr
    localPort  int
    remoteAddr netip.Addr
    remotePort int
    state      string
    uid        int
    inode      uint64
}

// SocketEntityFromPrimitives claves de ConnectionKeyEnum. Las IP en texto ("192.168.1.5", "::1"); una IPv4 mapeada
// en IPv6 se guarda como IPv4. Una IP que no se puede leer queda vacia (sin direccion).
func SocketEntityFromPrimitives(primitives map[string]any) *SocketEntity {
    localAddr, _ := netip.ParseAddr(application.GetTrimmedString(primitives, enums.ConnectionKeyEnum.LOCAL_IP))
    remoteAddr, _ := netip.ParseAddr(application.GetTrimmedString(primitives, enums.ConnectionKeyEnum.REMOTE_IP))
    return &SocketEntity{
        protocol:   application.GetTrimmedString(primitives, enums.ConnectionKeyEnum.PROTOCOL),
        ipVersion:  application.GetInt(primitives, enums.ConnectionKeyEnum.IP_VERSION),
        localAddr:  localAddr.Unmap(),
        localPort:  application.GetInt(primitives, enums.ConnectionKeyEnum.LOCAL_PORT),
        remoteAddr: remoteAddr.Unmap(),
        remotePort: application.GetInt(primitives, enums.ConnectionKeyEnum.REMOTE_PORT),
        state:      application.GetTrimmedString(primitives, enums.ConnectionKeyEnum.STATE),
        uid:        application.GetInt(primitives, enums.ConnectionKeyEnum.UID),
        inode:      uint64(application.GetInt(primitives, enums.ConnectionKeyEnum.INODE)),
    }
}

func (e *SocketEntity) GetProtocol() string {
    return e.protocol
}

func (e *SocketEntity) GetIpVersion() int {
    return e.ipVersion
}

func (e *SocketEntity) GetLocalAddr() netip.Addr {
    return e.localAddr
}

func (e *SocketEntity) GetLocalIp() string {
    return e.localAddr.String()
}

func (e *SocketEntity) GetLocalPort() int {
    return e.localPort
}

func (e *SocketEntity) GetRemoteAddr() netip.Addr {
    return e.remoteAddr
}

func (e *SocketEntity) GetRemoteIp() string {
    return e.remoteAddr.String()
}

func (e *SocketEntity) GetRemotePort() int {
    return e.remotePort
}

func (e *SocketEntity) GetState() string {
    return e.state
}

func (e *SocketEntity) GetUid() int {
    return e.uid
}

func (e *SocketEntity) GetInode() uint64 {
    return e.inode
}

// IsListening TCP en LISTEN, o UDP sin extremo remoto (un servicio UDP que recibe de cualquiera).
func (e *SocketEntity) IsListening() bool {
    if e.protocol == enums.ProtocolEnum.TCP {
        return e.state == enums.TcpStateEnum.LISTEN
    }
    return e.remotePort == 0 && (!e.remoteAddr.IsValid() || e.remoteAddr.IsUnspecified())
}

// IsLoopback los dos extremos en la propia maquina (127.0.0.0/8, ::1). Un socket en escucha cuenta por su IP local.
func (e *SocketEntity) IsLoopback() bool {
    if e.IsListening() {
        return e.localAddr.IsLoopback()
    }
    return e.localAddr.IsLoopback() && e.remoteAddr.IsLoopback()
}

// GetLocalEndpoint "192.168.1.5:52780" / "[::1]:631".
func (e *SocketEntity) GetLocalEndpoint() string {
    return netip.AddrPortFrom(e.localAddr, uint16(e.localPort)).String()
}

// GetRemoteEndpoint como GetLocalEndpoint; "*" para un socket en escucha.
func (e *SocketEntity) GetRemoteEndpoint() string {
    if e.IsListening() {
        return enums.ConnectionFormatEnum.ANY_ENDPOINT
    }
    return netip.AddrPortFrom(e.remoteAddr, uint16(e.remotePort)).String()
}

// GetKey identifica el socket entre dos fotos: protocolo + extremos + inodo (el inodo cambia si se cierra y se abre
// otro con los mismos extremos).
func (e *SocketEntity) GetKey() string {
    return strings.Join([]string{
        e.protocol, e.GetLocalEndpoint(), e.GetRemoteEndpoint(), strconv.FormatUint(e.inode, enums.ProcNetColumnEnum.DECIMAL_BASE),
    }, enums.ConnectionFormatEnum.KEY_SEPARATOR)
}
