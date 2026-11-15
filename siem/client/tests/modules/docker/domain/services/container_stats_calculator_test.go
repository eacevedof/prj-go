package services

import (
    "testing"

    "siem-client/src/modules/docker/domain/services"
    "siem-client/tests"
)

func getRawStats(cpuTotal float64, systemTotal float64) map[string]any {
    return map[string]any{
        "cpu_stats": map[string]any{
            "cpu_usage":        map[string]any{"total_usage": cpuTotal},
            "system_cpu_usage": systemTotal,
            "online_cpus":      float64(4),
        },
        "memory_stats": map[string]any{
            "usage": float64(300 << 20),
            "limit": float64(1 << 30),
            "stats": map[string]any{"inactive_file": float64(44 << 20)},
        },
        "networks": map[string]any{
            "eth0": map[string]any{"rx_bytes": float64(1000), "tx_bytes": float64(500)},
            "eth1": map[string]any{"rx_bytes": float64(24), "tx_bytes": float64(12)},
        },
    }
}

func TestFirstSampleHasNoCpuPercent(t *testing.T) {
    containerStatsEntity := services.ContainerStatsCalculatorGetInstance().GetContainerStatsEntity(getRawStats(5e9, 1e12), 0, 0)
    tests.AssertEquals(t, false, containerStatsEntity.HasCpuPercent())
    tests.AssertEquals(t, 5e9, containerStatsEntity.GetCpuTotal())
}

func TestCpuPercentAgainstPreviousSample(t *testing.T) {
    // el contenedor gasto 1e9 ns de cpu mientras el sistema (4 cpus) avanzo 1e10: 10% * 4 = 40%
    containerStatsEntity := services.ContainerStatsCalculatorGetInstance().GetContainerStatsEntity(getRawStats(6e9, 1.01e12), 5e9, 1e12)
    tests.AssertEquals(t, true, containerStatsEntity.HasCpuPercent())
    tests.AssertEquals(t, 40.0, containerStatsEntity.GetCpuPercent())
}

func TestMemoryWithoutInactiveFileAndNetworksAdded(t *testing.T) {
    containerStatsEntity := services.ContainerStatsCalculatorGetInstance().GetContainerStatsEntity(getRawStats(1, 1), 0, 0)
    // 300 MiB - 44 MiB de cache inactiva = 256 MiB de 1 GiB = 25%
    tests.AssertEquals(t, int64(256<<20), containerStatsEntity.GetMemoryUsageBytes())
    tests.AssertEquals(t, 25.0, containerStatsEntity.GetMemoryPercent())
    tests.AssertEquals(t, int64(1024), containerStatsEntity.GetNetRxBytes())
    tests.AssertEquals(t, int64(512), containerStatsEntity.GetNetTxBytes())
}
