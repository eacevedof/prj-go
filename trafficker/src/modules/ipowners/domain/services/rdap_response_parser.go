package services

import (
    "net/netip"
    "regexp"
    "strings"
    "time"

    "trafficker/src/modules/ipowners/domain/entities"
)

// un "registrante" que en realidad es un mantenedor de la base de datos de RIPE (TISCALI-ES-MNT), no una organizacion
var maintainerHandlePattern = regexp.MustCompile(`(?i)^[A-Z0-9_-]+-MNT$`)

// RdapResponseParser la respuesta RDAP de una red IP -> IpOwnerEntity. Cada registro regional rellena cosas
// distintas, por eso cada dato tiene su orden de preferencia:
//   - organizacion: el "fn" de la entidad registrante (ARIN, APNIC) salvo que sea un mantenedor (RIPE: XXX-MNT);
//     si no, la primera linea de la descripcion de la red (RIPE: "WWW Ibercom Net"); si no, el nombre de la red.
//   - pais: el "country" de la red (RIPE, APNIC); si no, la ultima linea de la direccion del registrante (ARIN: el
//     nombre del pais, "United States").
//   - rango: startAddress-endAddress; si falta, solo la IP preguntada.
type RdapResponseParser struct{}

func RdapResponseParserGetInstance() *RdapResponseParser {
    return &RdapResponseParser{}
}

// GetIpOwnerEntityOrNil nil si la respuesta no es una red IP (sin nombre ni organizacion).
func (p *RdapResponseParser) GetIpOwnerEntityOrNil(rdapResponse map[string]any, queriedAddr netip.Addr, now time.Time) *entities.IpOwnerEntity {
    network := strings.TrimSpace(getText(rdapResponse, "name"))
    registrant := p.getRegistrantOrNil(getList(rdapResponse, "entities"))
    org := p.getOrg(rdapResponse, registrant)
    if network == "" && org == "" {
        return nil
    }
    startAddr, endAddr := p.getRange(rdapResponse, queriedAddr)
    return entities.NewIpOwnerEntity(startAddr, endAddr, network, org, p.getCountry(rdapResponse, registrant), now)
}

func (p *RdapResponseParser) getOrg(rdapResponse map[string]any, registrant map[string]any) string {
    if registrant != nil {
        registrantName := p.getVcardText(registrant, "fn")
        if registrantName != "" && !maintainerHandlePattern.MatchString(registrantName) {
            return registrantName
        }
    }
    for _, rawRemark := range getList(rdapResponse, "remarks") {
        remark, _ := rawRemark.(map[string]any)
        title := strings.ToLower(getText(remark, "title"))
        if title != "" && title != "description" {
            continue
        }
        for _, rawLine := range getList(remark, "description") {
            if line, _ := rawLine.(string); strings.TrimSpace(line) != "" {
                return strings.TrimSpace(line)
            }
        }
    }
    return ""
}

func (p *RdapResponseParser) getCountry(rdapResponse map[string]any, registrant map[string]any) string {
    if country := strings.TrimSpace(getText(rdapResponse, "country")); country != "" {
        return strings.ToUpper(country)
    }
    if registrant == nil {
        return ""
    }
    addressLines := strings.Split(p.getVcardAddressLabel(registrant), "\n")
    countryName := strings.TrimSpace(addressLines[len(addressLines)-1])
    if countryCode, isKnown := arinCountryCodeByName[strings.ToLower(countryName)]; isKnown {
        return countryCode
    }
    return countryName
}

// arinCountryCodeByName ARIN no da el codigo de pais, solo el nombre en la direccion: los de su region a ISO 3166
// para que la columna sea homogenea con RIPE / APNIC / LACNIC / AFRINIC. Uno que no este aqui sale con su nombre.
var arinCountryCodeByName = map[string]string{
    "united states": "US", "canada": "CA", "puerto rico": "PR", "bahamas": "BS", "barbados": "BB", "bermuda": "BM",
    "cayman islands": "KY", "jamaica": "JM", "antigua and barbuda": "AG", "anguilla": "AI", "dominica": "DM",
    "grenada": "GD", "saint kitts and nevis": "KN", "saint lucia": "LC", "saint vincent and the grenadines": "VC",
    "turks and caicos islands": "TC", "virgin islands, british": "VG", "virgin islands, u.s.": "VI",
    "united states minor outlying islands": "UM", "saint pierre and miquelon": "PM", "antarctica": "AQ",
}

func (p *RdapResponseParser) getRange(rdapResponse map[string]any, queriedAddr netip.Addr) (netip.Addr, netip.Addr) {
    startAddr, errStart := netip.ParseAddr(getText(rdapResponse, "startAddress"))
    endAddr, errEnd := netip.ParseAddr(getText(rdapResponse, "endAddress"))
    queriedAddr = queriedAddr.Unmap()
    if errStart != nil || errEnd != nil || startAddr.BitLen() != queriedAddr.BitLen() || endAddr.Compare(startAddr) < 0 {
        return queriedAddr, queriedAddr
    }
    return startAddr, endAddr
}

// getRegistrantOrNil la primera entidad con rol registrant, primero las de arriba y luego las anidadas.
func (p *RdapResponseParser) getRegistrantOrNil(rawEntities []any) map[string]any {
    pending := rawEntities
    for len(pending) > 0 {
        nested := []any{}
        for _, rawEntity := range pending {
            entity, _ := rawEntity.(map[string]any)
            if entity == nil {
                continue
            }
            for _, rawRole := range getList(entity, "roles") {
                if role, _ := rawRole.(string); role == "registrant" {
                    return entity
                }
            }
            nested = append(nested, getList(entity, "entities")...)
        }
        pending = nested
    }
    return nil
}

// getVcardText vcardArray = ["vcard", [["fn", {}, "text", "Google LLC"], ...]]
func (p *RdapResponseParser) getVcardText(entity map[string]any, property string) string {
    for _, vcardProperty := range p.getVcardProperties(entity) {
        if len(vcardProperty) >= 4 && vcardProperty[0] == property {
            text, _ := vcardProperty[3].(string)
            return strings.TrimSpace(text)
        }
    }
    return ""
}

// getVcardAddressLabel ["adr", {"label": "1600 Amphitheatre Parkway\nMountain View\nCA\n94043\nUnited States"}, ...]
func (p *RdapResponseParser) getVcardAddressLabel(entity map[string]any) string {
    for _, vcardProperty := range p.getVcardProperties(entity) {
        if len(vcardProperty) >= 2 && vcardProperty[0] == "adr" {
            parameters, _ := vcardProperty[1].(map[string]any)
            return strings.TrimSpace(getText(parameters, "label"))
        }
    }
    return ""
}

func (p *RdapResponseParser) getVcardProperties(entity map[string]any) [][]any {
    vcardArray := getList(entity, "vcardArray")
    if len(vcardArray) < 2 {
        return nil
    }
    rawProperties, _ := vcardArray[1].([]any)
    properties := [][]any{}
    for _, rawProperty := range rawProperties {
        if property, isList := rawProperty.([]any); isList {
            properties = append(properties, property)
        }
    }
    return properties
}

func getText(data map[string]any, key string) string {
    if data == nil {
        return ""
    }
    text, _ := data[key].(string)
    return text
}

func getList(data map[string]any, key string) []any {
    if data == nil {
        return nil
    }
    list, _ := data[key].([]any)
    return list
}
