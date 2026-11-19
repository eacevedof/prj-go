package repositories

import (
    sharedrepositories "notifier-backend/src/modules/shared/infrastructure/repositories"

    "notifier-backend/src/modules/messages/domain/entities"
    "notifier-backend/src/modules/messages/domain/enums"
)

type MessagesWriterSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func MessagesWriterSqliteRepositoryGetInstance() *MessagesWriterSqliteRepository {
    return &MessagesWriterSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// InsertMessage lo encola: pending y con turno ya (next_attempt_at = created_at).
func (r *MessagesWriterSqliteRepository) InsertMessage(messageEntity *entities.MessageEntity) (int, error) {
    sql := `
        -- insertMessage
        INSERT INTO ` + "`app_messages`" + `
            (message_uuid, api_token_id, sender_id, channel, recipients, subject, body, body_format, message_status,
            attempts, max_attempts, next_attempt_at, idempotency_key, request_ip, created_at, updated_at)
        VALUES
            (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?, ?, ?, ?)
        `
    _, err := r.Execute(sql,
        messageEntity.GetMessageUuid(),
        messageEntity.GetApiTokenId(),
        messageEntity.GetSenderId(),
        messageEntity.GetChannel(),
        messageEntity.GetRecipients().ToJson(),
        messageEntity.GetSubject(),
        messageEntity.GetBody(),
        messageEntity.GetBodyFormat(),
        enums.MessageStatusEnum.PENDING,
        messageEntity.GetMaxAttempts(),
        messageEntity.GetCreatedAt(),
        r.getNullIfEmpty(messageEntity.GetIdempotencyKey()),
        r.getTruncated(messageEntity.GetRequestIp(), enums.MessagesLimitEnum.REQUEST_IP_MAX_LENGTH),
        messageEntity.GetCreatedAt(),
        messageEntity.GetCreatedAt(),
    )
    if err != nil {
        return 0, err
    }
    return r.GetLastInsertId(), nil
}

func (r *MessagesWriterSqliteRepository) InsertAttachment(attachmentEntity *entities.MessageAttachmentEntity) (int, error) {
    sql := `
        -- insertAttachment
        INSERT INTO ` + "`app_message_attachments`" + `
            (message_id, attachment_type, file_name, content, source_url)
        VALUES
            (?, ?, ?, ?, ?)
        `
    _, err := r.Execute(sql,
        attachmentEntity.GetMessageId(),
        attachmentEntity.GetAttachmentType(),
        attachmentEntity.GetFileName(),
        r.getNullIfEmpty(attachmentEntity.GetContent()),
        r.getNullIfEmpty(attachmentEntity.GetSourceUrl()),
    )
    if err != nil {
        return 0, err
    }
    return r.GetLastInsertId(), nil
}

// UpdatePendingAsSendingById reclama el mensaje: solo cambia si sigue pending. 0 filas = otro proceso (el worker o
// ./console message-queue-process) ya lo cogio.
func (r *MessagesWriterSqliteRepository) UpdatePendingAsSendingById(messageId int, updatedAt string) (int, error) {
    sql := `
        -- updatePendingAsSendingById
        UPDATE ` + "`app_messages`" + `
        SET message_status = ?, updated_at = ?
        WHERE 1=1
        AND id = ?
        AND message_status = ?
        `
    return r.Execute(sql, enums.MessageStatusEnum.SENDING, updatedAt, messageId, enums.MessageStatusEnum.PENDING)
}

func (r *MessagesWriterSqliteRepository) UpdateAsSentById(messageId int, providerMessageId string, sentAt string) (int, error) {
    sql := `
        -- updateAsSentById
        UPDATE ` + "`app_messages`" + `
        SET message_status = ?, attempts = attempts + 1, provider_message_id = ?, error_message = NULL,
            sent_at = ?, updated_at = ?
        WHERE 1=1
        AND id = ?
        `
    return r.Execute(sql, enums.MessageStatusEnum.SENT, providerMessageId, sentAt, sentAt, messageId)
}

// UpdateAfterFailureById un intento fallido mas: vuelve a pending con su siguiente turno o queda failed.
func (r *MessagesWriterSqliteRepository) UpdateAfterFailureById(messageId int, messageStatus string, errorMessage string, nextAttemptAt string, updatedAt string) (int, error) {
    sql := `
        -- updateAfterFailureById
        UPDATE ` + "`app_messages`" + `
        SET message_status = ?, attempts = attempts + 1, error_message = ?, next_attempt_at = ?, updated_at = ?
        WHERE 1=1
        AND id = ?
        `
    return r.Execute(sql, messageStatus, r.getTruncated(errorMessage, enums.MessagesLimitEnum.ERROR_MESSAGE_MAX_LENGTH), nextAttemptAt, updatedAt, messageId)
}

// UpdateNextAttemptAtById aplaza un pendiente sin gastar un intento (su perfil llego al limite de envios).
func (r *MessagesWriterSqliteRepository) UpdateNextAttemptAtById(messageId int, nextAttemptAt string, updatedAt string) (int, error) {
    sql := `
        -- updateNextAttemptAtById
        UPDATE ` + "`app_messages`" + `
        SET next_attempt_at = ?, updated_at = ?
        WHERE 1=1
        AND id = ?
        AND message_status = ?
        `
    return r.Execute(sql, nextAttemptAt, updatedAt, messageId, enums.MessageStatusEnum.PENDING)
}

// UpdateOrphanSendingAsPending los sending que llevan parados desde antes de staleBefore (la API murio a medio envio)
// vuelven a pending. Puede salir un duplicado si el proveedor llego a aceptarlo: es preferible a perderlo.
func (r *MessagesWriterSqliteRepository) UpdateOrphanSendingAsPending(staleBefore string, updatedAt string) (int, error) {
    sql := `
        -- updateOrphanSendingAsPending
        UPDATE ` + "`app_messages`" + `
        SET message_status = ?, updated_at = ?
        WHERE 1=1
        AND message_status = ?
        AND updated_at < ?
        `
    return r.Execute(sql, enums.MessageStatusEnum.PENDING, updatedAt, enums.MessageStatusEnum.SENDING, staleBefore)
}

// UpdateAttachmentMetadataById el tipo y el tamano reales, que se conocen al enviar.
func (r *MessagesWriterSqliteRepository) UpdateAttachmentMetadataById(attachmentId int, mimeType string, sizeBytes int) (int, error) {
    sql := `
        -- updateAttachmentMetadataById
        UPDATE ` + "`app_message_attachments`" + `
        SET mime_type = ?, size_bytes = ?
        WHERE 1=1
        AND id = ?
        `
    return r.Execute(sql, mimeType, sizeBytes, attachmentId)
}

func (r *MessagesWriterSqliteRepository) InsertAttempt(attemptEntity *entities.MessageAttemptEntity) (int, error) {
    sql := `
        -- insertAttempt
        INSERT INTO ` + "`app_message_attempts`" + `
            (message_id, attempt_number, attempt_status, provider_response, duration_ms, attempted_at)
        VALUES
            (?, ?, ?, ?, ?, ?)
        `
    return r.Execute(sql,
        attemptEntity.GetMessageId(),
        attemptEntity.GetAttemptNumber(),
        attemptEntity.GetAttemptStatus(),
        r.getTruncated(attemptEntity.GetProviderResponse(), enums.MessagesLimitEnum.PROVIDER_RESPONSE_MAX_LENGTH),
        attemptEntity.GetDurationMs(),
        attemptEntity.GetAttemptedAt(),
    )
}

func (r *MessagesWriterSqliteRepository) getNullIfEmpty(value string) any {
    if value == "" {
        return nil
    }
    return value
}

func (r *MessagesWriterSqliteRepository) getTruncated(text string, maxLength int) string {
    runes := []rune(text)
    if len(runes) <= maxLength {
        return text
    }
    return string(runes[:maxLength])
}
