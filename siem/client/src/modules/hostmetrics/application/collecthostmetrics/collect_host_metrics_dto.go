package collecthostmetrics

import "siem-client/src/modules/shared/application"

type CollectHostMetricsDto struct {
    // true = muestrea todos los monitores ya, toque o no (./console metrics-collect --force=1)
    isForced bool
}

func CollectHostMetricsDtoFromPrimitives(primitives map[string]any) *CollectHostMetricsDto {
    return &CollectHostMetricsDto{isForced: application.GetInt(primitives, "isForced") == 1}
}

func (d *CollectHostMetricsDto) IsForced() bool {
    return d.isForced
}
