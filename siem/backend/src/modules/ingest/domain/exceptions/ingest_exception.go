package exceptions

import "siem-backend/src/modules/shared/domain/enums"

// IngestException error de dominio de ingest con su codigo HTTP equivalente (400, 401).
type IngestException struct {
    message string
    code    int
}

func (e *IngestException) Error() string {
    return e.message
}

func (e *IngestException) GetCode() int {
    return e.code
}

func UnauthorizedCustom(message string) *IngestException {
    return &IngestException{message: message, code: enums.ExceptionEnum.CODE_UNAUTHORIZED}
}

func BadRequestCustom(message string) *IngestException {
    return &IngestException{message: message, code: enums.ExceptionEnum.CODE_BAD_REQUEST}
}
