package enums

// AppEnvKeyEnum claves del .env que la app consume (las que nadie lee no se mapean). Solo las usan
// EnvironmentReaderRawRepository (un getter semantico por clave) y el loader AppbootLoadenv():
// el resto del codigo nunca escribe un nombre de variable de entorno.
var AppEnvKeyEnum = struct {
    APP_ENV           string
    APP_DEBUG         string
    APP_DEFAULT_LANG  string
    APP_DOMAIN        string
    APP_SQLITE_PATH   string
    APP_EMAIL_TO      string
    APP_EMAIL_FROM    string
    APP_SENDMAIL_PATH string
    APP_LOG_PATH      string
    APP_NAME          string
    APP_BASE_URL      string

    APP_DASHBOARD_TOKEN               string
    APP_QUEUE_WORKER_INTERVAL_SECONDS string
    APP_EVENTS_RETENTION_DAYS         string
}{
    APP_ENV:           "APP_ENV",
    APP_DEBUG:         "APP_DEBUG",
    APP_DEFAULT_LANG:  "APP_DEFAULT_LANG",
    APP_DOMAIN:        "APP_DOMAIN",
    APP_SQLITE_PATH:   "APP_SQLITE_PATH",
    APP_EMAIL_TO:      "APP_EMAIL_TO",
    APP_EMAIL_FROM:    "APP_EMAIL_FROM",
    APP_SENDMAIL_PATH: "APP_SENDMAIL_PATH",
    APP_LOG_PATH:      "APP_LOG_PATH",
    APP_NAME:          "APP_NAME",
    APP_BASE_URL:      "APP_BASE_URL",

    APP_DASHBOARD_TOKEN:               "APP_DASHBOARD_TOKEN",
    APP_QUEUE_WORKER_INTERVAL_SECONDS: "APP_QUEUE_WORKER_INTERVAL_SECONDS",
    APP_EVENTS_RETENTION_DAYS:         "APP_EVENTS_RETENTION_DAYS",
}
