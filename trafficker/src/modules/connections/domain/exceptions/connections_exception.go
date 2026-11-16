package exceptions

import "trafficker/src/modules/shared/domain/enums"

// ConnectionsException no se pueden leer los sockets del sistema (/proc/net no existe o no se puede leer): 503.
type ConnectionsException struct {
    message string
    code    int
}

func (e *ConnectionsException) Error() string {
    return e.message
}

func (e *ConnectionsException) GetCode() int {
    return e.code
}

func ServiceUnavailableCustom(message string) *ConnectionsException {
    return &ConnectionsException{message: message, code: enums.ExceptionEnum.CODE_SERVICE_UNAVAILABLE}
}
