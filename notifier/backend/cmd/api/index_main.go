package main

import (
    "fmt"
    "net/http"
    "runtime/debug"
    "strings"

    "notifier-backend/src/core/boot"
    "notifier-backend/src/core/config"
    "notifier-backend/src/modules/shared/domain/enums"
    "notifier-backend/src/modules/shared/infrastructure/components/httprequester"
    "notifier-backend/src/modules/shared/infrastructure/components/httpresponder"
    "notifier-backend/src/modules/shared/infrastructure/components/uuid"
    "notifier-backend/src/modules/shared/infrastructure/repositories/configuration"
    "notifier-backend/src/modules/shared/infrastructure/controllers"
)

const maxRequestBodyBytes = 8 << 20

// IndexMain front controller de la API (el public/index.php de task-php). Sin sesion ni vistas: solo JSON.
type IndexMain struct {
    routes []boot.Route
}

func IndexMainGetInstance(routes []boot.Route) *IndexMain {
    return &IndexMain{routes: routes}
}

func (m *IndexMain) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
    // un mensaje con adjuntos de texto cabe en 8 MB (los grandes van por url): nadie puede llenar la memoria con un body enorme
    request.Body = http.MaxBytesReader(responseWriter, request.Body, maxRequestBodyBytes)
    httpRequester := httprequester.GetInstance(request)
    defer func() {
        if recovered := recover(); recovered != nil {
            err := fmt.Errorf("%v", recovered)
            m.logDebugError(httpRequester, err)
            m.getUnexpectedErrorResponse(httpRequester, err).EchoHeadersAndBody(responseWriter)
        }
    }()

    m.exec(httpRequester).EchoHeadersAndBody(responseWriter)
}

func (m *IndexMain) exec(httpRequester *httprequester.HttpRequester) *httpresponder.HttpResponder {
    requestMethod := httpRequester.GetRequestMethod()
    if requestMethod == enums.HttpEnum.METHOD_OPTIONS {
        return m.getWithCommonHeaders(httpRequester, httpresponder.NoContent())
    }

    runData := boot.ApiRouterGetInstance(m.routes).GetRunData(requestMethod, httpRequester.GetRequestUri())
    if runData.Route == nil && len(runData.Allowed) == 0 {
        return m.getWithCommonHeaders(httpRequester, m.getErrorResponse(
            enums.ExceptionEnum.CODE_NOT_FOUND,
            config.Tr(enums.SharedMessageEnum.RESOURCE_NOT_FOUND),
            map[string]string{},
        ))
    }

    if runData.Route == nil {
        return m.getWithCommonHeaders(httpRequester, m.getErrorResponse(
            enums.ExceptionEnum.CODE_METHOD_NOT_ALLOWED,
            config.Tr(enums.SharedMessageEnum.METHOD_NOT_ALLOWED),
            map[string]string{enums.HttpEnum.HEADER_ALLOW: strings.ToUpper(strings.Join(runData.Allowed, ", "))},
        ))
    }

    httpRequester.LoadJsonBodyIntoPost()
    return m.getWithCommonHeaders(httpRequester, runData.Route.Controller(httpRequester, runData.Args))
}

// getWithCommonHeaders toda respuesta sale por aqui: las cabeceras comunes se suman a las del HttpResponder y se
// emiten juntas en EchoHeadersAndBody(), el unico punto que escribe la respuesta.
func (m *IndexMain) getWithCommonHeaders(httpRequester *httprequester.HttpRequester, httpResponder *httpresponder.HttpResponder) *httpresponder.HttpResponder {
    commonHeaders := m.getCorsHeaders(httpRequester)
    for key, value := range m.getSecurityHeaders() {
        commonHeaders[key] = value
    }
    return httpResponder.GetWithHeaders(commonHeaders)
}

func (m *IndexMain) getCorsHeaders(httpRequester *httprequester.HttpRequester) map[string]string {
    corsHeaders := map[string]string{}
    if origin := httpRequester.GetHeaderValueByKey(enums.HttpEnum.HEADER_ORIGIN); origin != "" {
        corsHeaders[enums.HttpEnum.HEADER_ACCESS_CONTROL_ALLOW_ORIGIN] = origin
        corsHeaders[enums.HttpEnum.HEADER_ACCESS_CONTROL_MAX_AGE] = "86400"
    }
    if httpRequester.GetRequestMethod() == enums.HttpEnum.METHOD_OPTIONS {
        corsHeaders[enums.HttpEnum.HEADER_ACCESS_CONTROL_ALLOW_METHODS] = "GET, POST, PUT, DELETE, OPTIONS"
        corsHeaders[enums.HttpEnum.HEADER_ACCESS_CONTROL_ALLOW_HEADERS] = "Content-Type, Accept, Accept-Language, Authorization, X-Time-Zone"
    }
    return corsHeaders
}

func (m *IndexMain) getSecurityHeaders() map[string]string {
    return map[string]string{
        enums.HttpEnum.HEADER_X_CONTENT_TYPE_OPTIONS: "nosniff",
        enums.HttpEnum.HEADER_REFERRER_POLICY:        "no-referrer",
    }
}

// getUnexpectedErrorResponse el panic que nadie ha tratado: se loguea con un id y el cliente recibe ese id en un 500
// generico.
func (m *IndexMain) getUnexpectedErrorResponse(httpRequester *httprequester.HttpRequester, err error) *httpresponder.HttpResponder {
    errorUuid := uuid.GetInstance().GetUniqueIdWithPrefix("")
    sensitiveDataMasker := boot.SensitiveDataMaskerGetInstance()
    if post := httpRequester.GetPost(); len(post) > 0 {
        config.Lgerr(sensitiveDataMasker.GetMaskedArray(post), "index-exception "+errorUuid+" POST", "error")
    }
    if get := httpRequester.GetGet(); len(get) > 0 {
        config.Lgerr(sensitiveDataMasker.GetMaskedArray(get), "index-exception "+errorUuid+" GET", "error")
    }
    config.Lgerr(err.Error(), "index-exception "+errorUuid, "error")
    config.Lgerr(string(debug.Stack()), "index-exception "+errorUuid+" TRACE", "error")

    return m.getWithCommonHeaders(httpRequester, m.getErrorResponse(
        enums.ExceptionEnum.CODE_INTERNAL_SERVER_ERROR,
        config.Tr(enums.SharedMessageEnum.UNEXPECTED_ERROR_WITH_ID, errorUuid),
        map[string]string{},
    ))
}

func (m *IndexMain) logDebugError(httpRequester *httprequester.HttpRequester, err error) {
    environmentReaderRawRepository := configuration.GetInstance()
    if environmentReaderRawRepository.IsProduction() || !environmentReaderRawRepository.IsDebug() {
        return
    }

    config.Lgerr(map[string]any{
        "Exception": err.Error(),
        "POST":      httpRequester.GetPost(),
        "GET":       httpRequester.GetGet(),
    }, "index-debug", "debug")
}

func (m *IndexMain) getErrorResponse(code int, errorMessage string, headers map[string]string) *httpresponder.HttpResponder {
    return httpresponder.Json(
        controllers.GetJsonEncoded(controllers.JsonEnvelope{
            Code:   code,
            Status: false,
            Errors: []string{errorMessage},
            Data:   []any{},
        }),
        code,
        headers,
    )
}
