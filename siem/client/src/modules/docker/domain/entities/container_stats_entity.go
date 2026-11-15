package entities

// ContainerStatsEntity una muestra de recursos de un contenedor. cpuPercent -1 = no se sabe (primera muestra: el %
// de cpu sale de comparar con la anterior). cpuTotal y systemTotal se guardan para la siguiente.
type ContainerStatsEntity struct {
    cpuPercent       float64
    memoryUsageBytes int64
    memoryLimitBytes int64
    memoryPercent    float64
    netRxBytes       int64
    netTxBytes       int64
    cpuTotal         float64
    systemTotal      float64
}

func ContainerStatsEntityFromPrimitives(primitives map[string]float64) *ContainerStatsEntity {
    return &ContainerStatsEntity{
        cpuPercent:       primitives["cpuPercent"],
        memoryUsageBytes: int64(primitives["memoryUsageBytes"]),
        memoryLimitBytes: int64(primitives["memoryLimitBytes"]),
        memoryPercent:    primitives["memoryPercent"],
        netRxBytes:       int64(primitives["netRxBytes"]),
        netTxBytes:       int64(primitives["netTxBytes"]),
        cpuTotal:         primitives["cpuTotal"],
        systemTotal:      primitives["systemTotal"],
    }
}

func (e *ContainerStatsEntity) HasCpuPercent() bool {
    return e.cpuPercent >= 0
}

func (e *ContainerStatsEntity) GetCpuPercent() float64 {
    return e.cpuPercent
}

func (e *ContainerStatsEntity) GetMemoryUsageBytes() int64 {
    return e.memoryUsageBytes
}

func (e *ContainerStatsEntity) GetMemoryLimitBytes() int64 {
    return e.memoryLimitBytes
}

func (e *ContainerStatsEntity) GetMemoryPercent() float64 {
    return e.memoryPercent
}

func (e *ContainerStatsEntity) GetNetRxBytes() int64 {
    return e.netRxBytes
}

func (e *ContainerStatsEntity) GetNetTxBytes() int64 {
    return e.netTxBytes
}

func (e *ContainerStatsEntity) GetCpuTotal() float64 {
    return e.cpuTotal
}

func (e *ContainerStatsEntity) GetSystemTotal() float64 {
    return e.systemTotal
}
