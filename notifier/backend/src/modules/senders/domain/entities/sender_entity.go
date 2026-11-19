package entities

import (
    "encoding/json"

    "notifier-backend/src/modules/shared/application"

    "notifier-backend/src/modules/senders/domain/enums"
)

const SENDER_UUID_PREFIX = "sender-"

// SenderPrimitives representacion publica: sin id, user_id ni el secreto (solo si esta puesto).
type SenderPrimitives struct {
    SenderUuid       string         `json:"sender_uuid"`
    SenderName       string         `json:"sender_name"`
    Channel          string         `json:"channel"`
    ProviderCode     string         `json:"provider_code"`
    Config           map[string]any `json:"config"`
    HasSecret        bool           `json:"has_secret"`
    IsEnabled        bool           `json:"is_enabled"`
    DailyLimit       int            `json:"daily_limit"`
    MinuteLimit      int            `json:"minute_limit"`
    LastUsedAt       string         `json:"last_used_at"`
    LastErrorAt      string         `json:"last_error_at"`
    LastErrorMessage string         `json:"last_error_message"`
    CreatedAt        string         `json:"created_at"`
    UpdatedAt        string         `json:"updated_at"`
}

type SenderEntity struct {
    id               int
    userId           int
    senderUuid       string
    senderName       string
    channel          string
    providerCode     string
    config           map[string]any
    secretEncrypted  string
    isEnabled        bool
    // 0 = el limite por defecto del proveedor (NULL en BD)
    dailyLimit       int
    minuteLimit      int
    lastUsedAt       string
    lastErrorAt      string
    lastErrorMessage string
    createdAt        string
    updatedAt        string
}

// SenderEntityFromPrimitives desde una fila de app_senders: config llega como texto JSON (de la BD) o ya como map.
func SenderEntityFromPrimitives(primitives map[string]any) *SenderEntity {
    return &SenderEntity{
        id:               application.GetInt(primitives, enums.SendersKeyEnum.ID),
        userId:           application.GetInt(primitives, enums.SendersKeyEnum.USER_ID),
        senderUuid:       application.GetTrimmedString(primitives, enums.SendersKeyEnum.SENDER_UUID),
        senderName:       application.GetTrimmedString(primitives, enums.SendersKeyEnum.SENDER_NAME),
        channel:          application.GetTrimmedString(primitives, enums.SendersKeyEnum.CHANNEL),
        providerCode:     application.GetTrimmedString(primitives, enums.SendersKeyEnum.PROVIDER_CODE),
        config:           getConfig(primitives[enums.SendersKeyEnum.CONFIG]),
        secretEncrypted:  application.GetTrimmedString(primitives, enums.SendersKeyEnum.SECRET_ENCRYPTED),
        isEnabled:        application.GetInt(primitives, enums.SendersKeyEnum.IS_ENABLED) == 1,
        dailyLimit:       application.GetInt(primitives, enums.SendersKeyEnum.DAILY_LIMIT),
        minuteLimit:      application.GetInt(primitives, enums.SendersKeyEnum.MINUTE_LIMIT),
        lastUsedAt:       application.GetTrimmedString(primitives, enums.SendersKeyEnum.LAST_USED_AT),
        lastErrorAt:      application.GetTrimmedString(primitives, enums.SendersKeyEnum.LAST_ERROR_AT),
        lastErrorMessage: application.GetTrimmedString(primitives, enums.SendersKeyEnum.LAST_ERROR_MESSAGE),
        createdAt:        application.GetTrimmedString(primitives, enums.SendersKeyEnum.CREATED_AT),
        updatedAt:        application.GetTrimmedString(primitives, enums.SendersKeyEnum.UPDATED_AT),
    }
}

func getConfig(rawConfig any) map[string]any {
    if config, isMap := rawConfig.(map[string]any); isMap {
        return config
    }
    config := map[string]any{}
    _ = json.Unmarshal([]byte(application.ToString(rawConfig)), &config)
    return config
}

func (e *SenderEntity) ToPrimitives() SenderPrimitives {
    return SenderPrimitives{
        SenderUuid:       e.senderUuid,
        SenderName:       e.senderName,
        Channel:          e.channel,
        ProviderCode:     e.providerCode,
        Config:           e.config,
        HasSecret:        e.secretEncrypted != "",
        IsEnabled:        e.isEnabled,
        DailyLimit:       e.GetEffectiveDailyLimit(),
        MinuteLimit:      e.GetEffectiveMinuteLimit(),
        LastUsedAt:       e.lastUsedAt,
        LastErrorAt:      e.lastErrorAt,
        LastErrorMessage: e.lastErrorMessage,
        CreatedAt:        e.createdAt,
        UpdatedAt:        e.updatedAt,
    }
}

// ToPrimitivesOrEmpty una lectura que no encuentra nada da el perfil vacio.
func ToPrimitivesOrEmpty(senderEntity *SenderEntity) SenderPrimitives {
    if senderEntity == nil {
        return SenderPrimitives{}
    }
    return senderEntity.ToPrimitives()
}

func (e *SenderEntity) GetId() int {
    return e.id
}

func (e *SenderEntity) GetUserId() int {
    return e.userId
}

func (e *SenderEntity) GetSenderUuid() string {
    return e.senderUuid
}

func (e *SenderEntity) GetSenderName() string {
    return e.senderName
}

func (e *SenderEntity) GetChannel() string {
    return e.channel
}

func (e *SenderEntity) GetProviderCode() string {
    return e.providerCode
}

func (e *SenderEntity) GetConfig() map[string]any {
    return e.config
}

// GetConfigJson config serializado para la columna.
func (e *SenderEntity) GetConfigJson() string {
    configJson, err := json.Marshal(e.config)
    if err != nil {
        return "{}"
    }
    return string(configJson)
}

func (e *SenderEntity) GetConfigString(key string) string {
    return application.GetTrimmedString(e.config, key)
}

func (e *SenderEntity) GetConfigInt(key string) int {
    return application.GetInt(e.config, key)
}

func (e *SenderEntity) GetSecretEncrypted() string {
    return e.secretEncrypted
}

func (e *SenderEntity) IsEnabled() bool {
    return e.isEnabled
}

func (e *SenderEntity) GetDailyLimit() int {
    return e.dailyLimit
}

func (e *SenderEntity) GetMinuteLimit() int {
    return e.minuteLimit
}

// GetEffectiveDailyLimit el del perfil o, si no tiene, el del proveedor. 0 = sin limite.
func (e *SenderEntity) GetEffectiveDailyLimit() int {
    if e.dailyLimit > 0 {
        return e.dailyLimit
    }
    return enums.ProviderEnum.DEFINITIONS[e.providerCode].DailyLimit
}

// GetEffectiveMinuteLimit el del perfil o, si no tiene, el del proveedor. 0 = sin limite.
func (e *SenderEntity) GetEffectiveMinuteLimit() int {
    if e.minuteLimit > 0 {
        return e.minuteLimit
    }
    return enums.ProviderEnum.DEFINITIONS[e.providerCode].MinuteLimit
}

func (e *SenderEntity) GetCreatedAt() string {
    return e.createdAt
}
