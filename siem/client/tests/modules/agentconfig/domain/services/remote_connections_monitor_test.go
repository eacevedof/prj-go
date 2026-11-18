package services

import (
    "strings"
    "testing"

    "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/agentconfig/domain/services"
    "siem-client/tests"
)

func TestConnectionsMonitorIsOnePerAgent(t *testing.T) {
    agentConfigEntity := entities.AgentConfigEntityFromPrimitives(entities.AgentConfigPrimitives{})
    tests.AssertNil(t, toAnyMonitor(agentConfigEntity.GetConnectionsMonitorEntityOrNil()))

    services.RemoteMonitorsMergerGetInstance().Merge(agentConfigEntity, []*entities.RemoteMonitorEntity{
        getMonitorEntity("monitor-net", "connections", "conexiones", map[string]any{
            "is_loopback_included": true,
            "ignore_processes":     []any{"sshd", " ", "chronyd"},
            "ignore_ports":         []any{float64(53), float64(123)},
        }),
        getMonitorEntity("monitor-net-2", "connections", "conexiones-2", nil),
    })

    connectionsMonitorEntity := agentConfigEntity.GetConnectionsMonitorEntityOrNil()
    tests.AssertEquals(t, "monitor-net", connectionsMonitorEntity.GetMonitorUuid())
    tests.AssertEquals(t, true, connectionsMonitorEntity.IsConnections())
    tests.AssertEquals(t, false, connectionsMonitorEntity.IsMetric())
    tests.AssertEquals(t, true, connectionsMonitorEntity.GetConfigBool("is_loopback_included"))
    tests.AssertEquals(t, []string{"sshd", "chronyd"}, connectionsMonitorEntity.GetConfigStrings("ignore_processes"))
    tests.AssertEquals(t, []int{53, 123}, connectionsMonitorEntity.GetConfigInts("ignore_ports"))
    tests.AssertEquals(t, 0, len(agentConfigEntity.GetMetricMonitorEntities()))

    warnings := strings.Join(agentConfigEntity.GetRemoteConfigWarnings(), "\n")
    tests.AssertStringContainsString(t, "conexiones-2 ignored: this agent already has a connections monitor (conexiones)", warnings)
}

// toAnyMonitor un puntero nil tipado no es nil como any: AssertNil necesita el nil de verdad.
func toAnyMonitor(remoteMonitorEntity *entities.RemoteMonitorEntity) any {
    if remoteMonitorEntity == nil {
        return nil
    }
    return remoteMonitorEntity
}
