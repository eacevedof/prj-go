package entities

import (
    "net/netip"
    "strconv"

    "trafficker/src/modules/connections/domain/enums"
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

func NewSocketEntity(protocol string, ipVersion int, localAddr netip.Addr, localPort int, remoteAddr netip.Addr,
    remotePort int, state string, uid int, inode uint64) *SocketEntity {
    return &SocketEntity{
        protocol:   protocol,
        ipVersion:  ipVersion,
        localAddr:  localAddr.Unmap(),
        localPort:  localPort,
        remoteAddr: remoteAddr.Unmap(),
        remotePort: remotePort,
        state:      state,
        uid:        uid,
        inode:      inode,
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

// HasPublicRemote la IP remota es de internet: ni loopback, ni privada (RFC 1918 / fc00::/7), ni link-local, ni
// CGNAT (100.64.0.0/10), ni sin especificar. Son las que tiene sentido preguntar por whois.
func (e *SocketEntity) HasPublicRemote() bool {
    remoteAddr := e.remoteAddr
    if !remoteAddr.IsValid() || remoteAddr.IsUnspecified() || remoteAddr.IsLoopback() || remoteAddr.IsPrivate() ||
        remoteAddr.IsLinkLocalUnicast() || remoteAddr.IsLinkLocalMulticast() || remoteAddr.IsMulticast() {
        return false
    }
    return !cgnatPrefix.Contains(remoteAddr)
}

// GetLocalEndpoint "192.168.1.5:52780" / "[::1]:631".
func (e *SocketEntity) GetLocalEndpoint() string {
    return netip.AddrPortFrom(e.localAddr, uint16(e.localPort)).String()
}

// GetRemoteEndpoint como GetLocalEndpoint; "*" para un socket en escucha.
func (e *SocketEntity) GetRemoteEndpoint() string {
    if e.IsListening() {
        return "*"
    }
    return netip.AddrPortFrom(e.remoteAddr, uint16(e.remotePort)).String()
}

// GetKey identifica el socket entre dos fotos: protocolo + extremos + inodo (el inodo cambia si se cierra y se abre
// otro con los mismos extremos).
func (e *SocketEntity) GetKey() string {
    return e.protocol + "|" + e.GetLocalEndpoint() + "|" + e.GetRemoteEndpoint() + "|" + strconv.FormatUint(e.inode, 10)
}

var cgnatPrefix = netip.MustParsePrefix("100.64.0.0/10")
