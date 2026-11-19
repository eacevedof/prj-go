package entities

import (
    "notifier-backend/src/modules/shared/application"

    "notifier-backend/src/modules/messages/domain/enums"
)

// MessageAttemptPrimitives un intento de entrega tal cual se ensena.
type MessageAttemptPrimitives struct {
    AttemptNumber    int    `json:"attempt_number"`
    AttemptStatus    string `json:"attempt_status"`
    ProviderResponse string `json:"provider_response"`
    DurationMs       int    `json:"duration_ms"`
    AttemptedAt      string `json:"attempted_at"`
}

type MessageAttemptEntity struct {
    messageId        int
    attemptNumber    int
    attemptStatus    string
    providerResponse string
    durationMs       int
    attemptedAt      string
}

// MessageAttemptEntityFromPrimitives desde una fila de app_message_attempts.
func MessageAttemptEntityFromPrimitives(primitives map[string]any) *MessageAttemptEntity {
    keys := enums.MessagesKeyEnum
    return &MessageAttemptEntity{
        messageId:        application.GetInt(primitives, keys.MESSAGE_ID),
        attemptNumber:    application.GetInt(primitives, keys.ATTEMPT_NUMBER),
        attemptStatus:    application.GetTrimmedString(primitives, keys.ATTEMPT_STATUS),
        providerResponse: application.GetTrimmedString(primitives, keys.PROVIDER_RESPONSE),
        durationMs:       application.GetInt(primitives, keys.DURATION_MS),
        attemptedAt:      application.GetTrimmedString(primitives, keys.ATTEMPTED_AT),
    }
}

func (e *MessageAttemptEntity) ToPrimitives() MessageAttemptPrimitives {
    return MessageAttemptPrimitives{
        AttemptNumber:    e.attemptNumber,
        AttemptStatus:    e.attemptStatus,
        ProviderResponse: e.providerResponse,
        DurationMs:       e.durationMs,
        AttemptedAt:      e.attemptedAt,
    }
}

func (e *MessageAttemptEntity) GetMessageId() int {
    return e.messageId
}

func (e *MessageAttemptEntity) GetAttemptNumber() int {
    return e.attemptNumber
}

func (e *MessageAttemptEntity) GetAttemptStatus() string {
    return e.attemptStatus
}

func (e *MessageAttemptEntity) GetProviderResponse() string {
    return e.providerResponse
}

func (e *MessageAttemptEntity) GetDurationMs() int {
    return e.durationMs
}

func (e *MessageAttemptEntity) GetAttemptedAt() string {
    return e.attemptedAt
}
