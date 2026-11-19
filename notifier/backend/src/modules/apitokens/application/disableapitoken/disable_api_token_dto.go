package disableapitoken

import "notifier-backend/src/modules/shared/application"

type DisableApiTokenDto struct {
    application.AbstractInputDto
    tokenUuid string
}

func DisableApiTokenDtoFromPrimitives(primitives map[string]any) *DisableApiTokenDto {
    return &DisableApiTokenDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        tokenUuid:        application.GetTrimmedString(primitives, "tokenUuid"),
    }
}

func (d *DisableApiTokenDto) GetTokenUuid() string {
    return d.tokenUuid
}
