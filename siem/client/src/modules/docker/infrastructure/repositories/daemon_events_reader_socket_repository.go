package repositories

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "net/url"
    "time"

    "siem-client/src/modules/shared/infrastructure/components/dockersocket"
)

// DaemonEventsReaderSocketRepository eventos del daemon (GET /events). Con until, Docker devuelve los eventos del
// intervalo (un JSON por linea) y cierra: no hace falta mantener el stream abierto entre pasadas.
type DaemonEventsReaderSocketRepository struct{}

func DaemonEventsReaderSocketRepositoryGetInstance() *DaemonEventsReaderSocketRepository {
    return &DaemonEventsReaderSocketRepository{}
}

// GetRawDaemonEvents los eventos entre since y until, con nanosegundos: con segundos enteros, un evento en la fraccion
// entre dos ventanas no entraba en ninguna (se perdio un "start"). Una linea que no es JSON se ignora.
func (r *DaemonEventsReaderSocketRepository) GetRawDaemonEvents(socketPath string, since time.Time, until time.Time) ([]map[string]any, error) {
    query := url.Values{}
    query.Set("since", getUnixWithNanos(since))
    query.Set("until", getUnixWithNanos(until))
    body, err := dockersocket.GetInstance(socketPath).GetBody("/events?" + query.Encode())
    if err != nil {
        return nil, err
    }

    rawDaemonEvents := []map[string]any{}
    scanner := bufio.NewScanner(bytes.NewReader(body))
    scanner.Buffer(make([]byte, 64*1024), 1<<20)
    for scanner.Scan() {
        rawDaemonEvent := map[string]any{}
        if json.Unmarshal(scanner.Bytes(), &rawDaemonEvent) == nil {
            rawDaemonEvents = append(rawDaemonEvents, rawDaemonEvent)
        }
    }
    return rawDaemonEvents, nil
}

// getUnixWithNanos "1794768571.042000000": el formato de since/until que acepta Docker.
func getUnixWithNanos(moment time.Time) string {
    return fmt.Sprintf("%d.%09d", moment.Unix(), moment.Nanosecond())
}
