package enums

// UsersKeyEnum claves de la API y de consola de users, mas las columnas con las que se monta UserEntity.
var UsersKeyEnum = struct {
    USER_UUID  string
    USER_NAME  string
    USER_EMAIL string
    USER_ROLE  string
    IS_ENABLED string
    CREATED_AT string
    UPDATED_AT string

    PASSWORD  string
    API_TOKEN string
    USER      string

    // columnas internas de BD: nunca salen por la API
    ID             string
    PASSWORD_HASH  string
    API_TOKEN_HASH string
}{
    USER_UUID:  "user_uuid",
    USER_NAME:  "user_name",
    USER_EMAIL: "user_email",
    USER_ROLE:  "user_role",
    IS_ENABLED: "is_enabled",
    CREATED_AT: "created_at",
    UPDATED_AT: "updated_at",

    PASSWORD:  "password",
    API_TOKEN: "api_token",
    USER:      "user",

    ID:             "id",
    PASSWORD_HASH:  "password_hash",
    API_TOKEN_HASH: "api_token_hash",
}
