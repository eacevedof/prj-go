package state

import (
    "encoding/json"
    "errors"
    "io/fs"
    "os"
    "path/filepath"
)

type StateWriterJsonRepository struct{}

func StateWriterJsonRepositoryGetInstance() *StateWriterJsonRepository {
    return &StateWriterJsonRepository{}
}

// SaveState escritura atomica (fichero temporal + rename): un corte a mitad nunca deja el estado a medias.
func (r *StateWriterJsonRepository) SaveState(statePath string, stateFileName string, state map[string]any) error {
    if err := os.MkdirAll(statePath, 0o755); err != nil {
        return err
    }
    rawJson, err := json.MarshalIndent(state, "", "    ")
    if err != nil {
        return err
    }
    pathStateFile := filepath.Join(statePath, stateFileName)
    if err := os.WriteFile(pathStateFile+".tmp", rawJson, 0o644); err != nil {
        return err
    }
    return os.Rename(pathStateFile+".tmp", pathStateFile)
}

// DeleteState borra el fichero de estado (lo que vigilaba ya no esta en la config). Que no exista no es un error.
func (r *StateWriterJsonRepository) DeleteState(statePath string, stateFileName string) error {
    err := os.Remove(filepath.Join(statePath, stateFileName))
    if errors.Is(err, fs.ErrNotExist) {
        return nil
    }
    return err
}
