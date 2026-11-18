package services

import (
    "path/filepath"
    "slices"

    "siem-client/src/core/config"

    "siem-client/src/modules/agentconfig/domain/enums"
    "siem-client/src/modules/agentconfig/domain/entities"
)

// RemoteMonitorsMerger suma a la configuracion local (siem-client.json, si existe) los monitores que bajan del panel
// del backend. Un monitor remoto que no vale (nombre o tipo raros, path no absoluto, regex o glob invalidos) NO tumba
// el agente: se ignora y queda como aviso. Si su nombre choca con un watcher local, gana el local (y se avisa).
//   - folder / log_pattern -> un watcher mas (start_at end: no se reenvia el historico).
//   - cpu / memory / disk / gpu / network -> metricas para el modulo hostmetrics.
//   - connections -> el monitor del modulo connections (uno por agente: el segundo se ignora con aviso).
type RemoteMonitorsMerger struct{}

func RemoteMonitorsMergerGetInstance() *RemoteMonitorsMerger {
    return &RemoteMonitorsMerger{}
}

// Merge modifica agentConfigEntity (AddRemoteMonitors).
func (m *RemoteMonitorsMerger) Merge(agentConfigEntity *entities.AgentConfigEntity, remoteMonitorEntities []*entities.RemoteMonitorEntity) {
    usedNames := []string{}
    for _, watcherEntity := range agentConfigEntity.GetWatcherEntities() {
        usedNames = append(usedNames, watcherEntity.GetName())
    }

    watcherEntities := []*entities.WatcherEntity{}
    metricMonitorEntities := []*entities.RemoteMonitorEntity{}
    var connectionsMonitorEntity *entities.RemoteMonitorEntity
    warnings := []string{}
    for _, remoteMonitorEntity := range remoteMonitorEntities {
        monitorName := remoteMonitorEntity.GetMonitorName()
        if remoteMonitorEntity.GetMonitorUuid() == "" || !namePattern.MatchString(monitorName) {
            warnings = append(warnings, config.Tr(enums.AgentConfigMessageEnum.MONITOR_X_IS_NOT_VALID_X, monitorName, "monitor_uuid or monitor_name"))
            continue
        }
        if slices.Contains(usedNames, monitorName) {
            warnings = append(warnings, config.Tr(enums.AgentConfigMessageEnum.MONITOR_X_NAME_IS_ALREADY_USED, monitorName))
            continue
        }

        switch {
        case remoteMonitorEntity.IsMetric():
            metricMonitorEntities = append(metricMonitorEntities, remoteMonitorEntity)
        case remoteMonitorEntity.IsConnections():
            if connectionsMonitorEntity != nil {
                warnings = append(warnings, config.Tr(enums.AgentConfigMessageEnum.MONITOR_X_ONLY_ONE_CONNECTIONS, monitorName, connectionsMonitorEntity.GetMonitorName()))
                continue
            }
            connectionsMonitorEntity = remoteMonitorEntity
        case remoteMonitorEntity.IsLines():
            watcherEntity, problem := m.getWatcherEntityOrProblem(remoteMonitorEntity)
            if problem != "" {
                warnings = append(warnings, config.Tr(enums.AgentConfigMessageEnum.MONITOR_X_IS_NOT_VALID_X, monitorName, problem))
                continue
            }
            watcherEntities = append(watcherEntities, watcherEntity)
        default:
            warnings = append(warnings, config.Tr(enums.AgentConfigMessageEnum.MONITOR_X_IS_NOT_VALID_X, monitorName, "monitor_type "+remoteMonitorEntity.GetMonitorType()))
            continue
        }
        usedNames = append(usedNames, monitorName)
    }
    agentConfigEntity.AddRemoteMonitors(watcherEntities, metricMonitorEntities, warnings)
    agentConfigEntity.SetConnectionsMonitorEntity(connectionsMonitorEntity)
}

// getWatcherEntityOrProblem match y severity solo cuentan en log_pattern (un folder manda todas las lineas).
func (m *RemoteMonitorsMerger) getWatcherEntityOrProblem(remoteMonitorEntity *entities.RemoteMonitorEntity) (*entities.WatcherEntity, string) {
    watcherPrimitives := entities.WatcherPrimitives{
        Name:        remoteMonitorEntity.GetMonitorName(),
        Path:        remoteMonitorEntity.GetConfigString(enums.RemoteConfigKeyEnum.PATH),
        Pattern:     remoteMonitorEntity.GetConfigString(enums.RemoteConfigKeyEnum.PATTERN),
        IsRecursive: remoteMonitorEntity.GetConfigBool(enums.RemoteConfigKeyEnum.IS_RECURSIVE),
        StartAt:     enums.WatcherStartAtEnum.END,
        MonitorUuid: remoteMonitorEntity.GetMonitorUuid(),
        MonitorType: remoteMonitorEntity.GetMonitorType(),
    }
    if remoteMonitorEntity.GetMonitorType() == enums.MonitorTypeEnum.LOG_PATTERN {
        watcherPrimitives.Match = remoteMonitorEntity.GetConfigString(enums.RemoteConfigKeyEnum.MATCH)
        watcherPrimitives.Severity = remoteMonitorEntity.GetConfigString(enums.RemoteConfigKeyEnum.SEVERITY)
    }
    watcherEntity := entities.WatcherEntityFromPrimitives(watcherPrimitives)

    switch {
    case !filepath.IsAbs(watcherEntity.GetPath()):
        return nil, "path must be an absolute folder path"
    case !m.isGlobValid(watcherEntity.GetPattern()):
        return nil, "pattern is not a valid glob"
    case !watcherEntity.IsMatchValid():
        return nil, "match is not a valid regex"
    case remoteMonitorEntity.GetMonitorType() == enums.MonitorTypeEnum.LOG_PATTERN && watcherEntity.GetMatch() == "":
        return nil, "match is required"
    case watcherEntity.GetSeverity() != "" && !slices.Contains(enums.LogSeverityEnum.VALUES, watcherEntity.GetSeverity()):
        return nil, "severity is not valid"
    }
    return watcherEntity, ""
}

func (m *RemoteMonitorsMerger) isGlobValid(pattern string) bool {
    _, err := filepath.Match(pattern, "probe")
    return err == nil
}
