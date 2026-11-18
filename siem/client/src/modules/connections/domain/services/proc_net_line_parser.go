package services

import (
    "encoding/hex"
    "net/netip"
    "slices"
    "strconv"
    "strings"

    "siem-client/src/modules/connections/domain/enums"
    "siem-client/src/modules/connections/domain/entities"
)

// ProcNetLineParser convierte una linea de /proc/net/{tcp,tcp6,udp,udp6} en un SocketEntity (columnas en
// ProcNetColumnEnum). Las IP van en hex, en palabras de 32 bits en el orden de bytes del host (little-endian en
// amd64/arm64): 0100007F es 127.0.0.1. Los puertos van en hex normal (0277 = 631).
type ProcNetLineParser struct{}

func ProcNetLineParserGetInstance() *ProcNetLineParser {
    return &ProcNetLineParser{}
}

// GetSocketEntityOrNil nil para la cabecera y para cualquier linea que no se entienda (no rompe la foto entera).
func (p *ProcNetLineParser) GetSocketEntityOrNil(procNetFile string, line string) *entities.SocketEntity {
    columns := strings.Fields(line)
    if len(columns) < enums.ProcNetColumnEnum.MIN_COLUMNS || !strings.HasSuffix(columns[0], enums.ProcNetColumnEnum.SLOT_SUFFIX) {
        return nil
    }
    protocol, ipVersion := p.getProtocolAndIpVersion(procNetFile)
    if protocol == "" {
        return nil
    }
    localAddr, localPort, isLocalOk := p.getAddrAndPort(columns[enums.ProcNetColumnEnum.LOCAL_ADDRESS])
    remoteAddr, remotePort, isRemoteOk := p.getAddrAndPort(columns[enums.ProcNetColumnEnum.REMOTE_ADDRESS])
    if !isLocalOk || !isRemoteOk {
        return nil
    }
    uid, err := strconv.Atoi(columns[enums.ProcNetColumnEnum.UID])
    if err != nil {
        return nil
    }
    inode, err := strconv.ParseUint(columns[enums.ProcNetColumnEnum.INODE], enums.ProcNetColumnEnum.DECIMAL_BASE, enums.ProcNetColumnEnum.INODE_BITS)
    if err != nil {
        return nil
    }
    return entities.SocketEntityFromPrimitives(map[string]any{
        enums.ConnectionKeyEnum.PROTOCOL:    protocol,
        enums.ConnectionKeyEnum.IP_VERSION:  ipVersion,
        enums.ConnectionKeyEnum.LOCAL_IP:    localAddr.String(),
        enums.ConnectionKeyEnum.LOCAL_PORT:  localPort,
        enums.ConnectionKeyEnum.REMOTE_IP:   remoteAddr.String(),
        enums.ConnectionKeyEnum.REMOTE_PORT: remotePort,
        enums.ConnectionKeyEnum.STATE:       p.getState(protocol, strings.ToUpper(columns[enums.ProcNetColumnEnum.STATE]), remoteAddr, remotePort),
        enums.ConnectionKeyEnum.UID:         uid,
        enums.ConnectionKeyEnum.INODE:       int(inode),
    })
}

func (p *ProcNetLineParser) getProtocolAndIpVersion(procNetFile string) (string, int) {
    switch procNetFile {
    case enums.ProcNetFileEnum.TCP:
        return enums.ProtocolEnum.TCP, enums.IpVersionEnum.V4
    case enums.ProcNetFileEnum.TCP6:
        return enums.ProtocolEnum.TCP, enums.IpVersionEnum.V6
    case enums.ProcNetFileEnum.UDP:
        return enums.ProtocolEnum.UDP, enums.IpVersionEnum.V4
    case enums.ProcNetFileEnum.UDP6:
        return enums.ProtocolEnum.UDP, enums.IpVersionEnum.V6
    }
    return "", 0
}

// getState TCP: el nombre del estado. UDP no tiene estados: con extremo remoto es una "conexion" (ESTABLISHED) y sin
// el un socket abierto que recibe de cualquiera (UNCONN).
func (p *ProcNetLineParser) getState(protocol string, stateHex string, remoteAddr netip.Addr, remotePort int) string {
    if protocol == enums.ProtocolEnum.UDP {
        if remotePort == 0 && remoteAddr.IsUnspecified() {
            return enums.TcpStateEnum.UNCONN
        }
        return enums.TcpStateEnum.ESTABLISHED
    }
    if state, isKnown := enums.TcpStateByHex[stateHex]; isKnown {
        return state
    }
    return stateHex
}

// getAddrAndPort "0100007F:0277" -> 127.0.0.1, 631.
func (p *ProcNetLineParser) getAddrAndPort(rawEndpoint string) (netip.Addr, int, bool) {
    rawAddr, rawPort, isFound := strings.Cut(rawEndpoint, enums.ProcNetColumnEnum.ENDPOINT_SEPARATOR)
    if !isFound {
        return netip.Addr{}, 0, false
    }
    port, err := strconv.ParseUint(rawPort, enums.ProcNetColumnEnum.HEX_BASE, enums.ProcNetColumnEnum.PORT_BITS)
    if err != nil {
        return netip.Addr{}, 0, false
    }
    addrBytes, err := hex.DecodeString(rawAddr)
    if err != nil || (len(addrBytes) != enums.IpVersionEnum.V4_BYTES && len(addrBytes) != enums.IpVersionEnum.V6_BYTES) {
        return netip.Addr{}, 0, false
    }
    // cada palabra de 32 bits viene con los bytes al reves
    wordBytes := enums.ProcNetColumnEnum.ADDRESS_WORD_BYTES
    for wordStart := 0; wordStart < len(addrBytes); wordStart += wordBytes {
        slices.Reverse(addrBytes[wordStart : wordStart+wordBytes])
    }
    addr, isOk := netip.AddrFromSlice(addrBytes)
    if !isOk {
        return netip.Addr{}, 0, false
    }
    return addr, int(port), true
}
