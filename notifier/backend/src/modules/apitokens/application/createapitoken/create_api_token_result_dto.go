package createapitoken

import (
    "notifier-backend/src/modules/shared/application"

    "notifier-backend/src/modules/apitokens/domain/entities"
)

// CreateApiTokenPrimitives el token creado; api_token (en claro) solo sale aqui.
type CreateApiTokenPrimitives struct {
    Token    entities.ApiTokenPrimitives `json:"token"`
    ApiToken string                      `json:"api_token"`
}

type CreateApiTokenResultDto struct {
    token    entities.ApiTokenPrimitives
    apiToken string
}

func CreateApiTokenResultDtoFromPrimitives(primitives map[string]any) *CreateApiTokenResultDto {
    token, _ := primitives["token"].(entities.ApiTokenPrimitives)
    return &CreateApiTokenResultDto{token: token, apiToken: application.GetTrimmedString(primitives, "apiToken")}
}

func (d *CreateApiTokenResultDto) ToArray() CreateApiTokenPrimitives {
    return CreateApiTokenPrimitives{Token: d.token, ApiToken: d.apiToken}
}
