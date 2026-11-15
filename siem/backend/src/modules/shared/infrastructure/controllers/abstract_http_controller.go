package controllers

import (
    "bytes"
    "encoding/json"

    "siem-backend/src/modules/agents/application/authenticateagent"
    "siem-backend/src/modules/dashboard/application/authenticatedashboard"
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/components/httprequester"
    "siem-backend/src/modules/shared/infrastructure/traits"
)

// AbstractHttpController base de todo lo que entra por HTTP: la peticion (HttpRequester), log + mail y las dos
// autenticaciones por Bearer (un siem-client con el token de su agente; el front con el token del dashboard).
// AbstractRestController pone encima el envelope JSON y el mapeo de errores.
type AbstractHttpController struct {
    traits.LogTrait
    traits.MailTrait
    httpRequest *httprequester.HttpRequester
}

func NewAbstractHttpController(httpRequester *httprequester.HttpRequester) AbstractHttpController {
    logTrait := traits.NewLogTrait(httpRequester)
    return AbstractHttpController{
        LogTrait:    logTrait,
        MailTrait:   traits.NewMailTrait(logTrait),
        httpRequest: httpRequester,
    }
}

func (c *AbstractHttpController) GetHttpRequest() *httprequester.HttpRequester {
    return c.httpRequest
}

// GetAuthAgentOrFail se llama DENTRO del controller: el 401 (AgentsException) lo mapea su rama de error como
// cualquier 4xx. El token Bearer de la peticion -> authAgentId + authAgentUuid + authProjectId.
func (c *AbstractHttpController) GetAuthAgentOrFail() (map[string]any, error) {
    authenticateAgentResultDto, err := authenticateagent.AuthenticateAgentServiceGetInstance().Invoke(
        authenticateagent.AuthenticateAgentDtoFromPrimitives(application.MergePrimitives(
            map[string]any{"apiToken": c.httpRequest.GetBearerToken()},
            c.httpRequest.GetMetadata(),
        )),
    )
    if err != nil {
        return nil, err
    }
    return authenticateAgentResultDto.ToAuthPrimitives(), nil
}

// GetDashboardAuthOrFail el Bearer del front -> isDashboardAuthenticated (401 DashboardException si no vale).
func (c *AbstractHttpController) GetDashboardAuthOrFail() (map[string]any, error) {
    authenticateDashboardResultDto, err := authenticatedashboard.AuthenticateDashboardServiceGetInstance().Invoke(
        authenticatedashboard.AuthenticateDashboardDtoFromPrimitives(application.MergePrimitives(
            map[string]any{"dashboardToken": c.httpRequest.GetBearerToken()},
            c.httpRequest.GetMetadata(),
        )),
    )
    if err != nil {
        return nil, err
    }
    return authenticateDashboardResultDto.ToAuthPrimitives(), nil
}

// GetJsonEncoded json_encode con JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES: sin escapar <, > ni &.
func GetJsonEncoded(value any) string {
    var jsonBuffer bytes.Buffer
    jsonEncoder := json.NewEncoder(&jsonBuffer)
    jsonEncoder.SetEscapeHTML(false)
    if err := jsonEncoder.Encode(value); err != nil {
        return ""
    }
    return string(bytes.TrimRight(jsonBuffer.Bytes(), "\n"))
}
