package repositories

import (
    "fmt"

    "siem-backend/src/modules/events/domain/entities"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

type EventsWriterSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func EventsWriterSqliteRepositoryGetInstance() *EventsWriterSqliteRepository {
    return &EventsWriterSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// InsertEvent exit_code solo en lo que lo tiene (comando, estado de un contenedor, evento del daemon) y duration_ms solo
// en un comando: en el resto van NULL. Sin attributes, NULL.
func (r *EventsWriterSqliteRepository) InsertEvent(eventEntity *entities.EventEntity, sourceId int) (int, error) {
    var exitCode, durationMs, filePath, attributes any
    if eventEntity.HasExitCode() {
        exitCode = eventEntity.GetExitCode()
    }
    if eventEntity.IsCommandResult() {
        durationMs = eventEntity.GetDurationMs()
    }
    if eventEntity.GetFilePath() != "" {
        filePath = eventEntity.GetFilePath()
    }
    if eventEntity.GetAttributes() != "" {
        attributes = eventEntity.GetAttributes()
    }

    sql := fmt.Sprintf(`
        -- insertEvent
        INSERT INTO `+"`events`"+`
            (project_id, agent_id, source_id, event_type, severity, message, file_path, exit_code, duration_ms,
             attributes, occurred_at, received_at)
        VALUES
            (%d, %d, %d, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `, eventEntity.GetProjectId(), eventEntity.GetAgentId(), sourceId)
    return r.Execute(sql,
        eventEntity.GetEventType(),
        eventEntity.GetSeverity(),
        eventEntity.GetMessage(),
        filePath,
        exitCode,
        durationMs,
        attributes,
        eventEntity.GetOccurredAt(),
        eventEntity.GetReceivedAt(),
    )
}

// DeleteEventsReceivedBefore retencion: borra los eventos recibidos antes de receivedBefore ("Y-m-d H:i:s", UTC).
func (r *EventsWriterSqliteRepository) DeleteEventsReceivedBefore(receivedBefore string) (int, error) {
    sql := `
        -- deleteEventsReceivedBefore
        DELETE FROM ` + "`events`" + `
        WHERE 1=1
        AND received_at < ?
        `
    return r.Execute(sql, receivedBefore)
}
