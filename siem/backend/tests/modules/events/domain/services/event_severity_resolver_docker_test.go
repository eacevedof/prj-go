package services

import (
    "testing"

    "siem-backend/src/modules/events/domain/services"
    "siem-backend/tests"
)

func TestDockerStateSeverity(t *testing.T) {
    eventSeverityResolver := services.EventSeverityResolverGetInstance()
    for caseName, testCase := range map[string]struct {
        attributes       map[string]any
        expectedSeverity string
    }{
        "oom killed":         {map[string]any{"status": "exited", "exit_code": float64(137), "is_oom_killed": true}, "critical"},
        "dead":               {map[string]any{"status": "dead"}, "critical"},
        "exited with error":  {map[string]any{"status": "exited", "exit_code": float64(137)}, "error"},
        "unhealthy":          {map[string]any{"status": "running", "health": "unhealthy"}, "error"},
        "exited cleanly":     {map[string]any{"status": "exited", "exit_code": float64(0)}, "warning"},
        "restarting":         {map[string]any{"status": "restarting"}, "warning"},
        "restarted":          {map[string]any{"status": "running", "change_kind": "restarted"}, "warning"},
        "running and healthy": {map[string]any{"status": "running", "health": "healthy"}, "info"},
        "removed":            {map[string]any{"status": "removed", "change_kind": "removed"}, "info"},
    } {
        tests.AssertEquals(t, testCase.expectedSeverity, eventSeverityResolver.GetSeverity("docker_state", "x", 0, "", testCase.attributes), caseName)
    }
}

func TestDockerStatsSeverityByThresholds(t *testing.T) {
    eventSeverityResolver := services.EventSeverityResolverGetInstance()
    tests.AssertEquals(t, "debug", eventSeverityResolver.GetSeverity("docker_stats", "x", 0, "", map[string]any{"cpu_percent": 12.5, "memory_percent": 40.0}))
    tests.AssertEquals(t, "warning", eventSeverityResolver.GetSeverity("docker_stats", "x", 0, "", map[string]any{"cpu_percent": 95.0, "memory_percent": 40.0}))
    tests.AssertEquals(t, "warning", eventSeverityResolver.GetSeverity("docker_stats", "x", 0, "", map[string]any{"memory_percent": 88.0}))
    tests.AssertEquals(t, "error", eventSeverityResolver.GetSeverity("docker_stats", "x", 0, "", map[string]any{"memory_percent": 97.0}))
    // primera muestra: sin cpu_percent (no hay muestra anterior con la que comparar)
    tests.AssertEquals(t, "debug", eventSeverityResolver.GetSeverity("docker_stats", "x", 0, "", map[string]any{"memory_percent": 10.0}))
}

func TestDockerEventSeverityByAction(t *testing.T) {
    eventSeverityResolver := services.EventSeverityResolverGetInstance()
    tests.AssertEquals(t, "critical", eventSeverityResolver.GetSeverity("docker_event", "x", 0, "", map[string]any{"action": "oom"}))
    tests.AssertEquals(t, "error", eventSeverityResolver.GetSeverity("docker_event", "x", 0, "", map[string]any{"action": "die", "exit_code": float64(1)}))
    tests.AssertEquals(t, "info", eventSeverityResolver.GetSeverity("docker_event", "x", 0, "", map[string]any{"action": "die", "exit_code": float64(0)}))
    tests.AssertEquals(t, "error", eventSeverityResolver.GetSeverity("docker_event", "x", 0, "", map[string]any{"action": "health_status: unhealthy"}))
    tests.AssertEquals(t, "warning", eventSeverityResolver.GetSeverity("docker_event", "x", 0, "", map[string]any{"action": "kill"}))
    tests.AssertEquals(t, "info", eventSeverityResolver.GetSeverity("docker_event", "x", 0, "", map[string]any{"action": "start"}))
}

func TestDockerLogUsesKeywords(t *testing.T) {
    eventSeverityResolver := services.EventSeverityResolverGetInstance()
    tests.AssertEquals(t, "error", eventSeverityResolver.GetSeverity("docker_log", "ERROR connection refused", 0, "", nil))
    tests.AssertEquals(t, "info", eventSeverityResolver.GetSeverity("docker_log", "listening on :8080", 0, "", nil))
}
