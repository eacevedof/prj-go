package services

import (
    "encoding/binary"
    "testing"

    "siem-client/src/modules/docker/domain/services"
    "siem-client/tests"
)

// getFrame una trama del log multiplexado: [stream, 0, 0, 0, tamaño big-endian] + datos.
func getFrame(stream byte, payload string) []byte {
    header := make([]byte, 8)
    header[0] = stream
    binary.BigEndian.PutUint32(header[4:], uint32(len(payload)))
    return append(header, []byte(payload)...)
}

func TestMultiplexedLogIsSplitByStreamAndSortedByTime(t *testing.T) {
    body := append(getFrame(1, "2026-11-15T18:00:01.000000001Z listening on :80\n"), getFrame(2, "2026-11-15T18:00:00.5Z ERROR config missing\n")...)
    // una linea partida en dos tramas
    body = append(body, getFrame(1, "2026-11-15T18:00:02Z GET /hea")...)
    body = append(body, getFrame(1, "lth 200\n")...)

    lineEntities := services.DockerLogStreamParserGetInstance().GetContainerLogLineEntities(body, false)
    tests.AssertEquals(t, 3, len(lineEntities))
    tests.AssertEquals(t, "stderr", lineEntities[0].GetStream())
    tests.AssertEquals(t, "ERROR config missing", lineEntities[0].GetText())
    tests.AssertEquals(t, "listening on :80", lineEntities[1].GetText())
    tests.AssertEquals(t, 1, lineEntities[1].GetOccurredAt().Nanosecond())
    tests.AssertEquals(t, "GET /health 200", lineEntities[2].GetText())
}

func TestTtyLogIsPlainText(t *testing.T) {
    body := []byte("2026-11-15T18:00:00Z first\r\n2026-11-15T18:00:01Z second\n\n")
    lineEntities := services.DockerLogStreamParserGetInstance().GetContainerLogLineEntities(body, true)
    tests.AssertEquals(t, 2, len(lineEntities))
    tests.AssertEquals(t, "stdout", lineEntities[1].GetStream())
    tests.AssertEquals(t, "second", lineEntities[1].GetText())
}

func TestLineWithoutTimestampKeepsThePreviousOne(t *testing.T) {
    body := []byte("2026-11-15T18:00:00Z panic: boom\ngoroutine 1 [running]:\n")
    lineEntities := services.DockerLogStreamParserGetInstance().GetContainerLogLineEntities(body, true)
    tests.AssertEquals(t, 2, len(lineEntities))
    tests.AssertEquals(t, "goroutine 1 [running]:", lineEntities[1].GetText())
    tests.AssertEquals(t, lineEntities[0].GetOccurredAt(), lineEntities[1].GetOccurredAt())
}
