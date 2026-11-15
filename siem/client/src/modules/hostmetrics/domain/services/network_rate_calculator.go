package services

import (
    "strconv"
    "strings"
)

const (
    loopbackInterface = "lo"
    bitsPerByte       = 8
    bitsPerMegabit    = 1_000_000
)

// NetworkRateCalculator trafico de red del host a partir de /proc/net/dev. Cada linea ("eth0: rx_bytes rx_packets ...
// 8 campos de rx ... tx_bytes ...") lleva bytes acumulados desde el arranque: la velocidad es la diferencia entre dos
// muestras partida por el tiempo entre ellas.
type NetworkRateCalculator struct{}

func NetworkRateCalculatorGetInstance() *NetworkRateCalculator {
    return &NetworkRateCalculator{}
}

// GetInterfaceBytes bytes recibidos y enviados de una interfaz; con interfaceName vacio, la suma de todas menos lo.
// false si la interfaz pedida no existe (o si no hay ninguna que sumar).
func (c *NetworkRateCalculator) GetInterfaceBytes(procNetDevContent string, interfaceName string) (float64, float64, bool) {
    rxBytes, txBytes := 0.0, 0.0
    isFound := false
    for _, line := range strings.Split(procNetDevContent, "\n") {
        name, counters, hasColon := strings.Cut(line, ":")
        name = strings.TrimSpace(name)
        fields := strings.Fields(counters)
        // 16 contadores: 0 = rx_bytes, 8 = tx_bytes. Las dos lineas de cabecera no tienen ":" + numeros.
        if !hasColon || len(fields) < 9 {
            continue
        }
        if interfaceName == "" && name == loopbackInterface {
            continue
        }
        if interfaceName != "" && name != interfaceName {
            continue
        }
        lineRxBytes, rxErr := strconv.ParseFloat(fields[0], 64)
        lineTxBytes, txErr := strconv.ParseFloat(fields[8], 64)
        if rxErr != nil || txErr != nil {
            continue
        }
        rxBytes += lineRxBytes
        txBytes += lineTxBytes
        isFound = true
    }
    return rxBytes, txBytes, isFound
}

// GetMbitPerSecond de dos contadores de bytes y los milisegundos entre ellos; -1 si no se puede calcular (sin
// intervalo, o contador que retrocede: reinicio del host o de la interfaz).
func (c *NetworkRateCalculator) GetMbitPerSecond(previousBytes float64, bytes float64, elapsedMs float64) float64 {
    if elapsedMs <= 0 || bytes < previousBytes {
        return -1
    }
    return getRounded((bytes - previousBytes) * bitsPerByte / bitsPerMegabit / (elapsedMs / 1000))
}
