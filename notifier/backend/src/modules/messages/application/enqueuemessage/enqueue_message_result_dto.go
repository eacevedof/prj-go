package enqueuemessage

import "notifier-backend/src/modules/shared/application"

// EnqueueMessagePrimitives la respuesta 202 (claves en este orden).
type EnqueueMessagePrimitives struct {
    MessageUuid   string `json:"message_uuid"`
    Channel       string `json:"channel"`
    MessageStatus string `json:"message_status"`
    IsDuplicate   bool   `json:"is_duplicate"`
}

type EnqueueMessageResultDto struct {
    messageUuid   string
    channel       string
    messageStatus string
    isDuplicate   bool
}

func EnqueueMessageResultDtoFromPrimitives(primitives map[string]any) *EnqueueMessageResultDto {
    return &EnqueueMessageResultDto{
        messageUuid:   application.GetTrimmedString(primitives, "messageUuid"),
        channel:       application.GetTrimmedString(primitives, "channel"),
        messageStatus: application.GetTrimmedString(primitives, "messageStatus"),
        isDuplicate:   application.GetInt(primitives, "isDuplicate") == 1,
    }
}

func (d *EnqueueMessageResultDto) ToArray() EnqueueMessagePrimitives {
    return EnqueueMessagePrimitives{
        MessageUuid:   d.messageUuid,
        Channel:       d.channel,
        MessageStatus: d.messageStatus,
        IsDuplicate:   d.isDuplicate,
    }
}
