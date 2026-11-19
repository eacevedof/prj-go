package exceptions

import "errors"

// DeliveryException un fallo al entregar un mensaje: si es permanente (credenciales rechazadas, destinatario
// invalido, adjunto no permitido) no se reintenta; si es temporal, si.
type DeliveryException struct {
    message     string
    isPermanent bool
}

func (e *DeliveryException) Error() string {
    return e.message
}

func (e *DeliveryException) IsPermanent() bool {
    return e.isPermanent
}

func PermanentCustom(message string) *DeliveryException {
    return &DeliveryException{message: message, isPermanent: true}
}

func TemporaryCustom(message string) *DeliveryException {
    return &DeliveryException{message: message, isPermanent: false}
}

// IsPermanentError cualquier error que no sea un DeliveryException permanente cuenta como temporal: ante la duda,
// se reintenta.
func IsPermanentError(err error) bool {
    var deliveryException *DeliveryException
    return errors.As(err, &deliveryException) && deliveryException.IsPermanent()
}
