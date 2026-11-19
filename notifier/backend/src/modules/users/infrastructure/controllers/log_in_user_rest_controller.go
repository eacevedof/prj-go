package controllers

import (
    "notifier-backend/src/modules/shared/application"
    "notifier-backend/src/modules/shared/infrastructure/components/httprequester"
    "notifier-backend/src/modules/shared/infrastructure/components/httpresponder"
    sharedcontrollers "notifier-backend/src/modules/shared/infrastructure/controllers"

    "notifier-backend/src/modules/users/application/loginuser"
    "notifier-backend/src/modules/users/domain/enums"
)

// LogInUserRestController POST /api/v1/log-in {user_email, password} -> 200 {user, api_token}. El api_token es el
// Bearer del panel (fase 2); no sirve para enviar mensajes.
type LogInUserRestController struct {
    sharedcontrollers.AbstractRestController
}

func LogInUserRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *LogInUserRestController {
    return &LogInUserRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *LogInUserRestController) Invoke() *httpresponder.HttpResponder {
    httpRequest := c.GetHttpRequest()
    logInUserResultDto, err := loginuser.LogInUserServiceGetInstance().Invoke(
        loginuser.LogInUserDtoFromPrimitives(application.MergePrimitives(
            httpRequest.GetMetadata(),
            map[string]any{
                "userEmail": httpRequest.GetPostString(enums.UsersKeyEnum.USER_EMAIL),
                "password":  httpRequest.GetPostString(enums.UsersKeyEnum.PASSWORD),
            },
        )),
    )
    if err != nil {
        return c.RespondError(err, "LogInUserRestController.Invoke")
    }
    return c.RespondOk(logInUserResultDto.ToArray())
}
