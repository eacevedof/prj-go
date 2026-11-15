package exceptions

import "siem-backend/src/modules/shared/domain/enums"

// AgentsException error de dominio de agents con su codigo HTTP equivalente (401, 404, 409).
type AgentsException struct {
    message string
    code    int
}

func (e *AgentsException) Error() string {
    return e.message
}

func (e *AgentsException) GetCode() int {
    return e.code
}

func UnauthorizedCustom(message string) *AgentsException {
    return &AgentsException{message: message, code: enums.ExceptionEnum.CODE_UNAUTHORIZED}
}

func NotFoundCustom(message string) *AgentsException {
    return &AgentsException{message: message, code: enums.ExceptionEnum.CODE_NOT_FOUND}
}

func ConflictCustom(message string) *AgentsException {
    return &AgentsException{message: message, code: enums.ExceptionEnum.CODE_CONFLICT}
}
