package enums

// BodyFormatEnum formato del cuerpo. Mismos valores que el CHECK de app_messages.body_format.
var BodyFormatEnum = struct {
    TEXT   string
    HTML   string
    VALUES []string
}{
    TEXT:   "text",
    HTML:   "html",
    VALUES: []string{"text", "html"},
}

// AttachmentTypeEnum de donde sale un adjunto: el texto va en el mensaje (text) o lo descarga el worker (url).
var AttachmentTypeEnum = struct {
    TEXT string
    URL  string
}{
    TEXT: "text",
    URL:  "url",
}

// AttemptStatusEnum resultado de un intento de entrega. Mismos valores que el CHECK de
// app_message_attempts.attempt_status.
var AttemptStatusEnum = struct {
    SENT            string
    TEMPORARY_ERROR string
    PERMANENT_ERROR string
}{
    SENT:            "sent",
    TEMPORARY_ERROR: "temporary_error",
    PERMANENT_ERROR: "permanent_error",
}
