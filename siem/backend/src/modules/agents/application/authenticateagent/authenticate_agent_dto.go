package authenticateagent

import "siem-backend/src/modules/shared/application"

type AuthenticateAgentDto struct {
    application.AbstractInputDto
    apiToken string
}

func AuthenticateAgentDtoFromPrimitives(primitives map[string]any) *AuthenticateAgentDto {
    return &AuthenticateAgentDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        apiToken:         application.GetTrimmedString(primitives, "apiToken"),
    }
}

func (d *AuthenticateAgentDto) GetApiToken() string {
    return d.apiToken
}
