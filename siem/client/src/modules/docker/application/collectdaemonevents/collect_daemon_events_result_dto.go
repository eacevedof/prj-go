package collectdaemonevents

import "siem-client/src/modules/shared/application"

// CollectDaemonEventsPrimitives lo que imprime docker-collect de esta parte.
type CollectDaemonEventsPrimitives struct {
    IsDisabled     bool `json:"is_disabled"`
    DaemonEvents   int  `json:"daemon_events"`
    DroppedEvents  int  `json:"dropped_events"`
    SpooledBatches int  `json:"spooled_batches"`
}

type CollectDaemonEventsResultDto struct {
    isDisabled     bool
    daemonEvents   int
    droppedEvents  int
    spooledBatches int
}

func CollectDaemonEventsResultDtoFromPrimitives(primitives map[string]any) *CollectDaemonEventsResultDto {
    return &CollectDaemonEventsResultDto{
        isDisabled:     application.GetInt(primitives, "isDisabled") == 1,
        daemonEvents:   application.GetInt(primitives, "daemonEvents"),
        droppedEvents:  application.GetInt(primitives, "droppedEvents"),
        spooledBatches: application.GetInt(primitives, "spooledBatches"),
    }
}

func (d *CollectDaemonEventsResultDto) GetDaemonEvents() int {
    return d.daemonEvents
}

func (d *CollectDaemonEventsResultDto) GetDroppedEvents() int {
    return d.droppedEvents
}

func (d *CollectDaemonEventsResultDto) ToArray() CollectDaemonEventsPrimitives {
    return CollectDaemonEventsPrimitives{
        IsDisabled:     d.isDisabled,
        DaemonEvents:   d.daemonEvents,
        DroppedEvents:  d.droppedEvents,
        SpooledBatches: d.spooledBatches,
    }
}
