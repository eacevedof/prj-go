package services

import (
    "testing"
    "time"

    "siem-client/src/modules/docker/domain/entities"
    "siem-client/src/modules/docker/domain/services"
    "siem-client/tests"
)

func TestContainerStateEventCarriesWhatTheBackendNeeds(t *testing.T) {
    containerEntity := entities.ContainerEntityFromInspect(map[string]any{
        "Id": "0123456789abcdef", "Name": "/cont-kafka", "Config": map[string]any{"Image": "apache/kafka"},
        "State": map[string]any{"Status": "exited", "ExitCode": float64(137), "OOMKilled": true},
    })
    containerChangeEntity := entities.ContainerChangeEntityFromPrimitives(containerEntity, "status_changed", "running", "")
    event := services.DockerEventBuilderGetInstance().GetContainerStateEvent(containerChangeEntity, time.Date(2026, 11, 15, 18, 0, 0, 0, time.UTC))

    tests.AssertEquals(t, "docker_state", event["event_type"])
    tests.AssertEquals(t, "cont-kafka", event["source_name"])
    tests.AssertEquals(t, 137, event["exit_code"])
    tests.AssertEquals(t, "2026-11-15 18:00:00", event["occurred_at"])
    tests.AssertEquals(t, "cont-kafka: running -> exited (exit 137, OOM killed) · image apache/kafka", event["message"])
    attributes := event["attributes"].(map[string]any)
    tests.AssertEquals(t, "exited", attributes["status"])
    tests.AssertEquals(t, true, attributes["is_oom_killed"])
    tests.AssertEquals(t, "0123456789ab", attributes["container_id"])
}

func TestDaemonEventOfAContainerAndOfTheDaemon(t *testing.T) {
    dockerEventBuilder := services.DockerEventBuilderGetInstance()
    dieEvent := dockerEventBuilder.GetDaemonEvent(map[string]any{
        "Type": "container", "Action": "die", "time": float64(1794765668),
        "Actor": map[string]any{"ID": "0123456789abcdef", "Attributes": map[string]any{"name": "cont-web", "image": "nginx", "exitCode": "1"}},
    })
    tests.AssertEquals(t, "cont-web", dieEvent["source_name"])
    tests.AssertEquals(t, 1, dieEvent["exit_code"])
    tests.AssertEquals(t, "container die cont-web (exitCode=1) · image nginx", dieEvent["message"])
    tests.AssertEquals(t, "die", dieEvent["attributes"].(map[string]any)["action"])

    pullEvent := dockerEventBuilder.GetDaemonEvent(map[string]any{
        "Type": "image", "Action": "pull", "time": float64(1794765668), "Actor": map[string]any{"ID": "nginx:latest"},
    })
    tests.AssertEquals(t, "docker-daemon", pullEvent["source_name"])
    tests.AssertEquals(t, "image pull nginx:latest", pullEvent["message"])
}
