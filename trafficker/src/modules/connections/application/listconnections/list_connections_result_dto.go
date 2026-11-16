package listconnections

import (
    "trafficker/src/modules/connections/domain/entities"
    "trafficker/src/modules/connections/domain/enums"
)

// ListConnectionsPrimitives lo que imprime connections-list --format=json (claves en este orden).
type ListConnectionsPrimitives struct {
    Total               int              `json:"total"`
    Listening           int              `json:"listening"`
    Inbound             int              `json:"inbound"`
    Outbound            int              `json:"outbound"`
    UnreadableProcesses int              `json:"unreadable_processes"`
    Connections         []map[string]any `json:"connections"`
}

type ListConnectionsResultDto struct {
    connectionEntities  []*entities.ConnectionEntity
    unreadableProcesses int
}

func NewListConnectionsResultDto(connectionEntities []*entities.ConnectionEntity, unreadableProcesses int) *ListConnectionsResultDto {
    return &ListConnectionsResultDto{connectionEntities: connectionEntities, unreadableProcesses: unreadableProcesses}
}

func (d *ListConnectionsResultDto) GetConnectionEntities() []*entities.ConnectionEntity {
    return d.connectionEntities
}

// GetUnreadableProcesses procesos cuyos fd no se han podido leer: sus sockets salen con proceso "?".
func (d *ListConnectionsResultDto) GetUnreadableProcesses() int {
    return d.unreadableProcesses
}

func (d *ListConnectionsResultDto) GetCountByDirection(direction string) int {
    count := 0
    for _, connectionEntity := range d.connectionEntities {
        if connectionEntity.GetDirection() == direction {
            count++
        }
    }
    return count
}

func (d *ListConnectionsResultDto) ToArray() ListConnectionsPrimitives {
    connections := make([]map[string]any, 0, len(d.connectionEntities))
    for _, connectionEntity := range d.connectionEntities {
        connections = append(connections, connectionEntity.ToPrimitives())
    }
    return ListConnectionsPrimitives{
        Total:               len(d.connectionEntities),
        Listening:           d.GetCountByDirection(enums.DirectionEnum.LISTEN),
        Inbound:             d.GetCountByDirection(enums.DirectionEnum.INBOUND),
        Outbound:            d.GetCountByDirection(enums.DirectionEnum.OUTBOUND),
        UnreadableProcesses: d.unreadableProcesses,
        Connections:         connections,
    }
}
