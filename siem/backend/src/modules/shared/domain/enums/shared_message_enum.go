package enums

// SharedMessageEnum mensajes de error comunes (en ingles; los que se usan, el resto de task-php es herencia sin uso).
// Los {0}, {1}... los sustituye config.Tr; MAILER_SENDMAIL_FAILED es un formato de fmt.Errorf (envuelve el error).
var SharedMessageEnum = struct {
    AN_UNEXPECTED_ERROR_HAS_OCCURRED string
    FIELDS_VALIDATION_ERRORS         string
    RESOURCE_NOT_FOUND               string
    METHOD_NOT_ALLOWED               string
    UNEXPECTED_ERROR_WITH_ID         string
    CONSOLE_COMMAND_NOT_FOUND        string
    MISSING_ENV_SQLITE_PATH          string
    MAILER_FROM_AND_TO_REQUIRED      string
    MAILER_SENDMAIL_FAILED           string
}{
    AN_UNEXPECTED_ERROR_HAS_OCCURRED: "An unexpected error has occurred",
    FIELDS_VALIDATION_ERRORS:         "Fields validation errors",
    RESOURCE_NOT_FOUND:               "Resource not found",
    METHOD_NOT_ALLOWED:               "Method not allowed",
    UNEXPECTED_ERROR_WITH_ID:         "Sorry! but some unexpected error occurred ({0})",
    CONSOLE_COMMAND_NOT_FOUND:        "no class found for cmd {0}",
    MISSING_ENV_SQLITE_PATH:          "missing env APP_SQLITE_PATH",
    MAILER_FROM_AND_TO_REQUIRED:      "FuncMailer: email from and at least one email to are required",
    MAILER_SENDMAIL_FAILED:           "FuncMailer: mail() failed. %w %s",
}
