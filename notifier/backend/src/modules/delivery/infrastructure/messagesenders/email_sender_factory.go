package messagesenders

import (
    "notifier-backend/src/core/config"

    sendersenums "notifier-backend/src/modules/senders/domain/enums"

    "notifier-backend/src/modules/delivery/domain/enums"
    "notifier-backend/src/modules/delivery/domain/exceptions"
    "notifier-backend/src/modules/delivery/domain/interfaces"
)

// EmailSenderFactory el InterfaceEmailSender de cada proveedor. Un proveedor sin sender es un fallo permanente del
// mensaje (no deberia pasar: el alta del perfil ya lo impide).
type EmailSenderFactory struct{}

func EmailSenderFactoryGetInstance() *EmailSenderFactory {
    return &EmailSenderFactory{}
}

func (f *EmailSenderFactory) GetEmailSender(providerCode string) (interfaces.InterfaceEmailSender, error) {
    switch providerCode {
    case sendersenums.ProviderEnum.GMAIL, sendersenums.ProviderEnum.SMTP_CUSTOM:
        return SmtpMessageSenderGetInstance(), nil
    }
    return nil, exceptions.PermanentCustom(config.Tr(enums.DeliveryMessageEnum.PROVIDER_X_HAS_NO_MESSAGE_SENDER, providerCode))
}
