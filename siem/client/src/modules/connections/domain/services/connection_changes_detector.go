package services

import (
    "maps"
    "slices"

    "siem-client/src/modules/connections/domain/entities"
    "siem-client/src/modules/connections/domain/enums"
    "siem-client/src/modules/shared/application"
)

// ConnectionChangesDetector compara la foto de ahora con la anterior (state/connections.json) y saca SOLO los cambios,
// nunca la foto completa:
//   - una conexion (o escucha) que pasa el filtro y no estaba -> opened (listen_opened). Sin foto anterior (primera
//     pasada, o monitor nuevo) son todas, marcadas is_initial: el estado de partida.
//   - una que estaba y ya no existe -> closed (listen_closed), con su duracion.
//   - una que estaba, sigue existiendo pero ya no pasa el filtro (se cambio la config en el panel) -> se olvida sin
//     evento: no se ha cerrado.
//
// La foto siguiente guarda los primitivos de cada conexion y opened_at_unix (para la duracion del cierre).
type ConnectionChangesDetector struct{}

func ConnectionChangesDetectorGetInstance() *ConnectionChangesDetector {
    return &ConnectionChangesDetector{}
}

// GetChangesAndNextConnections previousConnections: clave -> primitivos de la foto anterior (nil = no hay foto).
// allConnectionEntities: todas las de ahora, sin filtrar (para distinguir "cerrada" de "ya no se vigila").
func (d *ConnectionChangesDetector) GetChangesAndNextConnections(
    previousConnections map[string]any,
    allConnectionEntities []*entities.ConnectionEntity,
    filteredConnectionEntities []*entities.ConnectionEntity,
    nowUnix int64,
) ([]*entities.ConnectionChangeEntity, map[string]any) {
    isInitial := previousConnections == nil
    changes := []*entities.ConnectionChangeEntity{}
    nextConnections := map[string]any{}

    existingKeys := map[string]bool{}
    for _, connectionEntity := range allConnectionEntities {
        existingKeys[connectionEntity.GetKey()] = true
    }

    // cierres primero (orden estable por clave): asi un puerto que cierra y abre en la misma pasada se lee bien
    for _, key := range slices.Sorted(maps.Keys(previousConnections)) {
        if existingKeys[key] {
            continue
        }
        connectionPrimitives, _ := previousConnections[key].(map[string]any)
        openedAtUnix := int64(application.GetInt(connectionPrimitives, enums.ConnectionsStateEnum.OPENED_AT_UNIX))
        changes = append(changes, entities.NewConnectionChangeEntity(
            d.getCloseAction(connectionPrimitives),
            d.getWithoutStateKeys(connectionPrimitives),
            int(max(0, nowUnix-openedAtUnix)),
            false,
        ))
    }

    for _, connectionEntity := range filteredConnectionEntities {
        key := connectionEntity.GetKey()
        if previousConnectionPrimitives, wasOpen := previousConnections[key].(map[string]any); wasOpen {
            nextConnections[key] = previousConnectionPrimitives
            continue
        }
        connectionPrimitives := connectionEntity.ToPrimitives()
        changes = append(changes, entities.NewConnectionChangeEntity(
            d.getOpenAction(connectionEntity.IsListening()),
            connectionPrimitives,
            0,
            isInitial,
        ))
        nextConnectionPrimitives := maps.Clone(connectionPrimitives)
        nextConnectionPrimitives[enums.ConnectionsStateEnum.OPENED_AT_UNIX] = nowUnix
        nextConnections[key] = nextConnectionPrimitives
    }
    return changes, nextConnections
}

func (d *ConnectionChangesDetector) getOpenAction(isListening bool) string {
    if isListening {
        return enums.ConnectionActionEnum.LISTEN_OPENED
    }
    return enums.ConnectionActionEnum.OPENED
}

func (d *ConnectionChangesDetector) getCloseAction(connectionPrimitives map[string]any) string {
    if application.GetTrimmedString(connectionPrimitives, enums.ConnectionKeyEnum.DIRECTION) == enums.DirectionEnum.LISTEN {
        return enums.ConnectionActionEnum.LISTEN_CLOSED
    }
    return enums.ConnectionActionEnum.CLOSED
}

// getWithoutStateKeys los primitivos guardados sin lo que es solo del estado (opened_at_unix).
func (d *ConnectionChangesDetector) getWithoutStateKeys(connectionPrimitives map[string]any) map[string]any {
    withoutStateKeys := maps.Clone(connectionPrimitives)
    if withoutStateKeys == nil {
        return map[string]any{}
    }
    delete(withoutStateKeys, enums.ConnectionsStateEnum.OPENED_AT_UNIX)
    return withoutStateKeys
}
