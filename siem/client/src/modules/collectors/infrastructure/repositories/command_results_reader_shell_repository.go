package repositories

import (
    "siem-client/src/modules/collectors/domain/entities"
    "siem-client/src/modules/shared/infrastructure/components/shellrunner"
)

// CommandResultsReaderShellRepository el resultado de un comando programado: la fuente de datos es la shell.
type CommandResultsReaderShellRepository struct {
    shellRunner *shellrunner.ShellRunner
}

func CommandResultsReaderShellRepositoryGetInstance() *CommandResultsReaderShellRepository {
    return &CommandResultsReaderShellRepository{shellRunner: shellrunner.GetInstance()}
}

// GetCommandResult error solo si el comando no llega a arrancar; un exit code distinto de 0 va en la entidad.
func (r *CommandResultsReaderShellRepository) GetCommandResult(commandLine string, timeoutSeconds int) (*entities.CommandResultEntity, error) {
    shellRunResult, err := r.shellRunner.Run(commandLine, timeoutSeconds)
    if err != nil {
        return nil, err
    }
    isTimedOut := 0
    if shellRunResult.IsTimedOut {
        isTimedOut = 1
    }
    return entities.CommandResultEntityFromPrimitives(map[string]any{
        "output":     shellRunResult.Output,
        "exitCode":   shellRunResult.ExitCode,
        "durationMs": shellRunResult.DurationMs,
        "isTimedOut": isTimedOut,
    }), nil
}
