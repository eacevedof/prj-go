package collectcontainerstats

import "siem-client/src/modules/shared/application"

type CollectContainerStatsDto struct {
    // true = toma la muestra ya, toque o no (./console docker-collect --force_stats=1)
    isForced bool
}

func CollectContainerStatsDtoFromPrimitives(primitives map[string]any) *CollectContainerStatsDto {
    return &CollectContainerStatsDto{isForced: application.GetInt(primitives, "isForced") == 1}
}

func (d *CollectContainerStatsDto) IsForced() bool {
    return d.isForced
}
