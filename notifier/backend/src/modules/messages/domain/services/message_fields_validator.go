package services

import (
    "slices"
    "strings"
    "unicode/utf8"

    "notifier-backend/src/core/config"

    sharedservices "notifier-backend/src/modules/shared/domain/services"

    "notifier-backend/src/modules/messages/domain/enums"
)

// EmailMessageFields lo que se valida de un mensaje de email, ya normalizado (MessagePayloadNormalizer).
type EmailMessageFields struct {
    To             []string
    Cc             []string
    Bcc            []string
    ReplyTo        string
    Subject        string
    Body           string
    BodyFormat     string
    IdempotencyKey string
    Attachments    []AttachmentInput
}

// MessageFieldsLimits los limites que vienen del .env.
type MessageFieldsLimits struct {
    MaxRecipients      int
    AttachmentMaxBytes int
}

// MessageFieldsValidator reglas de un mensaje de email. Devuelve un error "campo: mensaje" por regla incumplida.
type MessageFieldsValidator struct {
    emailAddressValidator *sharedservices.EmailAddressValidator
    attachmentUrlChecker  *AttachmentUrlChecker
}

func MessageFieldsValidatorGetInstance(attachmentUrlChecker *AttachmentUrlChecker) *MessageFieldsValidator {
    return &MessageFieldsValidator{
        emailAddressValidator: sharedservices.EmailAddressValidatorGetInstance(),
        attachmentUrlChecker:  attachmentUrlChecker,
    }
}

func (v *MessageFieldsValidator) GetEmailErrors(messageFields EmailMessageFields, limits MessageFieldsLimits) []string {
    keys := enums.MessagesKeyEnum
    errors := []string{}

    if len(messageFields.To) == 0 {
        errors = append(errors, v.getError(keys.TO, enums.MessagesMessageEnum.FIELD_IS_REQUIRED))
    }
    if len(messageFields.To)+len(messageFields.Cc)+len(messageFields.Bcc) > limits.MaxRecipients {
        errors = append(errors, v.getError(keys.TO, enums.MessagesMessageEnum.MAX_X_RECIPIENTS, limits.MaxRecipients))
    }
    // en este orden: el mensaje de error sale siempre igual
    emailLists := []struct {
        fieldName string
        emails    []string
    }{{keys.TO, messageFields.To}, {keys.CC, messageFields.Cc}, {keys.BCC, messageFields.Bcc}}
    for _, emailList := range emailLists {
        if invalidEmails := v.emailAddressValidator.GetInvalidOnes(emailList.emails); len(invalidEmails) > 0 {
            errors = append(errors, v.getError(emailList.fieldName, enums.MessagesMessageEnum.INVALID_EMAILS_X, strings.Join(invalidEmails, ", ")))
        }
    }
    if messageFields.ReplyTo != "" && !v.emailAddressValidator.IsValid(messageFields.ReplyTo) {
        errors = append(errors, v.getError(keys.REPLY_TO, enums.MessagesMessageEnum.INVALID_EMAILS_X, messageFields.ReplyTo))
    }

    switch {
    case messageFields.Subject == "":
        errors = append(errors, v.getError(keys.SUBJECT, enums.MessagesMessageEnum.FIELD_IS_REQUIRED))
    case utf8.RuneCountInString(messageFields.Subject) > enums.MessagesLimitEnum.SUBJECT_MAX_LENGTH:
        errors = append(errors, v.getError(keys.SUBJECT, enums.MessagesMessageEnum.MAX_LENGTH_IS_X, enums.MessagesLimitEnum.SUBJECT_MAX_LENGTH))
    case strings.ContainsAny(messageFields.Subject, "\r\n"):
        errors = append(errors, v.getError(keys.SUBJECT, enums.MessagesMessageEnum.NO_LINE_BREAKS))
    }

    if len(messageFields.Body) > enums.MessagesLimitEnum.BODY_MAX_BYTES {
        errors = append(errors, v.getError(keys.BODY, enums.MessagesMessageEnum.MAX_LENGTH_IS_X, enums.MessagesLimitEnum.BODY_MAX_BYTES))
    }
    if !slices.Contains(enums.BodyFormatEnum.VALUES, messageFields.BodyFormat) {
        errors = append(errors, v.getError(keys.BODY_FORMAT, enums.MessagesMessageEnum.MUST_BE_X, strings.Join(enums.BodyFormatEnum.VALUES, ", ")))
    }
    if utf8.RuneCountInString(messageFields.IdempotencyKey) > enums.MessagesLimitEnum.IDEMPOTENCY_KEY_MAX_LENGTH {
        errors = append(errors, v.getError(keys.IDEMPOTENCY_KEY, enums.MessagesMessageEnum.MAX_LENGTH_IS_X, enums.MessagesLimitEnum.IDEMPOTENCY_KEY_MAX_LENGTH))
    }
    return append(errors, v.getAttachmentsErrors(messageFields.Attachments, limits.AttachmentMaxBytes)...)
}

func (v *MessageFieldsValidator) getAttachmentsErrors(attachments []AttachmentInput, attachmentMaxBytes int) []string {
    keys := enums.MessagesKeyEnum
    if len(attachments) > enums.MessagesLimitEnum.MAX_ATTACHMENTS {
        return []string{v.getError(keys.ATTACHMENTS, enums.MessagesMessageEnum.MAX_X_ATTACHMENTS, enums.MessagesLimitEnum.MAX_ATTACHMENTS)}
    }

    errors := []string{}
    for i, attachment := range attachments {
        position := i + 1
        if attachment.HasContent == attachment.HasUrl {
            errors = append(errors, v.getError(keys.ATTACHMENTS, enums.MessagesMessageEnum.ATTACHMENT_X_NEEDS_CONTENT_OR_URL, position))
            continue
        }
        if !v.isFileNameValid(attachment.FileName) {
            errors = append(errors, v.getError(keys.ATTACHMENTS, enums.MessagesMessageEnum.ATTACHMENT_X_INVALID_FILE_NAME, position))
        }
        if attachment.HasContent && len(attachment.Content) > attachmentMaxBytes {
            errors = append(errors, v.getError(keys.ATTACHMENTS, enums.MessagesMessageEnum.ATTACHMENT_X_TOO_BIG, position, attachmentMaxBytes))
        }
        if attachment.HasUrl && (len(attachment.Url) > enums.MessagesLimitEnum.SOURCE_URL_MAX_LENGTH || !v.attachmentUrlChecker.IsAllowed(attachment.Url)) {
            errors = append(errors, v.getError(keys.ATTACHMENTS, enums.MessagesMessageEnum.ATTACHMENT_X_URL_NOT_ALLOWED, position))
        }
    }
    return errors
}

// isFileNameValid el nombre va en la cabecera del adjunto: ni rutas ni comillas ni saltos de linea.
func (v *MessageFieldsValidator) isFileNameValid(fileName string) bool {
    return fileName != "" &&
        utf8.RuneCountInString(fileName) <= enums.MessagesLimitEnum.FILE_NAME_MAX_LENGTH &&
        !strings.ContainsAny(fileName, "/\\\"\r\n\x00")
}

func (v *MessageFieldsValidator) getError(fieldName string, message string, args ...any) string {
    return fieldName + ": " + config.Tr(message, args...)
}
