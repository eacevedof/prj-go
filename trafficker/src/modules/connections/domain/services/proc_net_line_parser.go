package services

import (
    "encoding/hex"
    "net/netip"
    "strconv"
    "strings"

    "trafficker/src/modules/connections/domain/entities"
    "trafficker/src/modules/connections/domain/enums"
)

// ProcNetLineParser convierte una linea de /proc/net/{tcp,tcp6,udp,udp6} en un SocketEntity.
//
//   sl  local_address rem_address   st tx_queue:rx_queue tr:tm->when retrnsmt   uid  timeout inode
//    0: 0100007F:0277 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 22512 1 ...
//
// Las IP van en hex, en palabras de 32 bits en el orden de bytes del host (little-endian en amd64/arm64): 0100007F es
// 127.0.0.1. Los puertos van en hex normal (0277 = 631).
type ProcNetLineParser struct{}

func ProcNetLineParserGetInstance() *ProcNetLineParser {
    return &ProcNetLineParser{}
}

// GetSocketEntityOrNil nil para la cabecera y para cualquier linea que no se entienda (no rompe la foto entera).
func (p *ProcNetLineParser) GetSocketEntityOrNil(procNetFile string, line string) *entities.SocketEntity {
    fields := strings.Fields(line)
    if len(fields) < 10 || !strings.HasSuffix(fields[0], ":") {
        return nil
    }
    protocol, ipVersion := p.getProtocolAndIpVersion(procNetFile)
    if protocol == "" {
        return nil
    }
    localAddr, localPort, isLocalOk := p.getAddrAndPort(fields[1])
    remoteAddr, remotePort, isRemoteOk := p.getAddrAndPort(fields[2])
    if !isLocalOk || !isRemoteOk {
        return nil
    }
    uid, err := strconv.Atoi(fields[7])
    if err != nil {
        return nil
    }
    inode, err := strconv.ParseUint(fields[9], 10, 64)
    if err != nil {
        return nil
    }
    state := p.getState(protocol, strings.ToUpper(fields[3]), remoteAddr, remotePort)
    return entities.NewSocketEntity(protocol, ipVersion, localAddr, localPort, remoteAddr, remotePort, state, uid, inode)
}

func (p *ProcNetLineParser) getProtocolAndIpVersion(procNetFile string) (string, int) {
    switch procNetFile {
    case enums.ProcNetFileEnum.TCP:
        return enums.ProtocolEnum.TCP, 4
    case enums.ProcNetFileEnum.TCP6:
        return enums.ProtocolEnum.TCP, 6
    case enums.ProcNetFileEnum.UDP:
        return enums.ProtocolEnum.UDP, 4
    case enums.ProcNetFileEnum.UDP6:
        return enums.ProtocolEnum.UDP, 6
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
    rawAddr, rawPort, isFound := strings.Cut(rawEndpoint, ":")
    if !isFound {
        return netip.Addr{}, 0, false
    }
    port, err := strconv.ParseUint(rawPort, 16, 16)
    if err != nil {
        return netip.Addr{}, 0, false
    }
    addrBytes, err := hex.DecodeString(rawAddr)
    if err != nil || (len(addrBytes) != 4 && len(addrBytes) != 16) {
        return netip.Addr{}, 0, false
    }
    // cada palabra de 32 bits viene con los bytes al reves
    for wordStart := 0; wordStart < len(addrBytes); wordStart += 4 {
        word := addrBytes[wordStart : wordStart+4]
        word[0], word[1], word[2], word[3] = word[3], word[2], word[1], word[0]
    }
    addr, isOk := netip.AddrFromSlice(addrBytes)
    if !isOk {
        return netip.Addr{}, 0, false
    }
    return addr, int(port), true
}
