package services

import (
    "slices"
    "strings"

    "siem-client/src/modules/connections/domain/entities"
)

// ConnectionsFilter que conexiones se quedan segun la config del monitor connections del panel: fuera loopback (salvo
// is_loopback_included), los procesos ignorados (sin distinguir mayusculas) y los puertos ignorados (el local, y el
// remoto si no es una escucha).
type ConnectionsFilter struct{}

func ConnectionsFilterGetInstance() *ConnectionsFilter {
    return &ConnectionsFilter{}
}

func (f *ConnectionsFilter) GetFilteredConnectionEntities(
    connectionEntities []*entities.ConnectionEntity,
    isLoopbackIncluded bool,
    ignoredProcesses []string,
    ignoredPorts []int,
) []*entities.ConnectionEntity {
    filteredConnectionEntities := []*entities.ConnectionEntity{}
    for _, connectionEntity := range connectionEntities {
        if f.isIncluded(connectionEntity, isLoopbackIncluded, ignoredProcesses, ignoredPorts) {
            filteredConnectionEntities = append(filteredConnectionEntities, connectionEntity)
        }
    }
    return filteredConnectionEntities
}

func (f *ConnectionsFilter) isIncluded(
    connectionEntity *entities.ConnectionEntity,
    isLoopbackIncluded bool,
    ignoredProcesses []string,
    ignoredPorts []int,
) bool {
    socketEntity := connectionEntity.GetSocketEntity()
    if !isLoopbackIncluded && socketEntity.IsLoopback() {
        return false
    }
    processName := connectionEntity.GetProcessEntity().GetName()
    isIgnoredProcess := slices.ContainsFunc(ignoredProcesses, func(ignoredProcess string) bool {
        return strings.EqualFold(ignoredProcess, processName)
    })
    if isIgnoredProcess {
        return false
    }
    if slices.Contains(ignoredPorts, socketEntity.GetLocalPort()) {
        return false
    }
    if !socketEntity.IsListening() && slices.Contains(ignoredPorts, socketEntity.GetRemotePort()) {
        return false
    }
    return true
}
