package enums

// IpOwnersCacheEnum la cache de whois en disco (APP_STATE_PATH/ip-owners.json): un rango de red por entrada, asi una
// consulta sirve para todas las IP del mismo bloque.
var IpOwnersCacheEnum = struct {
    FILE_NAME string
    RANGES    string
    // dias que vale una respuesta buena (como el localhost-traffic de Rust) y horas que se recuerda un fallo
    TTL_DAYS          int
    FAILED_TTL_HOURS  int
}{
    FILE_NAME:        "ip-owners.json",
    RANGES:           "ranges",
    TTL_DAYS:         15,
    FAILED_TTL_HOURS: 1,
}

// IpOwnerKeyEnum claves de una entrada de la cache (y de los primitivos del dueño de una IP).
var IpOwnerKeyEnum = struct {
    START_IP    string
    END_IP      string
    NETWORK     string
    ORG         string
    COUNTRY     string
    SCOPE       string
    IS_FAILED   string
    RESOLVED_AT string
    EXPIRES_AT  string
}{
    START_IP:    "start_ip",
    END_IP:      "end_ip",
    NETWORK:     "network",
    ORG:         "org",
    COUNTRY:     "country",
    SCOPE:       "scope",
    IS_FAILED:   "is_failed",
    RESOLVED_AT: "resolved_at",
    EXPIRES_AT:  "expires_at",
}

// IpScopeEnum de que tipo es una IP. Solo las public se preguntan por whois; el resto se nombran en local.
var IpScopeEnum = struct {
    PUBLIC      string
    PRIVATE     string
    LOOPBACK    string
    LINK_LOCAL  string
    CGNAT       string
    MULTICAST   string
    UNSPECIFIED string
}{
    PUBLIC:      "public",
    PRIVATE:     "private",
    LOOPBACK:    "loopback",
    LINK_LOCAL:  "link-local",
    CGNAT:       "cgnat",
    MULTICAST:   "multicast",
    UNSPECIFIED: "unspecified",
}

// IpOwnersRdapEnum consulta RDAP (el whois moderno, en JSON). rdap.org redirige al registro regional que toca
// (ARIN, RIPE, APNIC, LACNIC, AFRINIC). User-Agent generico: nunca un email ni datos del usuario.
var IpOwnersRdapEnum = struct {
    BASE_URL        string
    USER_AGENT      string
    TIMEOUT_SECONDS int
    // cuantas IP se preguntan como mucho en una pasada y cuantas a la vez (rdap.org limita el ritmo)
    MAX_LOOKUPS_PER_RUN int
    PARALLEL_LOOKUPS    int
}{
    BASE_URL:            "https://rdap.org/ip/",
    USER_AGENT:          "trafficker/",
    TIMEOUT_SECONDS:     8,
    MAX_LOOKUPS_PER_RUN: 40,
    PARALLEL_LOOKUPS:    4,
}
