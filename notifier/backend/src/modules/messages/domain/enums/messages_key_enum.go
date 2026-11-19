package enums

// MessagesKeyEnum claves del contrato de POST /api/v1/messages y de las respuestas, mas las columnas de app_messages,
// app_message_attachments y app_message_attempts.
var MessagesKeyEnum = struct {
    API_TOKEN       string
    TO              string
    CC              string
    BCC             string
    REPLY_TO        string
    SUBJECT         string
    BODY            string
    BODY_FORMAT     string
    ATTACHMENTS     string
    IDEMPOTENCY_KEY string

    FILE_NAME       string
    CONTENT         string
    URL             string
    ATTACHMENT_TYPE string
    SOURCE_URL      string
    MIME_TYPE       string
    SIZE_BYTES      string

    MESSAGE_UUID        string
    CHANNEL             string
    MESSAGE_STATUS      string
    RECIPIENTS          string
    ATTEMPTS            string
    MAX_ATTEMPTS        string
    NEXT_ATTEMPT_AT     string
    ERROR_MESSAGE       string
    PROVIDER_MESSAGE_ID string
    REQUEST_IP          string
    CREATED_AT          string
    SENT_AT             string
    UPDATED_AT          string
    IS_DUPLICATE        string

    ATTEMPT_NUMBER    string
    ATTEMPT_STATUS    string
    PROVIDER_RESPONSE string
    DURATION_MS       string
    ATTEMPTED_AT      string

    // columnas internas de BD: nunca salen por la API
    ID           string
    API_TOKEN_ID string
    SENDER_ID    string
    MESSAGE_ID   string
}{
    API_TOKEN:       "api_token",
    TO:              "to",
    CC:              "cc",
    BCC:             "bcc",
    REPLY_TO:        "reply_to",
    SUBJECT:         "subject",
    BODY:            "body",
    BODY_FORMAT:     "body_format",
    ATTACHMENTS:     "attachments",
    IDEMPOTENCY_KEY: "idempotency_key",

    FILE_NAME:       "file_name",
    CONTENT:         "content",
    URL:             "url",
    ATTACHMENT_TYPE: "attachment_type",
    SOURCE_URL:      "source_url",
    MIME_TYPE:       "mime_type",
    SIZE_BYTES:      "size_bytes",

    MESSAGE_UUID:        "message_uuid",
    CHANNEL:             "channel",
    MESSAGE_STATUS:      "message_status",
    RECIPIENTS:          "recipients",
    ATTEMPTS:            "attempts",
    MAX_ATTEMPTS:        "max_attempts",
    NEXT_ATTEMPT_AT:     "next_attempt_at",
    ERROR_MESSAGE:       "error_message",
    PROVIDER_MESSAGE_ID: "provider_message_id",
    REQUEST_IP:          "request_ip",
    CREATED_AT:          "created_at",
    SENT_AT:             "sent_at",
    UPDATED_AT:          "updated_at",
    IS_DUPLICATE:        "is_duplicate",

    ATTEMPT_NUMBER:    "attempt_number",
    ATTEMPT_STATUS:    "attempt_status",
    PROVIDER_RESPONSE: "provider_response",
    DURATION_MS:       "duration_ms",
    ATTEMPTED_AT:      "attempted_at",

    ID:           "id",
    API_TOKEN_ID: "api_token_id",
    SENDER_ID:    "sender_id",
    MESSAGE_ID:   "message_id",
}
