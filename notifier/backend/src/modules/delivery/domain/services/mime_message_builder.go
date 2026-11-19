package services

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "mime"
    "mime/multipart"
    "mime/quotedprintable"
    "net/mail"
    "net/textproto"
    "strings"
    "time"

    "notifier-backend/src/modules/delivery/domain/entities"
)

const (
    crlf               = "\r\n"
    base64LineLength   = 76
    contentTypeText    = "text/plain; charset=UTF-8"
    contentTypeHtml    = "text/html; charset=UTF-8"
    defaultMimeType    = "application/octet-stream"
    encodingQuoted     = "quoted-printable"
    encodingBase64     = "base64"
)

// MimeMessageBuilder el email en formato RFC 5322 + MIME, listo para el DATA de SMTP: cabeceras con los textos no ASCII
// codificados (RFC 2047), cuerpo en quoted-printable y, si hay adjuntos, multipart/mixed con cada uno en base64.
// Bcc nunca va en las cabeceras. Todos los valores ya vienen validados (sin saltos de linea).
type MimeMessageBuilder struct{}

func MimeMessageBuilderGetInstance() *MimeMessageBuilder {
    return &MimeMessageBuilder{}
}

func (b *MimeMessageBuilder) GetMessage(outgoingEmail entities.OutgoingEmail) ([]byte, error) {
    var message bytes.Buffer
    b.writeHeaders(&message, outgoingEmail)

    if len(outgoingEmail.Attachments) == 0 {
        b.writeHeader(&message, "Content-Type", b.getBodyContentType(outgoingEmail))
        b.writeHeader(&message, "Content-Transfer-Encoding", encodingQuoted)
        message.WriteString(crlf)
        return message.Bytes(), b.writeQuotedPrintable(&message, outgoingEmail.Body)
    }

    multipartWriter := multipart.NewWriter(&message)
    b.writeHeader(&message, "Content-Type", mime.FormatMediaType("multipart/mixed", map[string]string{"boundary": multipartWriter.Boundary()}))
    message.WriteString(crlf)

    bodyPart, err := multipartWriter.CreatePart(textproto.MIMEHeader{
        "Content-Type":              {b.getBodyContentType(outgoingEmail)},
        "Content-Transfer-Encoding": {encodingQuoted},
    })
    if err != nil {
        return nil, err
    }
    if err := b.writeQuotedPrintable(bodyPart, outgoingEmail.Body); err != nil {
        return nil, err
    }
    for _, attachment := range outgoingEmail.Attachments {
        if err := b.writeAttachment(multipartWriter, attachment); err != nil {
            return nil, err
        }
    }
    if err := multipartWriter.Close(); err != nil {
        return nil, err
    }
    return message.Bytes(), nil
}

func (b *MimeMessageBuilder) writeHeaders(message *bytes.Buffer, outgoingEmail entities.OutgoingEmail) {
    date := outgoingEmail.Date
    if date.IsZero() {
        date = time.Now()
    }
    b.writeHeader(message, "Date", date.Format(time.RFC1123Z))
    b.writeHeader(message, "Message-ID", "<"+outgoingEmail.MessageId+">")
    b.writeHeader(message, "From", (&mail.Address{Name: outgoingEmail.FromName, Address: outgoingEmail.FromEmail}).String())
    b.writeHeader(message, "To", b.getAddressList(outgoingEmail.To))
    if len(outgoingEmail.Cc) > 0 {
        b.writeHeader(message, "Cc", b.getAddressList(outgoingEmail.Cc))
    }
    if outgoingEmail.ReplyTo != "" {
        b.writeHeader(message, "Reply-To", (&mail.Address{Address: outgoingEmail.ReplyTo}).String())
    }
    b.writeHeader(message, "Subject", mime.QEncoding.Encode("utf-8", outgoingEmail.Subject))
    b.writeHeader(message, "MIME-Version", "1.0")
}

func (b *MimeMessageBuilder) writeHeader(message *bytes.Buffer, name string, value string) {
    message.WriteString(name + ": " + value + crlf)
}

func (b *MimeMessageBuilder) getAddressList(emailAddresses []string) string {
    formattedAddresses := make([]string, 0, len(emailAddresses))
    for _, emailAddress := range emailAddresses {
        formattedAddresses = append(formattedAddresses, (&mail.Address{Address: emailAddress}).String())
    }
    return strings.Join(formattedAddresses, ", ")
}

func (b *MimeMessageBuilder) getBodyContentType(outgoingEmail entities.OutgoingEmail) string {
    if outgoingEmail.IsHtml {
        return contentTypeHtml
    }
    return contentTypeText
}

// writeQuotedPrintable los saltos de linea del cuerpo pasan a CRLF (lo que exige SMTP) antes de codificar.
func (b *MimeMessageBuilder) writeQuotedPrintable(writer interface{ Write([]byte) (int, error) }, body string) error {
    normalizedBody := strings.ReplaceAll(strings.ReplaceAll(body, "\r\n", "\n"), "\n", crlf)
    quotedPrintableWriter := quotedprintable.NewWriter(writer)
    if _, err := quotedPrintableWriter.Write([]byte(normalizedBody)); err != nil {
        return err
    }
    return quotedPrintableWriter.Close()
}

// writeAttachment el nombre va con FormatMediaType: si no es ASCII sale como filename*=utf-8''... (RFC 2231).
func (b *MimeMessageBuilder) writeAttachment(multipartWriter *multipart.Writer, attachment entities.OutgoingAttachment) error {
    // FormatMediaType da "" con un tipo mal formado: entonces va como binario generico
    contentType := mime.FormatMediaType(attachment.MimeType, map[string]string{"name": attachment.FileName})
    if contentType == "" {
        contentType = mime.FormatMediaType(defaultMimeType, map[string]string{"name": attachment.FileName})
    }
    attachmentPart, err := multipartWriter.CreatePart(textproto.MIMEHeader{
        "Content-Type":              {contentType},
        "Content-Disposition":       {mime.FormatMediaType("attachment", map[string]string{"filename": attachment.FileName})},
        "Content-Transfer-Encoding": {encodingBase64},
    })
    if err != nil {
        return err
    }
    encodedContent := base64.StdEncoding.EncodeToString(attachment.Content)
    for start := 0; start < len(encodedContent); start += base64LineLength {
        end := min(start+base64LineLength, len(encodedContent))
        if _, err := fmt.Fprint(attachmentPart, encodedContent[start:end]+crlf); err != nil {
            return err
        }
    }
    return nil
}
