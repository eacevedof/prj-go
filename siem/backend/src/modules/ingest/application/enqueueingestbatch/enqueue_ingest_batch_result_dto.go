package enqueueingestbatch

import "siem-backend/src/modules/shared/application"

// EnqueueIngestBatchPrimitives la respuesta al siem-client (claves en este orden).
type EnqueueIngestBatchPrimitives struct {
    BatchUuid    string `json:"batch_uuid"`
    QueuedEvents int    `json:"queued_events"`
    IsDuplicate  bool   `json:"is_duplicate"`
}

type EnqueueIngestBatchResultDto struct {
    batchUuid    string
    queuedEvents int
    isDuplicate  bool
}

func EnqueueIngestBatchResultDtoFromPrimitives(primitives map[string]any) *EnqueueIngestBatchResultDto {
    return &EnqueueIngestBatchResultDto{
        batchUuid:    application.GetTrimmedString(primitives, "batchUuid"),
        queuedEvents: application.GetInt(primitives, "queuedEvents"),
        isDuplicate:  application.GetInt(primitives, "isDuplicate") == 1,
    }
}

func (d *EnqueueIngestBatchResultDto) GetBatchUuid() string {
    return d.batchUuid
}

func (d *EnqueueIngestBatchResultDto) GetQueuedEvents() int {
    return d.queuedEvents
}

func (d *EnqueueIngestBatchResultDto) IsDuplicate() bool {
    return d.isDuplicate
}

func (d *EnqueueIngestBatchResultDto) ToArray() EnqueueIngestBatchPrimitives {
    return EnqueueIngestBatchPrimitives{BatchUuid: d.batchUuid, QueuedEvents: d.queuedEvents, IsDuplicate: d.isDuplicate}
}
