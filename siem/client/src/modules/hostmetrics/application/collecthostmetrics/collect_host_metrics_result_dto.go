package collecthostmetrics

import "siem-client/src/modules/shared/application"

// CollectHostMetricsPrimitives lo que imprime metrics-collect (claves en este orden).
type CollectHostMetricsPrimitives struct {
    MetricMonitors      int              `json:"metric_monitors"`
    SampledMonitors     int              `json:"sampled_monitors"`
    UnavailableMonitors int              `json:"unavailable_monitors"`
    SpooledBatches      int              `json:"spooled_batches"`
    Results             []map[string]any `json:"results"`
}

type CollectHostMetricsResultDto struct {
    // monitores de metricas en la config (les tocara o no en esta pasada)
    metricMonitors int
    // los que han dejado un evento con valor
    sampledMonitors int
    // los que han dejado el aviso (unico) de metrica no disponible
    unavailableMonitors int
    spooledBatches      int
    // los eventos host_metric de esta pasada (solo primitivos)
    results []map[string]any
}

func CollectHostMetricsResultDtoFromPrimitives(primitives map[string]any) *CollectHostMetricsResultDto {
    results, _ := primitives["results"].([]map[string]any)
    if results == nil {
        results = []map[string]any{}
    }
    return &CollectHostMetricsResultDto{
        metricMonitors:      application.GetInt(primitives, "metricMonitors"),
        sampledMonitors:     application.GetInt(primitives, "sampledMonitors"),
        unavailableMonitors: application.GetInt(primitives, "unavailableMonitors"),
        spooledBatches:      application.GetInt(primitives, "spooledBatches"),
        results:             results,
    }
}

func (d *CollectHostMetricsResultDto) GetSampledMonitors() int {
    return d.sampledMonitors
}

func (d *CollectHostMetricsResultDto) GetUnavailableMonitors() int {
    return d.unavailableMonitors
}

func (d *CollectHostMetricsResultDto) ToArray() CollectHostMetricsPrimitives {
    return CollectHostMetricsPrimitives{
        MetricMonitors:      d.metricMonitors,
        SampledMonitors:     d.sampledMonitors,
        UnavailableMonitors: d.unavailableMonitors,
        SpooledBatches:      d.spooledBatches,
        Results:             d.results,
    }
}
