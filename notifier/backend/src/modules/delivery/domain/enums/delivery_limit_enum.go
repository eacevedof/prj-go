package enums

// DeliveryLimitEnum tiempos y tamanos de la entrega.
var DeliveryLimitEnum = struct {
    // espera tras el intento N fallido (N = posicion + 1): 1 min, 5 min, 30 min, 2 h. El quinto fallo ya es failed.
    RETRY_DELAYS_SECONDS []int
    // mensajes por pasada del worker de la API
    MAX_MESSAGES_PER_WORKER_ROUND int
    // ./console message-queue-process --max_messages
    MAX_MESSAGES_DEFAULT int
    MAX_MESSAGES_MIN     int
    MAX_MESSAGES_MAX     int
    // un sending sin tocar desde hace mas de esto es un huerfano (la API murio a medio envio): vuelve a pending
    ORPHAN_SENDING_SECONDS int
    // un perfil que llego a su limite: el mensaje espera sin gastar intento
    MINUTE_LIMIT_POSTPONE_SECONDS int
    DAILY_LIMIT_POSTPONE_SECONDS  int
    // conexion SMTP completa (conectar, TLS, auth, DATA)
    SMTP_TIMEOUT_SECONDS int
    // redirecciones que sigue la descarga de un adjunto (cada una vuelve a pasar por la lista blanca)
    ATTACHMENT_MAX_REDIRECTS int
}{
    RETRY_DELAYS_SECONDS:          []int{60, 300, 1800, 7200},
    MAX_MESSAGES_PER_WORKER_ROUND: 20,
    MAX_MESSAGES_DEFAULT:          50,
    MAX_MESSAGES_MIN:              1,
    MAX_MESSAGES_MAX:              500,
    ORPHAN_SENDING_SECONDS:        600,
    MINUTE_LIMIT_POSTPONE_SECONDS: 60,
    DAILY_LIMIT_POSTPONE_SECONDS:  3600,
    SMTP_TIMEOUT_SECONDS:          60,
    ATTACHMENT_MAX_REDIRECTS:      3,
}
