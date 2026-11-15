package exceptions

import "siem-client/src/modules/shared/domain/enums"

// AgentConfigException la configuracion del agente (siem-client.json o .env) no vale: el agente no arranca.
type AgentConfigException struct {
    message string
    code    int
}

func (e *AgentConfigException) Error() string {
    return e.message
}

func (e *AgentConfigException) GetCode() int {
    return e.code
}

func NotFoundCustom(message string) *AgentConfigException {
    return &AgentConfigException{message: message, code: enums.ExceptionEnum.CODE_NOT_FOUND}
}

func BadRequestCustom(message string) *AgentConfigException {
    return &AgentConfigException{message: message, code: enums.ExceptionEnum.CODE_BAD_REQUEST}
}

// UnauthorizedCustom el backend rechaza el token al pedir la config remota.
func UnauthorizedCustom(message string) *AgentConfigException {
    return &AgentConfigException{message: message, code: enums.ExceptionEnum.CODE_UNAUTHORIZED}
}

// ServiceUnavailableCustom la config remota no se puede bajar ahora (backend caido o respuesta inesperada): el agente
// sigue con la ultima guardada.
func ServiceUnavailableCustom(message string) *AgentConfigException {
    return &AgentConfigException{message: message, code: enums.ExceptionEnum.CODE_SERVICE_UNAVAILABLE}
}
