package services

import (
    "sort"

    "siem-client/src/modules/docker/domain/entities"
    "siem-client/src/modules/docker/domain/enums"
    "siem-client/src/modules/shared/application"
)

// ContainerChangesDetector compara los contenedores de ahora con los de la pasada anterior (state) y devuelve que le ha
// pasado a cada uno. En la primera pasada todos salen como first_seen: es el estado de partida que pinta el dashboard.
// Por contenedor, un solo cambio, el mas importante: estado > reinicio > health.
type ContainerChangesDetector struct{}

func ContainerChangesDetectorGetInstance() *ContainerChangesDetector {
    return &ContainerChangesDetector{}
}

// GetContainerChangeEntities previousById: ToStatePrimitives de la pasada anterior por id de contenedor.
func (d *ContainerChangesDetector) GetContainerChangeEntities(previousById map[string]map[string]any, containerEntities []*entities.ContainerEntity, isFirstScan bool) []*entities.ContainerChangeEntity {
    containerChangeEntities := []*entities.ContainerChangeEntity{}
    currentIds := map[string]bool{}
    for _, containerEntity := range containerEntities {
        currentIds[containerEntity.GetId()] = true
        if containerChangeEntity := d.getContainerChangeEntityOrNil(previousById, containerEntity, isFirstScan); containerChangeEntity != nil {
            containerChangeEntities = append(containerChangeEntities, containerChangeEntity)
        }
    }

    // los que estaban y ya no: borrados (orden estable para que la salida no cambie de una pasada a otra)
    removedIds := []string{}
    for containerId := range previousById {
        if !currentIds[containerId] {
            removedIds = append(removedIds, containerId)
        }
    }
    sort.Strings(removedIds)
    for _, containerId := range removedIds {
        previousContainerEntity := entities.ContainerEntityFromStatePrimitives(containerId, previousById[containerId])
        containerChangeEntities = append(containerChangeEntities, entities.ContainerChangeEntityFromPrimitives(
            previousContainerEntity.GetWithStatus(enums.ContainerStatusEnum.REMOVED),
            enums.ContainerChangeKindEnum.REMOVED,
            previousContainerEntity.GetStatus(),
            previousContainerEntity.GetHealth(),
        ))
    }
    return containerChangeEntities
}

func (d *ContainerChangesDetector) getContainerChangeEntityOrNil(previousById map[string]map[string]any, containerEntity *entities.ContainerEntity, isFirstScan bool) *entities.ContainerChangeEntity {
    previous, wasSeen := previousById[containerEntity.GetId()]
    if !wasSeen {
        changeKind := enums.ContainerChangeKindEnum.APPEARED
        if isFirstScan {
            changeKind = enums.ContainerChangeKindEnum.FIRST_SEEN
        }
        return entities.ContainerChangeEntityFromPrimitives(containerEntity, changeKind, "", "")
    }

    previousStatus := application.GetTrimmedString(previous, enums.DockerAttributeEnum.STATUS)
    previousHealth := application.GetTrimmedString(previous, enums.DockerAttributeEnum.HEALTH)
    previousRestartCount := application.GetInt(previous, enums.DockerAttributeEnum.RESTART_COUNT)
    var changeKind string
    switch {
    case containerEntity.GetStatus() != previousStatus:
        changeKind = enums.ContainerChangeKindEnum.STATUS_CHANGED
    case containerEntity.GetRestartCount() > previousRestartCount:
        changeKind = enums.ContainerChangeKindEnum.RESTARTED
    case containerEntity.GetHealth() != previousHealth:
        changeKind = enums.ContainerChangeKindEnum.HEALTH_CHANGED
    default:
        return nil
    }
    return entities.ContainerChangeEntityFromPrimitives(containerEntity, changeKind, previousStatus, previousHealth)
}
