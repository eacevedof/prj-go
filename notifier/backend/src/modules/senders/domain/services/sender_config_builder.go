package services

import (
    "strings"

    "notifier-backend/src/modules/senders/domain/enums"
)

// SenderConfigBuilder aplica los datos por defecto del proveedor a lo que manda el usuario y monta el JSON de config
// (lo no secreto). En Gmail el host, el puerto y la seguridad son siempre los suyos; el usuario SMTP, si no se da, es
// el propio from_email.
type SenderConfigBuilder struct{}

func SenderConfigBuilderGetInstance() *SenderConfigBuilder {
    return &SenderConfigBuilder{}
}

func (b *SenderConfigBuilder) GetWithProviderDefaults(senderFields SenderFields) SenderFields {
    providerDefinition := enums.ProviderEnum.DEFINITIONS[senderFields.ProviderCode]
    if providerDefinition.SmtpHost != "" {
        senderFields.SmtpHost = providerDefinition.SmtpHost
        senderFields.SmtpPort = providerDefinition.SmtpPort
        senderFields.SmtpSecurity = providerDefinition.SmtpSecurity
    }
    if senderFields.SmtpSecurity == "" {
        senderFields.SmtpSecurity = enums.SmtpSecurityEnum.STARTTLS
    }
    senderFields.SmtpHost = strings.ToLower(senderFields.SmtpHost)
    senderFields.FromEmail = strings.ToLower(senderFields.FromEmail)
    if senderFields.Username == "" && senderFields.Secret != "" {
        senderFields.Username = senderFields.FromEmail
    }
    return senderFields
}

func (b *SenderConfigBuilder) GetEmailConfig(senderFields SenderFields) map[string]any {
    keys := enums.SendersKeyEnum
    return map[string]any{
        keys.FROM_EMAIL:    senderFields.FromEmail,
        keys.FROM_NAME:     senderFields.FromName,
        keys.SMTP_HOST:     senderFields.SmtpHost,
        keys.SMTP_PORT:     senderFields.SmtpPort,
        keys.SMTP_SECURITY: senderFields.SmtpSecurity,
        keys.USERNAME:      senderFields.Username,
    }
}
