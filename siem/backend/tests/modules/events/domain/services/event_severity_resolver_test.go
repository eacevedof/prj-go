package services

import (
    "testing"

    "siem-backend/tests"

    "siem-backend/src/modules/events/domain/services"
)

func TestDeclaredSeverityWins(t *testing.T) {
    severity := services.EventSeverityResolverGetInstance().GetSeverity("log_line", "ERROR boom", 0, "Warning", nil)
    tests.AssertEquals(t, "warning", severity)
}

func TestLogLineSeverityByKeywords(t *testing.T) {
    eventSeverityResolver := services.EventSeverityResolverGetInstance()
    for message, expectedSeverity := range map[string]string{
        "kernel: Out of memory: killed process 123": "critical",
        "PANIC: something exploded":                 "critical",
        "2026-11-15 ERROR db connection refused":    "error",
        "Permission denied (publickey)":             "error",
        "WARN cache miss":                           "warning",
        "request timed out after 30s":               "warning",
        "DEBUG payload dumped":                      "debug",
        "user logged in":                            "info",
        "terrorism is not a keyword":                "info",
    } {
        tests.AssertEquals(t, expectedSeverity, eventSeverityResolver.GetSeverity("log_line", message, 0, "", nil), message)
    }
}

func TestCommandResultSeverityByExitCode(t *testing.T) {
    eventSeverityResolver := services.EventSeverityResolverGetInstance()
    tests.AssertEquals(t, "info", eventSeverityResolver.GetSeverity("command_result", "ERROR in the output", 0, "", nil))
    tests.AssertEquals(t, "error", eventSeverityResolver.GetSeverity("command_result", "all good", 2, "", nil))
    tests.AssertEquals(t, "error", eventSeverityResolver.GetSeverity("command_result", "[timeout after 5 s]", -1, "", nil))
}
