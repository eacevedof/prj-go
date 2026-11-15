package enums

// QueueStatusEnum estado de un lote en ingest_queue. Mismos valores que el CHECK de ingest_queue.queue_status.
var QueueStatusEnum = struct {
    PENDING string
    DONE    string
    FAILED  string
}{
    PENDING: "pending",
    DONE:    "done",
    FAILED:  "failed",
}
