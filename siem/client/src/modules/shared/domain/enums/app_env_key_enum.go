package enums

// AppEnvKeyEnum claves del .env que la app consume. Solo las usan EnvironmentReaderRawRepository (un getter
// semantico por clave) y el loader AppbootLoadenv(): el resto del codigo nunca escribe un nombre de variable.
var AppEnvKeyEnum = struct {
    APP_ENV          string
    APP_DEBUG        string
    APP_NAME         string
    APP_DEFAULT_LANG string
    APP_LOG_PATH     string

    // a donde se envia: el backend y el api_token del agente (./console agent-create en el backend)
    APP_SIEM_ENDPOINT_URL string
    APP_SIEM_API_TOKEN    string
    // 1 (defecto) = lo que se vigila tambien baja del panel del backend (GET /api/v1/agent-config); 0 = solo siem-client.json
    APP_SIEM_REMOTE_CONFIG string

    // que se vigila (siem-client.json) y donde se guarda lo pendiente de enviar y lo ya leido
    APP_CONFIG_PATH string
    APP_SPOOL_PATH  string
    APP_STATE_PATH  string
}{
    APP_ENV:          "APP_ENV",
    APP_DEBUG:        "APP_DEBUG",
    APP_NAME:         "APP_NAME",
    APP_DEFAULT_LANG: "APP_DEFAULT_LANG",
    APP_LOG_PATH:     "APP_LOG_PATH",

    APP_SIEM_ENDPOINT_URL: "APP_SIEM_ENDPOINT_URL",
    APP_SIEM_API_TOKEN:    "APP_SIEM_API_TOKEN",
    APP_SIEM_REMOTE_CONFIG: "APP_SIEM_REMOTE_CONFIG",

    APP_CONFIG_PATH: "APP_CONFIG_PATH",
    APP_SPOOL_PATH:  "APP_SPOOL_PATH",
    APP_STATE_PATH:  "APP_STATE_PATH",
}

// AppEnvDefaultEnum valores por defecto cuando ni el OS ni el .env traen la clave (las rutas salen de BootPathsEnum).
var AppEnvDefaultEnum = struct {
    APP_ENV          string
    APP_DEBUG        string
    APP_NAME         string
    APP_DEFAULT_LANG string
    APP_SIEM_REMOTE_CONFIG string
}{
    APP_ENV:          "production",
    APP_DEBUG:        "0",
    APP_NAME:         "siem-client",
    APP_DEFAULT_LANG: "es",
    APP_SIEM_REMOTE_CONFIG: "1",
}
