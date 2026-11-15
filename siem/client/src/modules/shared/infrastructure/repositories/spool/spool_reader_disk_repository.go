package spool

import (
    "encoding/json"
    "errors"
    "io/fs"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

const SPOOL_FILE_EXTENSION = ".json"

// SpoolReaderDiskRepository el spool: los lotes de eventos recogidos y aun no entregados. Un fichero
// <batch_uuid>.json por lote (lista JSON de eventos). El batch_uuid empieza por la hora en hex: ordenar por nombre es
// ordenar por antiguedad.
type SpoolReaderDiskRepository struct{}

func SpoolReaderDiskRepositoryGetInstance() *SpoolReaderDiskRepository {
    return &SpoolReaderDiskRepository{}
}

// GetAllSpoolBatchUuids todos los lotes pendientes, del mas antiguo al mas nuevo.
func (r *SpoolReaderDiskRepository) GetAllSpoolBatchUuids(spoolPath string) ([]string, error) {
    dirEntries, err := os.ReadDir(spoolPath)
    if errors.Is(err, fs.ErrNotExist) {
        return []string{}, nil
    }
    if err != nil {
        return nil, err
    }

    batchUuids := []string{}
    for _, dirEntry := range dirEntries {
        if dirEntry.IsDir() || !strings.HasSuffix(dirEntry.Name(), SPOOL_FILE_EXTENSION) {
            continue
        }
        batchUuids = append(batchUuids, strings.TrimSuffix(dirEntry.Name(), SPOOL_FILE_EXTENSION))
    }
    sort.Strings(batchUuids)
    return batchUuids, nil
}

// GetSpoolBatchEvents los eventos de un lote.
func (r *SpoolReaderDiskRepository) GetSpoolBatchEvents(spoolPath string, batchUuid string) ([]map[string]any, error) {
    rawJson, err := os.ReadFile(filepath.Join(spoolPath, batchUuid+SPOOL_FILE_EXTENSION))
    if err != nil {
        return nil, err
    }
    events := []map[string]any{}
    if err := json.Unmarshal(rawJson, &events); err != nil {
        return nil, err
    }
    return events, nil
}
