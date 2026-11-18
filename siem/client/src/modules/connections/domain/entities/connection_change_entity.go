package entities

// ConnectionChangeEntity lo que ha cambiado entre dos pasadas para una conexion o un puerto en escucha: la accion
// (opened, closed, listen_opened, listen_closed), sus primitivos (los de ConnectionEntity.ToPrimitives: en un cierre,
// los guardados en el estado, porque la conexion ya no esta), cuanto duro (solo en los cierres) y si es el estado de
// partida (la primera pasada manda lo que ya estaba abierto).
type ConnectionChangeEntity struct {
    action               string
    connectionPrimitives map[string]any
    durationSeconds      int
    isInitial            bool
}

func NewConnectionChangeEntity(action string, connectionPrimitives map[string]any, durationSeconds int, isInitial bool) *ConnectionChangeEntity {
    return &ConnectionChangeEntity{
        action:               action,
        connectionPrimitives: connectionPrimitives,
        durationSeconds:      durationSeconds,
        isInitial:            isInitial,
    }
}

func (e *ConnectionChangeEntity) GetAction() string {
    return e.action
}

func (e *ConnectionChangeEntity) GetConnectionPrimitives() map[string]any {
    return e.connectionPrimitives
}

func (e *ConnectionChangeEntity) GetDurationSeconds() int {
    return e.durationSeconds
}

func (e *ConnectionChangeEntity) IsInitial() bool {
    return e.isInitial
}
