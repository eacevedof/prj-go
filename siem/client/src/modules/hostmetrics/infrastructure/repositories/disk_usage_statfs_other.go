//go:build !linux

package repositories

import (
    "errors"
    "runtime"
)

// getStatfsBytes fuera de Linux la metrica de disco no esta implementada (el agente la avisa una vez como no disponible).
func getStatfsBytes(path string) (float64, float64, float64, error) {
    return 0, 0, 0, errors.New("disk usage is not supported on " + runtime.GOOS)
}
