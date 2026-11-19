package controllers

import (
    "notifier-backend/src/modules/shared/application"
    "notifier-backend/src/modules/shared/infrastructure/components/httprequester"
    "notifier-backend/src/modules/shared/infrastructure/components/httpresponder"
    sharedcontrollers "notifier-backend/src/modules/shared/infrastructure/controllers"

    "notifier-backend/src/modules/messages/application/enqueuemessage"
    "notifier-backend/src/modules/messages/application/getmessage"
    "notifier-backend/src/modules/messages/domain/enums"
)

// EnqueueMessageRestController POST /api/v1/messages (token de envio en el body "api_token" o en Authorization: Bearer)
//
//    body: {to, cc, bcc, reply_to, subject, body, body_format, attachments: [{file_name, content|url}], idempotency_key}
//    202: {message_uuid, channel, message_status, is_duplicate}
type EnqueueMessageRestController struct {
    sharedcontrollers.AbstractRestController
}

func EnqueueMessageRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *EnqueueMessageRestController {
    return &EnqueueMessageRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *EnqueueMessageRestController) Invoke() *httpresponder.HttpResponder {
    authApiTokenPrimitives, err := c.GetAuthApiTokenOrFail(c.GetHttpRequest().GetPostString(enums.MessagesKeyEnum.API_TOKEN))
    if err != nil {
        return c.RespondError(err, "EnqueueMessageRestController.Invoke")
    }

    keys := enums.MessagesKeyEnum
    httpRequest := c.GetHttpRequest()
    enqueueMessageResultDto, err := enqueuemessage.EnqueueMessageServiceGetInstance().Invoke(
        enqueuemessage.EnqueueMessageDtoFromPrimitives(application.MergePrimitives(
            httpRequest.GetMetadata(),
            authApiTokenPrimitives,
            map[string]any{
                "toJson":          c.getJsonOrEmpty(keys.TO),
                "ccJson":          c.getJsonOrEmpty(keys.CC),
                "bccJson":         c.getJsonOrEmpty(keys.BCC),
                "replyTo":         httpRequest.GetPostString(keys.REPLY_TO),
                "subject":         httpRequest.GetPostString(keys.SUBJECT),
                "body":            application.ToString(httpRequest.GetPostValue(keys.BODY, "")),
                "bodyFormat":      httpRequest.GetPostString(keys.BODY_FORMAT),
                "attachmentsJson": c.getJsonOrEmpty(keys.ATTACHMENTS),
                "idempotencyKey":  httpRequest.GetPostString(keys.IDEMPOTENCY_KEY),
            },
        )),
    )
    if err != nil {
        return c.RespondError(err, "EnqueueMessageRestController.Invoke")
    }
    return c.RespondAccepted(enqueueMessageResultDto.ToArray())
}

// getJsonOrEmpty el campo del body tal cual (string, lista u objeto) serializado; "" si no viene.
func (c *EnqueueMessageRestController) getJsonOrEmpty(key string) string {
    rawValue := c.GetHttpRequest().GetPostValue(key, nil)
    if rawValue == nil {
        return ""
    }
    return sharedcontrollers.GetJsonEncoded(rawValue)
}

// GetMessageRestController GET /api/v1/messages/:message_uuid (Bearer: el mismo token de envio) -> 200 el mensaje con
// sus adjuntos e intentos.
type GetMessageRestController struct {
    sharedcontrollers.AbstractRestController
}

func GetMessageRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *GetMessageRestController {
    return &GetMessageRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *GetMessageRestController) Invoke(messageUuid string) *httpresponder.HttpResponder {
    authApiTokenPrimitives, err := c.GetAuthApiTokenOrFail("")
    if err != nil {
        return c.RespondError(err, "GetMessageRestController.Invoke")
    }
    getMessageResultDto, err := getmessage.GetMessageServiceGetInstance().Invoke(
        getmessage.GetMessageDtoFromPrimitives(application.MergePrimitives(
            c.GetHttpRequest().GetMetadata(), authApiTokenPrimitives, map[string]any{"messageUuid": messageUuid},
        )),
    )
    if err != nil {
        return c.RespondError(err, "GetMessageRestController.Invoke")
    }
    return c.RespondOk(getMessageResultDto.ToArray())
}
