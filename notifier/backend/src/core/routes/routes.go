package routes

import (
    "notifier-backend/src/core/boot"
    "notifier-backend/src/modules/shared/domain/enums"
    "notifier-backend/src/modules/shared/infrastructure/components/httprequester"
    "notifier-backend/src/modules/shared/infrastructure/components/httpresponder"
    sharedcontrollers "notifier-backend/src/modules/shared/infrastructure/controllers"

    messagescontrollers "notifier-backend/src/modules/messages/infrastructure/controllers"
    userscontrollers "notifier-backend/src/modules/users/infrastructure/controllers"
)

// GetRoutes rutas de la API (ApiRouter). Url con tags ":nombre" -> args con ese nombre para el controller.
// Allowed: metodos HTTP en minusculas (HEAD va implicito en GET). Una misma url puede repetirse con otro metodo.
// Auth: no lo decide el router. Cada controller autentica con GetAuthApiTokenOrFail() (token de envio: encolar y
// consultar un mensaje) o GetAuthUserOrFail() (Bearer del panel, fase 2).
func GetRoutes() []boot.Route {
    return []boot.Route{
        {Url: "/health", Allowed: []string{enums.HttpEnum.METHOD_GET}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return sharedcontrollers.HealthControllerGetInstance(httpRequester).Invoke()
        }},

        //@cliente: envio con su token
        {Url: "/api/v1/messages", Allowed: []string{enums.HttpEnum.METHOD_POST}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return messagescontrollers.EnqueueMessageRestControllerGetInstance(httpRequester).Invoke()
        }},
        {Url: "/api/v1/messages/:message_uuid", Allowed: []string{enums.HttpEnum.METHOD_GET}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return messagescontrollers.GetMessageRestControllerGetInstance(httpRequester).Invoke(args["message_uuid"])
        }},

        //@panel: acceso
        {Url: "/api/v1/log-in", Allowed: []string{enums.HttpEnum.METHOD_POST}, Controller: func(httpRequester *httprequester.HttpRequester, args map[string]string) *httpresponder.HttpResponder {
            return userscontrollers.LogInUserRestControllerGetInstance(httpRequester).Invoke()
        }},
    }
}
