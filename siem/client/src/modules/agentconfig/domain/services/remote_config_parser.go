package services

import (
    "encoding/json"
    "errors"
    "strings"

    "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/agentconfig/domain/enums"
    "siem-client/src/modules/shared/application"
)

const (
    ingestPathSuffix      = "/ingest-batches"
    agentConfigPathSuffix = "/agent-config"
)

// RemoteConfigParser lo que tiene de regla la config remota: de donde se baja y que forma tiene que traer.
type RemoteConfigParser struct{}

func RemoteConfigParserGetInstance() *RemoteConfigParser {
    return &RemoteConfigParser{}
}

// GetAgentConfigUrl la url de la config remota sale de la de ingesta (APP_SIEM_ENDPOINT_URL) cambiando el final
// /ingest-batches por /agent-config. Si la de ingesta no acaba asi, se cambia su ultimo tramo.
func (p *RemoteConfigParser) GetAgentConfigUrl(ingestEndpointUrl string) string {
    endpointUrl := strings.TrimRight(strings.TrimSpace(ingestEndpointUrl), "/")
    if strings.HasSuffix(endpointUrl, ingestPathSuffix) {
        return strings.TrimSuffix(endpointUrl, ingestPathSuffix) + agentConfigPathSuffix
    }
    if lastSlash := strings.LastIndex(endpointUrl, "/"); lastSlash > strings.Index(endpointUrl, "://")+2 {
        return endpointUrl[:lastSlash] + agentConfigPathSuffix
    }
    return endpointUrl + agentConfigPathSuffix
}

// GetRemoteConfigOrFail del body de la respuesta ({code, status, errors, data: {config_hash, monitors: [...]}}) a lo
// que se guarda en state/remote-config.json: {config_hash, monitors}. Error si no es JSON, si data no trae una lista
// de monitores o si alguno no es un objeto (una respuesta a medias no debe sustituir a la config buena guardada).
func (p *RemoteConfigParser) GetRemoteConfigOrFail(responseBody string) (map[string]any, error) {
    response := map[string]any{}
    if err := json.Unmarshal([]byte(responseBody), &response); err != nil {
        return nil, errors.New("response is not json")
    }
    data, isObject := response[enums.RemoteConfigKeyEnum.DATA].(map[string]any)
    if !isObject {
        return nil, errors.New("data is missing")
    }
    rawMonitors, isList := data[enums.RemoteConfigKeyEnum.MONITORS].([]any)
    if !isList {
        return nil, errors.New("data.monitors must be an array")
    }
    for _, rawMonitor := range rawMonitors {
        if _, isMonitorObject := rawMonitor.(map[string]any); !isMonitorObject {
            return nil, errors.New("data.monitors must contain objects")
        }
    }
    return map[string]any{
        enums.RemoteConfigKeyEnum.CONFIG_HASH: application.GetTrimmedString(data, enums.RemoteConfigKeyEnum.CONFIG_HASH),
        enums.RemoteConfigKeyEnum.MONITORS:    rawMonitors,
    }, nil
}

// GetRemoteMonitorEntities los monitores de una config remota guardada (state/remote-config.json ya decodificado).
func (p *RemoteConfigParser) GetRemoteMonitorEntities(remoteConfig map[string]any) []*entities.RemoteMonitorEntity {
    remoteMonitorEntities := []*entities.RemoteMonitorEntity{}
    rawMonitors, _ := remoteConfig[enums.RemoteConfigKeyEnum.MONITORS].([]any)
    for _, rawMonitor := range rawMonitors {
        if monitor, isObject := rawMonitor.(map[string]any); isObject {
            remoteMonitorEntities = append(remoteMonitorEntities, entities.RemoteMonitorEntityFromPrimitives(monitor))
        }
    }
    return remoteMonitorEntities
}
