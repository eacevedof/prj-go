package repositories

import (
    "bytes"
    "encoding/json"
    "os"

    "siem-client/src/modules/agentconfig/domain/entities"
)

// AgentConfigReaderJsonRepository lee siem-client.json. Las claves desconocidas son un error (una errata como
// "intervalo_seconds" no puede quedar en silencio con el valor por defecto).
type AgentConfigReaderJsonRepository struct{}

func AgentConfigReaderJsonRepositoryGetInstance() *AgentConfigReaderJsonRepository {
    return &AgentConfigReaderJsonRepository{}
}

func (r *AgentConfigReaderJsonRepository) DoesConfigFileExist(pathConfigFile string) bool {
    fileInfo, err := os.Stat(pathConfigFile)
    return err == nil && !fileInfo.IsDir()
}

// GetAgentConfigEntity error si el fichero no se lee o no es el JSON esperado (sintaxis, tipos, claves).
func (r *AgentConfigReaderJsonRepository) GetAgentConfigEntity(pathConfigFile string) (*entities.AgentConfigEntity, error) {
    rawJson, err := os.ReadFile(pathConfigFile)
    if err != nil {
        return nil, err
    }

    jsonDecoder := json.NewDecoder(bytes.NewReader(rawJson))
    jsonDecoder.DisallowUnknownFields()
    agentConfigPrimitives := entities.AgentConfigPrimitives{}
    if err := jsonDecoder.Decode(&agentConfigPrimitives); err != nil {
        return nil, err
    }
    return entities.AgentConfigEntityFromPrimitives(agentConfigPrimitives), nil
}
