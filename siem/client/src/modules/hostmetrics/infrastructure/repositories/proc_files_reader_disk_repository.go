package repositories

import (
    "os"
    "runtime"
)

const (
    PROC_STAT    = "/proc/stat"
    PROC_MEMINFO = "/proc/meminfo"
    PROC_NET_DEV = "/proc/net/dev"
    PROC_LOADAVG = "/proc/loadavg"
)

// ProcFilesReaderDiskRepository los contadores del kernel de Linux en /proc (cpu, memoria, red, carga). Dentro de un
// contenedor /proc/stat y /proc/meminfo son los del host (en Docker Desktop, los de su maquina virtual).
type ProcFilesReaderDiskRepository struct{}

func ProcFilesReaderDiskRepositoryGetInstance() *ProcFilesReaderDiskRepository {
    return &ProcFilesReaderDiskRepository{}
}

// IsSupported /proc solo existe en Linux: en Windows y macOS estas metricas no estan disponibles.
func (r *ProcFilesReaderDiskRepository) IsSupported() bool {
    return runtime.GOOS == "linux"
}

func (r *ProcFilesReaderDiskRepository) GetOsName() string {
    return runtime.GOOS
}

func (r *ProcFilesReaderDiskRepository) GetProcFileContent(procFilePath string) (string, error) {
    content, err := os.ReadFile(procFilePath)
    if err != nil {
        return "", err
    }
    return string(content), nil
}
