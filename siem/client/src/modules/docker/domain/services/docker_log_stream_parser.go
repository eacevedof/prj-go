package services

import (
    "bytes"
    "encoding/binary"
    "strings"
    "time"

    "siem-client/src/modules/docker/domain/entities"
)

const (
    STREAM_STDOUT     = "stdout"
    STREAM_STDERR     = "stderr"
    frameHeaderLength = 8
    // en la cabecera de cada trama: 1 = stdout, 2 = stderr
    frameStreamStderr = 2
)

// DockerLogStreamParser lee la respuesta de GET /containers/{id}/logs?timestamps=1. Un contenedor sin TTY la manda
// multiplexada en tramas ([stream, 0, 0, 0, tamaño uint32 big-endian] + datos); con TTY, texto tal cual (todo stdout).
// Cada linea empieza por la marca de tiempo RFC3339Nano y un espacio. Una linea sin marca valida se queda con la de la
// linea anterior (no se pierde).
type DockerLogStreamParser struct{}

func DockerLogStreamParserGetInstance() *DockerLogStreamParser {
    return &DockerLogStreamParser{}
}

func (p *DockerLogStreamParser) GetContainerLogLineEntities(body []byte, isTty bool) []*entities.ContainerLogLineEntity {
    if isTty || !p.IsMultiplexed(body) {
        return p.getLineEntities(STREAM_STDOUT, string(body))
    }

    // las tramas de un mismo stream se juntan antes de partir en lineas: una linea larga puede venir en varias tramas
    textByStream := map[string]*strings.Builder{STREAM_STDOUT: {}, STREAM_STDERR: {}}
    streamsInOrder := []string{}
    for len(body) >= frameHeaderLength {
        frameLength := int(binary.BigEndian.Uint32(body[4:frameHeaderLength]))
        if frameHeaderLength+frameLength > len(body) {
            frameLength = len(body) - frameHeaderLength
        }
        stream := STREAM_STDOUT
        if body[0] == frameStreamStderr {
            stream = STREAM_STDERR
        }
        if textByStream[stream].Len() == 0 {
            streamsInOrder = append(streamsInOrder, stream)
        }
        textByStream[stream].Write(body[frameHeaderLength : frameHeaderLength+frameLength])
        body = body[frameHeaderLength+frameLength:]
    }

    lineEntities := []*entities.ContainerLogLineEntity{}
    for _, stream := range streamsInOrder {
        lineEntities = append(lineEntities, p.getLineEntities(stream, textByStream[stream].String())...)
    }
    // stdout y stderr intercalados por su hora real
    p.sortByOccurredAt(lineEntities)
    return lineEntities
}

func (p *DockerLogStreamParser) getLineEntities(stream string, text string) []*entities.ContainerLogLineEntity {
    lineEntities := []*entities.ContainerLogLineEntity{}
    var lastOccurredAt time.Time
    for _, rawLine := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
        occurredAt, lineText := p.getTimestampAndText(rawLine)
        if occurredAt.IsZero() {
            occurredAt = lastOccurredAt
        }
        lastOccurredAt = occurredAt
        if strings.TrimSpace(lineText) == "" || occurredAt.IsZero() {
            continue
        }
        lineEntities = append(lineEntities, entities.ContainerLogLineEntityFromPrimitives(stream, occurredAt, strings.TrimRight(lineText, " \t")))
    }
    return lineEntities
}

func (p *DockerLogStreamParser) getTimestampAndText(rawLine string) (time.Time, string) {
    timestamp, lineText, hasSpace := strings.Cut(rawLine, " ")
    if !hasSpace {
        timestamp, lineText = rawLine, ""
    }
    occurredAt, err := time.Parse(time.RFC3339Nano, timestamp)
    if err != nil {
        return time.Time{}, rawLine
    }
    return occurredAt.UTC(), lineText
}

// sortByOccurredAt insercion estable: pocas lineas por pasada y casi siempre ya ordenadas.
func (p *DockerLogStreamParser) sortByOccurredAt(lineEntities []*entities.ContainerLogLineEntity) {
    for i := 1; i < len(lineEntities); i++ {
        for j := i; j > 0 && lineEntities[j].GetOccurredAt().Before(lineEntities[j-1].GetOccurredAt()); j-- {
            lineEntities[j], lineEntities[j-1] = lineEntities[j-1], lineEntities[j]
        }
    }
}

// IsMultiplexed true si el cuerpo empieza por una cabecera de trama (primer byte 1 o 2 y tres ceros): un cuerpo vacio
// o de texto se lee tal cual aunque el contenedor diga que no tiene TTY.
func (p *DockerLogStreamParser) IsMultiplexed(body []byte) bool {
    return len(body) >= frameHeaderLength && (body[0] == 1 || body[0] == frameStreamStderr) && bytes.Equal(body[1:4], []byte{0, 0, 0})
}
