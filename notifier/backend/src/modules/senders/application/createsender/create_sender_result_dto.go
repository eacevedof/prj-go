package createsender

import (
    "notifier-backend/src/modules/shared/application"

    apitokensentities "notifier-backend/src/modules/apitokens/domain/entities"

    "notifier-backend/src/modules/senders/domain/entities"
)

// CreateSenderPrimitives el perfil creado y su primer token de envio. api_token solo sale aqui: no se puede volver a
// ver.
type CreateSenderPrimitives struct {
    Sender   entities.SenderPrimitives            `json:"sender"`
    Token    apitokensentities.ApiTokenPrimitives `json:"token"`
    ApiToken string                               `json:"api_token"`
}

type CreateSenderResultDto struct {
    sender   entities.SenderPrimitives
    token    apitokensentities.ApiTokenPrimitives
    apiToken string
}

func CreateSenderResultDtoFromPrimitives(primitives map[string]any) *CreateSenderResultDto {
    sender, _ := primitives["sender"].(entities.SenderPrimitives)
    token, _ := primitives["token"].(apitokensentities.ApiTokenPrimitives)
    return &CreateSenderResultDto{sender: sender, token: token, apiToken: application.GetTrimmedString(primitives, "apiToken")}
}

func (d *CreateSenderResultDto) ToArray() CreateSenderPrimitives {
    return CreateSenderPrimitives{Sender: d.sender, Token: d.token, ApiToken: d.apiToken}
}
