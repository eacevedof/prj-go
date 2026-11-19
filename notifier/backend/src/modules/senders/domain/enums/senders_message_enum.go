package enums

// SendersMessageEnum mensajes de error de senders ({0}, {1}... los rellena config.Tr).
var SendersMessageEnum = struct {
    FIELD_IS_REQUIRED                string
    MIN_LENGTH_IS_X                  string
    MAX_LENGTH_IS_X                  string
    INVALID_EMAIL                    string
    INVALID_HOST                     string
    PORT_MUST_BE_1_TO_65535          string
    SECURITY_MUST_BE_X               string
    SECURITY_NONE_ONLY_FOR_CUSTOM    string
    LIMIT_MUST_BE_0_TO_X             string
    PROVIDER_X_UNKNOWN               string
    PROVIDER_X_NOT_IMPLEMENTED_YET   string
    SENDER_NAME_X_ALREADY_EXISTS     string
    USER_X_NOT_FOUND                 string
    SENDER_X_NOT_FOUND               string
}{
    FIELD_IS_REQUIRED:              "field is required",
    MIN_LENGTH_IS_X:                "min length is {0}",
    MAX_LENGTH_IS_X:                "max length is {0}",
    INVALID_EMAIL:                  "invalid email",
    INVALID_HOST:                   "invalid host name",
    PORT_MUST_BE_1_TO_65535:        "must be between 1 and 65535",
    SECURITY_MUST_BE_X:             "must be one of {0}",
    SECURITY_NONE_ONLY_FOR_CUSTOM:  "none is only allowed for smtp_custom",
    LIMIT_MUST_BE_0_TO_X:           "must be between 0 and {0} (0 = provider default)",
    PROVIDER_X_UNKNOWN:             "unknown provider {0}",
    PROVIDER_X_NOT_IMPLEMENTED_YET: "provider {0} is not available yet",
    SENDER_NAME_X_ALREADY_EXISTS:   "Sender name {0} already exists for this user",
    USER_X_NOT_FOUND:               "User {0} not found",
    SENDER_X_NOT_FOUND:             "Sender {0} not found",
}
