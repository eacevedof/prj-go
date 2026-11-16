package enums

// ExceptionEnum codigos HTTP equivalentes de los errores (misma lectura que task-go aunque aqui no haya API).
var ExceptionEnum = struct {
    CODE_OK                    int
    CODE_BAD_REQUEST           int
    CODE_UNAUTHORIZED          int
    CODE_NOT_FOUND             int
    CODE_CONFLICT              int
    CODE_INTERNAL_SERVER_ERROR int
    CODE_SERVICE_UNAVAILABLE   int
}{
    CODE_OK:                    200,
    CODE_BAD_REQUEST:           400,
    CODE_UNAUTHORIZED:          401,
    CODE_NOT_FOUND:             404,
    CODE_CONFLICT:              409,
    CODE_INTERNAL_SERVER_ERROR: 500,
    CODE_SERVICE_UNAVAILABLE:   503,
}
