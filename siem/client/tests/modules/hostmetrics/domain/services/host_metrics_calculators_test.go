package services

import (
    "testing"

    "siem-client/src/modules/hostmetrics/domain/services"
    "siem-client/tests"
)

const procStat = `cpu  1000 50 500 8000 200 0 50 0 300 0
cpu0 500 25 250 4000 100 0 25 0 150 0
intr 12345
`

func TestCpuTimesFromProcStat(t *testing.T) {
    total, idle, err := services.CpuUsageCalculatorGetInstance().GetCpuTimes(procStat)
    tests.AssertNil(t, err)
    // user nice system idle iowait irq softirq steal (guest ya va dentro de user)
    tests.AssertEquals(t, 9800.0, total)
    tests.AssertEquals(t, 8200.0, idle)

    _, _, err = services.CpuUsageCalculatorGetInstance().GetCpuTimes("intr 1 2 3\n")
    tests.AssertNotEquals(t, nil, err)
}

func TestCpuUsagePercentBetweenTwoSamples(t *testing.T) {
    cpuUsageCalculator := services.CpuUsageCalculatorGetInstance()
    // 1000 jiffies entre muestras, 750 de ellos idle: 25% de uso
    tests.AssertEquals(t, 25.0, cpuUsageCalculator.GetUsagePercent(9800, 8200, 10800, 8950))
    // sin muestra anterior, o contadores que retroceden (reinicio del host): no se puede calcular
    tests.AssertEquals(t, -1.0, cpuUsageCalculator.GetUsagePercent(0, 0, 10800, 8950))
    tests.AssertEquals(t, -1.0, cpuUsageCalculator.GetUsagePercent(9800, 8200, 500, 400))
}

func TestLoadAverages(t *testing.T) {
    loadAverages, hasLoadAverages := services.CpuUsageCalculatorGetInstance().GetLoadAverages("0.52 0.40 0.31 1/234 5678\n")
    tests.AssertEquals(t, true, hasLoadAverages)
    tests.AssertEquals(t, [3]float64{0.52, 0.40, 0.31}, loadAverages)
    _, hasLoadAverages = services.CpuUsageCalculatorGetInstance().GetLoadAverages("garbage")
    tests.AssertEquals(t, false, hasLoadAverages)
}

func TestMemoryUsedIsTotalMinusAvailable(t *testing.T) {
    memoryUsageCalculator := services.MemoryUsageCalculatorGetInstance()
    totalBytes, usedBytes, err := memoryUsageCalculator.GetMemoryBytes("MemTotal:        8000000 kB\nMemFree:          500000 kB\nMemAvailable:    6000000 kB\nBuffers:          100000 kB\nCached:          2000000 kB\n")
    tests.AssertNil(t, err)
    tests.AssertEquals(t, 8000000.0*1024, totalBytes)
    tests.AssertEquals(t, 2000000.0*1024, usedBytes)
    tests.AssertEquals(t, 25.0, memoryUsageCalculator.GetUsedPercent(totalBytes, usedBytes))

    // kernel antiguo sin MemAvailable: MemFree + Buffers + Cached
    _, usedBytes, err = memoryUsageCalculator.GetMemoryBytes("MemTotal: 1000 kB\nMemFree: 100 kB\nBuffers: 100 kB\nCached: 300 kB\n")
    tests.AssertNil(t, err)
    tests.AssertEquals(t, 500.0*1024, usedBytes)

    _, _, err = memoryUsageCalculator.GetMemoryBytes("nothing here")
    tests.AssertNotEquals(t, nil, err)
}

const procNetDev = `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo: 9000000   100    0    0    0     0          0         0  9000000   100    0    0    0     0       0          0
  eth0: 1000000   200    0    0    0     0          0         0   500000   150    0    0    0     0       0          0
  eth1:  250000    20    0    0    0     0          0         0   125000    10    0    0    0     0       0          0
`

func TestNetworkBytesOfOneInterfaceOrAllButLoopback(t *testing.T) {
    networkRateCalculator := services.NetworkRateCalculatorGetInstance()
    rxBytes, txBytes, isFound := networkRateCalculator.GetInterfaceBytes(procNetDev, "eth0")
    tests.AssertEquals(t, true, isFound)
    tests.AssertEquals(t, 1000000.0, rxBytes)
    tests.AssertEquals(t, 500000.0, txBytes)

    rxBytes, txBytes, isFound = networkRateCalculator.GetInterfaceBytes(procNetDev, "")
    tests.AssertEquals(t, true, isFound)
    tests.AssertEquals(t, 1250000.0, rxBytes)
    tests.AssertEquals(t, 625000.0, txBytes)

    _, _, isFound = networkRateCalculator.GetInterfaceBytes(procNetDev, "wlan9")
    tests.AssertEquals(t, false, isFound)
}

func TestNetworkMbitPerSecond(t *testing.T) {
    networkRateCalculator := services.NetworkRateCalculatorGetInstance()
    // 2.500.000 bytes en 10 s = 20.000.000 bits / 10 s = 2 Mbit/s
    tests.AssertEquals(t, 2.0, networkRateCalculator.GetMbitPerSecond(1000000, 3500000, 10000))
    // contador que retrocede (reinicio) o sin intervalo: no se puede calcular
    tests.AssertEquals(t, -1.0, networkRateCalculator.GetMbitPerSecond(3500000, 1000, 10000))
    tests.AssertEquals(t, -1.0, networkRateCalculator.GetMbitPerSecond(1000, 2000, 0))
}

func TestGpuStatsFromNvidiaSmiCsv(t *testing.T) {
    gpuStatsParser := services.GpuStatsParserGetInstance()
    gpuStats := gpuStatsParser.GetGpuStats("NVIDIA GeForce RTX 4090, 45, 3012, 24564, 61\nTesla T4, 80, [N/A], [N/A], 70\nbroken line\nGTX, [N/A], 1, 2, 3\n")
    tests.AssertEquals(t, 2, len(gpuStats))
    tests.AssertEquals(t, "NVIDIA GeForce RTX 4090", gpuStats[0].Name)
    tests.AssertEquals(t, 12.26, gpuStats[0].GetMemoryPercent())
    tests.AssertEquals(t, -1.0, gpuStats[1].GetMemoryPercent())

    busiestGpuStat, hasGpu := gpuStatsParser.GetBusiestGpuStat(gpuStats)
    tests.AssertEquals(t, true, hasGpu)
    tests.AssertEquals(t, "Tesla T4", busiestGpuStat.Name)
    _, hasGpu = gpuStatsParser.GetBusiestGpuStat(gpuStatsParser.GetGpuStats("nvidia-smi: command not found"))
    tests.AssertEquals(t, false, hasGpu)
}
