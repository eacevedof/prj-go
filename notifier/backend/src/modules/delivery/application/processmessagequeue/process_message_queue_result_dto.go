package processmessagequeue

import "notifier-backend/src/modules/shared/application"

// ProcessMessageQueuePrimitives lo que imprime ./console message-queue-process (claves en este orden).
type ProcessMessageQueuePrimitives struct {
    Channel           string `json:"channel"`
    SentMessages      int    `json:"sent_messages"`
    RetriedMessages   int    `json:"retried_messages"`
    FailedMessages    int    `json:"failed_messages"`
    PostponedMessages int    `json:"postponed_messages"`
    OrphansRecovered  int    `json:"orphans_recovered"`
}

type ProcessMessageQueueResultDto struct {
    channel           string
    sentMessages      int
    retriedMessages   int
    failedMessages    int
    postponedMessages int
    orphansRecovered  int
    dueMessages       int
}

func ProcessMessageQueueResultDtoFromPrimitives(primitives map[string]any) *ProcessMessageQueueResultDto {
    return &ProcessMessageQueueResultDto{
        channel:           application.GetTrimmedString(primitives, "channel"),
        sentMessages:      application.GetInt(primitives, "sentMessages"),
        retriedMessages:   application.GetInt(primitives, "retriedMessages"),
        failedMessages:    application.GetInt(primitives, "failedMessages"),
        postponedMessages: application.GetInt(primitives, "postponedMessages"),
        orphansRecovered:  application.GetInt(primitives, "orphansRecovered"),
        dueMessages:       application.GetInt(primitives, "dueMessages"),
    }
}

// GetDueMessages cuantos pendientes con turno se leyeron: si iguala al maximo pedido, puede quedar cola.
func (d *ProcessMessageQueueResultDto) GetDueMessages() int {
    return d.dueMessages
}

// GetHandledMessages los que cambiaron en esta pasada (enviados, reintentados, fallidos o aplazados). Los de un perfil
// deshabilitado no cuentan: siguen en la cola sin cambiar.
func (d *ProcessMessageQueueResultDto) GetHandledMessages() int {
    return d.sentMessages + d.retriedMessages + d.failedMessages + d.postponedMessages
}

func (d *ProcessMessageQueueResultDto) HasFailures() bool {
    return d.failedMessages > 0 || d.retriedMessages > 0
}

func (d *ProcessMessageQueueResultDto) ToArray() ProcessMessageQueuePrimitives {
    return ProcessMessageQueuePrimitives{
        Channel:           d.channel,
        SentMessages:      d.sentMessages,
        RetriedMessages:   d.retriedMessages,
        FailedMessages:    d.failedMessages,
        PostponedMessages: d.postponedMessages,
        OrphansRecovered:  d.orphansRecovered,
    }
}
