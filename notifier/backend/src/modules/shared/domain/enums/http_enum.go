package enums

// HttpEnum cabeceras, metodos, tipos y claves del envelope REST {code, status, errors, data}: sin literales sueltos en
// el transporte. Los metodos van en minusculas: asi los compara el router (Route.Allowed).
var HttpEnum = struct {
    HEADER_CONTENT_TYPE     string
    HEADER_WWW_AUTHENTICATE string
    HEADER_AUTHORIZATION    string
    HEADER_USER_AGENT       string
    HEADER_CLIENT_IP        string
    HEADER_X_FORWARDED_FOR  string
    HEADER_X_TIME_ZONE      string
    HEADER_ACCEPT_LANGUAGE  string
    HEADER_ORIGIN           string
    HEADER_ALLOW            string
    HEADER_CACHE_CONTROL    string

    HEADER_ACCESS_CONTROL_ALLOW_ORIGIN  string
    HEADER_ACCESS_CONTROL_MAX_AGE       string
    HEADER_ACCESS_CONTROL_ALLOW_METHODS string
    HEADER_ACCESS_CONTROL_ALLOW_HEADERS string
    HEADER_X_CONTENT_TYPE_OPTIONS       string
    HEADER_REFERRER_POLICY              string

    CONTENT_TYPE_JSON                string
    CONTENT_TYPE_HTML                string
    CONTENT_TYPE_FORM_URLENCODED     string
    CONTENT_TYPE_MULTIPART_FORM_DATA string

    AUTH_SCHEME_BEARER string

    METHOD_GET     string
    METHOD_HEAD    string
    METHOD_POST    string
    METHOD_PUT     string
    METHOD_DELETE  string
    METHOD_OPTIONS string

    DEFAULT_CLIENT_IP string

    KEY_CODE   string
    KEY_STATUS string
    KEY_ERRORS string
    KEY_DATA   string
}{
    HEADER_CONTENT_TYPE:     "Content-Type",
    HEADER_WWW_AUTHENTICATE: "WWW-Authenticate",
    HEADER_AUTHORIZATION:    "Authorization",
    HEADER_USER_AGENT:       "User-Agent",
    HEADER_CLIENT_IP:        "Client-Ip",
    HEADER_X_FORWARDED_FOR:  "X-Forwarded-For",
    HEADER_X_TIME_ZONE:      "X-Time-Zone",
    HEADER_ACCEPT_LANGUAGE:  "Accept-Language",
    HEADER_ORIGIN:           "Origin",
    HEADER_ALLOW:            "Allow",
    HEADER_CACHE_CONTROL:    "Cache-Control",

    HEADER_ACCESS_CONTROL_ALLOW_ORIGIN:  "Access-Control-Allow-Origin",
    HEADER_ACCESS_CONTROL_MAX_AGE:       "Access-Control-Max-Age",
    HEADER_ACCESS_CONTROL_ALLOW_METHODS: "Access-Control-Allow-Methods",
    HEADER_ACCESS_CONTROL_ALLOW_HEADERS: "Access-Control-Allow-Headers",
    HEADER_X_CONTENT_TYPE_OPTIONS:       "X-Content-Type-Options",
    HEADER_REFERRER_POLICY:              "Referrer-Policy",

    CONTENT_TYPE_JSON:                "application/json",
    CONTENT_TYPE_HTML:                "text/html; charset=utf-8",
    CONTENT_TYPE_FORM_URLENCODED:     "application/x-www-form-urlencoded",
    CONTENT_TYPE_MULTIPART_FORM_DATA: "multipart/form-data",

    AUTH_SCHEME_BEARER: "Bearer",

    METHOD_GET:     "get",
    METHOD_HEAD:    "head",
    METHOD_POST:    "post",
    METHOD_PUT:     "put",
    METHOD_DELETE:  "delete",
    METHOD_OPTIONS: "options",

    DEFAULT_CLIENT_IP: "127.0.0.1",

    KEY_CODE:   "code",
    KEY_STATUS: "status",
    KEY_ERRORS: "errors",
    KEY_DATA:   "data",
}
