package authenticateapitoken

import "notifier-backend/src/modules/shared/application"

type AuthenticateApiTokenDto struct {
    application.AbstractInputDto
    apiToken string
}

func AuthenticateApiTokenDtoFromPrimitives(primitives map[string]any) *AuthenticateApiTokenDto {
    return &AuthenticateApiTokenDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        apiToken:         application.GetTrimmedString(primitives, "apiToken"),
    }
}

func (d *AuthenticateApiTokenDto) GetApiToken() string {
    return d.apiToken
}
