package interfaces

import "notifier-backend/src/modules/delivery/domain/entities"

// InterfaceEmailSender un proveedor de email (SMTP hoy; Resend por HTTP en la fase 3). Devuelve lo que contesto el
// proveedor, para el log de intentos. Un error que deba no reintentarse es un DeliveryException permanente.
type InterfaceEmailSender interface {
    Send(outgoingEmail entities.OutgoingEmail, smtpSettings entities.SmtpSettings) (string, error)
}
