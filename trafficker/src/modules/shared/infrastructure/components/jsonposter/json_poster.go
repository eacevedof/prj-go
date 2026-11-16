package jsonposter

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "time"
)

const (
    requestTimeout       = 30 * time.Second
    maxResponseBodyBytes = 64 << 10
    // la config remota puede traer cientos de monitores
    maxGetResponseBodyBytes = 4 << 20
)

// JsonPoster POST de un JSON con Bearer y lectura de la respuesta. Un 4xx/5xx no es un error del componente: se
// devuelve el codigo y el body. Solo da error si no hay respuesta (red, DNS, timeout).
type JsonPoster struct {
    httpClient *http.Client
}

func GetInstance() *JsonPoster {
    return &JsonPoster{httpClient: &http.Client{Timeout: requestTimeout}}
}

// PostJson userAgent: un nombre generico de cliente ("trafficker/0.1.0"), nunca datos personales.
func (p *JsonPoster) PostJson(url string, bearerToken string, userAgent string, body any) (int, string, error) {
    jsonBody, err := json.Marshal(body)
    if err != nil {
        return 0, "", err
    }

    request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
    if err != nil {
        return 0, "", err
    }
    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Accept", "application/json")
    request.Header.Set("Authorization", "Bearer "+bearerToken)
    request.Header.Set("User-Agent", userAgent)

    response, err := p.httpClient.Do(request)
    if err != nil {
        return 0, "", err
    }
    defer response.Body.Close()
    responseBody, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBodyBytes))
    if err != nil {
        return response.StatusCode, "", err
    }
    return response.StatusCode, string(responseBody), nil
}

// GetJson GET con Bearer (la config remota del agente). Misma lectura que PostJson: codigo HTTP y body; error solo si
// no hay respuesta.
func (p *JsonPoster) GetJson(url string, bearerToken string, userAgent string) (int, string, error) {
    request, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return 0, "", err
    }
    request.Header.Set("Accept", "application/json")
    request.Header.Set("Authorization", "Bearer "+bearerToken)
    request.Header.Set("User-Agent", userAgent)

    response, err := p.httpClient.Do(request)
    if err != nil {
        return 0, "", err
    }
    defer response.Body.Close()
    responseBody, err := io.ReadAll(io.LimitReader(response.Body, maxGetResponseBodyBytes))
    if err != nil {
        return response.StatusCode, "", err
    }
    return response.StatusCode, string(responseBody), nil
}
