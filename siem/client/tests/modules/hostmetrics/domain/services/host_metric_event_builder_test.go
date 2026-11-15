package services

import (
    "testing"
    "time"

    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/hostmetrics/domain/services"
    "siem-client/tests"
)

var occurredAt = time.Date(2026, 11, 15, 18, 0, 0, 0, time.UTC)

func getRemoteMonitorEntity(monitorType string, monitorName string) *agentconfigentities.RemoteMonitorEntity {
    return agentconfigentities.RemoteMonitorEntityFromPrimitives(map[string]any{
        "monitor_uuid": "monitor-abc", "monitor_type": monitorType, "monitor_name": monitorName, "interval_seconds": float64(60),
    })
}

func TestDiskMetricEventFollowsTheContract(t *testing.T) {
    hostMetricSampleEntity := services.HostMetricSampleBuilderGetInstance().GetDiskSample("/", 81.2, 174376673280, 214748364800)
    event := services.HostMetricEventBuilderGetInstance().GetMetricEvent(getRemoteMonitorEntity("disk", "disco-raiz"), hostMetricSampleEntity, occurredAt)

    tests.AssertEquals(t, "host_metric", event["event_type"])
    tests.AssertEquals(t, "disco-raiz", event["source_name"])
    tests.AssertEquals(t, "disk / 81.2% used (162.4 GiB of 200.0 GiB)", event["message"])
    tests.AssertEquals(t, "2026-11-15 18:00:00", event["occurred_at"])
    attributes := event["attributes"].(map[string]any)
    tests.AssertEquals(t, "monitor-abc", attributes["monitor_uuid"])
    tests.AssertEquals(t, "disk", attributes["monitor_type"])
    tests.AssertEquals(t, 81.2, attributes["value"])
    tests.AssertEquals(t, "percent", attributes["unit"])
    tests.AssertEquals(t, "/", attributes["path"])
    tests.AssertEquals(t, int64(214748364800), attributes["total_bytes"])
    _, hasSeverity := event["severity"]
    tests.AssertEquals(t, false, hasSeverity)
}

func TestNetworkValueIsTheHighestOfRxAndTx(t *testing.T) {
    hostMetricSampleEntity := services.HostMetricSampleBuilderGetInstance().GetNetworkSample("", 1.2, 3.45)
    event := services.HostMetricEventBuilderGetInstance().GetMetricEvent(getRemoteMonitorEntity("network", "red"), hostMetricSampleEntity, occurredAt)
    attributes := event["attributes"].(map[string]any)
    tests.AssertEquals(t, 3.45, attributes["value"])
    tests.AssertEquals(t, "mbit_per_second", attributes["unit"])
    tests.AssertEquals(t, "network all rx 1.20 Mbit/s tx 3.45 Mbit/s", event["message"])
}

func TestCpuMemoryAndGpuMessages(t *testing.T) {
    hostMetricSampleBuilder := services.HostMetricSampleBuilderGetInstance()
    tests.AssertEquals(t, "cpu 12.3% used (load 0.52 0.40 0.31)", hostMetricSampleBuilder.GetCpuSample(12.34, [3]float64{0.52, 0.4, 0.31}, true).GetMessage())
    tests.AssertEquals(t, "cpu 12.3% used", hostMetricSampleBuilder.GetCpuSample(12.34, [3]float64{}, false).GetMessage())
    tests.AssertEquals(t, "memory 25.0% used (2.0 GiB of 8.0 GiB)", hostMetricSampleBuilder.GetMemorySample(25, 2<<30, 8<<30).GetMessage())

    gpuSample := hostMetricSampleBuilder.GetGpuSample(services.GpuStat{Name: "RTX 4090", UtilizationPercent: 45, MemoryUsedMib: 1000, MemoryTotalMib: 4000, TemperatureCelsius: 61}, 2)
    tests.AssertEquals(t, "gpu RTX 4090 45.0% used · memory 25.0% · 61 °C", gpuSample.GetMessage())
    tests.AssertEquals(t, 45.0, gpuSample.GetValue())
    tests.AssertEquals(t, 2, gpuSample.GetAttributes()["gpu_count"])
}

func TestUnavailableEventHasNoValue(t *testing.T) {
    event := services.HostMetricEventBuilderGetInstance().GetUnavailableEvent(getRemoteMonitorEntity("gpu", "gpu"), "nvidia-smi not available: not found", occurredAt)
    attributes := event["attributes"].(map[string]any)
    tests.AssertEquals(t, true, attributes["is_unavailable"])
    _, hasValue := attributes["value"]
    tests.AssertEquals(t, false, hasValue)
    tests.AssertEquals(t, "gpu not available: nvidia-smi not available: not found", event["message"])
    tests.AssertEquals(t, "gpu", attributes["monitor_type"])
}
