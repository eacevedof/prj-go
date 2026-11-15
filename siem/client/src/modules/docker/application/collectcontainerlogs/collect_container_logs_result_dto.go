package collectcontainerlogs

import "siem-client/src/modules/shared/application"

// CollectContainerLogsPrimitives lo que imprime docker-collect de esta parte.
type CollectContainerLogsPrimitives struct {
    IsDisabled     bool `json:"is_disabled"`
    ContainersRead int  `json:"containers_read"`
    LogLines       int  `json:"log_lines"`
    SpooledBatches int  `json:"spooled_batches"`
}

type CollectContainerLogsResultDto struct {
    isDisabled     bool
    containersRead int
    logLines       int
    spooledBatches int
}

func CollectContainerLogsResultDtoFromPrimitives(primitives map[string]any) *CollectContainerLogsResultDto {
    return &CollectContainerLogsResultDto{
        isDisabled:     application.GetInt(primitives, "isDisabled") == 1,
        containersRead: application.GetInt(primitives, "containersRead"),
        logLines:       application.GetInt(primitives, "logLines"),
        spooledBatches: application.GetInt(primitives, "spooledBatches"),
    }
}

func (d *CollectContainerLogsResultDto) GetLogLines() int {
    return d.logLines
}

func (d *CollectContainerLogsResultDto) ToArray() CollectContainerLogsPrimitives {
    return CollectContainerLogsPrimitives{
        IsDisabled:     d.isDisabled,
        ContainersRead: d.containersRead,
        LogLines:       d.logLines,
        SpooledBatches: d.spooledBatches,
    }
}
