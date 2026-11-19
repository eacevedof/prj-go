package exceptions

import "notifier-backend/src/modules/shared/domain/enums"

// MessagesException error de dominio de messages con su codigo HTTP equivalente (400, 404).
type MessagesException struct {
    message string
    code    int
}

func (e *MessagesException) Error() string {
    return e.message
}

func (e *MessagesException) GetCode() int {
    return e.code
}

func BadRequestCustom(message string) *MessagesException {
    return &MessagesException{message: message, code: enums.ExceptionEnum.CODE_BAD_REQUEST}
}

func NotFoundCustom(message string) *MessagesException {
    return &MessagesException{message: message, code: enums.ExceptionEnum.CODE_NOT_FOUND}
}

func UnauthorizedCustom(message string) *MessagesException {
    return &MessagesException{message: message, code: enums.ExceptionEnum.CODE_UNAUTHORIZED}
}
