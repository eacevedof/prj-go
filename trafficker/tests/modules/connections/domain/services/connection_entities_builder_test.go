package services

import (
    "testing"

    "trafficker/src/modules/connections/domain/entities"
    "trafficker/src/modules/connections/domain/services"
    "trafficker/tests"
)

const (
    // sshd escucha en 0.0.0.0:22 y alguien de la LAN se ha conectado
    sshListenLine  = "   0: 00000000:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 100 1"
    sshInboundLine = "   1: 0501A8C0:0016 0A01A8C0:D431 01 00000000:00000000 00:00000000 00000000     0        0 101 1"
)

func getSocketEntities(t *testing.T, lines ...string) []*entities.SocketEntity {
    t.Helper()
    parser := services.ProcNetLineParserGetInstance()
    socketEntities := []*entities.SocketEntity{}
    for _, line := range lines {
        socketEntity := parser.GetSocketEntityOrNil("tcp", line)
        if socketEntity == nil {
            t.Fatalf("line not parsed: %s", line)
        }
        socketEntities = append(socketEntities, socketEntity)
    }
    return socketEntities
}

func TestDirectionsAndOrder(t *testing.T) {
    socketEntities := getSocketEntities(t, tcpEstablishedLine, sshInboundLine, sshListenLine)
    processEntityByInode := map[uint64]*entities.ProcessEntity{
        100:   entities.NewProcessEntity(812, "sshd", "/usr/sbin/sshd"),
        101:   entities.NewProcessEntity(9001, "sshd", "/usr/sbin/sshd"),
        98765: entities.NewProcessEntity(4321, "firefox", "/usr/lib/firefox/firefox"),
    }
    connectionEntities := services.ConnectionEntitiesBuilderGetInstance().GetConnectionEntities(
        socketEntities, processEntityByInode, map[int]string{0: "root", 1000: "easu"},
    )

    tests.AssertEquals(t, 3, len(connectionEntities))
    tests.AssertEquals(t, "listen", connectionEntities[0].GetDirection())
    tests.AssertEquals(t, "inbound", connectionEntities[1].GetDirection())
    tests.AssertEquals(t, "sshd(9001)", connectionEntities[1].GetProcessLabel())
    tests.AssertEquals(t, "root", connectionEntities[1].GetUserName())
    tests.AssertEquals(t, "outbound", connectionEntities[2].GetDirection())
    tests.AssertEquals(t, "firefox(4321)", connectionEntities[2].GetProcessLabel())
    tests.AssertEquals(t, "easu", connectionEntities[2].GetUserName())
}

func TestTimeWaitIsDroppedAndUnknownProcessIsKept(t *testing.T) {
    socketEntities := getSocketEntities(t, tcpTimeWaitLine, tcpEstablishedLine)
    connectionEntities := services.ConnectionEntitiesBuilderGetInstance().GetConnectionEntities(
        socketEntities, map[uint64]*entities.ProcessEntity{}, map[int]string{},
    )
    tests.AssertEquals(t, 1, len(connectionEntities))
    tests.AssertEquals(t, "?", connectionEntities[0].GetProcessLabel())
    tests.AssertEquals(t, 0, connectionEntities[0].GetProcessEntity().GetPid())
    tests.AssertEquals(t, "1000", connectionEntities[0].GetUserName())
}

func TestPrimitivesOfAConnection(t *testing.T) {
    connectionEntities := services.ConnectionEntitiesBuilderGetInstance().GetConnectionEntities(
        getSocketEntities(t, tcpEstablishedLine),
        map[uint64]*entities.ProcessEntity{98765: entities.NewProcessEntity(4321, "firefox", "/usr/lib/firefox/firefox")},
        map[int]string{1000: "easu"},
    )
    tests.AssertEquals(t, map[string]any{
        "protocol":    "tcp",
        "ip_version":  4,
        "direction":   "outbound",
        "state":       "ESTABLISHED",
        "local_ip":    "192.168.1.5",
        "local_port":  52524,
        "remote_ip":   "142.250.250.1",
        "remote_port": 443,
        "pid":         4321,
        "process":     "firefox",
        "exe":         "/usr/lib/firefox/firefox",
        "user":        "easu",
    }, connectionEntities[0].ToPrimitives())
}
