package services

import (
    "testing"

    "siem-client/src/modules/connections/domain/entities"
    "siem-client/src/modules/connections/domain/services"
    "siem-client/tests"
)

func getConnectionEntities(t *testing.T) []*entities.ConnectionEntity {
    t.Helper()
    return services.ConnectionEntitiesBuilderGetInstance().GetConnectionEntities(
        getSocketEntities(t, tcpListenLine, sshListenLine, sshInboundLine, tcpEstablishedLine),
        map[uint64]*entities.ProcessEntity{
            22512: entities.ProcessEntityFromPrimitives(map[string]any{"pid": 500, "process": "cupsd", "exe": "/usr/sbin/cupsd"}),
            100:   entities.ProcessEntityFromPrimitives(map[string]any{"pid": 812, "process": "sshd", "exe": "/usr/sbin/sshd"}),
            101:   entities.ProcessEntityFromPrimitives(map[string]any{"pid": 9001, "process": "sshd", "exe": "/usr/sbin/sshd"}),
            98765: entities.ProcessEntityFromPrimitives(map[string]any{"pid": 4321, "process": "firefox", "exe": "/usr/lib/firefox/firefox"}),
        },
        map[int]string{0: "root", 1000: "easu"},
    )
}

func getProcessLabels(connectionEntities []*entities.ConnectionEntity) []string {
    processLabels := []string{}
    for _, connectionEntity := range connectionEntities {
        processLabels = append(processLabels, connectionEntity.GetProcessLabel())
    }
    return processLabels
}

func TestLoopbackIsOutByDefault(t *testing.T) {
    connectionsFilter := services.ConnectionsFilterGetInstance()
    filtered := connectionsFilter.GetFilteredConnectionEntities(getConnectionEntities(t), false, []string{}, []int{})
    tests.AssertEquals(t, []string{"sshd(812)", "sshd(9001)", "firefox(4321)"}, getProcessLabels(filtered))

    withLoopback := connectionsFilter.GetFilteredConnectionEntities(getConnectionEntities(t), true, []string{}, []int{})
    tests.AssertEquals(t, 4, len(withLoopback))
}

func TestIgnoredProcessesAndPorts(t *testing.T) {
    connectionsFilter := services.ConnectionsFilterGetInstance()
    withoutSshd := connectionsFilter.GetFilteredConnectionEntities(getConnectionEntities(t), false, []string{"sshd"}, []int{})
    tests.AssertEquals(t, []string{"firefox(4321)"}, getProcessLabels(withoutSshd))

    // el panel guarda el nombre como se escribio: se compara sin distinguir mayusculas
    withoutFirefox := connectionsFilter.GetFilteredConnectionEntities(getConnectionEntities(t), false, []string{"Firefox"}, []int{})
    tests.AssertEquals(t, []string{"sshd(812)", "sshd(9001)"}, getProcessLabels(withoutFirefox))

    withoutHttps := connectionsFilter.GetFilteredConnectionEntities(getConnectionEntities(t), false, []string{}, []int{443})
    tests.AssertEquals(t, []string{"sshd(812)", "sshd(9001)"}, getProcessLabels(withoutHttps))
}
