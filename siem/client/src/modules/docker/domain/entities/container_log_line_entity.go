package entities

import "time"

// ContainerLogLineEntity una linea del log de un contenedor: su stream (stdout / stderr), la marca de tiempo que pone
// Docker (UTC, con nanosegundos: sirve de cursor para no repetir lineas) y el texto.
type ContainerLogLineEntity struct {
    stream     string
    occurredAt time.Time
    text       string
}

func ContainerLogLineEntityFromPrimitives(stream string, occurredAt time.Time, text string) *ContainerLogLineEntity {
    return &ContainerLogLineEntity{stream: stream, occurredAt: occurredAt, text: text}
}

func (e *ContainerLogLineEntity) GetStream() string {
    return e.stream
}

func (e *ContainerLogLineEntity) GetOccurredAt() time.Time {
    return e.occurredAt
}

func (e *ContainerLogLineEntity) GetText() string {
    return e.text
}
