package repositories

import (
    "fmt"

    "siem-backend/src/modules/events/domain/entities"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

const lastMessageMaxLength = 500

type SourcesWriterSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func SourcesWriterSqliteRepositoryGetInstance() *SourcesWriterSqliteRepository {
    return &SourcesWriterSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// InsertSource una source nace con su primer evento.
func (r *SourcesWriterSqliteRepository) InsertSource(agentId int, sourceType string, sourceName string, createdAt string) (int, error) {
    sql := fmt.Sprintf(`
        -- insertSource
        INSERT INTO `+"`sources`"+`
            (agent_id, source_type, source_name, created_at, updated_at)
        VALUES
            (%d, ?, ?, ?, ?)
        `, agentId)
    if _, err := r.Execute(sql, sourceType, sourceName, createdAt, createdAt); err != nil {
        return 0, err
    }
    return r.GetLastInsertId(), nil
}

// UpdateLastEventById el ultimo estado conocido de la source = su evento mas reciente del lote. newEventsCount se
// suma a events_count.
func (r *SourcesWriterSqliteRepository) UpdateLastEventById(sourceId int, lastEventEntity *entities.EventEntity, newEventsCount int) (int, error) {
    var lastExitCode any
    if lastEventEntity.HasExitCode() {
        lastExitCode = lastEventEntity.GetExitCode()
    }
    lastMessage := []rune(lastEventEntity.GetMessage())
    if len(lastMessage) > lastMessageMaxLength {
        lastMessage = lastMessage[:lastMessageMaxLength]
    }

    sql := fmt.Sprintf(`
        -- updateLastEventById
        UPDATE `+"`sources`"+`
        SET events_count = events_count + %d, last_event_at = ?, last_severity = ?, last_exit_code = ?,
            last_message = ?, updated_at = ?
        WHERE 1=1
        AND id = %d
        `, newEventsCount, sourceId)
    return r.Execute(sql,
        lastEventEntity.GetOccurredAt(),
        lastEventEntity.GetSeverity(),
        lastExitCode,
        string(lastMessage),
        lastEventEntity.GetReceivedAt(),
    )
}
