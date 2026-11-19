package enums

// SenderFieldsLimitEnum limites de los campos de un perfil de envio (SenderFieldsValidator).
var SenderFieldsLimitEnum = struct {
    SENDER_NAME_MIN_LENGTH int
    SENDER_NAME_MAX_LENGTH int
    FROM_EMAIL_MAX_LENGTH  int
    FROM_NAME_MAX_LENGTH   int
    SMTP_HOST_MAX_LENGTH   int
    USERNAME_MAX_LENGTH    int
    SECRET_MAX_LENGTH      int
    LIMIT_MAX              int
    LAST_ERROR_MAX_LENGTH  int
}{
    SENDER_NAME_MIN_LENGTH: 3,
    SENDER_NAME_MAX_LENGTH: 100,
    FROM_EMAIL_MAX_LENGTH:  150,
    FROM_NAME_MAX_LENGTH:   100,
    SMTP_HOST_MAX_LENGTH:   255,
    USERNAME_MAX_LENGTH:    150,
    SECRET_MAX_LENGTH:      500,
    LIMIT_MAX:              100000,
    LAST_ERROR_MAX_LENGTH:  500,
}
