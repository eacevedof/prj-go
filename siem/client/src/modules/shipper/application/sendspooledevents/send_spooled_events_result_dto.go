package sendspooledevents

import "siem-client/src/modules/shared/application"

// SendSpooledEventsPrimitives lo que imprime events-send (claves en este orden).
type SendSpooledEventsPrimitives struct {
    SentBatches     int  `json:"sent_batches"`
    SentEvents      int  `json:"sent_events"`
    RejectedBatches int  `json:"rejected_batches"`
    DroppedBatches  int  `json:"dropped_batches"`
    PendingBatches  int  `json:"pending_batches"`
    IsHeartbeatSent bool `json:"is_heartbeat_sent"`
}

type SendSpooledEventsResultDto struct {
    sentBatches     int
    sentEvents      int
    rejectedBatches int
    droppedBatches  int
    pendingBatches  int
    isHeartbeatSent bool
}

func SendSpooledEventsResultDtoFromPrimitives(primitives map[string]any) *SendSpooledEventsResultDto {
    return &SendSpooledEventsResultDto{
        sentBatches:     application.GetInt(primitives, "sentBatches"),
        sentEvents:      application.GetInt(primitives, "sentEvents"),
        rejectedBatches: application.GetInt(primitives, "rejectedBatches"),
        droppedBatches:  application.GetInt(primitives, "droppedBatches"),
        pendingBatches:  application.GetInt(primitives, "pendingBatches"),
        isHeartbeatSent: application.GetInt(primitives, "isHeartbeatSent") == 1,
    }
}

func (d *SendSpooledEventsResultDto) GetSentBatches() int {
    return d.sentBatches
}

func (d *SendSpooledEventsResultDto) GetSentEvents() int {
    return d.sentEvents
}

func (d *SendSpooledEventsResultDto) GetRejectedBatches() int {
    return d.rejectedBatches
}

func (d *SendSpooledEventsResultDto) GetDroppedBatches() int {
    return d.droppedBatches
}

func (d *SendSpooledEventsResultDto) GetPendingBatches() int {
    return d.pendingBatches
}

func (d *SendSpooledEventsResultDto) ToArray() SendSpooledEventsPrimitives {
    return SendSpooledEventsPrimitives{
        SentBatches:     d.sentBatches,
        SentEvents:      d.sentEvents,
        RejectedBatches: d.rejectedBatches,
        DroppedBatches:  d.droppedBatches,
        PendingBatches:  d.pendingBatches,
        IsHeartbeatSent: d.isHeartbeatSent,
    }
}
