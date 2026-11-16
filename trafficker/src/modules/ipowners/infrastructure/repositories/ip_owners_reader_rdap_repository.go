package repositories

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/netip"
    "time"

    "trafficker/src/modules/ipowners/domain/enums"
)

// IpOwnersReaderRdapRepository pregunta por RDAP de quien es una IP. rdap.org redirige al registro regional; el
// cliente HTTP de la stdlib sigue la redireccion solo.
type IpOwnersReaderRdapRepository struct {
    httpClient *http.Client
    userAgent  string
}

func IpOwnersReaderRdapRepositoryGetInstance(appVersion string) *IpOwnersReaderRdapRepository {
    return &IpOwnersReaderRdapRepository{
        httpClient: &http.Client{Timeout: time.Duration(enums.IpOwnersRdapEnum.TIMEOUT_SECONDS) * time.Second},
        userAgent:  enums.IpOwnersRdapEnum.USER_AGENT + appVersion,
    }
}

// GetRdapResponse el JSON de la red que contiene la IP. Error si no responde, no es 200 (404 sin datos, 429 por
// ritmo) o no es JSON.
func (r *IpOwnersReaderRdapRepository) GetRdapResponse(addr netip.Addr) (map[string]any, error) {
    request, err := http.NewRequest(http.MethodGet, enums.IpOwnersRdapEnum.BASE_URL+addr.Unmap().String(), nil)
    if err != nil {
        return nil, err
    }
    request.Header.Set("User-Agent", r.userAgent)
    request.Header.Set("Accept", "application/rdap+json, application/json")

    response, err := r.httpClient.Do(request)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()
    if response.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("rdap %s: http %d", addr, response.StatusCode)
    }
    rawJson, err := io.ReadAll(io.LimitReader(response.Body, 2<<20))
    if err != nil {
        return nil, err
    }
    rdapResponse := map[string]any{}
    if err := json.Unmarshal(rawJson, &rdapResponse); err != nil {
        return nil, err
    }
    return rdapResponse, nil
}
