package enums

// AppEnvKeyEnum claves del .env que la app consume. Solo las usan EnvironmentReaderRawRepository (un getter
// semantico por clave) y el loader AppbootLoadenv(): el resto del codigo nunca escribe un nombre de variable.
var AppEnvKeyEnum = struct {
    APP_ENV          string
    APP_DEBUG        string
    APP_NAME         string
    APP_DEFAULT_LANG string
    APP_LOG_PATH     string

    // a donde se envia: el endpoint de ingesta del backend del siem y el api_token del agente
    APP_SIEM_ENDPOINT_URL string
    APP_SIEM_API_TOKEN    string

    // cada cuanto se hace la foto de los sockets y cada cuanto se envia el spool
    APP_COLLECT_INTERVAL_SECONDS string
    APP_SEND_INTERVAL_SECONDS    string
    // filtros: loopback (127.0.0.0/8, ::1) fuera salvo "1"; procesos y puertos a ignorar separados por comas
    APP_INCLUDE_LOOPBACK string
    APP_IGNORE_PROCESSES string
    APP_IGNORE_PORTS     string

    // donde se guarda lo pendiente de enviar (spool) y la ultima foto de los sockets (state)
    APP_SPOOL_PATH string
    APP_STATE_PATH string
    // raiz de /proc: solo cambia en los tests o al leer el /proc del host montado en un contenedor
    APP_PROC_PATH string
}{
    APP_ENV:          "APP_ENV",
    APP_DEBUG:        "APP_DEBUG",
    APP_NAME:         "APP_NAME",
    APP_DEFAULT_LANG: "APP_DEFAULT_LANG",
    APP_LOG_PATH:     "APP_LOG_PATH",

    APP_SIEM_ENDPOINT_URL: "APP_SIEM_ENDPOINT_URL",
    APP_SIEM_API_TOKEN:    "APP_SIEM_API_TOKEN",

    APP_COLLECT_INTERVAL_SECONDS: "APP_COLLECT_INTERVAL_SECONDS",
    APP_SEND_INTERVAL_SECONDS:    "APP_SEND_INTERVAL_SECONDS",
    APP_INCLUDE_LOOPBACK:         "APP_INCLUDE_LOOPBACK",
    APP_IGNORE_PROCESSES:         "APP_IGNORE_PROCESSES",
    APP_IGNORE_PORTS:             "APP_IGNORE_PORTS",

    APP_SPOOL_PATH: "APP_SPOOL_PATH",
    APP_STATE_PATH: "APP_STATE_PATH",
    APP_PROC_PATH:  "APP_PROC_PATH",
}

// AppEnvDefaultEnum valores por defecto cuando ni el OS ni el .env traen la clave (las rutas salen de BootPathsEnum).
var AppEnvDefaultEnum = struct {
    APP_ENV                      string
    APP_DEBUG                    string
    APP_NAME                     string
    APP_DEFAULT_LANG             string
    APP_COLLECT_INTERVAL_SECONDS string
    APP_SEND_INTERVAL_SECONDS    string
    APP_INCLUDE_LOOPBACK         string
    APP_PROC_PATH                string
}{
    APP_ENV:                      "production",
    APP_DEBUG:                    "0",
    APP_NAME:                     "trafficker",
    APP_DEFAULT_LANG:             "es",
    APP_COLLECT_INTERVAL_SECONDS: "5",
    APP_SEND_INTERVAL_SECONDS:    "30",
    APP_INCLUDE_LOOPBACK:         "0",
    APP_PROC_PATH:                "/proc",
}
