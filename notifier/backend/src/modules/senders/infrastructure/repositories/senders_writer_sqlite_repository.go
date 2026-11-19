package repositories

import (
    sharedrepositories "notifier-backend/src/modules/shared/infrastructure/repositories"

    "notifier-backend/src/modules/senders/domain/entities"
    "notifier-backend/src/modules/senders/domain/enums"
)

type SendersWriterSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func SendersWriterSqliteRepositoryGetInstance() *SendersWriterSqliteRepository {
    return &SendersWriterSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// InsertSender en el alta updated_at = created_at. Un limite 0 se guarda como NULL (el del proveedor).
func (r *SendersWriterSqliteRepository) InsertSender(senderEntity *entities.SenderEntity) (int, error) {
    sql := `
        -- insertSender
        INSERT INTO ` + "`app_senders`" + `
            (user_id, sender_uuid, sender_name, channel, provider_code, config, secret_encrypted, is_enabled,
            daily_limit, minute_limit, created_at, updated_at)
        VALUES
            (?, ?, ?, ?, ?, ?, ?, 1, ?, ?, ?, ?)
        `
    _, err := r.Execute(sql,
        senderEntity.GetUserId(),
        senderEntity.GetSenderUuid(),
        senderEntity.GetSenderName(),
        senderEntity.GetChannel(),
        senderEntity.GetProviderCode(),
        senderEntity.GetConfigJson(),
        senderEntity.GetSecretEncrypted(),
        r.getNullIfZero(senderEntity.GetDailyLimit()),
        r.getNullIfZero(senderEntity.GetMinuteLimit()),
        senderEntity.GetCreatedAt(),
        senderEntity.GetCreatedAt(),
    )
    if err != nil {
        return 0, err
    }
    return r.GetLastInsertId(), nil
}

// UpdateLastUsedAtById un envio correcto: limpia el ultimo error.
func (r *SendersWriterSqliteRepository) UpdateLastUsedAtById(senderId int, lastUsedAt string) (int, error) {
    sql := `
        -- updateLastUsedAtById
        UPDATE ` + "`app_senders`" + `
        SET last_used_at = ?, last_error_at = NULL, last_error_message = NULL
        WHERE 1=1
        AND id = ?
        `
    return r.Execute(sql, lastUsedAt, senderId)
}

// UpdateLastErrorById el ultimo fallo de entrega, para que el panel lo pinte en el perfil.
func (r *SendersWriterSqliteRepository) UpdateLastErrorById(senderId int, lastErrorMessage string, lastErrorAt string) (int, error) {
    errorRunes := []rune(lastErrorMessage)
    if len(errorRunes) > enums.SenderFieldsLimitEnum.LAST_ERROR_MAX_LENGTH {
        errorRunes = errorRunes[:enums.SenderFieldsLimitEnum.LAST_ERROR_MAX_LENGTH]
    }
    sql := `
        -- updateLastErrorById
        UPDATE ` + "`app_senders`" + `
        SET last_error_at = ?, last_error_message = ?
        WHERE 1=1
        AND id = ?
        `
    return r.Execute(sql, lastErrorAt, string(errorRunes), senderId)
}

func (r *SendersWriterSqliteRepository) getNullIfZero(value int) any {
    if value == 0 {
        return nil
    }
    return value
}
