package disableapitoken

import "notifier-backend/src/modules/apitokens/domain/entities"

type DisableApiTokenResultDto struct {
    token entities.ApiTokenPrimitives
}

func DisableApiTokenResultDtoFromPrimitives(primitives map[string]any) *DisableApiTokenResultDto {
    token, _ := primitives["token"].(entities.ApiTokenPrimitives)
    return &DisableApiTokenResultDto{token: token}
}

func (d *DisableApiTokenResultDto) ToArray() entities.ApiTokenPrimitives {
    return d.token
}
