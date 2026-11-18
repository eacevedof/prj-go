package services

import (
    "fmt"
    "strings"
    "testing"

    "siem-backend/tests"

    "siem-backend/src/modules/monitors/domain/services"
)

func TestConnectionsConfigIsNormalizedFromListsOrText(t *testing.T) {
    monitorFieldsValidator := services.MonitorFieldsValidatorGetInstance()
    tests.AssertEquals(t,
        map[string]any{"is_loopback_included": false, "ignore_processes": []string{}, "ignore_ports": []int{}},
        monitorFieldsValidator.GetNormalizedConfig("connections", map[string]any{"bogus": 1}),
    )
    // como llega del JSON del panel: arrays, numeros float64, repetidos y vacios
    tests.AssertEquals(t,
        map[string]any{"is_loopback_included": true, "ignore_processes": []string{"sshd", "nginx"}, "ignore_ports": []int{22, 443}},
        monitorFieldsValidator.GetNormalizedConfig("connections", map[string]any{
            "is_loopback_included": true,
            "ignore_processes":     []any{" sshd ", "nginx", "sshd", ""},
            "ignore_ports":         []any{float64(443), "22", float64(443)},
        }),
    )
    // texto separado por comas
    tests.AssertEquals(t,
        map[string]any{"is_loopback_included": false, "ignore_processes": []string{"chrome", "Code"}, "ignore_ports": []int{53, 5353}},
        monitorFieldsValidator.GetNormalizedConfig("connections", map[string]any{
            "ignore_processes": "chrome, Code,,",
            "ignore_ports":     "5353, 53",
        }),
    )
}

func TestValidConnectionsMonitorIgnoresThresholds(t *testing.T) {
    monitorFieldsValidator := services.MonitorFieldsValidatorGetInstance()
    tests.AssertEmpty(t, monitorFieldsValidator.GetErrors(
        "conexiones", "connections", 30,
        map[string]any{"ignore_processes": []any{"sshd"}, "ignore_ports": []any{float64(22)}},
        getThreshold(500), getThreshold(1),
    ))
}

func TestInvalidConnectionsConfig(t *testing.T) {
    monitorFieldsValidator := services.MonitorFieldsValidatorGetInstance()
    errors := strings.Join(monitorFieldsValidator.GetErrors(
        "conexiones", "connections", 5,
        map[string]any{"ignore_processes": []any{"my process"}, "ignore_ports": []any{"ssh", float64(70000)}},
        nil, nil,
    ), "\n")
    for _, expectedError := range []string{
        "interval_seconds: must be between 10 and 86400",
        "config.ignore_processes: up to 50 process names",
        "config.ignore_ports: up to 50 ports between 1 and 65535",
    } {
        tests.AssertStringContainsString(t, expectedError, errors)
    }

    tooManyPorts := []any{}
    tooManyProcesses := []any{}
    for item := 1; item <= 51; item++ {
        tooManyPorts = append(tooManyPorts, float64(item))
        tooManyProcesses = append(tooManyProcesses, fmt.Sprintf("proc-%d", item))
    }
    tooManyErrors := strings.Join(monitorFieldsValidator.GetErrors(
        "conexiones", "connections", 30,
        map[string]any{"ignore_processes": tooManyProcesses, "ignore_ports": tooManyPorts}, nil, nil,
    ), "\n")
    tests.AssertStringContainsString(t, "config.ignore_processes", tooManyErrors)
    tests.AssertStringContainsString(t, "config.ignore_ports", tooManyErrors)

    longNameErrors := monitorFieldsValidator.GetErrors(
        "conexiones", "connections", 30, map[string]any{"ignore_processes": []any{strings.Repeat("a", 65)}}, nil, nil,
    )
    tests.AssertStringContainsString(t, "config.ignore_processes", strings.Join(longNameErrors, "\n"))
}

func TestConnectionsIsOnePerAgentWhateverItsConfig(t *testing.T) {
    tests.AssertEquals(t,
        getResourceKey("connections", map[string]any{}),
        getResourceKey("connections", map[string]any{"is_loopback_included": true, "ignore_ports": []any{float64(22)}}),
    )
    tests.AssertNotEquals(t, getResourceKey("connections", map[string]any{}), getResourceKey("network", map[string]any{}))
}
