package sendspooledevents

import "siem-client/src/modules/shared/application"

const MAX_BATCHES_DEFAULT = 20

type SendSpooledEventsDto struct {
    // lotes como mucho por pasada: con el backend recien vuelto, el spool se vacia en varias pasadas sin saturarlo
    maxBatches int
}

func SendSpooledEventsDtoFromPrimitives(primitives map[string]any) *SendSpooledEventsDto {
    maxBatches := application.GetInt(primitives, "maxBatches")
    if maxBatches < 1 {
        maxBatches = MAX_BATCHES_DEFAULT
    }
    return &SendSpooledEventsDto{maxBatches: maxBatches}
}

func (d *SendSpooledEventsDto) GetMaxBatches() int {
    return d.maxBatches
}
