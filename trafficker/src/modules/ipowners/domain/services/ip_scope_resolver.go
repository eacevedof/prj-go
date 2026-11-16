package services

import (
    "net/netip"

    "trafficker/src/modules/ipowners/domain/enums"
)

var cgnatPrefix = netip.MustParsePrefix("100.64.0.0/10")

// IpScopeResolver de que tipo es una IP: solo las public tiene sentido preguntarlas por whois.
type IpScopeResolver struct{}

func IpScopeResolverGetInstance() *IpScopeResolver {
    return &IpScopeResolver{}
}

func (r *IpScopeResolver) GetScope(addr netip.Addr) string {
    addr = addr.Unmap()
    switch {
    case !addr.IsValid() || addr.IsUnspecified():
        return enums.IpScopeEnum.UNSPECIFIED
    case addr.IsLoopback():
        return enums.IpScopeEnum.LOOPBACK
    case addr.IsPrivate():
        return enums.IpScopeEnum.PRIVATE
    case addr.IsLinkLocalUnicast():
        return enums.IpScopeEnum.LINK_LOCAL
    case addr.IsMulticast() || addr.IsLinkLocalMulticast():
        return enums.IpScopeEnum.MULTICAST
    case cgnatPrefix.Contains(addr):
        return enums.IpScopeEnum.CGNAT
    }
    return enums.IpScopeEnum.PUBLIC
}
