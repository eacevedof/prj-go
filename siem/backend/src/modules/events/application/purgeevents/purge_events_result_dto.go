package purgeevents

import "siem-backend/src/modules/shared/application"

// PurgeEventsPrimitives lo que imprime ./console events-purge (claves en este orden).
type PurgeEventsPrimitives struct {
    RetentionDays  int    `json:"retention_days"`
    DeletedBefore  string `json:"deleted_before"`
    DeletedEvents  int    `json:"deleted_events"`
    DeletedBatches int    `json:"deleted_batches"`
}

type PurgeEventsResultDto struct {
    retentionDays  int
    deletedBefore  string
    deletedEvents  int
    deletedBatches int
}

func PurgeEventsResultDtoFromPrimitives(primitives map[string]any) *PurgeEventsResultDto {
    return &PurgeEventsResultDto{
        retentionDays:  application.GetInt(primitives, "retentionDays"),
        deletedBefore:  application.GetTrimmedString(primitives, "deletedBefore"),
        deletedEvents:  application.GetInt(primitives, "deletedEvents"),
        deletedBatches: application.GetInt(primitives, "deletedBatches"),
    }
}

func (d *PurgeEventsResultDto) GetRetentionDays() int {
    return d.retentionDays
}

func (d *PurgeEventsResultDto) GetDeletedBefore() string {
    return d.deletedBefore
}

func (d *PurgeEventsResultDto) GetDeletedEvents() int {
    return d.deletedEvents
}

func (d *PurgeEventsResultDto) GetDeletedBatches() int {
    return d.deletedBatches
}

func (d *PurgeEventsResultDto) ToArray() PurgeEventsPrimitives {
    return PurgeEventsPrimitives{
        RetentionDays:  d.retentionDays,
        DeletedBefore:  d.deletedBefore,
        DeletedEvents:  d.deletedEvents,
        DeletedBatches: d.deletedBatches,
    }
}
