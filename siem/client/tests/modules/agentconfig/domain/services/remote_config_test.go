package services

import (
    "strings"
    "testing"

    "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/agentconfig/domain/services"
    "siem-client/tests"
)

func TestAgentConfigUrlComesFromTheIngestUrl(t *testing.T) {
    remoteConfigParser := services.RemoteConfigParserGetInstance()
    tests.AssertEquals(t, "https://siem.example.com/api/v1/agent-config", remoteConfigParser.GetAgentConfigUrl("https://siem.example.com/api/v1/ingest-batches"))
    tests.AssertEquals(t, "http://back:8080/api/v1/agent-config", remoteConfigParser.GetAgentConfigUrl(" http://back:8080/api/v1/ingest-batches/ "))
    // una url de ingesta con otro final: se cambia su ultimo tramo
    tests.AssertEquals(t, "http://back:8080/api/v1/agent-config", remoteConfigParser.GetAgentConfigUrl("http://back:8080/api/v1/ingest"))
    tests.AssertEquals(t, "http://back:8080/agent-config", remoteConfigParser.GetAgentConfigUrl("http://back:8080"))
}

func TestRemoteConfigResponseIsValidated(t *testing.T) {
    remoteConfigParser := services.RemoteConfigParserGetInstance()
    remoteConfig, err := remoteConfigParser.GetRemoteConfigOrFail(`{"code":200,"status":true,"errors":[],"data":{"config_hash":"abc","monitors":[
        {"monitor_uuid":"monitor-1","monitor_type":"disk","monitor_name":"disco","interval_seconds":120,"config":{"path":"/data"}}
    ]}}`)
    tests.AssertNil(t, err)
    tests.AssertEquals(t, "abc", remoteConfig["config_hash"])
    remoteMonitorEntities := remoteConfigParser.GetRemoteMonitorEntities(remoteConfig)
    tests.AssertEquals(t, 1, len(remoteMonitorEntities))
    tests.AssertEquals(t, "disco", remoteMonitorEntities[0].GetMonitorName())
    tests.AssertEquals(t, 120, remoteMonitorEntities[0].GetIntervalSeconds())
    tests.AssertEquals(t, "/data", remoteMonitorEntities[0].GetConfigString("path"))
    tests.AssertEquals(t, true, remoteMonitorEntities[0].IsMetric())

    // respuestas que NO deben sustituir a la config guardada
    for _, badResponse := range []string{"<html>502</html>", `{"data":null}`, `{"data":{"monitors":"x"}}`, `{"data":{"monitors":[1]}}`} {
        _, err := remoteConfigParser.GetRemoteConfigOrFail(badResponse)
        tests.AssertNotEquals(t, nil, err)
    }
}

func getMonitorEntity(monitorUuid string, monitorType string, monitorName string, monitorConfig map[string]any) *entities.RemoteMonitorEntity {
    return entities.RemoteMonitorEntityFromPrimitives(map[string]any{
        "monitor_uuid": monitorUuid, "monitor_type": monitorType, "monitor_name": monitorName, "config": monitorConfig,
    })
}

