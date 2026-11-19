package enums

// MessagesMessageEnum mensajes de error de messages ({0}, {1}... los rellena config.Tr).
var MessagesMessageEnum = struct {
    FIELD_IS_REQUIRED                 string
    MUST_BE_A_STRING_OR_A_LIST        string
    MUST_BE_A_LIST                    string
    MAX_LENGTH_IS_X                   string
    MAX_X_RECIPIENTS                  string
    INVALID_EMAILS_X                  string
    NO_LINE_BREAKS                    string
    MUST_BE_X                         string
    MAX_X_ATTACHMENTS                 string
    ATTACHMENT_X_NEEDS_CONTENT_OR_URL string
    ATTACHMENT_X_INVALID_FILE_NAME    string
    ATTACHMENT_X_TOO_BIG              string
    ATTACHMENT_X_URL_NOT_ALLOWED      string
    CHANNEL_X_NOT_AVAILABLE_YET       string
    MESSAGE_X_NOT_FOUND               string
    INVALID_OR_MISSING_API_TOKEN      string
}{
    FIELD_IS_REQUIRED:                 "field is required",
    MUST_BE_A_STRING_OR_A_LIST:        "must be a string or a list of strings",
    MUST_BE_A_LIST:                    "must be a list",
    MAX_LENGTH_IS_X:                   "max length is {0}",
    MAX_X_RECIPIENTS:                  "max {0} recipients between to, cc and bcc",
    INVALID_EMAILS_X:                  "invalid emails: {0}",
    NO_LINE_BREAKS:                    "line breaks are not allowed",
    MUST_BE_X:                         "must be one of {0}",
    MAX_X_ATTACHMENTS:                 "max {0} attachments",
    ATTACHMENT_X_NEEDS_CONTENT_OR_URL: "attachment {0}: needs content or url (not both)",
    ATTACHMENT_X_INVALID_FILE_NAME:    "attachment {0}: file_name max 255 chars, without / \\ \" or line breaks",
    ATTACHMENT_X_TOO_BIG:              "attachment {0}: content is bigger than {1} bytes",
    ATTACHMENT_X_URL_NOT_ALLOWED:      "attachment {0}: url must be https and its host one of APP_ATTACHMENT_ALLOWED_HOSTS",
    CHANNEL_X_NOT_AVAILABLE_YET:       "channel {0} is not available yet",
    MESSAGE_X_NOT_FOUND:               "Message {0} not found",
    INVALID_OR_MISSING_API_TOKEN:      "Invalid or missing api token",
}
