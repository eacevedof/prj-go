package runduecommands

import "siem-client/src/modules/shared/application"

type RunDueCommandsDto struct {
    // true = lanza todos los comandos ya, toque o no (./console commands-run --force=1)
    isForced bool
}

func RunDueCommandsDtoFromPrimitives(primitives map[string]any) *RunDueCommandsDto {
    return &RunDueCommandsDto{isForced: application.GetInt(primitives, "isForced") == 1}
}

func (d *RunDueCommandsDto) IsForced() bool {
    return d.isForced
}
