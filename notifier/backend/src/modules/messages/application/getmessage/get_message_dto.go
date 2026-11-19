package getmessage

import "notifier-backend/src/modules/shared/application"

type GetMessageDto struct {
    application.AbstractInputDto
    authApiTokenId int
    messageUuid    string
}

func GetMessageDtoFromPrimitives(primitives map[string]any) *GetMessageDto {
    return &GetMessageDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        authApiTokenId:   application.GetInt(primitives, "authApiTokenId"),
        messageUuid:      application.GetTrimmedString(primitives, "messageUuid"),
    }
}

func (d *GetMessageDto) GetAuthApiTokenId() int {
    return d.authApiTokenId
}

func (d *GetMessageDto) GetMessageUuid() string {
    return d.messageUuid
}
