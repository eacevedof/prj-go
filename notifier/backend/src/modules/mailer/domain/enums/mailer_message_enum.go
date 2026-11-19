package enums

// MailerMessageEnum mensajes de error de mailer ({0}... los rellena config.Tr).
var MailerMessageEnum = struct {
    UNKNOWN_ALERT_TYPE_X string
    SUBJECT_IS_REQUIRED  string
}{
    UNKNOWN_ALERT_TYPE_X: "SendAlertEmail: unknown alert type \"{0}\"",
    SUBJECT_IS_REQUIRED:  "SendAlertEmail: subject is required",
}
