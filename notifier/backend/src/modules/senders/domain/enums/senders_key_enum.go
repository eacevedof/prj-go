package enums

// SendersKeyEnum claves de la API y de consola de senders, las del JSON de config y las columnas de app_senders.
var SendersKeyEnum = struct {
    SENDER_UUID        string
    SENDER_NAME        string
    CHANNEL            string
    PROVIDER_CODE      string
    CONFIG             string
    HAS_SECRET         string
    IS_ENABLED         string
    DAILY_LIMIT        string
    MINUTE_LIMIT       string
    LAST_USED_AT       string
    LAST_ERROR_AT      string
    LAST_ERROR_MESSAGE string
    CREATED_AT         string
    UPDATED_AT         string

    // config (JSON, lo no secreto)
    FROM_EMAIL    string
    FROM_NAME     string
    SMTP_HOST     string
    SMTP_PORT     string
    SMTP_SECURITY string
    USERNAME      string

    // entrada: el secreto en claro (nunca sale) y el dueno
    SECRET     string
    USER_UUID  string
    TOKEN_NAME string
    SENDER     string
    API_TOKEN  string

    // columnas internas de BD: nunca salen por la API
    ID               string
    USER_ID          string
    SECRET_ENCRYPTED string
}{
    SENDER_UUID:        "sender_uuid",
    SENDER_NAME:        "sender_name",
    CHANNEL:            "channel",
    PROVIDER_CODE:      "provider_code",
    CONFIG:             "config",
    HAS_SECRET:         "has_secret",
    IS_ENABLED:         "is_enabled",
    DAILY_LIMIT:        "daily_limit",
    MINUTE_LIMIT:       "minute_limit",
    LAST_USED_AT:       "last_used_at",
    LAST_ERROR_AT:      "last_error_at",
    LAST_ERROR_MESSAGE: "last_error_message",
    CREATED_AT:         "created_at",
    UPDATED_AT:         "updated_at",

    FROM_EMAIL:    "from_email",
    FROM_NAME:     "from_name",
    SMTP_HOST:     "smtp_host",
    SMTP_PORT:     "smtp_port",
    SMTP_SECURITY: "smtp_security",
    USERNAME:      "username",

    SECRET:     "secret",
    USER_UUID:  "user_uuid",
    TOKEN_NAME: "token_name",
    SENDER:     "sender",
    API_TOKEN:  "api_token",

    ID:               "id",
    USER_ID:          "user_id",
    SECRET_ENCRYPTED: "secret_encrypted",
}
