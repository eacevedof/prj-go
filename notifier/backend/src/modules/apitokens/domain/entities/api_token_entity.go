package entities

import (
    "encoding/json"

    "notifier-backend/src/modules/shared/application"

    "notifier-backend/src/modules/apitokens/domain/enums"
)

const TOKEN_UUID_PREFIX = "token-"

// ApiTokenPrimitives representacion publica: nunca el hash, solo el prefijo visible.
type ApiTokenPrimitives struct {
    TokenUuid   string   `json:"token_uuid"`
    TokenName   string   `json:"token_name"`
    TokenPrefix string   `json:"token_prefix"`
    DefaultTo   []string `json:"default_to"`
    IsEnabled   bool     `json:"is_enabled"`
    LastUsedAt  string   `json:"last_used_at"`
    CreatedAt   string   `json:"created_at"`
}

type ApiTokenEntity struct {
    id           int
    senderId     int
    tokenUuid    string
    tokenName    string
    apiTokenHash string
    tokenPrefix  string
    defaultTo    []string
    isEnabled    bool
    lastUsedAt   string
    createdAt    string
}

// ApiTokenEntityFromPrimitives desde una fila de app_api_tokens: default_to llega como texto JSON o como []string.
func ApiTokenEntityFromPrimitives(primitives map[string]any) *ApiTokenEntity {
    return &ApiTokenEntity{
        id:           application.GetInt(primitives, enums.ApiTokensKeyEnum.ID),
        senderId:     application.GetInt(primitives, enums.ApiTokensKeyEnum.SENDER_ID),
        tokenUuid:    application.GetTrimmedString(primitives, enums.ApiTokensKeyEnum.TOKEN_UUID),
        tokenName:    application.GetTrimmedString(primitives, enums.ApiTokensKeyEnum.TOKEN_NAME),
        apiTokenHash: application.GetTrimmedString(primitives, enums.ApiTokensKeyEnum.API_TOKEN_HASH),
        tokenPrefix:  application.GetTrimmedString(primitives, enums.ApiTokensKeyEnum.TOKEN_PREFIX),
        defaultTo:    GetDefaultTo(primitives[enums.ApiTokensKeyEnum.DEFAULT_TO]),
        isEnabled:    application.GetInt(primitives, enums.ApiTokensKeyEnum.IS_ENABLED) == 1,
        lastUsedAt:   application.GetTrimmedString(primitives, enums.ApiTokensKeyEnum.LAST_USED_AT),
        createdAt:    application.GetTrimmedString(primitives, enums.ApiTokensKeyEnum.CREATED_AT),
    }
}

// GetDefaultTo la columna default_to (JSON) o una lista ya montada.
func GetDefaultTo(rawDefaultTo any) []string {
    if defaultTo, isList := rawDefaultTo.([]string); isList {
        return defaultTo
    }
    defaultTo := []string{}
    _ = json.Unmarshal([]byte(application.ToString(rawDefaultTo)), &defaultTo)
    return defaultTo
}

func (e *ApiTokenEntity) ToPrimitives() ApiTokenPrimitives {
    return ApiTokenPrimitives{
        TokenUuid:   e.tokenUuid,
        TokenName:   e.tokenName,
        TokenPrefix: e.tokenPrefix,
        DefaultTo:   e.defaultTo,
        IsEnabled:   e.isEnabled,
        LastUsedAt:  e.lastUsedAt,
        CreatedAt:   e.createdAt,
    }
}

// ToPrimitivesOrEmpty una lectura que no encuentra nada da el token vacio.
func ToPrimitivesOrEmpty(apiTokenEntity *ApiTokenEntity) ApiTokenPrimitives {
    if apiTokenEntity == nil {
        return ApiTokenPrimitives{DefaultTo: []string{}}
    }
    return apiTokenEntity.ToPrimitives()
}

func (e *ApiTokenEntity) GetId() int {
    return e.id
}

func (e *ApiTokenEntity) GetSenderId() int {
    return e.senderId
}

func (e *ApiTokenEntity) GetTokenUuid() string {
    return e.tokenUuid
}

func (e *ApiTokenEntity) GetTokenName() string {
    return e.tokenName
}

func (e *ApiTokenEntity) GetApiTokenHash() string {
    return e.apiTokenHash
}

func (e *ApiTokenEntity) GetTokenPrefix() string {
    return e.tokenPrefix
}

func (e *ApiTokenEntity) GetDefaultTo() []string {
    return e.defaultTo
}

// GetDefaultToJson para la columna: NULL si no hay destinatarios por defecto.
func (e *ApiTokenEntity) GetDefaultToJson() any {
    if len(e.defaultTo) == 0 {
        return nil
    }
    defaultToJson, _ := json.Marshal(e.defaultTo)
    return string(defaultToJson)
}

func (e *ApiTokenEntity) IsEnabled() bool {
    return e.isEnabled
}

func (e *ApiTokenEntity) GetCreatedAt() string {
    return e.createdAt
}
