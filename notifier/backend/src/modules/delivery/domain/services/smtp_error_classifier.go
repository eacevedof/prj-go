package services

import (
    "errors"
    "net/textproto"

    "notifier-backend/src/modules/delivery/domain/exceptions"
)

// SmtpErrorClassifier un error del servidor SMTP es permanente si su codigo es 5xx (535 credenciales, 550 buzon que no
// existe, 552 demasiado grande...): reintentar no lo arregla. Un 4xx (421, 450, 451...) o un fallo de red es temporal.
type SmtpErrorClassifier struct{}

func SmtpErrorClassifierGetInstance() *SmtpErrorClassifier {
    return &SmtpErrorClassifier{}
}

func (c *SmtpErrorClassifier) GetDeliveryException(err error) *exceptions.DeliveryException {
    var deliveryException *exceptions.DeliveryException
    if errors.As(err, &deliveryException) {
        return deliveryException
    }
    var smtpError *textproto.Error
    if errors.As(err, &smtpError) && smtpError.Code >= 500 && smtpError.Code < 600 {
        return exceptions.PermanentCustom(err.Error())
    }
    return exceptions.TemporaryCustom(err.Error())
}
