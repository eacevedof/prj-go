package collectcontainerstats

import "siem-client/src/modules/shared/application"

// CollectContainerStatsPrimitives lo que imprime docker-collect de esta parte.
type CollectContainerStatsPrimitives struct {
    IsDisabled        bool             `json:"is_disabled"`
    IsNotDue          bool             `json:"is_not_due"`
    SampledContainers int              `json:"sampled_containers"`
    SpooledBatches    int              `json:"spooled_batches"`
    Results           []map[string]any `json:"results"`
}

type CollectContainerStatsResultDto struct {
    isDisabled        bool
    isNotDue          bool
    sampledContainers int
    spooledBatches    int
    results           []map[string]any
}

func CollectContainerStatsResultDtoFromPrimitives(primitives map[string]any) *CollectContainerStatsResultDto {
    results, _ := primitives["results"].([]map[string]any)
    if results == nil {
        results = []map[string]any{}
    }
    return &CollectContainerStatsResultDto{
        isDisabled:        application.GetInt(primitives, "isDisabled") == 1,
        isNotDue:          application.GetInt(primitives, "isNotDue") == 1,
        sampledContainers: application.GetInt(primitives, "sampledContainers"),
        spooledBatches:    application.GetInt(primitives, "spooledBatches"),
        results:           results,
    }
}

func (d *CollectContainerStatsResultDto) GetSampledContainers() int {
    return d.sampledContainers
}

func (d *CollectContainerStatsResultDto) ToArray() CollectContainerStatsPrimitives {
    return CollectContainerStatsPrimitives{
        IsDisabled:        d.isDisabled,
        IsNotDue:          d.isNotDue,
        SampledContainers: d.sampledContainers,
        SpooledBatches:    d.spooledBatches,
        Results:           d.results,
    }
}
