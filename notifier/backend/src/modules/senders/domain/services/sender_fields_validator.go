package services

import (
    "regexp"
    "slices"
    "strings"
    "unicode/utf8"

    "notifier-backend/src/core/config"

    sharedservices "notifier-backend/src/modules/shared/domain/services"

    "notifier-backend/src/modules/senders/domain/enums"
)

var (
    // un nombre de host (o una ip): letras, digitos, puntos y guiones
    smtpHostPattern = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9.-]*[A-Za-z0-9])?$`)
)

// SenderFields lo que se valida de un perfil de envio de email (canal email). Ya con los valores por defecto del
// proveedor aplicados (SenderConfigBuilder).
type SenderFields struct {
    SenderName   string
    ProviderCode string
    FromEmail    string
    FromName     string
    SmtpHost     string
    SmtpPort     int
    SmtpSecurity string
    Username     string
    Secret       string
    DailyLimit   int
    MinuteLimit  int
}

// SenderFieldsValidator reglas de un perfil de envio. Devuelve un error "campo: mensaje" por regla incumplida.
type SenderFieldsValidator struct{}

func SenderFieldsValidatorGetInstance() *SenderFieldsValidator {
    return &SenderFieldsValidator{}
}

// GetProviderErrors el proveedor va primero: si no existe o no esta disponible, el resto no se puede validar.
func (v *SenderFieldsValidator) GetProviderErrors(providerCode string) []string {
    if providerCode == "" {
        return []string{v.getError(enums.SendersKeyEnum.PROVIDER_CODE, enums.SendersMessageEnum.FIELD_IS_REQUIRED)}
    }
    providerDefinition, isKnown := enums.ProviderEnum.DEFINITIONS[providerCode]
    if !isKnown {
        return []string{v.getError(enums.SendersKeyEnum.PROVIDER_CODE, enums.SendersMessageEnum.PROVIDER_X_UNKNOWN, providerCode)}
    }
    if !providerDefinition.IsImplemented {
        return []string{v.getError(enums.SendersKeyEnum.PROVIDER_CODE, enums.SendersMessageEnum.PROVIDER_X_NOT_IMPLEMENTED_YET, providerCode)}
    }
    return nil
}

func (v *SenderFieldsValidator) GetEmailSenderErrors(senderFields SenderFields) []string {
    errors := []string{}
    limits := enums.SenderFieldsLimitEnum
    keys := enums.SendersKeyEnum

    errors = append(errors, v.getLengthErrors(keys.SENDER_NAME, senderFields.SenderName, limits.SENDER_NAME_MIN_LENGTH, limits.SENDER_NAME_MAX_LENGTH)...)

    switch {
    case senderFields.FromEmail == "":
        errors = append(errors, v.getError(keys.FROM_EMAIL, enums.SendersMessageEnum.FIELD_IS_REQUIRED))
    case utf8.RuneCountInString(senderFields.FromEmail) > limits.FROM_EMAIL_MAX_LENGTH:
        errors = append(errors, v.getError(keys.FROM_EMAIL, enums.SendersMessageEnum.MAX_LENGTH_IS_X, limits.FROM_EMAIL_MAX_LENGTH))
    case !sharedservices.EmailAddressValidatorGetInstance().IsValid(senderFields.FromEmail):
        errors = append(errors, v.getError(keys.FROM_EMAIL, enums.SendersMessageEnum.INVALID_EMAIL))
    }
    // el nombre va en la cabecera From: sin saltos de linea (inyeccion de cabeceras)
    if utf8.RuneCountInString(senderFields.FromName) > limits.FROM_NAME_MAX_LENGTH || strings.ContainsAny(senderFields.FromName, "\r\n") {
        errors = append(errors, v.getError(keys.FROM_NAME, enums.SendersMessageEnum.MAX_LENGTH_IS_X, limits.FROM_NAME_MAX_LENGTH))
    }

    switch {
    case senderFields.SmtpHost == "":
        errors = append(errors, v.getError(keys.SMTP_HOST, enums.SendersMessageEnum.FIELD_IS_REQUIRED))
    case len(senderFields.SmtpHost) > limits.SMTP_HOST_MAX_LENGTH || !smtpHostPattern.MatchString(senderFields.SmtpHost):
        errors = append(errors, v.getError(keys.SMTP_HOST, enums.SendersMessageEnum.INVALID_HOST))
    }
    if senderFields.SmtpPort < 1 || senderFields.SmtpPort > 65535 {
        errors = append(errors, v.getError(keys.SMTP_PORT, enums.SendersMessageEnum.PORT_MUST_BE_1_TO_65535))
    }
    switch {
    case !slices.Contains(enums.SmtpSecurityEnum.VALUES, senderFields.SmtpSecurity):
        errors = append(errors, v.getError(keys.SMTP_SECURITY, enums.SendersMessageEnum.SECURITY_MUST_BE_X, strings.Join(enums.SmtpSecurityEnum.VALUES, ", ")))
    case senderFields.SmtpSecurity == enums.SmtpSecurityEnum.NONE && senderFields.ProviderCode != enums.ProviderEnum.SMTP_CUSTOM:
        errors = append(errors, v.getError(keys.SMTP_SECURITY, enums.SendersMessageEnum.SECURITY_NONE_ONLY_FOR_CUSTOM))
    }

    if utf8.RuneCountInString(senderFields.Username) > limits.USERNAME_MAX_LENGTH {
        errors = append(errors, v.getError(keys.USERNAME, enums.SendersMessageEnum.MAX_LENGTH_IS_X, limits.USERNAME_MAX_LENGTH))
    }
    // sin secreto solo un smtp_custom (Mailpit no pide autenticacion)
    if senderFields.Secret == "" && senderFields.ProviderCode != enums.ProviderEnum.SMTP_CUSTOM {
        errors = append(errors, v.getError(keys.SECRET, enums.SendersMessageEnum.FIELD_IS_REQUIRED))
    }
    if utf8.RuneCountInString(senderFields.Secret) > limits.SECRET_MAX_LENGTH {
        errors = append(errors, v.getError(keys.SECRET, enums.SendersMessageEnum.MAX_LENGTH_IS_X, limits.SECRET_MAX_LENGTH))
    }

    if senderFields.DailyLimit < 0 || senderFields.DailyLimit > limits.LIMIT_MAX {
        errors = append(errors, v.getError(keys.DAILY_LIMIT, enums.SendersMessageEnum.LIMIT_MUST_BE_0_TO_X, limits.LIMIT_MAX))
    }
    if senderFields.MinuteLimit < 0 || senderFields.MinuteLimit > limits.LIMIT_MAX {
        errors = append(errors, v.getError(keys.MINUTE_LIMIT, enums.SendersMessageEnum.LIMIT_MUST_BE_0_TO_X, limits.LIMIT_MAX))
    }
    return errors
}

func (v *SenderFieldsValidator) getLengthErrors(fieldName string, value string, minLength int, maxLength int) []string {
    valueLength := utf8.RuneCountInString(value)
    switch {
    case value == "":
        return []string{v.getError(fieldName, enums.SendersMessageEnum.FIELD_IS_REQUIRED)}
    case valueLength < minLength:
        return []string{v.getError(fieldName, enums.SendersMessageEnum.MIN_LENGTH_IS_X, minLength)}
    case valueLength > maxLength:
        return []string{v.getError(fieldName, enums.SendersMessageEnum.MAX_LENGTH_IS_X, maxLength)}
    }
    return nil
}

func (v *SenderFieldsValidator) getError(fieldName string, message string, args ...any) string {
    return fieldName + ": " + config.Tr(message, args...)
}
