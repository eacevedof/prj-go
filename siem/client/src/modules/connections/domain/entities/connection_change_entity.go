package entities

import (
    "siem-client/src/modules/shared/application"

    "siem-client/src/modules/connections/domain/enums"
)

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

// ConnectionChangeEntityFromPrimitives action, duration_seconds e is_initial (NetConnectionAttributeEnum) y connection
// (ConnectionKeyEnum: los primitivos de la conexion).
func ConnectionChangeEntityFromPrimitives(primitives map[string]any) *ConnectionChangeEntity {
    connectionPrimitives, _ := primitives[enums.ConnectionKeyEnum.CONNECTION].(map[string]any)
    if connectionPrimitives == nil {
        connectionPrimitives = map[string]any{}
    }
    isInitial, _ := primitives[enums.NetConnectionAttributeEnum.IS_INITIAL].(bool)
    return &ConnectionChangeEntity{
        action:               application.GetTrimmedString(primitives, enums.NetConnectionAttributeEnum.ACTION),
        connectionPrimitives: connectionPrimitives,
        durationSeconds:      max(0, application.GetInt(primitives, enums.NetConnectionAttributeEnum.DURATION_SECONDS)),
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

// IsClose closed o listen_closed (llevan la duracion).
func (e *ConnectionChangeEntity) IsClose() bool {
    return e.action == enums.ConnectionActionEnum.CLOSED || e.action == enums.ConnectionActionEnum.LISTEN_CLOSED
}

// IsListen listen_opened o listen_closed (un puerto en escucha: sin extremo remoto).
func (e *ConnectionChangeEntity) IsListen() bool {
    return e.action == enums.ConnectionActionEnum.LISTEN_OPENED || e.action == enums.ConnectionActionEnum.LISTEN_CLOSED
}
