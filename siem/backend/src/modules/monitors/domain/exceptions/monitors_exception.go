package exceptions

import "siem-backend/src/modules/shared/domain/enums"

// MonitorsException fallos controlados de los monitores: no existe (404) o nombre repetido en el proyecto (409).
type MonitorsException struct {
    message string
    code    int
}

func (e *MonitorsException) Error() string {
    return e.message
}

func (e *MonitorsException) GetCode() int {
    return e.code
}

func NotFoundCustom(message string) *MonitorsException {
    return &MonitorsException{message: message, code: enums.ExceptionEnum.CODE_NOT_FOUND}
}

func ConflictCustom(message string) *MonitorsException {
    return &MonitorsException{message: message, code: enums.ExceptionEnum.CODE_CONFLICT}
}
