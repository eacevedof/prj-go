package collectcontainerstates

import "siem-client/src/modules/shared/application"

// CollectContainerStatesPrimitives lo que imprime docker-collect de esta parte (claves en este orden).
type CollectContainerStatesPrimitives struct {
    IsDisabled        bool             `json:"is_disabled"`
    Containers        int              `json:"containers"`
    RunningContainers int              `json:"running_containers"`
    Changes           int              `json:"changes"`
    SpooledBatches    int              `json:"spooled_batches"`
    Results           []map[string]any `json:"results"`
}

type CollectContainerStatesResultDto struct {
    isDisabled        bool
    containers        int
    runningContainers int
    changes           int
    spooledBatches    int
    // los eventos docker_state de esta pasada (solo primitivos)
    results []map[string]any
}

func CollectContainerStatesResultDtoFromPrimitives(primitives map[string]any) *CollectContainerStatesResultDto {
    results, _ := primitives["results"].([]map[string]any)
    if results == nil {
        results = []map[string]any{}
    }
    return &CollectContainerStatesResultDto{
        isDisabled:        application.GetInt(primitives, "isDisabled") == 1,
        containers:        application.GetInt(primitives, "containers"),
        runningContainers: application.GetInt(primitives, "runningContainers"),
        changes:           application.GetInt(primitives, "changes"),
        spooledBatches:    application.GetInt(primitives, "spooledBatches"),
        results:           results,
    }
}

func (d *CollectContainerStatesResultDto) IsDisabled() bool {
    return d.isDisabled
}

func (d *CollectContainerStatesResultDto) GetContainers() int {
    return d.containers
}

func (d *CollectContainerStatesResultDto) GetRunningContainers() int {
    return d.runningContainers
}

func (d *CollectContainerStatesResultDto) GetChanges() int {
    return d.changes
}

func (d *CollectContainerStatesResultDto) ToArray() CollectContainerStatesPrimitives {
    return CollectContainerStatesPrimitives{
        IsDisabled:        d.isDisabled,
        Containers:        d.containers,
        RunningContainers: d.runningContainers,
        Changes:           d.changes,
        SpooledBatches:    d.spooledBatches,
        Results:           d.results,
    }
}
