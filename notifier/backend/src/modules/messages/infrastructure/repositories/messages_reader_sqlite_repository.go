package repositories

import (
    "fmt"

    "notifier-backend/src/modules/shared/application"
    sharedrepositories "notifier-backend/src/modules/shared/infrastructure/repositories"

    "notifier-backend/src/modules/messages/domain/entities"
    "notifier-backend/src/modules/messages/domain/enums"
)

const messageColumns = `id, message_uuid, api_token_id, sender_id, channel, recipients, subject, body, body_format, message_status,
            attempts, max_attempts, next_attempt_at, error_message, provider_message_id, idempotency_key, request_ip,
            created_at, sent_at, updated_at`

type MessagesReaderSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func MessagesReaderSqliteRepositoryGetInstance() *MessagesReaderSqliteRepository {
    return &MessagesReaderSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

func (r *MessagesReaderSqliteRepository) GetMessageById(messageId int) (*entities.MessageEntity, error) {
    sql := fmt.Sprintf(`
        -- getMessageById
        SELECT `+messageColumns+`
        FROM `+"`app_messages`"+`
        WHERE 1=1
        AND id = %d
        `, messageId)
    return r.getFirstMessageEntity(r.Query(sql))
}

// GetMessageByApiTokenIdAndMessageUuid un cliente solo ve los mensajes que encolo con su token.
func (r *MessagesReaderSqliteRepository) GetMessageByApiTokenIdAndMessageUuid(apiTokenId int, messageUuid string) (*entities.MessageEntity, error) {
    sql := `
        -- getMessageByApiTokenIdAndMessageUuid
        SELECT ` + messageColumns + `
        FROM ` + "`app_messages`" + `
        WHERE 1=1
        AND api_token_id = ?
        AND message_uuid = ?
        `
    return r.getFirstMessageEntity(r.Query(sql, apiTokenId, messageUuid))
}

// GetMessageByApiTokenIdAndIdempotencyKey el mensaje que ya se encolo con esa clave (un reenvio del cliente).
func (r *MessagesReaderSqliteRepository) GetMessageByApiTokenIdAndIdempotencyKey(apiTokenId int, idempotencyKey string) (*entities.MessageEntity, error) {
    sql := `
        -- getMessageByApiTokenIdAndIdempotencyKey
        SELECT ` + messageColumns + `
        FROM ` + "`app_messages`" + `
        WHERE 1=1
        AND api_token_id = ?
        AND idempotency_key = ?
        `
    return r.getFirstMessageEntity(r.Query(sql, apiTokenId, idempotencyKey))
}

// GetDueMessagesByChannel los pendientes de un canal cuyo turno ha llegado, del mas antiguo al mas nuevo (FIFO).
func (r *MessagesReaderSqliteRepository) GetDueMessagesByChannel(channel string, now string, limit int) ([]*entities.MessageEntity, error) {
    sql := fmt.Sprintf(`
        -- getDueMessagesByChannel
        SELECT `+messageColumns+`
        FROM `+"`app_messages`"+`
        WHERE 1=1
        AND channel = ?
        AND message_status = ?
        AND next_attempt_at <= ?
        ORDER BY id ASC
        LIMIT %d
        `, limit)
    return r.getMessageEntities(r.Query(sql, channel, enums.MessageStatusEnum.PENDING, now))
}

// CountSentBySenderIdSince envios correctos de un perfil desde una fecha: los limites por minuto y por dia.
func (r *MessagesReaderSqliteRepository) CountSentBySenderIdSince(senderId int, sentSince string) (int, error) {
    sql := `
        -- countSentBySenderIdSince
        SELECT COUNT(*) AS n_rows
        FROM ` + "`app_messages`" + `
        WHERE 1=1
        AND sender_id = ?
        AND message_status = ?
        AND sent_at >= ?
        `
    rows, err := r.Query(sql, senderId, enums.MessageStatusEnum.SENT, sentSince)
    if err != nil || len(rows) == 0 {
        return 0, err
    }
    return application.ToInt(rows[0]["n_rows"]), nil
}

func (r *MessagesReaderSqliteRepository) GetAttachmentsByMessageId(messageId int) ([]*entities.MessageAttachmentEntity, error) {
    sql := fmt.Sprintf(`
        -- getAttachmentsByMessageId
        SELECT id, message_id, attachment_type, file_name, content, source_url, mime_type, size_bytes
        FROM `+"`app_message_attachments`"+`
        WHERE 1=1
        AND message_id = %d
        ORDER BY id ASC
        `, messageId)
    rows, err := r.Query(sql)
    if err != nil {
        return nil, err
    }
    attachmentEntities := make([]*entities.MessageAttachmentEntity, 0, len(rows))
    for _, row := range rows {
        attachmentEntities = append(attachmentEntities, entities.MessageAttachmentEntityFromPrimitives(row))
    }
    return attachmentEntities, nil
}

func (r *MessagesReaderSqliteRepository) GetAttemptsByMessageId(messageId int) ([]*entities.MessageAttemptEntity, error) {
    sql := fmt.Sprintf(`
        -- getAttemptsByMessageId
        SELECT message_id, attempt_number, attempt_status, provider_response, duration_ms, attempted_at
        FROM `+"`app_message_attempts`"+`
        WHERE 1=1
        AND message_id = %d
        ORDER BY attempt_number ASC
        `, messageId)
    rows, err := r.Query(sql)
    if err != nil {
        return nil, err
    }
    attemptEntities := make([]*entities.MessageAttemptEntity, 0, len(rows))
    for _, row := range rows {
        attemptEntities = append(attemptEntities, entities.MessageAttemptEntityFromPrimitives(row))
    }
    return attemptEntities, nil
}

func (r *MessagesReaderSqliteRepository) getFirstMessageEntity(rows []map[string]any, err error) (*entities.MessageEntity, error) {
    if err != nil || len(rows) == 0 {
        return nil, err
    }
    return entities.MessageEntityFromPrimitives(rows[0]), nil
}

func (r *MessagesReaderSqliteRepository) getMessageEntities(rows []map[string]any, err error) ([]*entities.MessageEntity, error) {
    if err != nil {
        return nil, err
    }
    messageEntities := make([]*entities.MessageEntity, 0, len(rows))
    for _, row := range rows {
        messageEntities = append(messageEntities, entities.MessageEntityFromPrimitives(row))
    }
    return messageEntities, nil
}
