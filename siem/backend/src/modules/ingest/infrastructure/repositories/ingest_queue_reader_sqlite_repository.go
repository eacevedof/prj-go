package repositories

import (
    "fmt"

    "siem-backend/src/modules/ingest/domain/enums"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

type IngestQueueReaderSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func IngestQueueReaderSqliteRepositoryGetInstance() *IngestQueueReaderSqliteRepository {
    return &IngestQueueReaderSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// GetPendingBatches los lotes por procesar, del mas antiguo al mas nuevo (FIFO), con el proyecto de su agente.
func (r *IngestQueueReaderSqliteRepository) GetPendingBatches(limit int) ([]map[string]any, error) {
    sql := fmt.Sprintf(`
        -- getPendingBatches
        SELECT
            q.id, q.agent_id, a.project_id, q.batch_uuid, q.payload, q.attempts, q.received_at
        FROM `+"`ingest_queue`"+` q
        INNER JOIN `+"`agents`"+` a ON a.id = q.agent_id
        WHERE 1=1
        AND q.queue_status = ?
        ORDER BY q.id ASC
        LIMIT %d
        `, limit)
    return r.Query(sql, enums.QueueStatusEnum.PENDING)
}
