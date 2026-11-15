package state

import (
    "encoding/json"
    "errors"
    "io/fs"
    "os"
    "path/filepath"
)

// StateReaderJsonRepository el estado del agente que sobrevive a un reinicio (hasta donde se leyo cada fichero,
// cuando se lanzo cada comando, cuando fue el ultimo envio). Un fichero JSON por cosa dentro de APP_STATE_PATH.
type StateReaderJsonRepository struct{}

func StateReaderJsonRepositoryGetInstance() *StateReaderJsonRepository {
    return &StateReaderJsonRepository{}
}

// GetState el contenido del fichero de estado; vacio si aun no existe (primer arranque).
func (r *StateReaderJsonRepository) GetState(statePath string, stateFileName string) (map[string]any, error) {
    rawJson, err := os.ReadFile(filepath.Join(statePath, stateFileName))
    if errors.Is(err, fs.ErrNotExist) {
        return map[string]any{}, nil
    }
    if err != nil {
        return nil, err
    }
    state := map[string]any{}
    if err := json.Unmarshal(rawJson, &state); err != nil {
        return nil, err
    }
    return state, nil
}
