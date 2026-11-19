package services

import (
    "slices"
    "unicode/utf8"

    "notifier-backend/src/core/config"

    sharedservices "notifier-backend/src/modules/shared/domain/services"

    "notifier-backend/src/modules/users/domain/enums"
)

// UserFieldsValidator reglas de escritura de un usuario. Devuelve un error "campo: mensaje" por regla incumplida;
// vacio = valido.
type UserFieldsValidator struct{}

func UserFieldsValidatorGetInstance() *UserFieldsValidator {
    return &UserFieldsValidator{}
}

func (v *UserFieldsValidator) GetErrors(userName string, userEmail string, userRole string, password string) []string {
    errors := []string{}
    limits := enums.UserFieldsLimitEnum

    userNameLength := utf8.RuneCountInString(userName)
    switch {
    case userName == "":
        errors = append(errors, v.getError(enums.UsersKeyEnum.USER_NAME, enums.UsersMessageEnum.FIELD_IS_REQUIRED))
    case userNameLength < limits.USER_NAME_MIN_LENGTH:
        errors = append(errors, v.getError(enums.UsersKeyEnum.USER_NAME, enums.UsersMessageEnum.MIN_LENGTH_IS_X, limits.USER_NAME_MIN_LENGTH))
    case userNameLength > limits.USER_NAME_MAX_LENGTH:
        errors = append(errors, v.getError(enums.UsersKeyEnum.USER_NAME, enums.UsersMessageEnum.MAX_LENGTH_IS_X, limits.USER_NAME_MAX_LENGTH))
    }

    switch {
    case userEmail == "":
        errors = append(errors, v.getError(enums.UsersKeyEnum.USER_EMAIL, enums.UsersMessageEnum.FIELD_IS_REQUIRED))
    case utf8.RuneCountInString(userEmail) > limits.USER_EMAIL_MAX_LENGTH:
        errors = append(errors, v.getError(enums.UsersKeyEnum.USER_EMAIL, enums.UsersMessageEnum.MAX_LENGTH_IS_X, limits.USER_EMAIL_MAX_LENGTH))
    case !sharedservices.EmailAddressValidatorGetInstance().IsValid(userEmail):
        errors = append(errors, v.getError(enums.UsersKeyEnum.USER_EMAIL, enums.UsersMessageEnum.INVALID_EMAIL))
    }

    if !slices.Contains(enums.UserRoleEnum.VALUES, userRole) {
        errors = append(errors, v.getError(enums.UsersKeyEnum.USER_ROLE, enums.UsersMessageEnum.INVALID_ROLE))
    }

    switch {
    case password == "":
        errors = append(errors, v.getError(enums.UsersKeyEnum.PASSWORD, enums.UsersMessageEnum.FIELD_IS_REQUIRED))
    case utf8.RuneCountInString(password) < limits.PASSWORD_MIN_LENGTH:
        errors = append(errors, v.getError(enums.UsersKeyEnum.PASSWORD, enums.UsersMessageEnum.MIN_LENGTH_IS_X, limits.PASSWORD_MIN_LENGTH))
    case len(password) > limits.PASSWORD_MAX_BYTES:
        errors = append(errors, v.getError(enums.UsersKeyEnum.PASSWORD, enums.UsersMessageEnum.MAX_LENGTH_IS_X, limits.PASSWORD_MAX_BYTES))
    }
    return errors
}

func (v *UserFieldsValidator) getError(fieldName string, message string, args ...any) string {
    return fieldName + ": " + config.Tr(message, args...)
}
