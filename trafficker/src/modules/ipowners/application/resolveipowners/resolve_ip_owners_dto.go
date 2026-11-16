package resolveipowners

import "net/netip"

// ResolveIpOwnersDto las IP remotas de las que se quiere saber el dueño (repetidas o no validas se ignoran).
type ResolveIpOwnersDto struct {
    addrs []netip.Addr
}

func ResolveIpOwnersDtoFromPrimitives(primitives map[string]any) *ResolveIpOwnersDto {
    rawIps, _ := primitives["ips"].([]string)
    seen := map[netip.Addr]bool{}
    addrs := []netip.Addr{}
    for _, rawIp := range rawIps {
        addr, err := netip.ParseAddr(rawIp)
        if err != nil {
            continue
        }
        addr = addr.Unmap()
        if seen[addr] {
            continue
        }
        seen[addr] = true
        addrs = append(addrs, addr)
    }
    return &ResolveIpOwnersDto{addrs: addrs}
}

func (d *ResolveIpOwnersDto) GetAddrs() []netip.Addr {
    return d.addrs
}
