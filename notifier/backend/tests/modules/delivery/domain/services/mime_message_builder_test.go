package services

import (
    "io"
    "net/mail"
    "strings"
    "testing"
    "time"

    "mime"
    "mime/multipart"

    "notifier-backend/tests"

    "notifier-backend/src/modules/delivery/domain/entities"
    "notifier-backend/src/modules/delivery/domain/services"
)

func getOutgoingEmail() entities.OutgoingEmail {
    return entities.OutgoingEmail{
        FromEmail: "yo@x.com",
        FromName:  "Avisos ñ",
        To:        []string{"a@x.com", "b@x.com"},
        Cc:        []string{"c@x.com"},
        Bcc:       []string{"oculto@x.com"},
        ReplyTo:   "r@x.com",
        Subject:   "Factura de café",
        Body:      "línea 1\nlínea 2",
        MessageId: "message-1@x.com",
        Date:      time.Date(2026, 11, 18, 10, 0, 0, 0, time.UTC),
    }
}

func TestPlainEmailHeadersAndBody(t *testing.T) {
    rawMessage, err := services.MimeMessageBuilderGetInstance().GetMessage(getOutgoingEmail())
    tests.AssertNil(t, err)

    message, err := mail.ReadMessage(strings.NewReader(string(rawMessage)))
    tests.AssertNil(t, err)
    decoder := new(mime.WordDecoder)
    subject, _ := decoder.DecodeHeader(message.Header.Get("Subject"))
    tests.AssertEquals(t, "Factura de café", subject)
    from, _ := message.Header.AddressList("From")
    tests.AssertEquals(t, "Avisos ñ", from[0].Name)
    tests.AssertEquals(t, "<a@x.com>, <b@x.com>", message.Header.Get("To"))
    tests.AssertEquals(t, "<c@x.com>", message.Header.Get("Cc"))
    tests.AssertEquals(t, "", message.Header.Get("Bcc"))
    tests.AssertEquals(t, "<message-1@x.com>", message.Header.Get("Message-Id"))
    tests.AssertStringContainsString(t, "text/plain", message.Header.Get("Content-Type"))
    tests.AssertEquals(t, false, strings.Contains(string(rawMessage), "oculto@x.com"))
}

func TestEmailWithAttachmentsIsMultipart(t *testing.T) {
    outgoingEmail := getOutgoingEmail()
    outgoingEmail.IsHtml = true
    outgoingEmail.Attachments = []entities.OutgoingAttachment{
        {FileName: "nota ñ.txt", MimeType: "text/plain", Content: []byte("hola")},
        {FileName: "datos.bin", MimeType: "mal formado;;", Content: []byte{0, 1, 2}},
    }
    rawMessage, err := services.MimeMessageBuilderGetInstance().GetMessage(outgoingEmail)
    tests.AssertNil(t, err)

    message, _ := mail.ReadMessage(strings.NewReader(string(rawMessage)))
    mediaType, params, err := mime.ParseMediaType(message.Header.Get("Content-Type"))
    tests.AssertNil(t, err)
    tests.AssertEquals(t, "multipart/mixed", mediaType)

    multipartReader := multipart.NewReader(message.Body, params["boundary"])
    parts := []*multipart.Part{}
    fileNames := []string{}
    for {
        part, err := multipartReader.NextPart()
        if err == io.EOF {
            break
        }
        tests.AssertNil(t, err)
        parts = append(parts, part)
        fileNames = append(fileNames, part.FileName())
    }
    tests.AssertEquals(t, 3, len(parts))
    tests.AssertEquals(t, []string{"", "nota ñ.txt", "datos.bin"}, fileNames)
    tests.AssertStringContainsString(t, "text/html", parts[0].Header.Get("Content-Type"))
    tests.AssertStringContainsString(t, "application/octet-stream", parts[2].Header.Get("Content-Type"))
}
