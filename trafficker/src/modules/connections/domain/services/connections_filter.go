package services

import (
    "slices"
    "strings"

    "trafficker/src/modules/connections/domain/entities"
    "trafficker/src/modules/connections/domain/enums"
)

// ConnectionsFilter que conexiones se quedan.
//   - Siempre (config del .env): fuera loopback (salvo APP_INCLUDE_LOOPBACK), procesos y puertos ignorados.
//   - connections-list --filter="established,remote": terminos separados por comas, TODOS se tienen que cumplir (AND).
//     Cada termino se busca como texto en protocolo, estado, direccion, extremos, proceso, ejecutable y usuario;
//     "remote" deja solo las de IP remota publica.
type ConnectionsFilter struct{}

func ConnectionsFilterGetInstance() *ConnectionsFilter {
    return &ConnectionsFilter{}
}

func (f *ConnectionsFilter) GetFilteredConnectionEntities(
    connectionEntities []*entities.ConnectionEntity,
    isLoopbackIncluded bool,
    ignoredProcesses []string,
    ignoredPorts []int,
    filterTerms []string,
) []*entities.ConnectionEntity {
    filteredConnectionEntities := []*entities.ConnectionEntity{}
    for _, connectionEntity := range connectionEntities {
        if !f.isIncludedByConfig(connectionEntity, isLoopbackIncluded, ignoredProcesses, ignoredPorts) {
            continue
        }
        if !f.matchesAllTerms(connectionEntity, filterTerms) {
            continue
        }
        filteredConnectionEntities = append(filteredConnectionEntities, connectionEntity)
    }
    return filteredConnectionEntities
}

// GetFilterTerms "Established, remote" -> [established remote].
func (f *ConnectionsFilter) GetFilterTerms(rawFilter string) []string {
    filterTerms := []string{}
    for _, filterTerm := range strings.Split(rawFilter, ",") {
        filterTerm = strings.ToLower(strings.TrimSpace(filterTerm))
        if filterTerm != "" {
            filterTerms = append(filterTerms, filterTerm)
        }
    }
    return filterTerms
}

func (f *ConnectionsFilter) isIncludedByConfig(
    connectionEntity *entities.ConnectionEntity,
    isLoopbackIncluded bool,
    ignoredProcesses []string,
    ignoredPorts []int,
) bool {
    socketEntity := connectionEntity.GetSocketEntity()
    if !isLoopbackIncluded && socketEntity.IsLoopback() {
        return false
    }
    if slices.Contains(ignoredProcesses, strings.ToLower(connectionEntity.GetProcessEntity().GetName())) {
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

func (f *ConnectionsFilter) matchesAllTerms(connectionEntity *entities.ConnectionEntity, filterTerms []string) bool {
    searchableText := connectionEntity.GetSearchableText()
    for _, filterTerm := range filterTerms {
        if filterTerm == enums.FilterTermEnum.REMOTE {
            if !connectionEntity.GetSocketEntity().HasPublicRemote() {
                return false
            }
            continue
        }
        if !strings.Contains(searchableText, filterTerm) {
            return false
        }
    }
    return true
}
