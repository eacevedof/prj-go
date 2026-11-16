package services

import (
    "net/netip"
    "time"

    "trafficker/src/modules/ipowners/domain/entities"
)

// IpOwnerFinder busca una IP en los rangos ya conocidos (cache): el rango vigente mas pequeño que la contiene (un
// bloque asignado a un cliente dentro del bloque grande de su proveedor dice mas).
type IpOwnerFinder struct{}

func IpOwnerFinderGetInstance() *IpOwnerFinder {
    return &IpOwnerFinder{}
}

// GetIpOwnerEntityOrNil nil si ningun rango vigente la contiene: hay que preguntarla.
func (f *IpOwnerFinder) GetIpOwnerEntityOrNil(ipOwnerEntities []*entities.IpOwnerEntity, addr netip.Addr, now time.Time) *entities.IpOwnerEntity {
    var bestIpOwnerEntity *entities.IpOwnerEntity
    for _, ipOwnerEntity := range ipOwnerEntities {
        if ipOwnerEntity.IsExpired(now) || !ipOwnerEntity.Contains(addr) {
            continue
        }
        if bestIpOwnerEntity == nil || ipOwnerEntity.GetSize() < bestIpOwnerEntity.GetSize() {
            bestIpOwnerEntity = ipOwnerEntity
        }
    }
    return bestIpOwnerEntity
}

// GetValidIpOwnerEntities los que no han caducado (lo que se vuelve a guardar en la cache).
func (f *IpOwnerFinder) GetValidIpOwnerEntities(ipOwnerEntities []*entities.IpOwnerEntity, now time.Time) []*entities.IpOwnerEntity {
    validIpOwnerEntities := []*entities.IpOwnerEntity{}
    for _, ipOwnerEntity := range ipOwnerEntities {
        if !ipOwnerEntity.IsExpired(now) {
            validIpOwnerEntities = append(validIpOwnerEntities, ipOwnerEntity)
        }
    }
    return validIpOwnerEntities
}
