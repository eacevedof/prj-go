package services

import (
    "fmt"

    "siem-client/src/modules/hostmetrics/domain/entities"
    "siem-client/src/modules/hostmetrics/domain/enums"
)

// HostMetricSampleBuilder de las medidas crudas de cada tipo a la muestra que viaja: el value que el backend compara
// con los umbrales, un mensaje legible para la tabla de eventos y los datos sueltos del tipo.
type HostMetricSampleBuilder struct{}

func HostMetricSampleBuilderGetInstance() *HostMetricSampleBuilder {
    return &HostMetricSampleBuilder{}
}

// GetCpuSample "cpu 12.3% used (load 0.52 0.40 0.31)".
func (b *HostMetricSampleBuilder) GetCpuSample(usagePercent float64, loadAverages [3]float64, hasLoadAverages bool) *entities.HostMetricSampleEntity {
    attributes := map[string]any{}
    message := fmt.Sprintf("cpu %.1f%% used", usagePercent)
    if hasLoadAverages {
        attributes[enums.HostMetricAttributeEnum.LOAD_1M] = loadAverages[0]
        attributes[enums.HostMetricAttributeEnum.LOAD_5M] = loadAverages[1]
        attributes[enums.HostMetricAttributeEnum.LOAD_15M] = loadAverages[2]
        message += fmt.Sprintf(" (load %.2f %.2f %.2f)", loadAverages[0], loadAverages[1], loadAverages[2])
    }
    return entities.HostMetricSampleEntityFromPrimitives(usagePercent, enums.MetricUnitEnum.PERCENT, message, attributes)
}

// GetMemorySample "memory 63.1% used (4.8 GiB of 7.7 GiB)".
func (b *HostMetricSampleBuilder) GetMemorySample(usedPercent float64, usedBytes float64, totalBytes float64) *entities.HostMetricSampleEntity {
    return entities.HostMetricSampleEntityFromPrimitives(
        usedPercent, enums.MetricUnitEnum.PERCENT,
        fmt.Sprintf("memory %.1f%% used (%s of %s)", usedPercent, b.getBytesLabel(usedBytes), b.getBytesLabel(totalBytes)),
        map[string]any{
            enums.HostMetricAttributeEnum.USED_BYTES:  int64(usedBytes),
            enums.HostMetricAttributeEnum.TOTAL_BYTES: int64(totalBytes),
        },
    )
}

// GetDiskSample "disk / 81.2% used (162.4 GiB of 200.0 GiB)".
func (b *HostMetricSampleBuilder) GetDiskSample(path string, usedPercent float64, usedBytes float64, totalBytes float64) *entities.HostMetricSampleEntity {
    return entities.HostMetricSampleEntityFromPrimitives(
        usedPercent, enums.MetricUnitEnum.PERCENT,
        fmt.Sprintf("disk %s %.1f%% used (%s of %s)", path, usedPercent, b.getBytesLabel(usedBytes), b.getBytesLabel(totalBytes)),
        map[string]any{
            enums.HostMetricAttributeEnum.PATH:        path,
            enums.HostMetricAttributeEnum.USED_BYTES:  int64(usedBytes),
            enums.HostMetricAttributeEnum.TOTAL_BYTES: int64(totalBytes),
        },
    )
}

// GetNetworkSample value = max(rx, tx) en Mbit/s. "network eth0 rx 1.20 Mbit/s tx 0.30 Mbit/s" ("all" = todas menos lo).
func (b *HostMetricSampleBuilder) GetNetworkSample(interfaceName string, rxMbit float64, txMbit float64) *entities.HostMetricSampleEntity {
    interfaceLabel := interfaceName
    if interfaceLabel == "" {
        interfaceLabel = "all"
    }
    return entities.HostMetricSampleEntityFromPrimitives(
        max(rxMbit, txMbit), enums.MetricUnitEnum.MBIT_PER_SECOND,
        fmt.Sprintf("network %s rx %.2f Mbit/s tx %.2f Mbit/s", interfaceLabel, rxMbit, txMbit),
        map[string]any{
            enums.HostMetricAttributeEnum.INTERFACE: interfaceName,
            enums.HostMetricAttributeEnum.RX_MBIT:   rxMbit,
            enums.HostMetricAttributeEnum.TX_MBIT:   txMbit,
        },
    )
}

// GetGpuSample value = % de uso de la GPU mas cargada. "gpu NVIDIA RTX 4090 45.0% used · memory 12.3% · 61 °C".
func (b *HostMetricSampleBuilder) GetGpuSample(busiestGpuStat GpuStat, gpuCount int) *entities.HostMetricSampleEntity {
    attributes := map[string]any{
        enums.HostMetricAttributeEnum.GPU_NAME:  busiestGpuStat.Name,
        enums.HostMetricAttributeEnum.GPU_COUNT: gpuCount,
    }
    message := fmt.Sprintf("gpu %s %.1f%% used", busiestGpuStat.Name, busiestGpuStat.UtilizationPercent)
    if memoryPercent := busiestGpuStat.GetMemoryPercent(); memoryPercent >= 0 {
        attributes[enums.HostMetricAttributeEnum.GPU_MEMORY_PERCENT] = memoryPercent
        message += fmt.Sprintf(" · memory %.1f%%", memoryPercent)
    }
    if busiestGpuStat.TemperatureCelsius >= 0 {
        attributes[enums.HostMetricAttributeEnum.GPU_TEMPERATURE] = busiestGpuStat.TemperatureCelsius
        message += fmt.Sprintf(" · %.0f °C", busiestGpuStat.TemperatureCelsius)
    }
    return entities.HostMetricSampleEntityFromPrimitives(busiestGpuStat.UtilizationPercent, enums.MetricUnitEnum.PERCENT, message, attributes)
}

// getBytesLabel en binario: 512 B, 1.5 KiB, 256.0 MiB, 7.7 GiB.
func (b *HostMetricSampleBuilder) getBytesLabel(bytes float64) string {
    units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
    unitIndex := 0
    for bytes >= 1024 && unitIndex < len(units)-1 {
        bytes /= 1024
        unitIndex++
    }
    if unitIndex == 0 {
        return fmt.Sprintf("%.0f B", bytes)
    }
    return fmt.Sprintf("%.1f %s", bytes, units[unitIndex])
}
