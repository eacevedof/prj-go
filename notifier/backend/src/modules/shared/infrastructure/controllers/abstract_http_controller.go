package controllers

import (
    "bytes"
    "encoding/json"

    "notifier-backend/src/modules/apitokens/application/authenticateapitoken"
    "notifier-backend/src/modules/shared/application"
    "notifier-backend/src/modules/shared/infrastructure/components/httprequester"
    "notifier-backend/src/modules/shared/infrastructure/traits"
    "notifier-backend/src/modules/users/application/authenticateuser"
)

// AbstractHttpController base de todo lo que entra por HTTP: la peticion (HttpRequester), log + mail y las dos
// autenticaciones (el panel con el Bearer de su sesion; un cliente que envia con su token de envio).
// AbstractRestController pone encima el envelope JSON y el mapeo de errores.
type AbstractHttpController struct {
    traits.LogTrait
    traits.MailTrait
    httpRequester *httprequester.HttpRequester
}

func NewAbstractHttpController(httpRequester *httprequester.HttpRequester) AbstractHttpController {
    logTrait := traits.NewLogTrait(httpRequester)
    return AbstractHttpController{
        LogTrait:      logTrait,
        MailTrait:     traits.NewMailTrait(logTrait),
        httpRequester: httpRequester,
    }
}

func (c *AbstractHttpController) GetHttpRequest() *httprequester.HttpRequester {
    return c.httpRequester
}

// GetAuthUserOrFail el Bearer del panel -> authUserId + authUserUuid + authUserRole. Se llama DENTRO del controller: el
// 401 (UsersException) lo mapea su rama de error como cualquier 4xx.
func (c *AbstractHttpController) GetAuthUserOrFail() (map[string]any, error) {
    authenticateUserResultDto, err := authenticateuser.AuthenticateUserServiceGetInstance().Invoke(
        authenticateuser.AuthenticateUserDtoFromPrimitives(application.MergePrimitives(
            map[string]any{"apiToken": c.httpRequester.GetBearerToken()},
            c.httpRequester.GetMetadata(),
        )),
    )
    if err != nil {
        return nil, err
    }
    return authenticateUserResultDto.ToAuthPrimitives(), nil
}

// GetAuthApiTokenOrFail un token de envio -> authApiTokenId + authSenderId + authChannel + authDefaultTo (401
// ApiTokensException si no vale). El token viaja en el body (bodyApiToken, "api_token") o en Authorization: Bearer.
func (c *AbstractHttpController) GetAuthApiTokenOrFail(bodyApiToken string) (map[string]any, error) {
    apiToken := bodyApiToken
    if apiToken == "" {
        apiToken = c.httpRequester.GetBearerToken()
    }
    authenticateApiTokenResultDto, err := authenticateapitoken.AuthenticateApiTokenServiceGetInstance().Invoke(
        authenticateapitoken.AuthenticateApiTokenDtoFromPrimitives(application.MergePrimitives(
            map[string]any{"apiToken": apiToken},
            c.httpRequester.GetMetadata(),
        )),
    )
    if err != nil {
        return nil, err
    }
    return authenticateApiTokenResultDto.ToAuthPrimitives(), nil
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
