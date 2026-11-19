package exceptions

import "notifier-backend/src/modules/shared/domain/enums"

// ApiTokensException error de dominio de apitokens con su codigo HTTP equivalente (401, 404, 409).
type ApiTokensException struct {
    message string
    code    int
}

func (e *ApiTokensException) Error() string {
    return e.message
}

func (e *ApiTokensException) GetCode() int {
    return e.code
}

func UnauthorizedCustom(message string) *ApiTokensException {
    return &ApiTokensException{message: message, code: enums.ExceptionEnum.CODE_UNAUTHORIZED}
}

func NotFoundCustom(message string) *ApiTokensException {
    return &ApiTokensException{message: message, code: enums.ExceptionEnum.CODE_NOT_FOUND}
}

func ConflictCustom(message string) *ApiTokensException {
    return &ApiTokensException{message: message, code: enums.ExceptionEnum.CODE_CONFLICT}
}
