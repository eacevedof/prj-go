package repositories

import (
    "fmt"

    "siem-backend/src/modules/ingest/domain/enums"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

const errorMessageMaxLength = 500

type IngestQueueWriterSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func IngestQueueWriterSqliteRepositoryGetInstance() *IngestQueueWriterSqliteRepository {
    return &IngestQueueWriterSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// InsertBatchIfNew encola un lote. Devuelve 0 filas si (agent_id, batch_uuid) ya estaba: un reenvio del cliente
// (no le llego el 202) no duplica.
func (r *IngestQueueWriterSqliteRepository) InsertBatchIfNew(agentId int, batchUuid string, payload string, eventsCount int, receivedAt string) (int, error) {
    sql := fmt.Sprintf(`
        -- insertBatchIfNew
        INSERT OR IGNORE INTO `+"`app_ingest_queue`"+`
            (agent_id, batch_uuid, payload, events_count, queue_status, received_at)
        VALUES
            (%d, ?, ?, %d, ?, ?)
        `, agentId, eventsCount)
    return r.Execute(sql, batchUuid, payload, enums.QueueStatusEnum.PENDING, receivedAt)
}

// UpdatePendingBatchAsDoneById reclama el lote: solo cambia si sigue pending. 0 filas = otro proceso (el worker de
// la API o ./console ingest-queue-process) ya lo proceso. El payload se vacia: ya esta en events.
func (r *IngestQueueWriterSqliteRepository) UpdatePendingBatchAsDoneById(batchId int, processedAt string) (int, error) {
    sql := fmt.Sprintf(`
        -- updatePendingBatchAsDoneById
        UPDATE `+"`app_ingest_queue`"+`
        SET queue_status = ?, payload = '', processed_at = ?, error_message = NULL
        WHERE 1=1
        AND id = %d
        AND queue_status = ?
        `, batchId)
    return r.Execute(sql, enums.QueueStatusEnum.DONE, processedAt, enums.QueueStatusEnum.PENDING)
}

// UpdateBatchAfterFailureById un intento fallido mas: queueStatus sigue pending (se reintenta) o pasa a failed.
func (r *IngestQueueWriterSqliteRepository) UpdateBatchAfterFailureById(batchId int, queueStatus string, errorMessage string, processedAt string) (int, error) {
    errorMessageRunes := []rune(errorMessage)
    if len(errorMessageRunes) > errorMessageMaxLength {
        errorMessageRunes = errorMessageRunes[:errorMessageMaxLength]
    }
    sql := fmt.Sprintf(`
        -- updateBatchAfterFailureById
        UPDATE `+"`app_ingest_queue`"+`
        SET queue_status = ?, attempts = attempts + 1, error_message = ?, processed_at = ?
        WHERE 1=1
        AND id = %d
        `, batchId)
    return r.Execute(sql, queueStatus, string(errorMessageRunes), processedAt)
}

// DeleteDoneBatchesProcessedBefore limpia los lotes ya procesados (solo se guardan para reconocer reenvios).
func (r *IngestQueueWriterSqliteRepository) DeleteDoneBatchesProcessedBefore(processedBefore string) (int, error) {
    sql := `
        -- deleteDoneBatchesProcessedBefore
        DELETE FROM ` + "`app_ingest_queue`" + `
        WHERE 1=1
        AND queue_status = ?
        AND processed_at < ?
        `
    return r.Execute(sql, enums.QueueStatusEnum.DONE, processedBefore)
}
