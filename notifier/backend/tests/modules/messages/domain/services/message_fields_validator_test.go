package services

import (
    "strings"
    "testing"

    "notifier-backend/tests"

    "notifier-backend/src/modules/messages/domain/services"
)

var defaultLimits = services.MessageFieldsLimits{MaxRecipients: 3, AttachmentMaxBytes: 10}

func getValidator() *services.MessageFieldsValidator {
    return services.MessageFieldsValidatorGetInstance(services.AttachmentUrlCheckerFromAllowedHosts([]string{"cdn.x.com"}))
}

func getValidFields() services.EmailMessageFields {
    return services.EmailMessageFields{To: []string{"a@x.com"}, Subject: "hola", Body: "cuerpo", BodyFormat: "text"}
}

func TestValidEmailHasNoErrors(t *testing.T) {
    tests.AssertEmpty(t, getValidator().GetEmailErrors(getValidFields(), defaultLimits))
}

func TestToSubjectAndFormatAreChecked(t *testing.T) {
    messageFields := getValidFields()
    messageFields.To = []string{}
    messageFields.Subject = "linea\r\nBcc: x@y.com"
    messageFields.BodyFormat = "markdown"
    errors := getValidator().GetEmailErrors(messageFields, defaultLimits)
    tests.AssertEquals(t, 3, len(errors))
    tests.AssertStringContainsString(t, "to: field is required", errors[0])
    tests.AssertStringContainsString(t, "subject: line breaks", errors[1])
    tests.AssertStringContainsString(t, "body_format", errors[2])
}

func TestRecipientsLimitAndInvalidEmails(t *testing.T) {
    messageFields := getValidFields()
    messageFields.To = []string{"a@x.com", "no-es-email"}
    messageFields.Cc = []string{"c@x.com"}
    messageFields.Bcc = []string{"d@x.com"}
    errors := strings.Join(getValidator().GetEmailErrors(messageFields, defaultLimits), " | ")
    tests.AssertStringContainsString(t, "max 3 recipients", errors)
    tests.AssertStringContainsString(t, "invalid emails: no-es-email", errors)
}

func TestAttachmentsRules(t *testing.T) {
    messageFields := getValidFields()
    messageFields.Attachments = []services.AttachmentInput{
        {FileName: "ok.txt", Content: "corto", HasContent: true},
        {FileName: "../etc/passwd", Content: "x", HasContent: true},
        {FileName: "grande.txt", Content: "mas de diez bytes", HasContent: true},
        {FileName: "a.pdf", Url: "http://cdn.x.com/a.pdf", HasUrl: true},
        {FileName: "b.pdf", Url: "https://evil.com/b.pdf", HasUrl: true},
        {FileName: "c.pdf", Url: "https://cdn.x.com/c.pdf", HasUrl: true},
        {FileName: "d.txt"},
    }
    errors := getValidator().GetEmailErrors(messageFields, defaultLimits)
    tests.AssertEquals(t, 5, len(errors))
    tests.AssertStringContainsString(t, "attachment 2: file_name", errors[0])
    tests.AssertStringContainsString(t, "attachment 3: content is bigger", errors[1])
    tests.AssertStringContainsString(t, "attachment 4: url", errors[2])
    tests.AssertStringContainsString(t, "attachment 5: url", errors[3])
    tests.AssertStringContainsString(t, "attachment 7: needs content or url", errors[4])
}

func TestAttachmentUrlChecker(t *testing.T) {
    checker := services.AttachmentUrlCheckerFromAllowedHosts([]string{"cdn.x.com"})
    tests.AssertEquals(t, true, checker.IsAllowed("https://CDN.x.com/a.pdf"))
    tests.AssertEquals(t, true, checker.IsAllowed("https://cdn.x.com:443/a.pdf"))
    tests.AssertEquals(t, false, checker.IsAllowed("https://cdn.x.com:8443/a.pdf"))
    tests.AssertEquals(t, false, checker.IsAllowed("https://user:pw@cdn.x.com/a.pdf"))
    tests.AssertEquals(t, false, checker.IsAllowed("https://cdn.x.com.evil.com/a.pdf"))
    tests.AssertEquals(t, false, checker.IsAllowed("ftp://cdn.x.com/a.pdf"))
}
