package services

import (
    "testing"

    "trafficker/src/modules/connections/domain/entities"
    "trafficker/src/modules/connections/domain/services"
    "trafficker/tests"
)

func getConnectionEntities(t *testing.T) []*entities.ConnectionEntity {
    t.Helper()
    return services.ConnectionEntitiesBuilderGetInstance().GetConnectionEntities(
        getSocketEntities(t, tcpListenLine, sshListenLine, sshInboundLine, tcpEstablishedLine),
        map[uint64]*entities.ProcessEntity{
            22512: entities.NewProcessEntity(500, "cupsd", "/usr/sbin/cupsd"),
            100:   entities.NewProcessEntity(812, "sshd", "/usr/sbin/sshd"),
            101:   entities.NewProcessEntity(9001, "sshd", "/usr/sbin/sshd"),
            98765: entities.NewProcessEntity(4321, "firefox", "/usr/lib/firefox/firefox"),
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
    filtered := connectionsFilter.GetFilteredConnectionEntities(getConnectionEntities(t), false, []string{}, []int{}, []string{})
    tests.AssertEquals(t, []string{"sshd(812)", "sshd(9001)", "firefox(4321)"}, getProcessLabels(filtered))

    withLoopback := connectionsFilter.GetFilteredConnectionEntities(getConnectionEntities(t), true, []string{}, []int{}, []string{})
    tests.AssertEquals(t, 4, len(withLoopback))
}

func TestIgnoredProcessesAndPorts(t *testing.T) {
    connectionsFilter := services.ConnectionsFilterGetInstance()
    withoutSshd := connectionsFilter.GetFilteredConnectionEntities(getConnectionEntities(t), false, []string{"sshd"}, []int{}, []string{})
    tests.AssertEquals(t, []string{"firefox(4321)"}, getProcessLabels(withoutSshd))

    withoutHttps := connectionsFilter.GetFilteredConnectionEntities(getConnectionEntities(t), false, []string{}, []int{443}, []string{})
    tests.AssertEquals(t, []string{"sshd(812)", "sshd(9001)"}, getProcessLabels(withoutHttps))
}

func TestFilterTermsAreAnded(t *testing.T) {
    connectionsFilter := services.ConnectionsFilterGetInstance()
    filterTerms := connectionsFilter.GetFilterTerms(" Established , remote ")
    tests.AssertEquals(t, []string{"established", "remote"}, filterTerms)

    filtered := connectionsFilter.GetFilteredConnectionEntities(getConnectionEntities(t), false, []string{}, []int{}, filterTerms)
    tests.AssertEquals(t, []string{"firefox(4321)"}, getProcessLabels(filtered))

    byProcess := connectionsFilter.GetFilteredConnectionEntities(getConnectionEntities(t), false, []string{}, []int{},
        connectionsFilter.GetFilterTerms("sshd,inbound"))
    tests.AssertEquals(t, []string{"sshd(9001)"}, getProcessLabels(byProcess))
}
