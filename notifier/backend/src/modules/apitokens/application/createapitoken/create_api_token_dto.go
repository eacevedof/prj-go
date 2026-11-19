package createapitoken

import "notifier-backend/src/modules/shared/application"

type CreateApiTokenDto struct {
    application.AbstractInputDto
    senderUuid   string
    tokenName    string
    defaultToCsv string
}

func CreateApiTokenDtoFromPrimitives(primitives map[string]any) *CreateApiTokenDto {
    return &CreateApiTokenDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        senderUuid:       application.GetTrimmedString(primitives, "senderUuid"),
        tokenName:        application.GetTrimmedString(primitives, "tokenName"),
        defaultToCsv:     application.GetTrimmedString(primitives, "defaultToCsv"),
    }
}

func (d *CreateApiTokenDto) GetSenderUuid() string {
    return d.senderUuid
}

func (d *CreateApiTokenDto) GetTokenName() string {
    return d.tokenName
}

// GetDefaultToCsv destinatarios por defecto separados por comas ("" = ninguno).
func (d *CreateApiTokenDto) GetDefaultToCsv() string {
    return d.defaultToCsv
}
