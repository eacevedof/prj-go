package services

import (
    "math"

    "siem-client/src/modules/docker/domain/entities"
)

// ContainerStatsCalculator de la respuesta de GET /containers/{id}/stats?stream=false&one-shot=true (ya decodificada)
// a una muestra de recursos. one-shot no trae la medida anterior de cpu, asi que el % de cpu se calcula contra la
// muestra que guardo el agente en la pasada anterior (como hace docker stats con precpu_stats).
type ContainerStatsCalculator struct{}

func ContainerStatsCalculatorGetInstance() *ContainerStatsCalculator {
    return &ContainerStatsCalculator{}
}

// GetContainerStatsEntity previousCpuTotal y previousSystemTotal a 0 = no hay muestra anterior (cpu desconocida).
func (c *ContainerStatsCalculator) GetContainerStatsEntity(rawStats map[string]any, previousCpuTotal float64, previousSystemTotal float64) *entities.ContainerStatsEntity {
    cpuStats := getMap(rawStats, "cpu_stats")
    cpuTotal := getNumber(getMap(cpuStats, "cpu_usage"), "total_usage")
    systemTotal := getNumber(cpuStats, "system_cpu_usage")

    cpuPercent := -1.0
    cpuDelta := cpuTotal - previousCpuTotal
    systemDelta := systemTotal - previousSystemTotal
    if previousSystemTotal > 0 && systemDelta > 0 && cpuDelta >= 0 {
        cpuPercent = getRounded(cpuDelta / systemDelta * c.getOnlineCpus(cpuStats) * 100)
    }

    memoryStats := getMap(rawStats, "memory_stats")
    memoryUsageBytes := c.getMemoryUsageBytes(memoryStats)
    memoryLimitBytes := getNumber(memoryStats, "limit")
    memoryPercent := 0.0
    if memoryLimitBytes > 0 {
        memoryPercent = getRounded(memoryUsageBytes / memoryLimitBytes * 100)
    }

    netRxBytes, netTxBytes := 0.0, 0.0
    for _, rawNetwork := range getMap(rawStats, "networks") {
        network, _ := rawNetwork.(map[string]any)
        netRxBytes += getNumber(network, "rx_bytes")
        netTxBytes += getNumber(network, "tx_bytes")
    }

    return entities.ContainerStatsEntityFromPrimitives(map[string]float64{
        "cpuPercent":       cpuPercent,
        "memoryUsageBytes": memoryUsageBytes,
        "memoryLimitBytes": memoryLimitBytes,
        "memoryPercent":    memoryPercent,
        "netRxBytes":       netRxBytes,
        "netTxBytes":       netTxBytes,
        "cpuTotal":         cpuTotal,
        "systemTotal":      systemTotal,
    })
}

// getOnlineCpus online_cpus; en daemons antiguos, el numero de cpus de percpu_usage; si no, 1.
func (c *ContainerStatsCalculator) getOnlineCpus(cpuStats map[string]any) float64 {
    if onlineCpus := getNumber(cpuStats, "online_cpus"); onlineCpus > 0 {
        return onlineCpus
    }
    if perCpuUsage, isList := getMap(cpuStats, "cpu_usage")["percpu_usage"].([]any); isList && len(perCpuUsage) > 0 {
        return float64(len(perCpuUsage))
    }
    return 1
}

// getMemoryUsageBytes como docker stats: el uso menos la cache de ficheros inactiva (cgroup v2: inactive_file;
// v1: total_inactive_file), que el kernel libera si hace falta.
func (c *ContainerStatsCalculator) getMemoryUsageBytes(memoryStats map[string]any) float64 {
    usage := getNumber(memoryStats, "usage")
    detailedStats := getMap(memoryStats, "stats")
    inactiveFile := getNumber(detailedStats, "inactive_file")
    if inactiveFile == 0 {
        inactiveFile = getNumber(detailedStats, "total_inactive_file")
    }
    if inactiveFile > 0 && inactiveFile < usage {
        return usage - inactiveFile
    }
    return usage
}

func getMap(primitives map[string]any, key string) map[string]any {
    value, _ := primitives[key].(map[string]any)
    return value
}

func getNumber(primitives map[string]any, key string) float64 {
    value, _ := primitives[key].(float64)
    return value
}

// getRounded 2 decimales: suficiente para un %, y el JSON no arrastra 15 decimales.
func getRounded(value float64) float64 {
    return math.Round(value*100) / 100
}
