package enums

// DeliveryMessageEnum mensajes de la entrega ({0}, {1}... los rellena config.Tr). Son los que quedan en
// app_messages.error_message y en app_message_attempts.provider_response.
var DeliveryMessageEnum = struct {
    MAX_MESSAGES_MUST_BE_1_TO_500     string
    CHANNEL_X_HAS_NO_WORKER           string
    SENDER_NOT_FOUND                  string
    PROVIDER_X_HAS_NO_MESSAGE_SENDER  string
    NO_RECIPIENTS                     string
    ATTACHMENT_X_URL_NOT_ALLOWED      string
    ATTACHMENT_X_HOST_IS_PRIVATE      string
    ATTACHMENT_X_HTTP_STATUS_Y        string
    ATTACHMENT_X_TOO_BIG              string
    ATTACHMENT_X_DOWNLOAD_FAILED      string
    SMTP_SERVER_HAS_NO_STARTTLS       string
    SMTP_AUTH_NEEDS_TLS               string
    SENT_WITH_MESSAGE_ID_X            string
}{
    MAX_MESSAGES_MUST_BE_1_TO_500:    "max_messages: must be between 1 and 500",
    CHANNEL_X_HAS_NO_WORKER:          "channel: {0} has no worker yet",
    SENDER_NOT_FOUND:                 "sender not found",
    PROVIDER_X_HAS_NO_MESSAGE_SENDER: "provider {0} has no message sender yet",
    NO_RECIPIENTS:                    "the message has no recipients",
    ATTACHMENT_X_URL_NOT_ALLOWED:     "attachment {0}: url not allowed",
    ATTACHMENT_X_HOST_IS_PRIVATE:     "attachment {0}: host resolves to a private address",
    ATTACHMENT_X_HTTP_STATUS_Y:       "attachment {0}: download answered HTTP {1}",
    ATTACHMENT_X_TOO_BIG:             "attachment {0}: bigger than {1} bytes",
    ATTACHMENT_X_DOWNLOAD_FAILED:     "attachment {0}: download failed: {1}",
    SMTP_SERVER_HAS_NO_STARTTLS:      "smtp server does not offer STARTTLS",
    SMTP_AUTH_NEEDS_TLS:              "smtp auth without TLS is not allowed (use starttls or ssl)",
    SENT_WITH_MESSAGE_ID_X:           "accepted by the smtp server, Message-ID {0}",
}
