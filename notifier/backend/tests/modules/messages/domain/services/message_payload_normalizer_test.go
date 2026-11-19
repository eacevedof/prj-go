package services

import (
    "testing"

    "notifier-backend/tests"

    "notifier-backend/src/modules/messages/domain/services"
)

func TestEmailListAcceptsStringCsvOrList(t *testing.T) {
    normalizer := services.MessagePayloadNormalizerGetInstance()

    emails, isValid := normalizer.GetEmailList(`"Ana@Example.com"`)
    tests.AssertEquals(t, true, isValid)
    tests.AssertEquals(t, []string{"ana@example.com"}, emails)

    emails, _ = normalizer.GetEmailList(`"a@x.com, b@x.com ,a@x.com"`)
    tests.AssertEquals(t, []string{"a@x.com", "b@x.com"}, emails)

    emails, _ = normalizer.GetEmailList(`["a@x.com", " B@x.com "]`)
    tests.AssertEquals(t, []string{"a@x.com", "b@x.com"}, emails)

    emails, isValid = normalizer.GetEmailList("")
    tests.AssertEquals(t, true, isValid)
    tests.AssertEquals(t, []string{}, emails)
}

func TestEmailListRejectsOtherTypes(t *testing.T) {
    normalizer := services.MessagePayloadNormalizerGetInstance()
    _, isValid := normalizer.GetEmailList(`[1, 2]`)
    tests.AssertEquals(t, false, isValid)
    _, isValid = normalizer.GetEmailList(`{"a": "b"}`)
    tests.AssertEquals(t, false, isValid)
}

func TestAttachmentsAcceptObjectsAndLooseStrings(t *testing.T) {
    attachments, isValid := services.MessagePayloadNormalizerGetInstance().GetAttachments(
        `[{"file_name": "nota.txt", "content": "hola"}, {"url": "https://cdn.x.com/docs/factura.pdf?v=2"}, "texto suelto", "https://cdn.x.com/a.png"]`,
    )
    tests.AssertEquals(t, true, isValid)
    tests.AssertEquals(t, 4, len(attachments))
    tests.AssertEquals(t, services.AttachmentInput{FileName: "nota.txt", Content: "hola", HasContent: true}, attachments[0])
    tests.AssertEquals(t, "factura.pdf", attachments[1].FileName)
    tests.AssertEquals(t, true, attachments[1].HasUrl)
    tests.AssertEquals(t, services.AttachmentInput{FileName: "attachment.txt", Content: "texto suelto", HasContent: true}, attachments[2])
    tests.AssertEquals(t, "a.png", attachments[3].FileName)
}

func TestAttachmentsMustBeAList(t *testing.T) {
    _, isValid := services.MessagePayloadNormalizerGetInstance().GetAttachments(`"https://x.com/a"`)
    tests.AssertEquals(t, false, isValid)
}
