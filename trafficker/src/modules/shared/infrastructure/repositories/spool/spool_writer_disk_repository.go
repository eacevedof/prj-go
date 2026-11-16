package spool

import (
    "encoding/json"
    "os"
    "path/filepath"
)

const REJECTED_FILE_EXTENSION = ".rejected"

type SpoolWriterDiskRepository struct{}

func SpoolWriterDiskRepositoryGetInstance() *SpoolWriterDiskRepository {
    return &SpoolWriterDiskRepository{}
}

// SaveSpoolBatch escritura atomica (temporal + rename): el lector nunca ve un lote a medio escribir.
func (r *SpoolWriterDiskRepository) SaveSpoolBatch(spoolPath string, batchUuid string, events []map[string]any) error {
    if err := os.MkdirAll(spoolPath, 0o755); err != nil {
        return err
    }
    rawJson, err := json.Marshal(events)
    if err != nil {
        return err
    }
    pathSpoolFile := filepath.Join(spoolPath, batchUuid+SPOOL_FILE_EXTENSION)
    if err := os.WriteFile(pathSpoolFile+".tmp", rawJson, 0o644); err != nil {
        return err
    }
    return os.Rename(pathSpoolFile+".tmp", pathSpoolFile)
}

// DeleteSpoolBatch el backend ya tiene el lote (202).
func (r *SpoolWriterDiskRepository) DeleteSpoolBatch(spoolPath string, batchUuid string) error {
    return os.Remove(filepath.Join(spoolPath, batchUuid+SPOOL_FILE_EXTENSION))
}

// RenameSpoolBatchAsRejected el backend lo rechaza (400) o no se puede leer: deja de reintentarse pero queda en
// disco (<batch_uuid>.rejected) para mirarlo.
func (r *SpoolWriterDiskRepository) RenameSpoolBatchAsRejected(spoolPath string, batchUuid string) error {
    return os.Rename(
        filepath.Join(spoolPath, batchUuid+SPOOL_FILE_EXTENSION),
        filepath.Join(spoolPath, batchUuid+REJECTED_FILE_EXTENSION),
    )
}
