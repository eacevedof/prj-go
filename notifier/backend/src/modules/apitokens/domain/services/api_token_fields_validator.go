package services

import (
    "strings"
    "unicode/utf8"

    "notifier-backend/src/core/config"

    sharedservices "notifier-backend/src/modules/shared/domain/services"

    "notifier-backend/src/modules/apitokens/domain/enums"
)

const (
    TOKEN_NAME_MIN_LENGTH = 2
    TOKEN_NAME_MAX_LENGTH = 100
)

// ApiTokenFieldsValidator reglas de un token de envio: nombre (la app que lo usa) y destinatarios por defecto.
type ApiTokenFieldsValidator struct {
    emailAddressValidator *sharedservices.EmailAddressValidator
}

func ApiTokenFieldsValidatorGetInstance() *ApiTokenFieldsValidator {
    return &ApiTokenFieldsValidator{emailAddressValidator: sharedservices.EmailAddressValidatorGetInstance()}
}

func (v *ApiTokenFieldsValidator) GetErrors(tokenName string, defaultTo []string, maxRecipients int) []string {
    errors := []string{}
    keys := enums.ApiTokensKeyEnum
    tokenNameLength := utf8.RuneCountInString(tokenName)
    switch {
    case tokenName == "":
        errors = append(errors, keys.TOKEN_NAME+": "+config.Tr(enums.ApiTokensMessageEnum.FIELD_IS_REQUIRED))
    case tokenNameLength < TOKEN_NAME_MIN_LENGTH:
        errors = append(errors, keys.TOKEN_NAME+": "+config.Tr(enums.ApiTokensMessageEnum.MIN_LENGTH_IS_X, TOKEN_NAME_MIN_LENGTH))
    case tokenNameLength > TOKEN_NAME_MAX_LENGTH:
        errors = append(errors, keys.TOKEN_NAME+": "+config.Tr(enums.ApiTokensMessageEnum.MAX_LENGTH_IS_X, TOKEN_NAME_MAX_LENGTH))
    }

    if len(defaultTo) > maxRecipients {
        errors = append(errors, keys.DEFAULT_TO+": "+config.Tr(enums.ApiTokensMessageEnum.MAX_DEFAULT_TO_IS_X, maxRecipients))
    }
    if invalidEmails := v.emailAddressValidator.GetInvalidOnes(defaultTo); len(invalidEmails) > 0 {
        errors = append(errors, keys.DEFAULT_TO+": "+config.Tr(enums.ApiTokensMessageEnum.INVALID_EMAILS_X, strings.Join(invalidEmails, ", ")))
    }
    return errors
}

// GetDefaultToFromCsv "a@x.com, b@x.com" -> [a@x.com b@x.com] normalizadas y sin repetir (la consola).
func (v *ApiTokenFieldsValidator) GetDefaultToFromCsv(defaultToCsv string) []string {
    defaultTo := []string{}
    seen := map[string]bool{}
    for _, emailAddress := range strings.Split(defaultToCsv, ",") {
        emailAddress = v.emailAddressValidator.GetNormalized(emailAddress)
        if emailAddress == "" || seen[emailAddress] {
            continue
        }
        seen[emailAddress] = true
        defaultTo = append(defaultTo, emailAddress)
    }
    return defaultTo
}

// GetVisiblePrefix los primeros caracteres del token en claro, para reconocerlo en el panel.
func (v *ApiTokenFieldsValidator) GetVisiblePrefix(apiToken string) string {
    if len(apiToken) <= enums.ApiTokenFormatEnum.VISIBLE_LENGTH {
        return apiToken
    }
    return apiToken[:enums.ApiTokenFormatEnum.VISIBLE_LENGTH]
}
