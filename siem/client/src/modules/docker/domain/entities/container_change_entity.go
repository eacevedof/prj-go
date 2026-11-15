package entities

// ContainerChangeEntity lo que le ha pasado a un contenedor desde la pasada anterior: su estado actual, el tipo de
// cambio (ContainerChangeKindEnum) y como estaba antes.
type ContainerChangeEntity struct {
    containerEntity *ContainerEntity
    changeKind      string
    previousStatus  string
    previousHealth  string
}

func ContainerChangeEntityFromPrimitives(containerEntity *ContainerEntity, changeKind string, previousStatus string, previousHealth string) *ContainerChangeEntity {
    return &ContainerChangeEntity{
        containerEntity: containerEntity,
        changeKind:      changeKind,
        previousStatus:  previousStatus,
        previousHealth:  previousHealth,
    }
}

func (e *ContainerChangeEntity) GetContainerEntity() *ContainerEntity {
    return e.containerEntity
}

func (e *ContainerChangeEntity) GetChangeKind() string {
    return e.changeKind
}

func (e *ContainerChangeEntity) GetPreviousStatus() string {
    return e.previousStatus
}

func (e *ContainerChangeEntity) GetPreviousHealth() string {
    return e.previousHealth
}
