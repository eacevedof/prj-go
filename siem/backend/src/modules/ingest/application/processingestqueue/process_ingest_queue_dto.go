package processingestqueue

import "siem-backend/src/modules/shared/application"

const MAX_BATCHES_DEFAULT = 100

type ProcessIngestQueueDto struct {
    application.AbstractInputDto
    maxBatches int
}

func ProcessIngestQueueDtoFromPrimitives(primitives map[string]any) *ProcessIngestQueueDto {
    maxBatches := MAX_BATCHES_DEFAULT
    // "" (opcion de consola sin pasar) cuenta como ausente
    if application.GetTrimmedString(primitives, "maxBatches") != "" {
        maxBatches = application.GetInt(primitives, "maxBatches")
    }
    return &ProcessIngestQueueDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        maxBatches:       maxBatches,
    }
}

func (d *ProcessIngestQueueDto) GetMaxBatches() int {
    return d.maxBatches
}
