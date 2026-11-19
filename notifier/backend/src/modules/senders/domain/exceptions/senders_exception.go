package exceptions

import "notifier-backend/src/modules/shared/domain/enums"

// SendersException error de dominio de senders con su codigo HTTP equivalente (404, 409).
type SendersException struct {
    message string
    code    int
}

func (e *SendersException) Error() string {
    return e.message
}

func (e *SendersException) GetCode() int {
    return e.code
}

func NotFoundCustom(message string) *SendersException {
    return &SendersException{message: message, code: enums.ExceptionEnum.CODE_NOT_FOUND}
}

func ConflictCustom(message string) *SendersException {
    return &SendersException{message: message, code: enums.ExceptionEnum.CODE_CONFLICT}
}
