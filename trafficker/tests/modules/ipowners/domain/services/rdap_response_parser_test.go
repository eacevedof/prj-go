package services

import (
    "encoding/json"
    "net/netip"
    "testing"
    "time"

    "trafficker/src/modules/ipowners/domain/services"
    "trafficker/tests"
)

// respuestas reales de rdap.org recortadas a lo que se usa (sin las personas de contacto)
const (
    arinResponse = `{"name": "GOGL", "startAddress": "8.8.8.0", "endAddress": "8.8.8.255",
        "entities": [{"roles": ["registrant"], "vcardArray": ["vcard", [["version", {}, "text", "4.0"],
            ["fn", {}, "text", "Google LLC"],
            ["adr", {"label": "1600 Amphitheatre Parkway\nMountain View\nCA\n94043\nUnited States"}, "text", ["", "", "", "", "", "", ""]]]],
            "entities": [{"roles": ["abuse"], "vcardArray": ["vcard", [["fn", {}, "text", "Abuse"]]]}]}]}`
    ripeResponse = `{"name": "IBERCOMNET", "country": "ES", "startAddress": "213.94.1.0", "endAddress": "213.94.63.255",
        "remarks": [{"description": ["WWW Ibercom Net", "Spain"]}],
        "entities": [{"roles": ["abuse"], "vcardArray": ["vcard", [["fn", {}, "text", "MasMovil Abuse Team"]]]},
            {"roles": ["registrant"], "vcardArray": ["vcard", [["fn", {}, "text", "TISCALI-ES-MNT"]]]}]}`
    apnicResponse = `{"name": "APNIC-LABS", "country": "au", "startAddress": "1.1.1.0", "endAddress": "1.1.1.255",
        "remarks": [{"title": "remarks", "description": ["---------------"]},
            {"title": "description", "description": ["APNIC and Cloudflare DNS Resolver project"]}],
        "entities": [{"roles": ["registrant"], "vcardArray": ["vcard", [["fn", {}, "text", "APNIC Research and Development"]]]}]}`
)

var parserNow = time.Date(2026, 11, 16, 10, 0, 0, 0, time.UTC)

func getRdapResponse(t *testing.T, rawJson string) map[string]any {
    t.Helper()
    rdapResponse := map[string]any{}
    if err := json.Unmarshal([]byte(rawJson), &rdapResponse); err != nil {
        t.Fatal(err)
    }
    return rdapResponse
}

func TestArinOrgFromRegistrantAndCountryCodeFromAddress(t *testing.T) {
    ipOwnerEntity := services.RdapResponseParserGetInstance().GetIpOwnerEntityOrNil(
        getRdapResponse(t, arinResponse), netip.MustParseAddr("8.8.8.8"), parserNow,
    )
    tests.AssertEquals(t, "GOGL", ipOwnerEntity.GetNetwork())
    tests.AssertEquals(t, "Google LLC", ipOwnerEntity.GetOrg())
    tests.AssertEquals(t, "US", ipOwnerEntity.GetCountry())
    tests.AssertEquals(t, true, ipOwnerEntity.Contains(netip.MustParseAddr("8.8.8.200")))
    tests.AssertEquals(t, false, ipOwnerEntity.Contains(netip.MustParseAddr("8.8.9.1")))
}

func TestRipeMaintainerIsSkippedForTheDescription(t *testing.T) {
    ipOwnerEntity := services.RdapResponseParserGetInstance().GetIpOwnerEntityOrNil(
        getRdapResponse(t, ripeResponse), netip.MustParseAddr("213.94.49.55"), parserNow,
    )
    tests.AssertEquals(t, "WWW Ibercom Net", ipOwnerEntity.GetOrg())
    tests.AssertEquals(t, "ES", ipOwnerEntity.GetCountry())
    tests.AssertEquals(t, "WWW Ibercom Net", ipOwnerEntity.GetLabel())
}

func TestApnicRegistrantAndUppercasedCountry(t *testing.T) {
    ipOwnerEntity := services.RdapResponseParserGetInstance().GetIpOwnerEntityOrNil(
        getRdapResponse(t, apnicResponse), netip.MustParseAddr("1.1.1.1"), parserNow,
    )
    tests.AssertEquals(t, "APNIC Research and Development", ipOwnerEntity.GetOrg())
    tests.AssertEquals(t, "AU", ipOwnerEntity.GetCountry())
}

func TestNetworkNameWhenThereIsNoOrg(t *testing.T) {
    ipOwnerEntity := services.RdapResponseParserGetInstance().GetIpOwnerEntityOrNil(
        map[string]any{"name": "SOME-NET"}, netip.MustParseAddr("5.5.5.5"), parserNow,
    )
    tests.AssertEquals(t, "", ipOwnerEntity.GetOrg())
    tests.AssertEquals(t, "SOME-NET", ipOwnerEntity.GetLabel())
    // sin rango en la respuesta: solo la IP preguntada
    tests.AssertEquals(t, true, ipOwnerEntity.Contains(netip.MustParseAddr("5.5.5.5")))
    tests.AssertEquals(t, false, ipOwnerEntity.Contains(netip.MustParseAddr("5.5.5.6")))
}

func TestEmptyResponseIsNil(t *testing.T) {
    ipOwnerEntity := services.RdapResponseParserGetInstance().GetIpOwnerEntityOrNil(
        map[string]any{"errorCode": 404.0}, netip.MustParseAddr("5.5.5.5"), parserNow,
    )
    if ipOwnerEntity != nil {
        t.Fatalf("expected nil, got %#v", ipOwnerEntity)
    }
}
