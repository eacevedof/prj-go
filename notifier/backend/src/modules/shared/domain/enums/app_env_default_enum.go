package enums

// AppEnvDefaultEnum valores que toma la app cuando una clave opcional del .env falta o no es valida.
var AppEnvDefaultEnum = struct {
    QUEUE_WORKER_INTERVAL_SECONDS int
    SENDMAIL_PATH                 string
    MESSAGE_MAX_RECIPIENTS        int
    ATTACHMENT_MAX_BYTES          int
    ATTACHMENT_TIMEOUT_SECONDS    int
}{
    QUEUE_WORKER_INTERVAL_SECONDS: 5,
    SENDMAIL_PATH:                 "/usr/sbin/sendmail -t -i",
    MESSAGE_MAX_RECIPIENTS:        50,
    ATTACHMENT_MAX_BYTES:          10 << 20,
    ATTACHMENT_TIMEOUT_SECONDS:    20,
}
