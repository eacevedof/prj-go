package getmessage

import "notifier-backend/src/modules/messages/domain/entities"

type GetMessageResultDto struct {
    message entities.MessagePrimitives
}

func GetMessageResultDtoFromPrimitives(primitives map[string]any) *GetMessageResultDto {
    message, _ := primitives["message"].(entities.MessagePrimitives)
    return &GetMessageResultDto{message: message}
}

func (d *GetMessageResultDto) ToArray() entities.MessagePrimitives {
    return d.message
}
