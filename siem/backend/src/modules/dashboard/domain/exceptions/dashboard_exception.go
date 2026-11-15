package exceptions

import (
    "siem-backend/src/modules/dashboard/domain/enums"
    sharedenums "siem-backend/src/modules/shared/domain/enums"
)

// DashboardException error de dominio del dashboard con su codigo HTTP equivalente (401, 404).
type DashboardException struct {
    message string
    code    int
}

func (e *DashboardException) Error() string {
    return e.message
}

func (e *DashboardException) GetCode() int {
    return e.code
}

// UnauthorizedInvalidToken mismo 401 para todo fallo de autenticacion del dashboard (sin token, token distinto,
// APP_DASHBOARD_TOKEN sin configurar).
func UnauthorizedInvalidToken() *DashboardException {
    return &DashboardException{
        message: enums.DashboardMessageEnum.INVALID_OR_MISSING_DASHBOARD_TOKEN,
        code:    sharedenums.ExceptionEnum.CODE_UNAUTHORIZED,
    }
}

func NotFoundCustom(message string) *DashboardException {
    return &DashboardException{message: message, code: sharedenums.ExceptionEnum.CODE_NOT_FOUND}
}
