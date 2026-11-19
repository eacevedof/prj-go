package enums

// UserFieldsLimitEnum limites de los campos de un usuario (UserFieldsValidator).
var UserFieldsLimitEnum = struct {
    USER_NAME_MIN_LENGTH  int
    USER_NAME_MAX_LENGTH  int
    USER_EMAIL_MAX_LENGTH int
    // bcrypt solo usa los primeros 72 bytes: mas largo daria por buena una password que no coincide entera
    PASSWORD_MIN_LENGTH int
    PASSWORD_MAX_BYTES  int
}{
    USER_NAME_MIN_LENGTH:  3,
    USER_NAME_MAX_LENGTH:  100,
    USER_EMAIL_MAX_LENGTH: 150,
    PASSWORD_MIN_LENGTH:   10,
    PASSWORD_MAX_BYTES:    72,
}
