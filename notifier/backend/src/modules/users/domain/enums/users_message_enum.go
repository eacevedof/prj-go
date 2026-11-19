package enums

// UsersMessageEnum mensajes de error de users ({0}, {1}... los rellena config.Tr).
var UsersMessageEnum = struct {
    FIELD_IS_REQUIRED             string
    MIN_LENGTH_IS_X               string
    MAX_LENGTH_IS_X               string
    INVALID_EMAIL                 string
    INVALID_ROLE                  string
    USER_EMAIL_X_ALREADY_EXISTS   string
    ROOT_ALREADY_EXISTS           string
    USER_X_NOT_FOUND              string
    INVALID_EMAIL_OR_PASSWORD     string
    MISSING_OR_INVALID_API_TOKEN  string
}{
    FIELD_IS_REQUIRED:            "field is required",
    MIN_LENGTH_IS_X:              "min length is {0}",
    MAX_LENGTH_IS_X:              "max length is {0}",
    INVALID_EMAIL:                "invalid email",
    INVALID_ROLE:                 "must be one of root, admin, owner",
    USER_EMAIL_X_ALREADY_EXISTS:  "User email {0} already exists",
    ROOT_ALREADY_EXISTS:          "There is already a root user",
    USER_X_NOT_FOUND:             "User {0} not found",
    INVALID_EMAIL_OR_PASSWORD:    "Invalid email or password",
    MISSING_OR_INVALID_API_TOKEN: "Missing or invalid api token",
}
