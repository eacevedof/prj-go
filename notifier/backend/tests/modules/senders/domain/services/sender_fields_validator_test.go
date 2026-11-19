package services

import (
    "strings"
    "testing"

    "notifier-backend/tests"

    "notifier-backend/src/modules/senders/domain/services"
)

func TestGmailGetsItsServerAndUsernameFromEmail(t *testing.T) {
    senderFields := services.SenderConfigBuilderGetInstance().GetWithProviderDefaults(services.SenderFields{
        SenderName: "avisos", ProviderCode: "gmail", FromEmail: "Yo@Gmail.com", SmtpHost: "otro.com", SmtpPort: 25, Secret: "app pass",
    })
    tests.AssertEquals(t, "smtp.gmail.com", senderFields.SmtpHost)
    tests.AssertEquals(t, 587, senderFields.SmtpPort)
    tests.AssertEquals(t, "starttls", senderFields.SmtpSecurity)
    tests.AssertEquals(t, "yo@gmail.com", senderFields.Username)
    tests.AssertEmpty(t, services.SenderFieldsValidatorGetInstance().GetEmailSenderErrors(senderFields))
}

func TestProviderMustExistAndBeImplemented(t *testing.T) {
    validator := services.SenderFieldsValidatorGetInstance()
    tests.AssertStringContainsString(t, "unknown provider", validator.GetProviderErrors("aol")[0])
    tests.AssertStringContainsString(t, "not available yet", validator.GetProviderErrors("telegram_bot")[0])
    tests.AssertEmpty(t, validator.GetProviderErrors("smtp_custom"))
}

func TestGmailNeedsSecretAndNoneIsOnlyForCustom(t *testing.T) {
    errors := strings.Join(services.SenderFieldsValidatorGetInstance().GetEmailSenderErrors(services.SenderFields{
        SenderName: "avisos", ProviderCode: "gmail", FromEmail: "yo@gmail.com", FromName: "a\nb",
        SmtpHost: "smtp.gmail.com", SmtpPort: 587, SmtpSecurity: "none",
    }), " | ")
    tests.AssertStringContainsString(t, "secret: field is required", errors)
    tests.AssertStringContainsString(t, "none is only allowed", errors)
    tests.AssertStringContainsString(t, "from_name", errors)
}
