package services

import (
    "testing"

    "siem-backend/tests"

    "siem-backend/src/modules/events/domain/services"
)

func TestNetConnectionSeverity(t *testing.T) {
    eventSeverityResolver := services.EventSeverityResolverGetInstance()
    for caseName, testCase := range map[string]struct {
        attributes       map[string]any
        expectedSeverity string
    }{
        "no soportado en el sistema":        {map[string]any{"is_unavailable": float64(1)}, "warning"},
        "entrante desde ip publica":         {map[string]any{"action": "opened", "direction": "inbound", "remote_ip": "8.8.8.8"}, "info"},
        "entrante desde ipv6 publica":       {map[string]any{"action": "opened", "direction": "inbound", "remote_ip": "2001:4860:4860::8888"}, "info"},
        "entrante ipv4 mapeada en ipv6":     {map[string]any{"action": "opened", "direction": "inbound", "remote_ip": "::ffff:8.8.8.8"}, "info"},
        "entrante desde la red local":       {map[string]any{"action": "opened", "direction": "inbound", "remote_ip": "192.168.1.20"}, "debug"},
        "entrante desde cgnat":              {map[string]any{"action": "opened", "direction": "inbound", "remote_ip": "100.64.1.1"}, "debug"},
        "entrante con ip ilegible":          {map[string]any{"action": "opened", "direction": "inbound", "remote_ip": "?"}, "debug"},
        "saliente a ip publica":             {map[string]any{"action": "opened", "direction": "outbound", "remote_ip": "142.250.1.1"}, "debug"},
        "cierre de una entrante publica":    {map[string]any{"action": "closed", "direction": "inbound", "remote_ip": "8.8.8.8"}, "debug"},
        "puerto nuevo en todas las ifaces":  {map[string]any{"action": "listen_opened", "local_ip": "0.0.0.0"}, "info"},
        "puerto nuevo en ::":                {map[string]any{"action": "listen_opened", "local_ip": "::"}, "info"},
        "puerto nuevo solo en loopback":     {map[string]any{"action": "listen_opened", "local_ip": "127.0.0.1"}, "debug"},
        "puerto nuevo solo en ::1":          {map[string]any{"action": "listen_opened", "local_ip": "::1"}, "debug"},
        "puerto que se cierra":              {map[string]any{"action": "listen_closed", "local_ip": "0.0.0.0"}, "debug"},
    } {
        tests.AssertEquals(t, testCase.expectedSeverity, eventSeverityResolver.GetSeverity("net_connection", "x", 0, "", testCase.attributes), caseName)
    }
}

func TestNetConnectionGoesToNetworkSourceWithItsMonitor(t *testing.T) {
    eventEntity := services.EventNormalizerGetInstance().GetEventEntityOrNil(map[string]any{
        "event_type":  "net_connection",
        "source_name": "firefox",
        "message":     "opened tcp outbound firefox(4321) 192.168.1.5:52780 -> 142.250.1.1:443",
        "attributes":  map[string]any{"monitor_uuid": "mon-1", "action": "opened", "direction": "outbound", "remote_ip": "142.250.1.1"},
    }, 1, 1, "2026-11-18 10:00:00")
    tests.AssertEquals(t, "network", eventEntity.GetSourceType())
    tests.AssertEquals(t, "debug", eventEntity.GetSeverity())
    tests.AssertEquals(t, "mon-1", eventEntity.GetMonitorUuid())
}
