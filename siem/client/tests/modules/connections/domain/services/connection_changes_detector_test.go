package services

import (
    "encoding/json"
    "testing"
    "time"

    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/connections/domain/entities"
    "siem-client/src/modules/connections/domain/services"
    "siem-client/tests"
)

const firstPassUnix = int64(1_000_000)

// getStateRoundTrip como queda la foto tras guardarla y releerla de state/connections.json (numeros como float64).
func getStateRoundTrip(t *testing.T, connections map[string]any) map[string]any {
    t.Helper()
    rawJson, err := json.Marshal(connections)
    if err != nil {
        t.Fatal(err)
    }
    roundTrip := map[string]any{}
    if err := json.Unmarshal(rawJson, &roundTrip); err != nil {
        t.Fatal(err)
    }
    return roundTrip
}

func getActions(connectionChangeEntities []*entities.ConnectionChangeEntity) []string {
    actions := []string{}
    for _, connectionChangeEntity := range connectionChangeEntities {
        actions = append(actions, connectionChangeEntity.GetAction()+" "+connectionChangeEntity.GetConnectionPrimitives()["process"].(string))
    }
    return actions
}

func TestFirstPassSendsWhatIsOpenAsInitial(t *testing.T) {
    connectionEntities := getConnectionEntities(t)
    filtered := services.ConnectionsFilterGetInstance().GetFilteredConnectionEntities(connectionEntities, false, []string{}, []int{})
    changes, nextConnections := services.ConnectionChangesDetectorGetInstance().GetChangesAndNextConnections(nil, connectionEntities, filtered, firstPassUnix)

    tests.AssertEquals(t, []string{"listen_opened sshd", "opened sshd", "opened firefox"}, getActions(changes))
    tests.AssertEquals(t, true, changes[0].IsInitial())
    tests.AssertEquals(t, 3, len(nextConnections))
}

func TestOnlyChangesAfterTheFirstPass(t *testing.T) {
    detector := services.ConnectionChangesDetectorGetInstance()
    connectionEntities := getConnectionEntities(t)
    filtered := services.ConnectionsFilterGetInstance().GetFilteredConnectionEntities(connectionEntities, false, []string{}, []int{})
    _, firstConnections := detector.GetChangesAndNextConnections(nil, connectionEntities, filtered, firstPassUnix)

    // nada cambia: ningun evento
    sameChanges, _ := detector.GetChangesAndNextConnections(getStateRoundTrip(t, firstConnections), connectionEntities, filtered, firstPassUnix+30)
    tests.AssertEquals(t, 0, len(sameChanges))

    // firefox cierra 45 s despues y el sshd en escucha se para
    withoutFirefoxAndListen := []*entities.ConnectionEntity{connectionEntities[0], connectionEntities[2]}
    changes, nextConnections := detector.GetChangesAndNextConnections(
        getStateRoundTrip(t, firstConnections), withoutFirefoxAndListen,
        services.ConnectionsFilterGetInstance().GetFilteredConnectionEntities(withoutFirefoxAndListen, false, []string{}, []int{}),
        firstPassUnix+45,
    )
    tests.AssertEquals(t, 2, len(changes))
    actionsByProcess := map[string]*entities.ConnectionChangeEntity{}
    for _, change := range changes {
        actionsByProcess[change.GetAction()] = change
    }
    tests.AssertEquals(t, 45, actionsByProcess["closed"].GetDurationSeconds())
    tests.AssertEquals(t, "firefox", actionsByProcess["closed"].GetConnectionPrimitives()["process"])
    tests.AssertNil(t, actionsByProcess["closed"].GetConnectionPrimitives()["opened_at_unix"])
    tests.AssertEquals(t, "listen_closed", actionsByProcess["listen_closed"].GetAction())
    tests.AssertEquals(t, false, actionsByProcess["closed"].IsInitial())
    tests.AssertEquals(t, 1, len(nextConnections))
}

func TestIgnoredInThePanelIsForgottenWithoutClosing(t *testing.T) {
    detector := services.ConnectionChangesDetectorGetInstance()
    connectionsFilter := services.ConnectionsFilterGetInstance()
    connectionEntities := getConnectionEntities(t)
    _, firstConnections := detector.GetChangesAndNextConnections(
        nil, connectionEntities, connectionsFilter.GetFilteredConnectionEntities(connectionEntities, false, []string{}, []int{}), firstPassUnix,
    )

    // en el panel se ignora sshd: sus conexiones siguen abiertas, no son un cierre
    changes, nextConnections := detector.GetChangesAndNextConnections(
        getStateRoundTrip(t, firstConnections), connectionEntities,
        connectionsFilter.GetFilteredConnectionEntities(connectionEntities, false, []string{"sshd"}, []int{}), firstPassUnix+30,
    )
    tests.AssertEquals(t, 0, len(changes))
    tests.AssertEquals(t, 1, len(nextConnections))

    // y si se deja de ignorar, vuelven como abiertas (no como estado de partida)
    backChanges, _ := detector.GetChangesAndNextConnections(
        getStateRoundTrip(t, nextConnections), connectionEntities,
        connectionsFilter.GetFilteredConnectionEntities(connectionEntities, false, []string{}, []int{}), firstPassUnix+60,
    )
    tests.AssertEquals(t, []string{"listen_opened sshd", "opened sshd"}, getActions(backChanges))
    tests.AssertEquals(t, false, backChanges[0].IsInitial())
}

