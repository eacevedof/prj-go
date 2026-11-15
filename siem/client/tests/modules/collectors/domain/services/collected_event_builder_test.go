package services

import (
    "testing"

    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/collectors/domain/services"
    "siem-client/tests"
)

func TestLocalWatcherLineHasNoSeverityNorAttributes(t *testing.T) {
    watcherEntity := agentconfigentities.WatcherEntityFromPrimitives(agentconfigentities.WatcherPrimitives{Name: "app-logs", Path: "/var/log/app"})
    event := services.CollectedEventBuilderGetInstance().GetWatcherLogLineEvent(watcherEntity, "/var/log/app/a.log", "hello", "2026-11-15 18:00:00")
    tests.AssertEquals(t, "log_line", event["event_type"])
    tests.AssertEquals(t, "app-logs", event["source_name"])
    _, hasSeverity := event["severity"]
    _, hasAttributes := event["attributes"]
    tests.AssertEquals(t, false, hasSeverity)
    tests.AssertEquals(t, false, hasAttributes)
}

func TestMonitorWatcherLineCarriesMonitorAndDeclaredSeverity(t *testing.T) {
    watcherEntity := agentconfigentities.WatcherEntityFromPrimitives(agentconfigentities.WatcherPrimitives{
        Name: "ssh-fallos", Path: "/var/log", Match: "Failed", Severity: "error", MonitorUuid: "monitor-3", MonitorType: "log_pattern",
    })
    event := services.CollectedEventBuilderGetInstance().GetWatcherLogLineEvent(watcherEntity, "/var/log/auth.log", "Failed password", "2026-11-15 18:00:00")
    tests.AssertEquals(t, "ssh-fallos", event["source_name"])
    tests.AssertEquals(t, "/var/log/auth.log", event["file_path"])
    tests.AssertEquals(t, "error", event["severity"])
    tests.AssertEquals(t, map[string]any{"monitor_uuid": "monitor-3", "monitor_type": "log_pattern"}, event["attributes"])
}
