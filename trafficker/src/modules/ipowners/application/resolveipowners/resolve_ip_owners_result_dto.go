package resolveipowners

import (
    "net/netip"

    "trafficker/src/modules/ipowners/domain/entities"
)

type ResolveIpOwnersResultDto struct {
    ipOwnerEntityByIp map[string]*entities.IpOwnerEntity
    fromCache         int
    lookedUp          int
    failed            int
    // IP publicas que no se han preguntado en esta pasada (mas de MAX_LOOKUPS_PER_RUN): salen sin dueño
    pending int
}

func (d *ResolveIpOwnersResultDto) GetIpOwnerEntityOrNil(ip string) *entities.IpOwnerEntity {
    addr, err := netip.ParseAddr(ip)
    if err != nil {
        return nil
    }
    return d.ipOwnerEntityByIp[addr.Unmap().String()]
}

func (d *ResolveIpOwnersResultDto) GetFromCache() int {
    return d.fromCache
}

func (d *ResolveIpOwnersResultDto) GetLookedUp() int {
    return d.lookedUp
}

func (d *ResolveIpOwnersResultDto) GetFailed() int {
    return d.failed
}

func (d *ResolveIpOwnersResultDto) GetPending() int {
    return d.pending
}
