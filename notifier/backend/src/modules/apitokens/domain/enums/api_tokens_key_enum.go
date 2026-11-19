package enums

// ApiTokensKeyEnum claves de la API y de consola de apitokens, mas las columnas de app_api_tokens.
var ApiTokensKeyEnum = struct {
    TOKEN_UUID   string
    TOKEN_NAME   string
    TOKEN_PREFIX string
    DEFAULT_TO   string
    IS_ENABLED   string
    LAST_USED_AT string
    CREATED_AT   string
    UPDATED_AT   string

    API_TOKEN   string
    SENDER_UUID string
    TOKEN       string

    // columnas internas de BD: nunca salen por la API
    ID             string
    SENDER_ID      string
    API_TOKEN_HASH string

    // de la autenticacion de un token de envio (join con app_senders y app_users)
    CHANNEL        string
    PROVIDER_CODE  string
    USER_ID        string
    SENDER_ENABLED string
    USER_ENABLED   string
}{
    TOKEN_UUID:   "token_uuid",
    TOKEN_NAME:   "token_name",
    TOKEN_PREFIX: "token_prefix",
    DEFAULT_TO:   "default_to",
    IS_ENABLED:   "is_enabled",
    LAST_USED_AT: "last_used_at",
    CREATED_AT:   "created_at",
    UPDATED_AT:   "updated_at",

    API_TOKEN:   "api_token",
    SENDER_UUID: "sender_uuid",
    TOKEN:       "token",

    ID:             "id",
    SENDER_ID:      "sender_id",
    API_TOKEN_HASH: "api_token_hash",

    CHANNEL:        "channel",
    PROVIDER_CODE:  "provider_code",
    USER_ID:        "user_id",
    SENDER_ENABLED: "sender_enabled",
    USER_ENABLED:   "user_enabled",
}
