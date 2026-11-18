package collectconnectionchanges

import "siem-client/src/modules/shared/application"

type CollectConnectionChangesDto struct {
    // true = muestrea ya, toque o no (./console connections-collect --force=1)
    isForced bool
}

func CollectConnectionChangesDtoFromPrimitives(primitives map[string]any) *CollectConnectionChangesDto {
    return &CollectConnectionChangesDto{isForced: application.GetInt(primitives, "isForced") == 1}
}

func (d *CollectConnectionChangesDto) IsForced() bool {
    return d.isForced
}
