package runduecommands

import "siem-client/src/modules/shared/application"

// RunDueCommandsPrimitives lo que imprime commands-run (claves en este orden).
type RunDueCommandsPrimitives struct {
    RanCommands    int              `json:"ran_commands"`
    FailedCommands int              `json:"failed_commands"`
    SpooledBatches int              `json:"spooled_batches"`
    Results        []map[string]any `json:"results"`
}

type RunDueCommandsResultDto struct {
    ranCommands    int
    failedCommands int
    spooledBatches int
    // los eventos command_result de esta pasada (solo primitivos)
    results []map[string]any
}

func RunDueCommandsResultDtoFromPrimitives(primitives map[string]any) *RunDueCommandsResultDto {
    results, _ := primitives["results"].([]map[string]any)
    if results == nil {
        results = []map[string]any{}
    }
    return &RunDueCommandsResultDto{
        ranCommands:    application.GetInt(primitives, "ranCommands"),
        failedCommands: application.GetInt(primitives, "failedCommands"),
        spooledBatches: application.GetInt(primitives, "spooledBatches"),
        results:        results,
    }
}

func (d *RunDueCommandsResultDto) GetRanCommands() int {
    return d.ranCommands
}

func (d *RunDueCommandsResultDto) GetFailedCommands() int {
    return d.failedCommands
}

func (d *RunDueCommandsResultDto) ToArray() RunDueCommandsPrimitives {
    return RunDueCommandsPrimitives{
        RanCommands:    d.ranCommands,
        FailedCommands: d.failedCommands,
        SpooledBatches: d.spooledBatches,
        Results:        d.results,
    }
}