func getConnectionsMonitorEntity() *agentconfigentities.RemoteMonitorEntity {
    return agentconfigentities.RemoteMonitorEntityFromPrimitives(map[string]any{
        "monitor_uuid": "monitor-net", "monitor_type": "connections", "monitor_name": "conexiones", "interval_seconds": float64(30),
    })
}

func TestNetConnectionEventsFollowTheContract(t *testing.T) {
    occurredAt := time.Date(2026, 11, 18, 10, 0, 0, 0, time.UTC)
    eventBuilder := services.ConnectionEventBuilderGetInstance()
    connectionEntities := getConnectionEntities(t)

    opened := eventBuilder.GetChangeEvent(getConnectionsMonitorEntity(), entities.NewConnectionChangeEntity("opened", connectionEntities[3].ToPrimitives(), 0, false), occurredAt)
    tests.AssertEquals(t, "net_connection", opened["event_type"])
    tests.AssertEquals(t, "firefox", opened["source_name"])
    tests.AssertEquals(t, "opened tcp outbound firefox(4321) 192.168.1.5:52524 -> 142.250.250.1:443", opened["message"])
    tests.AssertEquals(t, "2026-11-18 10:00:00", opened["occurred_at"])
    attributes := opened["attributes"].(map[string]any)
    tests.AssertEquals(t, "monitor-net", attributes["monitor_uuid"])
    tests.AssertEquals(t, "opened", attributes["action"])
    tests.AssertEquals(t, "outbound", attributes["direction"])
    tests.AssertEquals(t, "142.250.250.1", attributes["remote_ip"])
    tests.AssertNil(t, attributes["duration_seconds"])

    listen := eventBuilder.GetChangeEvent(getConnectionsMonitorEntity(), entities.NewConnectionChangeEntity("listen_opened", connectionEntities[1].ToPrimitives(), 0, true), occurredAt)
    tests.AssertEquals(t, "listen_opened tcp sshd(812) 0.0.0.0:22 (already open)", listen["message"])
    tests.AssertEquals(t, true, listen["attributes"].(map[string]any)["is_initial"])

    // un cierre sale de la foto guardada (numeros float64) y lleva su duracion
    closedPrimitives := getStateRoundTrip(t, map[string]any{"c": connectionEntities[3].ToPrimitives()})["c"].(map[string]any)
    closed := eventBuilder.GetChangeEvent(getConnectionsMonitorEntity(), entities.NewConnectionChangeEntity("closed", closedPrimitives, 45, false), occurredAt)
    tests.AssertEquals(t, "closed tcp outbound firefox(4321) 192.168.1.5:52524 -> 142.250.250.1:443 after 45 s", closed["message"])
    tests.AssertEquals(t, 45, closed["attributes"].(map[string]any)["duration_seconds"])

    unknown := eventBuilder.GetChangeEvent(getConnectionsMonitorEntity(), entities.NewConnectionChangeEntity("opened", map[string]any{
        "protocol": "udp", "direction": "outbound", "local_ip": "::1", "local_port": 5353, "remote_ip": "ff02::fb", "remote_port": 5353, "pid": 0, "process": "?",
    }, 0, false), occurredAt)
    tests.AssertEquals(t, "?", unknown["source_name"])
    tests.AssertEquals(t, "opened udp outbound ? [::1]:5353 -> [ff02::fb]:5353", unknown["message"])
}

func TestUnavailableEventIsFromTheMonitor(t *testing.T) {
    event := services.ConnectionEventBuilderGetInstance().GetUnavailableEvent(
        getConnectionsMonitorEntity(), "OS windows is not supported (only Linux: /proc/net)", time.Date(2026, 11, 18, 10, 0, 0, 0, time.UTC),
    )
    tests.AssertEquals(t, "net_connection", event["event_type"])
    tests.AssertEquals(t, "conexiones", event["source_name"])
    tests.AssertEquals(t, "connections not available: OS windows is not supported (only Linux: /proc/net)", event["message"])
    attributes := event["attributes"].(map[string]any)
    tests.AssertEquals(t, true, attributes["is_unavailable"])
    tests.AssertEquals(t, "monitor-net", attributes["monitor_uuid"])
}
