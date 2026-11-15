package services

import (
    "testing"

    "siem-client/src/modules/docker/domain/entities"
    "siem-client/src/modules/docker/domain/services"
    "siem-client/tests"
)

func getContainerEntity(id string, name string, status string, health string, exitCode float64, restartCount float64) *entities.ContainerEntity {
    return entities.ContainerEntityFromInspect(map[string]any{
        "Id":           id,
        "Name":         "/" + name,
        "RestartCount": restartCount,
        "Config":       map[string]any{"Image": "nginx:1.27"},
        "State":        map[string]any{"Status": status, "ExitCode": exitCode, "Health": map[string]any{"Status": health}},
    })
}

func getChangeKinds(containerChangeEntities []*entities.ContainerChangeEntity) map[string]string {
    changeKindByName := map[string]string{}
    for _, containerChangeEntity := range containerChangeEntities {
        changeKindByName[containerChangeEntity.GetContainerEntity().GetName()] = containerChangeEntity.GetChangeKind()
    }
    return changeKindByName
}

func TestFirstScanReportsEveryContainer(t *testing.T) {
    containerEntities := []*entities.ContainerEntity{
        getContainerEntity("aaa", "web", "running", "healthy", 0, 0),
        getContainerEntity("bbb", "kafka", "exited", "", 137, 0),
    }
    changeKinds := getChangeKinds(services.ContainerChangesDetectorGetInstance().GetContainerChangeEntities(map[string]map[string]any{}, containerEntities, true))
    tests.AssertEquals(t, map[string]string{"web": "first_seen", "kafka": "first_seen"}, changeKinds)
}

func TestNoChangesNoEvents(t *testing.T) {
    webEntity := getContainerEntity("aaa", "web", "running", "healthy", 0, 1)
    previousById := map[string]map[string]any{"aaa": webEntity.ToStatePrimitives()}
    tests.AssertEmpty(t, services.ContainerChangesDetectorGetInstance().GetContainerChangeEntities(previousById, []*entities.ContainerEntity{webEntity}, false))
}

func TestEachKindOfChange(t *testing.T) {
    previousById := map[string]map[string]any{
        "aaa": getContainerEntity("aaa", "web", "running", "", 0, 0).ToStatePrimitives(),
        "bbb": getContainerEntity("bbb", "api", "running", "", 0, 2).ToStatePrimitives(),
        "ccc": getContainerEntity("ccc", "db", "running", "healthy", 0, 0).ToStatePrimitives(),
        "ddd": getContainerEntity("ddd", "old", "exited", "", 0, 0).ToStatePrimitives(),
    }
    containerEntities := []*entities.ContainerEntity{
        getContainerEntity("aaa", "web", "exited", "", 137, 0),
        // se cayo y la restart policy lo levanto entre dos pasadas: mismo estado, un reinicio mas
        getContainerEntity("bbb", "api", "running", "", 0, 3),
        getContainerEntity("ccc", "db", "running", "unhealthy", 0, 0),
        getContainerEntity("eee", "new", "running", "", 0, 0),
    }
    containerChangeEntities := services.ContainerChangesDetectorGetInstance().GetContainerChangeEntities(previousById, containerEntities, false)
    tests.AssertEquals(t, map[string]string{
        "web": "status_changed", "api": "restarted", "db": "health_changed", "new": "appeared", "old": "removed",
    }, getChangeKinds(containerChangeEntities))

    for _, containerChangeEntity := range containerChangeEntities {
        if containerChangeEntity.GetChangeKind() == "removed" {
            tests.AssertEquals(t, "removed", containerChangeEntity.GetContainerEntity().GetStatus())
            tests.AssertEquals(t, "exited", containerChangeEntity.GetPreviousStatus())
        }
        if containerChangeEntity.GetChangeKind() == "status_changed" {
            tests.AssertEquals(t, "running", containerChangeEntity.GetPreviousStatus())
        }
    }
}
