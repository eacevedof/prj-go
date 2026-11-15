package services

import (
    "errors"
    "strconv"
    "strings"
)

const bytesPerKibibyte = 1024

// MemoryUsageCalculator memoria usada del host a partir de /proc/meminfo: MemTotal - MemAvailable (lo que el kernel
// estima que puede dar a un proceso nuevo sin tirar de swap: cuenta la cache liberable, como free -m en "available").
type MemoryUsageCalculator struct{}

func MemoryUsageCalculatorGetInstance() *MemoryUsageCalculator {
    return &MemoryUsageCalculator{}
}

// GetMemoryBytes total y usada, en bytes. Sin MemAvailable (kernels < 3.14): MemFree + Buffers + Cached.
func (c *MemoryUsageCalculator) GetMemoryBytes(procMeminfoContent string) (float64, float64, error) {
    kibibytesByKey := map[string]float64{}
    for _, line := range strings.Split(procMeminfoContent, "\n") {
        key, rawValue, hasColon := strings.Cut(line, ":")
        fields := strings.Fields(rawValue)
        if !hasColon || len(fields) == 0 {
            continue
        }
        if kibibytes, err := strconv.ParseFloat(fields[0], 64); err == nil {
            kibibytesByKey[strings.TrimSpace(key)] = kibibytes
        }
    }

    total, hasTotal := kibibytesByKey["MemTotal"]
    if !hasTotal || total <= 0 {
        return 0, 0, errors.New("MemTotal not found")
    }
    available, hasAvailable := kibibytesByKey["MemAvailable"]
    if !hasAvailable {
        available = kibibytesByKey["MemFree"] + kibibytesByKey["Buffers"] + kibibytesByKey["Cached"]
    }
    used := total - available
    if used < 0 {
        used = 0
    }
    return total * bytesPerKibibyte, used * bytesPerKibibyte, nil
}

// GetUsedPercent 0 si total no es positivo.
func (c *MemoryUsageCalculator) GetUsedPercent(totalBytes float64, usedBytes float64) float64 {
    if totalBytes <= 0 {
        return 0
    }
    return getRounded(usedBytes / totalBytes * 100)
}
