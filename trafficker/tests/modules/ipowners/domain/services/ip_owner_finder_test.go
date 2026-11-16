package services

import (
    "net/netip"
    "testing"
    "time"

    "trafficker/src/modules/ipowners/domain/entities"
    "trafficker/src/modules/ipowners/domain/services"
    "trafficker/tests"
)

func getIpOwnerEntity(start string, end string, org string, resolvedAt time.Time) *entities.IpOwnerEntity {
    return entities.NewIpOwnerEntity(netip.MustParseAddr(start), netip.MustParseAddr(end), "NET", org, "ES", resolvedAt)
}

func TestSmallestValidRangeWins(t *testing.T) {
    ipOwnerEntities := []*entities.IpOwnerEntity{
        getIpOwnerEntity("213.94.0.0", "213.94.255.255", "Provider", parserNow),
        getIpOwnerEntity("213.94.49.0", "213.94.49.255", "Customer", parserNow),
    }
    found := services.IpOwnerFinderGetInstance().GetIpOwnerEntityOrNil(ipOwnerEntities, netip.MustParseAddr("213.94.49.55"), parserNow)
    tests.AssertEquals(t, "Customer", found.GetOrg())
}

func TestExpiredRangesAreIgnoredAndDropped(t *testing.T) {
    ipOwnerFinder := services.IpOwnerFinderGetInstance()
    oldIpOwnerEntity := getIpOwnerEntity("8.8.8.0", "8.8.8.255", "Google LLC", parserNow.AddDate(0, 0, -16))
    if found := ipOwnerFinder.GetIpOwnerEntityOrNil([]*entities.IpOwnerEntity{oldIpOwnerEntity}, netip.MustParseAddr("8.8.8.8"), parserNow); found != nil {
        t.Fatalf("expected nil, the range is 16 days old")
    }
    freshIpOwnerEntity := getIpOwnerEntity("1.1.1.0", "1.1.1.255", "APNIC", parserNow.AddDate(0, 0, -14))
    valid := ipOwnerFinder.GetValidIpOwnerEntities([]*entities.IpOwnerEntity{oldIpOwnerEntity, freshIpOwnerEntity}, parserNow)
    tests.AssertEquals(t, 1, len(valid))
    tests.AssertEquals(t, "APNIC", valid[0].GetOrg())
}

func TestFailedLookupIsRememberedOneHour(t *testing.T) {
    failedIpOwnerEntity := entities.FailedIpOwnerEntity(netip.MustParseAddr("5.5.5.5"), parserNow)
    tests.AssertEquals(t, "?", failedIpOwnerEntity.GetLabel())
    tests.AssertEquals(t, false, failedIpOwnerEntity.IsExpired(parserNow.Add(59*time.Minute)))
    tests.AssertEquals(t, true, failedIpOwnerEntity.IsExpired(parserNow.Add(time.Hour)))
}

func TestIpv4NeverMatchesAnIpv6Range(t *testing.T) {
    ipv6IpOwnerEntity := getIpOwnerEntity("::", "::ffff:ffff:ffff", "IPv6 range", parserNow)
    tests.AssertEquals(t, false, ipv6IpOwnerEntity.Contains(netip.MustParseAddr("8.8.8.8")))
}

func TestCachePrimitivesRoundTrip(t *testing.T) {
    ipOwnerEntity := getIpOwnerEntity("8.8.8.0", "8.8.8.255", "Google LLC", parserNow)
    restored := entities.IpOwnerEntityFromCachePrimitives(toJsonPrimitives(t, ipOwnerEntity.ToCachePrimitives()))
    tests.AssertEquals(t, ipOwnerEntity.ToCachePrimitives(), restored.ToCachePrimitives())
    tests.AssertEquals(t, true, restored.Contains(netip.MustParseAddr("8.8.8.8")))
}

func TestScopes(t *testing.T) {
    ipScopeResolver := services.IpScopeResolverGetInstance()
    tests.AssertEquals(t, "private", ipScopeResolver.GetScope(netip.MustParseAddr("10.0.0.191")))
    tests.AssertEquals(t, "loopback", ipScopeResolver.GetScope(netip.MustParseAddr("::1")))
    tests.AssertEquals(t, "link-local", ipScopeResolver.GetScope(netip.MustParseAddr("169.254.169.254")))
    tests.AssertEquals(t, "cgnat", ipScopeResolver.GetScope(netip.MustParseAddr("100.70.1.1")))
    tests.AssertEquals(t, "private", ipScopeResolver.GetScope(netip.MustParseAddr("::ffff:192.168.1.5")))
    tests.AssertEquals(t, "public", ipScopeResolver.GetScope(netip.MustParseAddr("160.79.104.10")))
}
