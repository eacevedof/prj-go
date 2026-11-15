package services

import (
    "strconv"
    "strings"
)

// GpuStat una GPU de la salida de nvidia-smi. Un dato que la tarjeta no da ("[N/A]") queda en -1.
type GpuStat struct {
    Name               string
    UtilizationPercent float64
    MemoryUsedMib      float64
    MemoryTotalMib     float64
    TemperatureCelsius float64
}

// GetMemoryPercent -1 si la tarjeta no informa de su memoria.
func (g GpuStat) GetMemoryPercent() float64 {
    if g.MemoryUsedMib < 0 || g.MemoryTotalMib <= 0 {
        return -1
    }
    return getRounded(g.MemoryUsedMib / g.MemoryTotalMib * 100)
}

// GpuStatsParser lee la salida de
// nvidia-smi --query-gpu=name,utilization.gpu,memory.used,memory.total,temperature.gpu --format=csv,noheader,nounits
// (una linea por GPU: "NVIDIA GeForce RTX 4090, 45, 3012, 24564, 61").
type GpuStatsParser struct{}

func GpuStatsParserGetInstance() *GpuStatsParser {
    return &GpuStatsParser{}
}

// GetGpuStats las GPUs con un % de uso legible; una linea que no tiene los 5 campos o cuyo uso no es un numero se
// ignora.
func (p *GpuStatsParser) GetGpuStats(nvidiaSmiOutput string) []GpuStat {
    gpuStats := []GpuStat{}
    for _, line := range strings.Split(nvidiaSmiOutput, "\n") {
        fields := strings.Split(line, ",")
        if len(fields) < 5 {
            continue
        }
        utilizationPercent := p.getNumberOrMinusOne(fields[1])
        if utilizationPercent < 0 {
            continue
        }
        gpuStats = append(gpuStats, GpuStat{
            Name:               strings.TrimSpace(fields[0]),
            UtilizationPercent: utilizationPercent,
            MemoryUsedMib:      p.getNumberOrMinusOne(fields[2]),
            MemoryTotalMib:     p.getNumberOrMinusOne(fields[3]),
            TemperatureCelsius: p.getNumberOrMinusOne(fields[4]),
        })
    }
    return gpuStats
}

// GetBusiestGpuStat la mas cargada (su % es el value del evento); false si no hay ninguna.
func (p *GpuStatsParser) GetBusiestGpuStat(gpuStats []GpuStat) (GpuStat, bool) {
    if len(gpuStats) == 0 {
        return GpuStat{}, false
    }
    busiestGpuStat := gpuStats[0]
    for _, gpuStat := range gpuStats[1:] {
        if gpuStat.UtilizationPercent > busiestGpuStat.UtilizationPercent {
            busiestGpuStat = gpuStat
        }
    }
    return busiestGpuStat, true
}

func (p *GpuStatsParser) getNumberOrMinusOne(field string) float64 {
    number, err := strconv.ParseFloat(strings.TrimSpace(field), 64)
    if err != nil || number < 0 {
        return -1
    }
    return number
}
