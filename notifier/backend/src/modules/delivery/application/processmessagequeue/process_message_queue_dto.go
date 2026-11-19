package processmessagequeue

import (
    "notifier-backend/src/modules/shared/application"

    "notifier-backend/src/modules/delivery/domain/enums"
)

type ProcessMessageQueueDto struct {
    application.AbstractInputDto
    channel     string
    maxMessages int
}

func ProcessMessageQueueDtoFromPrimitives(primitives map[string]any) *ProcessMessageQueueDto {
    maxMessages := enums.DeliveryLimitEnum.MAX_MESSAGES_DEFAULT
    // "" (opcion de consola sin pasar) cuenta como ausente
    if application.GetTrimmedString(primitives, "maxMessages") != "" {
        maxMessages = application.GetInt(primitives, "maxMessages")
    }
    return &ProcessMessageQueueDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        channel:          application.GetTrimmedString(primitives, "channel"),
        maxMessages:      maxMessages,
    }
}

func (d *ProcessMessageQueueDto) GetChannel() string {
    return d.channel
}

func (d *ProcessMessageQueueDto) GetMaxMessages() int {
    return d.maxMessages
}
