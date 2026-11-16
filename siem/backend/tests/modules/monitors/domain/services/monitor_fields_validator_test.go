package services

import (
    "strings"
    "testing"

    "siem-backend/tests"

    "siem-backend/src/modules/monitors/domain/services"
)

func getThreshold(value float64) *float64 {
    return &value
}

func TestValidMonitorsOfEveryType(t *testing.T) {
    monitorFieldsValidator := services.MonitorFieldsValidatorGetInstance()
    for monitorType, monitorConfig := range map[string]map[string]any{
        "cpu":         {},
        "memory":      {},
        "disk":        {"path": "/var"},
        "gpu":         {},
        "network":     {"interface": "eth0"},
        "folder":      {"path": "/var/log/nginx", "pattern": "*.log", "is_recursive": true},
        "log_pattern": {"path": "C:\\logs", "match": "(?i)failed password", "severity": "critical"},
    } {
        errors := monitorFieldsValidator.GetErrors("monitor-1", monitorType, 60, monitorConfig, getThreshold(80), getThreshold(95))
        tests.AssertEmpty(t, errors)
    }
}

func TestConfigIsNormalizedWithDefaultsAndWithoutUnknownKeys(t *testing.T) {
    monitorFieldsValidator := services.MonitorFieldsValidatorGetInstance()
    tests.AssertEquals(t, map[string]any{"path": "/"}, monitorFieldsValidator.GetNormalizedConfig("disk", map[string]any{"bogus": 1}))
    tests.AssertEquals(t, map[string]any{}, monitorFieldsValidator.GetNormalizedConfig("cpu", map[string]any{"path": "/"}))
    tests.AssertEquals(t,
        map[string]any{"path": "/var/log", "pattern": "*", "is_recursive": false},
        monitorFieldsValidator.GetNormalizedConfig("folder", map[string]any{"path": " /var/log "}),
    )
}

func TestInvalidMonitorReportsEveryError(t *testing.T) {
    monitorFieldsValidator := services.MonitorFieldsValidatorGetInstance()
    errors := strings.Join(monitorFieldsValidator.GetErrors(
        "bad name!", "log_pattern", 5,
        map[string]any{"path": "relative/path", "pattern": "[", "match": "(", "severity": "fatal"}, nil, nil,
    ), "\n")
    for _, expectedError := range []string{
        "monitor_name: required",
        "interval_seconds: must be between 10 and 86400",
        "config.path: must be an absolute path",
        "config.pattern: is not a valid glob",
        "config.match: is not a valid regex",
        "config.severity: must be empty or one of",
    } {
        tests.AssertStringContainsString(t, expectedError, errors)
    }

    unknownTypeErrors := monitorFieldsValidator.GetErrors("ok", "temperature", 60, map[string]any{}, nil, nil)
    tests.AssertStringContainsString(t, "monitor_type: must be one of", strings.Join(unknownTypeErrors, "\n"))
}

func TestThresholdRules(t *testing.T) {
    monitorFieldsValidator := services.MonitorFieldsValidatorGetInstance()
    percentErrors := strings.Join(monitorFieldsValidator.GetErrors("disk-1", "disk", 60, map[string]any{}, getThreshold(95), getThreshold(120)), "\n")
    tests.AssertStringContainsString(t, "critical_threshold: must be between 0 and 100", percentErrors)

    orderErrors := strings.Join(monitorFieldsValidator.GetErrors("disk-1", "disk", 60, map[string]any{}, getThreshold(90), getThreshold(80)), "\n")
    tests.AssertStringContainsString(t, "warning_threshold: must be lower than critical_threshold", orderErrors)

    // la red no es un %: 500 Mbit/s vale; 0 o negativo, no
    tests.AssertEmpty(t, monitorFieldsValidator.GetErrors("net-1", "network", 60, map[string]any{}, getThreshold(500), nil))
    networkErrors := strings.Join(monitorFieldsValidator.GetErrors("net-1", "network", 60, map[string]any{}, getThreshold(0), nil), "\n")
    tests.AssertStringContainsString(t, "warning_threshold: must be greater than 0", networkErrors)

    // folder no tiene valor numerico: sus umbrales se ignoran
    tests.AssertEmpty(t, monitorFieldsValidator.GetErrors("logs", "folder", 60, map[string]any{"path": "/var/log"}, getThreshold(500), getThreshold(1)))
}
