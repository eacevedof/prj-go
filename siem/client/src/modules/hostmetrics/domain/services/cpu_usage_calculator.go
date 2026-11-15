package services

import (
    "errors"
    "math"
    "strconv"
    "strings"
)

// CpuUsageCalculator % de cpu del host a partir de /proc/stat. La linea "cpu" lleva jiffies acumulados desde el
// arranque (user nice system idle iowait irq softirq steal guest guest_nice): el uso es lo que no fue idle ni iowait
// entre dos muestras, asi que hace falta la anterior.
type CpuUsageCalculator struct{}

func CpuUsageCalculatorGetInstance() *CpuUsageCalculator {
    return &CpuUsageCalculator{}
}

// GetCpuTimes total e idle (idle + iowait) de la linea agregada "cpu". guest y guest_nice ya van dentro de user y
// nice: no se suman otra vez.
func (c *CpuUsageCalculator) GetCpuTimes(procStatContent string) (float64, float64, error) {
    for _, line := range strings.Split(procStatContent, "\n") {
        fields := strings.Fields(line)
        if len(fields) < 5 || fields[0] != "cpu" {
            continue
        }
        total, idle := 0.0, 0.0
        for index, field := range fields[1:] {
            if index >= 8 {
                break
            }
            jiffies, err := strconv.ParseFloat(field, 64)
            if err != nil {
                return 0, 0, errors.New("cpu line is not numeric")
            }
            total += jiffies
            // 3 = idle, 4 = iowait
            if index == 3 || index == 4 {
                idle += jiffies
            }
        }
        return total, idle, nil
    }
    return 0, 0, errors.New("cpu line not found")
}

// GetUsagePercent entre dos muestras; -1 si no se puede calcular (sin muestra anterior, o contadores que retroceden:
// reinicio del host).
func (c *CpuUsageCalculator) GetUsagePercent(previousTotal float64, previousIdle float64, total float64, idle float64) float64 {
    totalDelta := total - previousTotal
    idleDelta := idle - previousIdle
    if previousTotal <= 0 || totalDelta <= 0 || idleDelta < 0 || idleDelta > totalDelta {
        return -1
    }
    return getRounded((totalDelta - idleDelta) / totalDelta * 100)
}

// GetLoadAverages 1, 5 y 15 minutos de /proc/loadavg ("0.52 0.40 0.31 1/234 5678"); false si no tiene ese formato.
func (c *CpuUsageCalculator) GetLoadAverages(procLoadavgContent string) ([3]float64, bool) {
    loadAverages := [3]float64{}
    fields := strings.Fields(procLoadavgContent)
    if len(fields) < 3 {
        return loadAverages, false
    }
    for index := range loadAverages {
        loadAverage, err := strconv.ParseFloat(fields[index], 64)
        if err != nil {
            return [3]float64{}, false
        }
        loadAverages[index] = loadAverage
    }
    return loadAverages, true
}

// getRounded 2 decimales: suficiente para un % o unos Mbit/s, y el JSON no arrastra 15 decimales.
func getRounded(value float64) float64 {
    return math.Round(value*100) / 100
}
