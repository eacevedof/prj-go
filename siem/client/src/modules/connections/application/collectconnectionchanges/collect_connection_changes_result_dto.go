package collectconnectionchanges

import "siem-client/src/modules/shared/application"

// CollectConnectionChangesPrimitives lo que imprime connections-collect (claves en este orden).
type CollectConnectionChangesPrimitives struct {
    IsDisabled          bool             `json:"is_disabled"`
    IsSampled           bool             `json:"is_sampled"`
    IsUnavailable       bool             `json:"is_unavailable"`
    Connections         int              `json:"connections"`
    Opened              int              `json:"opened"`
    Closed              int              `json:"closed"`
    UnreadableProcesses int              `json:"unreadable_processes"`
    SpooledBatches      int              `json:"spooled_batches"`
    Results             []map[string]any `json:"results"`
}

type CollectConnectionChangesResultDto struct {
    // sin monitor connections en la config: no se muestrea (y el estado se ha borrado)
    isDisabled bool
    // se ha tomado la foto en esta pasada (si no tocaba por el intervalo, false)
    isSampled bool
    // se ha mandado el aviso (unico) de que en este sistema no se puede muestrear
    isUnavailable bool
    // conexiones vigiladas despues de esta pasada (las que pasan el filtro)
    connections int
    opened      int
    closed      int
    // procesos cuyos sockets no se han podido leer (sin root): salen con proceso "?"
    unreadableProcesses int
    spooledBatches      int
    // los eventos net_connection de esta pasada (solo primitivos)
    results []map[string]any
}

func CollectConnectionChangesResultDtoFromPrimitives(primitives map[string]any) *CollectConnectionChangesResultDto {
    results, _ := primitives["results"].([]map[string]any)
    if results == nil {
        results = []map[string]any{}
    }
    return &CollectConnectionChangesResultDto{
        isDisabled:          application.GetInt(primitives, "isDisabled") == 1,
        isSampled:           application.GetInt(primitives, "isSampled") == 1,
        isUnavailable:       application.GetInt(primitives, "isUnavailable") == 1,
        connections:         application.GetInt(primitives, "connections"),
        opened:              application.GetInt(primitives, "opened"),
        closed:              application.GetInt(primitives, "closed"),
        unreadableProcesses: application.GetInt(primitives, "unreadableProcesses"),
        spooledBatches:      application.GetInt(primitives, "spooledBatches"),
        results:             results,
    }
}

func (d *CollectConnectionChangesResultDto) IsDisabled() bool {
    return d.isDisabled
}

func (d *CollectConnectionChangesResultDto) IsSampled() bool {
    return d.isSampled
}

func (d *CollectConnectionChangesResultDto) IsUnavailable() bool {
    return d.isUnavailable
}

func (d *CollectConnectionChangesResultDto) GetConnections() int {
    return d.connections
}

func (d *CollectConnectionChangesResultDto) GetOpened() int {
    return d.opened
}

func (d *CollectConnectionChangesResultDto) GetClosed() int {
    return d.closed
}

func (d *CollectConnectionChangesResultDto) GetUnreadableProcesses() int {
    return d.unreadableProcesses
}

func (d *CollectConnectionChangesResultDto) ToArray() CollectConnectionChangesPrimitives {
    return CollectConnectionChangesPrimitives{
        IsDisabled:          d.isDisabled,
        IsSampled:           d.isSampled,
        IsUnavailable:       d.isUnavailable,
        Connections:         d.connections,
        Opened:              d.opened,
        Closed:              d.closed,
        UnreadableProcesses: d.unreadableProcesses,
        SpooledBatches:      d.spooledBatches,
        Results:             d.results,
    }
}
