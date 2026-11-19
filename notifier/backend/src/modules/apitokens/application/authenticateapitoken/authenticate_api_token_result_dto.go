package authenticateapitoken

import "notifier-backend/src/modules/shared/application"

// AuthenticateApiTokenResultDto lo que sabe el token de envio de su perfil: con esto se encola sin volver a leerlo.
type AuthenticateApiTokenResultDto struct {
    authApiTokenId  int
    authSenderId    int
    authChannel     string
    authDefaultToJson string
}

func AuthenticateApiTokenResultDtoFromPrimitives(primitives map[string]any) *AuthenticateApiTokenResultDto {
    return &AuthenticateApiTokenResultDto{
        authApiTokenId:  application.GetInt(primitives, "authApiTokenId"),
        authSenderId:    application.GetInt(primitives, "authSenderId"),
        authChannel:     application.GetTrimmedString(primitives, "authChannel"),
        authDefaultToJson: application.GetTrimmedString(primitives, "authDefaultTo"),
    }
}

// ToAuthPrimitives lo que el controller mezcla en los primitivos del dto del caso de uso que viene despues.
func (d *AuthenticateApiTokenResultDto) ToAuthPrimitives() map[string]any {
    return map[string]any{
        "authApiTokenId": d.authApiTokenId,
        "authSenderId":   d.authSenderId,
        "authChannel":    d.authChannel,
        "authDefaultTo":  d.authDefaultToJson,
    }
}
