package authenticateuser

import "notifier-backend/src/modules/shared/application"

type AuthenticateUserResultDto struct {
    authUserId   int
    authUserUuid string
    authUserRole string
}

func AuthenticateUserResultDtoFromPrimitives(primitives map[string]any) *AuthenticateUserResultDto {
    return &AuthenticateUserResultDto{
        authUserId:   application.GetInt(primitives, "authUserId"),
        authUserUuid: application.GetTrimmedString(primitives, "authUserUuid"),
        authUserRole: application.GetTrimmedString(primitives, "authUserRole"),
    }
}

// ToAuthPrimitives lo que el controller mezcla en los primitivos del dto del caso de uso que viene despues.
func (d *AuthenticateUserResultDto) ToAuthPrimitives() map[string]any {
    return map[string]any{
        "authUserId":   d.authUserId,
        "authUserUuid": d.authUserUuid,
        "authUserRole": d.authUserRole,
    }
}
