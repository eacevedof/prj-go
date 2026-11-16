package entities

import (
    "net/netip"
    "time"

    "trafficker/src/modules/ipowners/domain/enums"
)

const dateTimeFormat = "2006-01-02 15:04:05"

// IpOwnerEntity de quien es un rango de IP: nombre de la red, organizacion y pais (lo que dice el registro regional
// por RDAP). Para IP no publicas solo lleva el scope (private, loopback...) y no se pregunta a nadie.
type IpOwnerEntity struct {
    startAddr  netip.Addr
    endAddr    netip.Addr
    network    string
    org        string
    country    string
    scope      string
    isFailed   bool
    resolvedAt time.Time
    expiresAt  time.Time
}

func NewIpOwnerEntity(startAddr netip.Addr, endAddr netip.Addr, network string, org string, country string,
    resolvedAt time.Time) *IpOwnerEntity {
    return &IpOwnerEntity{
        startAddr:  startAddr.Unmap(),
        endAddr:    endAddr.Unmap(),
        network:    network,
        org:        org,
        country:    country,
        scope:      enums.IpScopeEnum.PUBLIC,
        resolvedAt: resolvedAt.UTC(),
        expiresAt:  resolvedAt.UTC().AddDate(0, 0, enums.IpOwnersCacheEnum.TTL_DAYS),
    }
}

// FailedIpOwnerEntity la consulta de esa IP ha fallado: se recuerda poco tiempo para no insistir en cada pasada.
func FailedIpOwnerEntity(addr netip.Addr, resolvedAt time.Time) *IpOwnerEntity {
    return &IpOwnerEntity{
        startAddr:  addr.Unmap(),
        endAddr:    addr.Unmap(),
        scope:      enums.IpScopeEnum.PUBLIC,
        isFailed:   true,
        resolvedAt: resolvedAt.UTC(),
        expiresAt:  resolvedAt.UTC().Add(time.Duration(enums.IpOwnersCacheEnum.FAILED_TTL_HOURS) * time.Hour),
    }
}

// LocalIpOwnerEntity IP que no es de internet (privada, loopback...): no se cachea ni se pregunta.
func LocalIpOwnerEntity(addr netip.Addr, scope string) *IpOwnerEntity {
    return &IpOwnerEntity{startAddr: addr.Unmap(), endAddr: addr.Unmap(), scope: scope}
}

// IpOwnerEntityFromCachePrimitives nil si la entrada no se entiende (se ignora y se vuelve a preguntar).
func IpOwnerEntityFromCachePrimitives(primitives map[string]any) *IpOwnerEntity {
    startAddr, errStart := netip.ParseAddr(getString(primitives, enums.IpOwnerKeyEnum.START_IP))
    endAddr, errEnd := netip.ParseAddr(getString(primitives, enums.IpOwnerKeyEnum.END_IP))
    resolvedAt, errResolved := time.Parse(dateTimeFormat, getString(primitives, enums.IpOwnerKeyEnum.RESOLVED_AT))
    expiresAt, errExpires := time.Parse(dateTimeFormat, getString(primitives, enums.IpOwnerKeyEnum.EXPIRES_AT))
    if errStart != nil || errEnd != nil || errResolved != nil || errExpires != nil {
        return nil
    }
    isFailed, _ := primitives[enums.IpOwnerKeyEnum.IS_FAILED].(bool)
    return &IpOwnerEntity{
        startAddr:  startAddr,
        endAddr:    endAddr,
        network:    getString(primitives, enums.IpOwnerKeyEnum.NETWORK),
        org:        getString(primitives, enums.IpOwnerKeyEnum.ORG),
        country:    getString(primitives, enums.IpOwnerKeyEnum.COUNTRY),
        scope:      enums.IpScopeEnum.PUBLIC,
        isFailed:   isFailed,
        resolvedAt: resolvedAt,
        expiresAt:  expiresAt,
    }
}

func (e *IpOwnerEntity) GetNetwork() string {
    return e.network
}

func (e *IpOwnerEntity) GetOrg() string {
    return e.org
}

func (e *IpOwnerEntity) GetCountry() string {
    return e.country
}

func (e *IpOwnerEntity) GetScope() string {
    return e.scope
}

func (e *IpOwnerEntity) IsFailed() bool {
    return e.isFailed
}

func (e *IpOwnerEntity) IsPublic() bool {
    return e.scope == enums.IpScopeEnum.PUBLIC
}

// Contains la IP cae en el rango (misma familia: una IPv4 nunca cae en un rango IPv6).
func (e *IpOwnerEntity) Contains(addr netip.Addr) bool {
    addr = addr.Unmap()
    return addr.BitLen() == e.startAddr.BitLen() && e.startAddr.Compare(addr) <= 0 && addr.Compare(e.endAddr) <= 0
}

func (e *IpOwnerEntity) IsExpired(now time.Time) bool {
    return !now.UTC().Before(e.expiresAt)
}

// GetSize cuantas IP hay en el rango, aproximado (para preferir el rango mas pequeño si hay varios que la contienen).
func (e *IpOwnerEntity) GetSize() float64 {
    startBytes, endBytes := e.startAddr.AsSlice(), e.endAddr.AsSlice()
    size := 0.0
    for index := range startBytes {
        size = size*256 + float64(int(endBytes[index])-int(startBytes[index]))
    }
    return size
}

// GetLabel lo que se ensenya en la tabla: la organizacion, el scope para IP locales, "?" si fallo.
func (e *IpOwnerEntity) GetLabel() string {
    if !e.IsPublic() {
        return e.scope
    }
    if e.isFailed {
        return "?"
    }
    if e.org != "" {
        return e.org
    }
    return e.network
}

// ToCachePrimitives lo que se guarda en ip-owners.json.
func (e *IpOwnerEntity) ToCachePrimitives() map[string]any {
    return map[string]any{
        enums.IpOwnerKeyEnum.START_IP:    e.startAddr.String(),
        enums.IpOwnerKeyEnum.END_IP:      e.endAddr.String(),
        enums.IpOwnerKeyEnum.NETWORK:     e.network,
        enums.IpOwnerKeyEnum.ORG:         e.org,
        enums.IpOwnerKeyEnum.COUNTRY:     e.country,
        enums.IpOwnerKeyEnum.IS_FAILED:   e.isFailed,
        enums.IpOwnerKeyEnum.RESOLVED_AT: e.resolvedAt.Format(dateTimeFormat),
        enums.IpOwnerKeyEnum.EXPIRES_AT:  e.expiresAt.Format(dateTimeFormat),
    }
}

// ToPrimitives lo que se anyade a cada conexion en connections-list --format=json.
func (e *IpOwnerEntity) ToPrimitives() map[string]any {
    return map[string]any{
        "remote_scope":   e.scope,
        "remote_network": e.network,
        "remote_org":     e.org,
        "remote_country": e.country,
    }
}

func getString(primitives map[string]any, key string) string {
    value, _ := primitives[key].(string)
    return value
}
