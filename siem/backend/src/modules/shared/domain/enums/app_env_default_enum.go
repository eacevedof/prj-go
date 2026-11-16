package enums

// AppEnvDefaultEnum valores que toma la app cuando una clave opcional del .env falta o no es valida.
var AppEnvDefaultEnum = struct {
    QUEUE_WORKER_INTERVAL_SECONDS int
    EVENTS_RETENTION_DAYS         int
    SENDMAIL_PATH                 string
}{
    QUEUE_WORKER_INTERVAL_SECONDS: 2,
    EVENTS_RETENTION_DAYS:         30,
    SENDMAIL_PATH:                 "/usr/sbin/sendmail -t -i",
}
