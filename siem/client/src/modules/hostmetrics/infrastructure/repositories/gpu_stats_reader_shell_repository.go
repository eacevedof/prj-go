package repositories

import (
    "errors"
    "strconv"
    "strings"

    "siem-client/src/modules/shared/infrastructure/components/shellrunner"
)

const (
    NVIDIA_SMI_COMMAND         = "nvidia-smi --query-gpu=name,utilization.gpu,memory.used,memory.total,temperature.gpu --format=csv,noheader,nounits"
    nvidiaSmiTimeoutSeconds    = 10
    nvidiaSmiErrorMaxCharacters = 200
)

// GpuStatsReaderShellRepository las GPUs NVIDIA del host por nvidia-smi (el unico origen soportado). Sin la
// herramienta (host sin GPU o sin drivers) la shell devuelve un codigo distinto de 0: es un error.
type GpuStatsReaderShellRepository struct {
    shellRunner *shellrunner.ShellRunner
}

func GpuStatsReaderShellRepositoryGetInstance() *GpuStatsReaderShellRepository {
    return &GpuStatsReaderShellRepository{shellRunner: shellrunner.GetInstance()}
}

// GetNvidiaSmiOutput una linea csv por GPU. Error (con la primera linea de lo que dijo la shell) si no arranca, tarda
// mas de 10 s o sale con codigo distinto de 0.
func (r *GpuStatsReaderShellRepository) GetNvidiaSmiOutput() (string, error) {
    shellRunResult, err := r.shellRunner.Run(NVIDIA_SMI_COMMAND, nvidiaSmiTimeoutSeconds)
    if err != nil {
        return "", err
    }
    if shellRunResult.IsTimedOut {
        return "", errors.New("timeout")
    }
    if shellRunResult.ExitCode != 0 {
        firstLine, _, _ := strings.Cut(strings.TrimSpace(shellRunResult.Output), "\n")
        if len(firstLine) > nvidiaSmiErrorMaxCharacters {
            firstLine = firstLine[:nvidiaSmiErrorMaxCharacters]
        }
        if firstLine == "" {
            firstLine = "exit code " + strconv.Itoa(shellRunResult.ExitCode)
        }
        return "", errors.New(strings.TrimSpace(firstLine))
    }
    return shellRunResult.Output, nil
}
