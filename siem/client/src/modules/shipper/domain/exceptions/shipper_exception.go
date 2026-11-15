package exceptions

import "siem-client/src/modules/shared/domain/enums"

// ShipperException el envio no se puede hacer ahora: token rechazado (401) o backend caido / respuesta inesperada
// (503). Los lotes siguen en el spool y se reintentan en la siguiente pasada.
type ShipperException struct {
    message string
    code    int
}

func (e *ShipperException) Error() string {
    return e.message
}

func (e *ShipperException) GetCode() int {
    return e.code
}

func UnauthorizedCustom(message string) *ShipperException {
    return &ShipperException{message: message, code: enums.ExceptionEnum.CODE_UNAUTHORIZED}
}

func ServiceUnavailableCustom(message string) *ShipperException {
    return &ShipperException{message: message, code: enums.ExceptionEnum.CODE_SERVICE_UNAVAILABLE}
}
