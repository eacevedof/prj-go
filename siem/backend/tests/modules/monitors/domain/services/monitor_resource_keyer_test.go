package services

import (
    "testing"

    "siem-backend/src/modules/monitors/domain/services"
    "siem-backend/tests"
)

func getResourceKey(monitorType string, monitorConfig map[string]any) string {
    normalizedConfig := services.MonitorFieldsValidatorGetInstance().GetNormalizedConfig(monitorType, monitorConfig)
    return services.MonitorResourceKeyerGetInstance().GetResourceKey(monitorType, normalizedConfig)
}

func TestMetricsWithoutConfigAreOnePerType(t *testing.T) {
    tests.AssertEquals(t, getResourceKey("cpu", map[string]any{}), getResourceKey("cpu", map[string]any{"path": "/x"}))
    tests.AssertNotEquals(t, getResourceKey("cpu", map[string]any{}), getResourceKey("memory", map[string]any{}))
}

func TestDiskIsTheSameWithDefaultOrTrailingSeparator(t *testing.T) {
    tests.AssertEquals(t, getResourceKey("disk", map[string]any{}), getResourceKey("disk", map[string]any{"path": "/"}))
    tests.AssertEquals(t, getResourceKey("disk", map[string]any{"path": "/var"}), getResourceKey("disk", map[string]any{"path": " /var/ "}))
    tests.AssertNotEquals(t, getResourceKey("disk", map[string]any{"path": "/"}), getResourceKey("disk", map[string]any{"path": "/var"}))
}

func TestWindowsPathsIgnoreCaseAndKeepTheirRoot(t *testing.T) {
    tests.AssertEquals(t, getResourceKey("disk", map[string]any{"path": "C:\\"}), getResourceKey("disk", map[string]any{"path": "c:\\"}))
    tests.AssertNotEquals(t, getResourceKey("disk", map[string]any{"path": "C:\\"}), getResourceKey("disk", map[string]any{"path": "C:"}))
    tests.AssertEquals(t,
        getResourceKey("folder", map[string]any{"path": `C:\Logs\`}),
        getResourceKey("folder", map[string]any{"path": `c:\logs`}),
    )
    // en Linux las mayusculas cuentan
    tests.AssertNotEquals(t, getResourceKey("folder", map[string]any{"path": "/Logs"}), getResourceKey("folder", map[string]any{"path": "/logs"}))
}

func TestNetworkInterfaceAndFolderPatternMakeTheDifference(t *testing.T) {
    tests.AssertEquals(t, getResourceKey("network", map[string]any{"interface": "ETH0"}), getResourceKey("network", map[string]any{"interface": "eth0"}))
    tests.AssertNotEquals(t, getResourceKey("network", map[string]any{}), getResourceKey("network", map[string]any{"interface": "eth0"}))
    tests.AssertEquals(t,
        getResourceKey("folder", map[string]any{"path": "/var/log"}),
        getResourceKey("folder", map[string]any{"path": "/var/log", "pattern": "*", "is_recursive": true}),
    )
    tests.AssertNotEquals(t,
        getResourceKey("folder", map[string]any{"path": "/var/log", "pattern": "*.log"}),
        getResourceKey("folder", map[string]any{"path": "/var/log", "pattern": "*.gz"}),
    )
}

func TestLogPatternAlsoComparesTheRegex(t *testing.T) {
    tests.AssertEquals(t,
        getResourceKey("log_pattern", map[string]any{"path": "/var/log", "match": "failed", "severity": "error"}),
        getResourceKey("log_pattern", map[string]any{"path": "/var/log/", "match": "failed", "severity": "critical"}),
    )
    tests.AssertNotEquals(t,
        getResourceKey("log_pattern", map[string]any{"path": "/var/log", "match": "failed"}),
        getResourceKey("log_pattern", map[string]any{"path": "/var/log", "match": "invalid user"}),
    )
}
