package services

import (
    "regexp"
    "strings"
)

// el FILTER_VALIDATE_EMAIL de PHP: parte local de atomos separados por puntos y dominio con al menos un punto
var emailAddressPattern = regexp.MustCompile(
    "^[A-Za-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[A-Za-z0-9!#$%&'*+/=?^_`{|}~-]+)*" +
        `@(?:[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?\.)+[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?$`,
)

const emailAddressMaxLength = 254

// EmailAddressValidator una direccion de email suelta (sin "Nombre <...>"): la usan users, senders, apitokens y
// messages. Al no admitir espacios ni saltos de linea, una direccion valida no puede inyectar cabeceras.
type EmailAddressValidator struct{}

func EmailAddressValidatorGetInstance() *EmailAddressValidator {
    return &EmailAddressValidator{}
}

func (v *EmailAddressValidator) IsValid(emailAddress string) bool {
    return len(emailAddress) <= emailAddressMaxLength && emailAddressPattern.MatchString(emailAddress)
}

// GetInvalidOnes las que no son validas, para el mensaje de error.
func (v *EmailAddressValidator) GetInvalidOnes(emailAddresses []string) []string {
    invalidOnes := []string{}
    for _, emailAddress := range emailAddresses {
        if !v.IsValid(emailAddress) {
            invalidOnes = append(invalidOnes, emailAddress)
        }
    }
    return invalidOnes
}

// GetNormalized en minusculas y sin espacios alrededor: dos formas de escribir la misma direccion cuentan como una.
func (v *EmailAddressValidator) GetNormalized(emailAddress string) string {
    return strings.ToLower(strings.TrimSpace(emailAddress))
}
