package listsenders

import "notifier-backend/src/modules/senders/domain/entities"

type ListSendersResultDto struct {
    senders []entities.SenderPrimitives
}

func ListSendersResultDtoFromPrimitives(primitives map[string]any) *ListSendersResultDto {
    senders, _ := primitives["senders"].([]entities.SenderPrimitives)
    if senders == nil {
        senders = []entities.SenderPrimitives{}
    }
    return &ListSendersResultDto{senders: senders}
}

func (d *ListSendersResultDto) ToArray() []entities.SenderPrimitives {
    return d.senders
}
