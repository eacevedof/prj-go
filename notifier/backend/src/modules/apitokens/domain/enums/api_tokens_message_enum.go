package enums

// ApiTokensMessageEnum mensajes de error de apitokens ({0}, {1}... los rellena config.Tr).
var ApiTokensMessageEnum = struct {
    FIELD_IS_REQUIRED             string
    MIN_LENGTH_IS_X               string
    MAX_LENGTH_IS_X               string
    INVALID_EMAILS_X              string
    MAX_DEFAULT_TO_IS_X           string
    TOKEN_NAME_X_ALREADY_EXISTS   string
    SENDER_X_NOT_FOUND            string
    TOKEN_X_NOT_FOUND             string
    INVALID_OR_MISSING_API_TOKEN  string
}{
    FIELD_IS_REQUIRED:            "field is required",
    MIN_LENGTH_IS_X:              "min length is {0}",
    MAX_LENGTH_IS_X:              "max length is {0}",
    INVALID_EMAILS_X:             "invalid emails: {0}",
    MAX_DEFAULT_TO_IS_X:          "max {0} recipients",
    TOKEN_NAME_X_ALREADY_EXISTS:  "Token name {0} already exists for this sender",
    SENDER_X_NOT_FOUND:           "Sender {0} not found",
    TOKEN_X_NOT_FOUND:            "Token {0} not found",
    INVALID_OR_MISSING_API_TOKEN: "Invalid or missing api token",
}
