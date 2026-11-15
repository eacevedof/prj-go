package dockersocket

import (
    "context"
    "fmt"
    "io"
    "net"
    "net/http"
    "time"
)

const (
    // el host da igual: la conexion va al socket; la API sin version (/containers/json) usa la del daemon
    baseUrl        = "http://docker"
    requestTimeout = 20 * time.Second
    // un log o una lista de eventos enormes no deben tumbar el agente
    maxResponseBytes = 16 << 20
)

// DockerSocketClient GET a la API de Docker Engine por su socket unix (/var/run/docker.sock), solo con la stdlib.
// Solo lectura: el agente nunca crea, para ni borra nada. Una instancia por socket (reutiliza conexiones).
type DockerSocketClient struct {
    httpClient *http.Client
}

var clientsBySocketPath = map[string]*DockerSocketClient{}

func GetInstance(socketPath string) *DockerSocketClient {
    if dockerSocketClient, isCreated := clientsBySocketPath[socketPath]; isCreated {
        return dockerSocketClient
    }
    dialer := &net.Dialer{Timeout: 5 * time.Second}
    dockerSocketClient := &DockerSocketClient{
        httpClient: &http.Client{
            Timeout: requestTimeout,
            Transport: &http.Transport{
                DialContext: func(ctx context.Context, _ string, _ string) (net.Conn, error) {
                    return dialer.DialContext(ctx, "unix", socketPath)
                },
                MaxIdleConns: 4,
            },
        },
    }
    clientsBySocketPath[socketPath] = dockerSocketClient
    return dockerSocketClient
}

// GetBody el cuerpo de GET path (p. ej. /containers/json?all=1). Error si no conecta o si no responde 200.
func (c *DockerSocketClient) GetBody(path string) ([]byte, error) {
    response, err := c.httpClient.Get(baseUrl + path)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes))
    if err != nil {
        return nil, err
    }
    if response.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("GET %s answered %d: %s", path, response.StatusCode, string(body))
    }
    return body, nil
}
