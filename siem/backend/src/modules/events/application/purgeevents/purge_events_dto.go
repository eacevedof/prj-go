package purgeevents

import "siem-backend/src/modules/shared/application"

type PurgeEventsDto struct {
    application.AbstractInputDto
    // 0 = no se pasa: se usa APP_EVENTS_RETENTION_DAYS
    retentionDays int
}

func PurgeEventsDtoFromPrimitives(primitives map[string]any) *PurgeEventsDto {
    return &PurgeEventsDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        retentionDays:    application.GetInt(primitives, "retentionDays"),
    }
}

func (d *PurgeEventsDto) GetRetentionDays() int {
    return d.retentionDays
}
