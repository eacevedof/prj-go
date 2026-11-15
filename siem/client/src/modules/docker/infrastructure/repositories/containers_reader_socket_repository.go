package repositories

import (
    "encoding/json"
    "fmt"
    "net/url"
    "strconv"
    "time"

    "siem-client/src/modules/docker/domain/entities"
    "siem-client/src/modules/shared/infrastructure/components/dockersocket"
)

// ContainersReaderSocketRepository contenedores, sus recursos y sus logs, de la API de Docker por el socket.
type ContainersReaderSocketRepository struct{}

func ContainersReaderSocketRepositoryGetInstance() *ContainersReaderSocketRepository {
    return &ContainersReaderSocketRepository{}
}

// GetAllContainerEntities todos los contenedores (tambien los parados). La lista no trae health, OOM, reinicios ni
// TTY: se pide el inspect de cada uno. Uno que desaparece entre la lista y su inspect se ignora.
func (r *ContainersReaderSocketRepository) GetAllContainerEntities(socketPath string) ([]*entities.ContainerEntity, error) {
    dockerSocketClient := dockersocket.GetInstance(socketPath)
    body, err := dockerSocketClient.GetBody("/containers/json?all=1")
    if err != nil {
        return nil, err
    }
    containerSummaries := []map[string]any{}
    if err := json.Unmarshal(body, &containerSummaries); err != nil {
        return nil, err
    }

    containerEntities := make([]*entities.ContainerEntity, 0, len(containerSummaries))
    for _, containerSummary := range containerSummaries {
        containerId, _ := containerSummary["Id"].(string)
        inspectBody, err := dockerSocketClient.GetBody("/containers/" + url.PathEscape(containerId) + "/json")
        if err != nil {
            continue
        }
        inspect := map[string]any{}
        if err := json.Unmarshal(inspectBody, &inspect); err != nil {
            continue
        }
        containerEntities = append(containerEntities, entities.ContainerEntityFromInspect(inspect))
    }
    return containerEntities, nil
}

// GetContainerRawStats una muestra (one-shot: no espera al segundo de muestreo de docker stats).
func (r *ContainersReaderSocketRepository) GetContainerRawStats(socketPath string, containerId string) (map[string]any, error) {
    body, err := dockersocket.GetInstance(socketPath).GetBody("/containers/" + url.PathEscape(containerId) + "/stats?stream=false&one-shot=true")
    if err != nil {
        return nil, err
    }
    rawStats := map[string]any{}
    if err := json.Unmarshal(body, &rawStats); err != nil {
        return nil, err
    }
    return rawStats, nil
}

// GetContainerLogsBody las lineas desde since (incluida: Docker la repite; el servicio descarta la ya enviada) y como
// mucho las ultimas maxLines. since cero = desde el principio. Con marcas de tiempo (sirven de cursor).
func (r *ContainersReaderSocketRepository) GetContainerLogsBody(socketPath string, containerId string, since time.Time, maxLines int) ([]byte, error) {
    query := url.Values{}
    query.Set("stdout", "1")
    query.Set("stderr", "1")
    query.Set("timestamps", "1")
    query.Set("tail", strconv.Itoa(maxLines))
    if !since.IsZero() {
        query.Set("since", fmt.Sprintf("%d.%09d", since.Unix(), since.Nanosecond()))
    }
    return dockersocket.GetInstance(socketPath).GetBody("/containers/" + url.PathEscape(containerId) + "/logs?" + query.Encode())
}
