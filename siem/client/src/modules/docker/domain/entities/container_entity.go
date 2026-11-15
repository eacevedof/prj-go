package entities

import (
    "strings"

    "siem-client/src/modules/docker/domain/enums"
    "siem-client/src/modules/shared/application"
)

const SHORT_ID_LENGTH = 12

// ContainerEntity un contenedor tal como lo ve el agente en una pasada (de GET /containers/{id}/json). Lo que se
// guarda entre pasadas para detectar cambios es ToStatePrimitives.
type ContainerEntity struct {
    id           string
    name         string
    image        string
    status       string
    health       string
    exitCode     int
    isOomKilled  bool
    restartCount int
    isTty        bool
    finishedAt   string
}

// ContainerEntityFromInspect el JSON de inspect ya decodificado (numeros como float64).
func ContainerEntityFromInspect(inspect map[string]any) *ContainerEntity {
    state, _ := inspect["State"].(map[string]any)
    health, _ := state["Health"].(map[string]any)
    containerConfig, _ := inspect["Config"].(map[string]any)
    isOomKilled, _ := state["OOMKilled"].(bool)
    isTty, _ := containerConfig["Tty"].(bool)
    return &ContainerEntity{
        id:           application.GetTrimmedString(inspect, "Id"),
        name:         strings.TrimPrefix(application.GetTrimmedString(inspect, "Name"), "/"),
        image:        application.GetTrimmedString(containerConfig, "Image"),
        status:       application.GetTrimmedString(state, "Status"),
        health:       application.GetTrimmedString(health, "Status"),
        exitCode:     application.GetInt(state, "ExitCode"),
        isOomKilled:  isOomKilled,
        restartCount: application.GetInt(inspect, "RestartCount"),
        isTty:        isTty,
        finishedAt:   application.GetTrimmedString(state, "FinishedAt"),
    }
}

// ContainerEntityFromStatePrimitives lo guardado en la pasada anterior (para pintar un contenedor que ya no esta).
func ContainerEntityFromStatePrimitives(containerId string, primitives map[string]any) *ContainerEntity {
    isOomKilled, _ := primitives[enums.DockerAttributeEnum.IS_OOM_KILLED].(bool)
    return &ContainerEntity{
        id:           containerId,
        name:         application.GetTrimmedString(primitives, enums.DockerAttributeEnum.NAME),
        image:        application.GetTrimmedString(primitives, enums.DockerAttributeEnum.IMAGE),
        status:       application.GetTrimmedString(primitives, enums.DockerAttributeEnum.STATUS),
        health:       application.GetTrimmedString(primitives, enums.DockerAttributeEnum.HEALTH),
        exitCode:     application.GetInt(primitives, enums.DockerAttributeEnum.EXIT_CODE),
        isOomKilled:  isOomKilled,
        restartCount: application.GetInt(primitives, enums.DockerAttributeEnum.RESTART_COUNT),
    }
}

func (e *ContainerEntity) ToStatePrimitives() map[string]any {
    return map[string]any{
        enums.DockerAttributeEnum.NAME:          e.name,
        enums.DockerAttributeEnum.IMAGE:         e.image,
        enums.DockerAttributeEnum.STATUS:        e.status,
        enums.DockerAttributeEnum.HEALTH:        e.health,
        enums.DockerAttributeEnum.EXIT_CODE:     e.exitCode,
        enums.DockerAttributeEnum.IS_OOM_KILLED: e.isOomKilled,
        enums.DockerAttributeEnum.RESTART_COUNT: e.restartCount,
    }
}

func (e *ContainerEntity) GetId() string {
    return e.id
}

func (e *ContainerEntity) GetShortId() string {
    if len(e.id) <= SHORT_ID_LENGTH {
        return e.id
    }
    return e.id[:SHORT_ID_LENGTH]
}

func (e *ContainerEntity) GetName() string {
    if e.name == "" {
        return e.GetShortId()
    }
    return e.name
}

func (e *ContainerEntity) GetImage() string {
    return e.image
}

func (e *ContainerEntity) GetStatus() string {
    return e.status
}

func (e *ContainerEntity) IsRunning() bool {
    return e.status == enums.ContainerStatusEnum.RUNNING
}

func (e *ContainerEntity) IsExited() bool {
    return e.status == enums.ContainerStatusEnum.EXITED
}

func (e *ContainerEntity) GetHealth() string {
    return e.health
}

func (e *ContainerEntity) GetExitCode() int {
    return e.exitCode
}

func (e *ContainerEntity) IsOomKilled() bool {
    return e.isOomKilled
}

func (e *ContainerEntity) GetRestartCount() int {
    return e.restartCount
}

func (e *ContainerEntity) IsTty() bool {
    return e.isTty
}

// GetFinishedAt RFC3339Nano de cuando se paro por ultima vez ("0001-01-01T00:00:00Z" si nunca).
func (e *ContainerEntity) GetFinishedAt() string {
    return e.finishedAt
}

// GetWithStatus una copia con otro estado (removed: el contenedor ya no existe).
func (e *ContainerEntity) GetWithStatus(status string) *ContainerEntity {
    containerEntity := *e
    containerEntity.status = status
    return &containerEntity
}
