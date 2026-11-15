package exceptions

import "siem-client/src/modules/shared/domain/enums"

// DockerException no se puede hablar con Docker ahora (socket no montado, daemon parado, respuesta inesperada): 503.
// El agente lo avisa una vez y sigue con lo demas (ficheros, comandos, envio).
type DockerException struct {
    message string
    code    int
}

func (e *DockerException) Error() string {
    return e.message
}

func (e *DockerException) GetCode() int {
    return e.code
}

func ServiceUnavailableCustom(message string) *DockerException {
    return &DockerException{message: message, code: enums.ExceptionEnum.CODE_SERVICE_UNAVAILABLE}
}
