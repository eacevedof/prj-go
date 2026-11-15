package services

// EventsBatcher parte los eventos de una pasada en lotes de como mucho maxEventsPerBatch (cada lote = un fichero del
// spool = un POST al backend).
type EventsBatcher struct{}

func EventsBatcherGetInstance() *EventsBatcher {
    return &EventsBatcher{}
}

func (b *EventsBatcher) GetBatches(events []map[string]any, maxEventsPerBatch int) [][]map[string]any {
    batches := [][]map[string]any{}
    for start := 0; start < len(events); start += maxEventsPerBatch {
        end := min(start+maxEventsPerBatch, len(events))
        batches = append(batches, events[start:end])
    }
    return batches
}
