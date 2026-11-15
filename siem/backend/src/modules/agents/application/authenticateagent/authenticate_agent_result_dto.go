package authenticateagent

import "siem-backend/src/modules/shared/application"

type AuthenticateAgentResultDto struct {
    authAgentId   int
    authAgentUuid string
    authProjectId int
}

func AuthenticateAgentResultDtoFromPrimitives(primitives map[string]any) *AuthenticateAgentResultDto {
    return &AuthenticateAgentResultDto{
        authAgentId:   application.GetInt(primitives, "authAgentId"),
        authAgentUuid: application.GetTrimmedString(primitives, "authAgentUuid"),
        authProjectId: application.GetInt(primitives, "authProjectId"),
    }
}

func (d *AuthenticateAgentResultDto) GetAuthAgentId() int {
    return d.authAgentId
}

func (d *AuthenticateAgentResultDto) GetAuthAgentUuid() string {
    return d.authAgentUuid
}

func (d *AuthenticateAgentResultDto) GetAuthProjectId() int {
    return d.authProjectId
}

// ToAuthPrimitives lo que el controller mezcla en los primitivos del dto del caso de uso que viene despues.
func (d *AuthenticateAgentResultDto) ToAuthPrimitives() map[string]any {
    return map[string]any{
        "authAgentId":   d.authAgentId,
        "authAgentUuid": d.authAgentUuid,
        "authProjectId": d.authProjectId,
    }
}
