package enums

// ApiTokenFormatEnum forma de un token de envio: "ntf-" + 64 hex. El prefijo deja reconocerlo en un log o en un
// escaner de secretos; VISIBLE_LENGTH caracteres (prefijo incluido) se guardan en claro para el panel.
var ApiTokenFormatEnum = struct {
    PREFIX         string
    VISIBLE_LENGTH int
    DEFAULT_NAME   string
}{
    PREFIX:         "ntf-",
    VISIBLE_LENGTH: 10,
    DEFAULT_NAME:   "default",
}
