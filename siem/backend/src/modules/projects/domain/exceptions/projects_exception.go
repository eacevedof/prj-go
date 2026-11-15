package exceptions

import "siem-backend/src/modules/shared/domain/enums"

// ProjectsException error de dominio de projects con su codigo HTTP equivalente (404, 409).
type ProjectsException struct {
    message string
    code    int
}

func (e *ProjectsException) Error() string {
    return e.message
}

func (e *ProjectsException) GetCode() int {
    return e.code
}

func NotFoundCustom(message string) *ProjectsException {
    return &ProjectsException{message: message, code: enums.ExceptionEnum.CODE_NOT_FOUND}
}

func ConflictCustom(message string) *ProjectsException {
    return &ProjectsException{message: message, code: enums.ExceptionEnum.CODE_CONFLICT}
}
