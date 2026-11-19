package entities

import "time"

// OutgoingAttachment un adjunto listo para ir en el email (el de url ya descargado).
type OutgoingAttachment struct {
    FileName string
    MimeType string
    Content  []byte
}

// OutgoingEmail todo lo que hace falta para montar y entregar un email, sin nada de la BD.
type OutgoingEmail struct {
    FromEmail   string
    FromName    string
    To          []string
    Cc          []string
    Bcc         []string
    ReplyTo     string
    Subject     string
    Body        string
    IsHtml      bool
    MessageId   string
    Date        time.Time
    Attachments []OutgoingAttachment
}

// GetEnvelopeRecipients a quien se entrega de verdad (RCPT TO): to + cc + bcc.
func (e OutgoingEmail) GetEnvelopeRecipients() []string {
    envelopeRecipients := make([]string, 0, len(e.To)+len(e.Cc)+len(e.Bcc))
    envelopeRecipients = append(envelopeRecipients, e.To...)
    envelopeRecipients = append(envelopeRecipients, e.Cc...)
    return append(envelopeRecipients, e.Bcc...)
}

// SmtpSettings con que servidor y credenciales sale (el secreto ya descifrado: vive solo en memoria).
type SmtpSettings struct {
    Host     string
    Port     int
    Security string
    Username string
    Password string
}
