package entities

import (
    "notifier-backend/src/modules/shared/application"

    "notifier-backend/src/modules/messages/domain/enums"
)

const MESSAGE_UUID_PREFIX = "message-"

// MessagePrimitives lo que ve el cliente de un mensaje: su estado, sin el cuerpo ni los ids internos.
type MessagePrimitives struct {
    MessageUuid       string                        `json:"message_uuid"`
    Channel           string                        `json:"channel"`
    MessageStatus     string                        `json:"message_status"`
    To                []string                      `json:"to"`
    Cc                []string                      `json:"cc"`
    Bcc               []string                      `json:"bcc"`
    ReplyTo           string                        `json:"reply_to"`
    Subject           string                        `json:"subject"`
    BodyFormat        string                        `json:"body_format"`
    Attempts          int                           `json:"attempts"`
    MaxAttempts       int                           `json:"max_attempts"`
    NextAttemptAt     string                        `json:"next_attempt_at"`
    ErrorMessage      string                        `json:"error_message"`
    ProviderMessageId string                        `json:"provider_message_id"`
    CreatedAt         string                        `json:"created_at"`
    SentAt            string                        `json:"sent_at"`
    Attachments       []MessageAttachmentPrimitives `json:"attachments"`
    AttemptsLog       []MessageAttemptPrimitives    `json:"attempts_log"`
}

type MessageEntity struct {
    id                int
    messageUuid       string
    apiTokenId        int
    senderId          int
    channel           string
    recipients        Recipients
    subject           string
    body              string
    bodyFormat        string
    messageStatus     string
    attempts          int
    maxAttempts       int
    nextAttemptAt     string
    errorMessage      string
    providerMessageId string
    idempotencyKey    string
    requestIp         string
    createdAt         string
    sentAt            string
    updatedAt         string
}

// MessageEntityFromPrimitives desde una fila de app_messages: recipients llega como texto JSON o como Recipients.
func MessageEntityFromPrimitives(primitives map[string]any) *MessageEntity {
    keys := enums.MessagesKeyEnum
    recipients, isRecipients := primitives[keys.RECIPIENTS].(Recipients)
    if !isRecipients {
        recipients = RecipientsFromJson(application.ToString(primitives[keys.RECIPIENTS]))
    }
    return &MessageEntity{
        id:                application.GetInt(primitives, keys.ID),
        messageUuid:       application.GetTrimmedString(primitives, keys.MESSAGE_UUID),
        apiTokenId:        application.GetInt(primitives, keys.API_TOKEN_ID),
        senderId:          application.GetInt(primitives, keys.SENDER_ID),
        channel:           application.GetTrimmedString(primitives, keys.CHANNEL),
        recipients:        recipients,
        subject:           application.GetTrimmedString(primitives, keys.SUBJECT),
        body:              application.ToString(primitives[keys.BODY]),
        bodyFormat:        application.GetTrimmedString(primitives, keys.BODY_FORMAT),
        messageStatus:     application.GetTrimmedString(primitives, keys.MESSAGE_STATUS),
        attempts:          application.GetInt(primitives, keys.ATTEMPTS),
        maxAttempts:       application.GetIntOrDefault(primitives, keys.MAX_ATTEMPTS, enums.MessagesLimitEnum.MAX_ATTEMPTS),
        nextAttemptAt:     application.GetTrimmedString(primitives, keys.NEXT_ATTEMPT_AT),
        errorMessage:      application.GetTrimmedString(primitives, keys.ERROR_MESSAGE),
        providerMessageId: application.GetTrimmedString(primitives, keys.PROVIDER_MESSAGE_ID),
        idempotencyKey:    application.GetTrimmedString(primitives, keys.IDEMPOTENCY_KEY),
        requestIp:         application.GetTrimmedString(primitives, keys.REQUEST_IP),
        createdAt:         application.GetTrimmedString(primitives, keys.CREATED_AT),
        sentAt:            application.GetTrimmedString(primitives, keys.SENT_AT),
        updatedAt:         application.GetTrimmedString(primitives, keys.UPDATED_AT),
    }
}

// ToPrimitives con sus adjuntos y sus intentos (las dos listas pueden ir vacias).
func (e *MessageEntity) ToPrimitives(attachments []MessageAttachmentPrimitives, attemptsLog []MessageAttemptPrimitives) MessagePrimitives {
    if attachments == nil {
        attachments = []MessageAttachmentPrimitives{}
    }
    if attemptsLog == nil {
        attemptsLog = []MessageAttemptPrimitives{}
    }
    recipients := e.recipients.getWithoutNilLists()
    return MessagePrimitives{
        MessageUuid:       e.messageUuid,
        Channel:           e.channel,
        MessageStatus:     e.messageStatus,
        To:                recipients.To,
        Cc:                recipients.Cc,
        Bcc:               recipients.Bcc,
        ReplyTo:           recipients.ReplyTo,
        Subject:           e.subject,
        BodyFormat:        e.bodyFormat,
        Attempts:          e.attempts,
        MaxAttempts:       e.maxAttempts,
        NextAttemptAt:     e.nextAttemptAt,
        ErrorMessage:      e.errorMessage,
        ProviderMessageId: e.providerMessageId,
        CreatedAt:         e.createdAt,
        SentAt:            e.sentAt,
        Attachments:       attachments,
        AttemptsLog:       attemptsLog,
    }
}

func (e *MessageEntity) GetId() int {
    return e.id
}

func (e *MessageEntity) GetMessageUuid() string {
    return e.messageUuid
}

func (e *MessageEntity) GetApiTokenId() int {
    return e.apiTokenId
}

func (e *MessageEntity) GetSenderId() int {
    return e.senderId
}

func (e *MessageEntity) GetChannel() string {
    return e.channel
}

func (e *MessageEntity) GetRecipients() Recipients {
    return e.recipients
}

func (e *MessageEntity) GetSubject() string {
    return e.subject
}

func (e *MessageEntity) GetBody() string {
    return e.body
}

func (e *MessageEntity) GetBodyFormat() string {
    return e.bodyFormat
}

func (e *MessageEntity) GetMessageStatus() string {
    return e.messageStatus
}

func (e *MessageEntity) GetAttempts() int {
    return e.attempts
}

func (e *MessageEntity) GetMaxAttempts() int {
    return e.maxAttempts
}

func (e *MessageEntity) GetNextAttemptAt() string {
    return e.nextAttemptAt
}

func (e *MessageEntity) GetIdempotencyKey() string {
    return e.idempotencyKey
}

func (e *MessageEntity) GetRequestIp() string {
    return e.requestIp
}

func (e *MessageEntity) GetCreatedAt() string {
    return e.createdAt
}
