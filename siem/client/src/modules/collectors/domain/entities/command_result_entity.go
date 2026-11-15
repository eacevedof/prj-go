package entities

import "siem-client/src/modules/shared/application"

// CommandResultEntity lo que dejo una ejecucion de un comando programado.
type CommandResultEntity struct {
    output     string
    exitCode   int
    durationMs int
    isTimedOut bool
}

func CommandResultEntityFromPrimitives(primitives map[string]any) *CommandResultEntity {
    return &CommandResultEntity{
        // sin trim aqui: lo hace quien monta el evento (la salida puede empezar por espacios significativos)
        output:     application.ToString(primitives["output"]),
        exitCode:   application.GetInt(primitives, "exitCode"),
        durationMs: application.GetInt(primitives, "durationMs"),
        isTimedOut: application.GetInt(primitives, "isTimedOut") == 1,
    }
}

func (e *CommandResultEntity) GetOutput() string {
    return e.output
}

func (e *CommandResultEntity) GetExitCode() int {
    return e.exitCode
}

func (e *CommandResultEntity) GetDurationMs() int {
    return e.durationMs
}

func (e *CommandResultEntity) IsTimedOut() bool {
    return e.isTimedOut
}
