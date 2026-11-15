package shellrunner

import (
    "context"
    "errors"
    "os/exec"
    "runtime"
    "time"
)

// ShellRunResult lo que deja un comando: salida (stdout + stderr juntos), codigo de salida, duracion y si se corto
// por tiempo (entonces ExitCode = -1).
type ShellRunResult struct {
    Output     string
    ExitCode   int
    DurationMs int
    IsTimedOut bool
}

// ShellRunner lanza una linea de comando en la shell del sistema (sh -c en Linux/macOS, cmd /C en Windows) con un
// tiempo maximo. Un codigo de salida distinto de 0 no es un error del componente: va en el resultado. Solo devuelve
// error si el comando ni siquiera arranca.
type ShellRunner struct{}

func GetInstance() *ShellRunner {
    return &ShellRunner{}
}

func (s *ShellRunner) Run(commandLine string, timeoutSeconds int) (ShellRunResult, error) {
    timeoutContext, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
    defer cancel()

    command := exec.CommandContext(timeoutContext, "sh", "-c", commandLine)
    if runtime.GOOS == "windows" {
        command = exec.CommandContext(timeoutContext, "cmd", "/C", commandLine)
    }
    // si el comando deja hijos con la salida abierta, no esperar por ellos mas de 2 s tras matarlo
    command.WaitDelay = 2 * time.Second

    startedAt := time.Now()
    output, err := command.CombinedOutput()
    shellRunResult := ShellRunResult{
        Output:     string(output),
        DurationMs: int(time.Since(startedAt).Milliseconds()),
    }

    if errors.Is(timeoutContext.Err(), context.DeadlineExceeded) {
        shellRunResult.ExitCode = -1
        shellRunResult.IsTimedOut = true
        return shellRunResult, nil
    }
    var exitError *exec.ExitError
    if errors.As(err, &exitError) {
        shellRunResult.ExitCode = exitError.ExitCode()
        return shellRunResult, nil
    }
    return shellRunResult, err
}