func TestRemoteMonitorsBecomeWatchersAndMetrics(t *testing.T) {
    agentConfigEntity := entities.AgentConfigEntityFromPrimitives(entities.AgentConfigPrimitives{
        Watchers: []entities.WatcherPrimitives{{Name: "app-logs", Path: getAbsolutePath()}},
    })
    services.RemoteMonitorsMergerGetInstance().Merge(agentConfigEntity, []*entities.RemoteMonitorEntity{
        getMonitorEntity("monitor-1", "cpu", "cpu", nil),
        getMonitorEntity("monitor-2", "folder", "nginx", map[string]any{"path": getAbsolutePath(), "pattern": "*.log", "is_recursive": true}),
        getMonitorEntity("monitor-3", "log_pattern", "ssh-fallos", map[string]any{"path": getAbsolutePath(), "match": "Failed password", "severity": "ERROR"}),
        // todo lo que sigue se ignora con aviso
        getMonitorEntity("monitor-4", "folder", "app-logs", map[string]any{"path": getAbsolutePath()}),
        getMonitorEntity("monitor-5", "folder", "relativa", map[string]any{"path": "var/log"}),
        getMonitorEntity("monitor-6", "log_pattern", "regex-rota", map[string]any{"path": getAbsolutePath(), "match": "("}),
        getMonitorEntity("monitor-7", "log_pattern", "sin-match", map[string]any{"path": getAbsolutePath()}),
        getMonitorEntity("monitor-8", "teletransporte", "raro", nil),
        getMonitorEntity("", "cpu", "sin-uuid", nil),
        getMonitorEntity("monitor-9", "memory", "cpu", nil),
    })

    tests.AssertEquals(t, 1, len(agentConfigEntity.GetMetricMonitorEntities()))
    watcherEntities := agentConfigEntity.GetWatcherEntities()
    tests.AssertEquals(t, 3, len(watcherEntities))

    folderWatcher := watcherEntities[1]
    tests.AssertEquals(t, "nginx", folderWatcher.GetName())
    tests.AssertEquals(t, "monitor-2", folderWatcher.GetMonitorUuid())
    tests.AssertEquals(t, "folder", folderWatcher.GetMonitorType())
    tests.AssertEquals(t, true, folderWatcher.IsRecursive())
    tests.AssertEquals(t, true, folderWatcher.IsStartAtEnd())
    tests.AssertEquals(t, true, folderWatcher.IsLineMatching("cualquier linea"))

    patternWatcher := watcherEntities[2]
    tests.AssertEquals(t, "error", patternWatcher.GetSeverity())
    tests.AssertEquals(t, "*", patternWatcher.GetPattern())
    tests.AssertEquals(t, true, patternWatcher.IsLineMatching("sshd: Failed password for root"))
    tests.AssertEquals(t, false, patternWatcher.IsLineMatching("sshd: Accepted publickey"))

    warnings := strings.Join(agentConfigEntity.GetRemoteConfigWarnings(), "\n")
    tests.AssertEquals(t, 7, len(agentConfigEntity.GetRemoteConfigWarnings()))
    for _, expectedWarning := range []string{
        "remote monitor app-logs ignored: its name is already used",
        "remote monitor relativa ignored: path must be an absolute folder path",
        "remote monitor regex-rota ignored: match is not a valid regex",
        "remote monitor sin-match ignored: match is required",
        "remote monitor raro ignored: monitor_type teletransporte",
        "remote monitor sin-uuid ignored: monitor_uuid or monitor_name",
        "remote monitor cpu ignored: its name is already used",
    } {
        tests.AssertStringContainsString(t, expectedWarning, warnings)
    }
}

func TestNothingToWatchIsAllowedOnlyWithRemoteConfig(t *testing.T) {
    validator := services.AgentConfigValidatorGetInstance()
    emptyConfig := entities.AgentConfigEntityFromPrimitives(entities.AgentConfigPrimitives{})
    tests.AssertEmpty(t, validator.GetErrorsWithRemoteConfig(emptyConfig, endpointUrl, apiToken, true))
    tests.AssertEquals(t, []string{"watchers, commands or docker: at least one is required"}, validator.GetErrorsWithRemoteConfig(emptyConfig, endpointUrl, apiToken, false))
}

func TestLocalWatcherMatchAndSeverityAreValidated(t *testing.T) {
    agentConfigEntity := entities.AgentConfigEntityFromPrimitives(entities.AgentConfigPrimitives{
        Watchers: []entities.WatcherPrimitives{{Name: "logs", Path: getAbsolutePath(), Match: "[", Severity: "grave"}},
    })
    errors := strings.Join(services.AgentConfigValidatorGetInstance().GetErrors(agentConfigEntity, endpointUrl, apiToken), "\n")
    tests.AssertStringContainsString(t, "watcher logs: match is not a valid regex", errors)
    tests.AssertStringContainsString(t, "watcher logs: severity must be debug, info, warning, error or critical", errors)
}
