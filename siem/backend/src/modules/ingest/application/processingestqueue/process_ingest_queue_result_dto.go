package processingestqueue

import "siem-backend/src/modules/shared/application"

// ProcessIngestQueuePrimitives lo que imprime ./console ingest-queue-process (claves en este orden).
type ProcessIngestQueuePrimitives struct {
    ProcessedBatches int `json:"processed_batches"`
    FailedBatches    int `json:"failed_batches"`
    InsertedEvents   int `json:"inserted_events"`
    SkippedEvents    int `json:"skipped_events"`
}

type ProcessIngestQueueResultDto struct {
    processedBatches int
    failedBatches    int
    insertedEvents   int
    skippedEvents    int
}

func ProcessIngestQueueResultDtoFromPrimitives(primitives map[string]any) *ProcessIngestQueueResultDto {
    return &ProcessIngestQueueResultDto{
        processedBatches: application.GetInt(primitives, "processedBatches"),
        failedBatches:    application.GetInt(primitives, "failedBatches"),
        insertedEvents:   application.GetInt(primitives, "insertedEvents"),
        skippedEvents:    application.GetInt(primitives, "skippedEvents"),
    }
}

func (d *ProcessIngestQueueResultDto) GetProcessedBatches() int {
    return d.processedBatches
}

func (d *ProcessIngestQueueResultDto) GetFailedBatches() int {
    return d.failedBatches
}

func (d *ProcessIngestQueueResultDto) GetInsertedEvents() int {
    return d.insertedEvents
}

func (d *ProcessIngestQueueResultDto) GetSkippedEvents() int {
    return d.skippedEvents
}

// GetTotalBatches lotes que salieron de pending en esta pasada (bien o mal): si iguala al maximo pedido, puede
// quedar cola.
func (d *ProcessIngestQueueResultDto) GetTotalBatches() int {
    return d.processedBatches + d.failedBatches
}

func (d *ProcessIngestQueueResultDto) ToArray() ProcessIngestQueuePrimitives {
    return ProcessIngestQueuePrimitives{
        ProcessedBatches: d.processedBatches,
        FailedBatches:    d.failedBatches,
        InsertedEvents:   d.insertedEvents,
        SkippedEvents:    d.skippedEvents,
    }
}
