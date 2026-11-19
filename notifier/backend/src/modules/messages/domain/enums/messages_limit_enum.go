package enums

// MessagesLimitEnum limites de un mensaje. El tope de destinatarios y el de tamano de cada adjunto salen del .env
// (APP_MESSAGE_MAX_RECIPIENTS, APP_ATTACHMENT_MAX_BYTES).
var MessagesLimitEnum = struct {
    SUBJECT_MAX_LENGTH         int
    BODY_MAX_BYTES             int
    MAX_ATTACHMENTS            int
    FILE_NAME_MAX_LENGTH       int
    SOURCE_URL_MAX_LENGTH      int
    IDEMPOTENCY_KEY_MAX_LENGTH int
    REQUEST_IP_MAX_LENGTH      int
    ERROR_MESSAGE_MAX_LENGTH   int
    PROVIDER_RESPONSE_MAX_LENGTH int
    // reintentos: 5 intentos con espera creciente
    MAX_ATTEMPTS int
    // un adjunto en texto sin nombre
    DEFAULT_TEXT_FILE_NAME string
}{
    SUBJECT_MAX_LENGTH:           255,
    BODY_MAX_BYTES:               1 << 20,
    MAX_ATTACHMENTS:              10,
    FILE_NAME_MAX_LENGTH:         255,
    SOURCE_URL_MAX_LENGTH:        1000,
    IDEMPOTENCY_KEY_MAX_LENGTH:   100,
    REQUEST_IP_MAX_LENGTH:        45,
    ERROR_MESSAGE_MAX_LENGTH:     500,
    PROVIDER_RESPONSE_MAX_LENGTH: 1000,
    MAX_ATTEMPTS:                 5,
    DEFAULT_TEXT_FILE_NAME:       "attachment.txt",
}
