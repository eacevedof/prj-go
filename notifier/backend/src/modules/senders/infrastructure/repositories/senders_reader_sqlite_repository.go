package repositories

import (
    "fmt"

    "notifier-backend/src/modules/shared/application"
    sharedrepositories "notifier-backend/src/modules/shared/infrastructure/repositories"

    "notifier-backend/src/modules/senders/domain/entities"
)

const senderColumns = `id, user_id, sender_uuid, sender_name, channel, provider_code, config, secret_encrypted, is_enabled,
            daily_limit, minute_limit, last_used_at, last_error_at, last_error_message, created_at, updated_at`

type SendersReaderSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func SendersReaderSqliteRepositoryGetInstance() *SendersReaderSqliteRepository {
    return &SendersReaderSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

func (r *SendersReaderSqliteRepository) GetSenderById(senderId int) (*entities.SenderEntity, error) {
    sql := fmt.Sprintf(`
        -- getSenderById
        SELECT `+senderColumns+`
        FROM `+"`app_senders`"+`
        WHERE 1=1
        AND id = %d
        `, senderId)
    return r.getFirstSenderEntity(r.Query(sql))
}

func (r *SendersReaderSqliteRepository) GetSenderBySenderUuid(senderUuid string) (*entities.SenderEntity, error) {
    sql := `
        -- getSenderBySenderUuid
        SELECT ` + senderColumns + `
        FROM ` + "`app_senders`" + `
        WHERE 1=1
        AND sender_uuid = ?
        `
    return r.getFirstSenderEntity(r.Query(sql, senderUuid))
}

// GetSenders los perfiles de un usuario (userId > 0) o todos (0), por nombre.
func (r *SendersReaderSqliteRepository) GetSenders(userId int) ([]*entities.SenderEntity, error) {
    sql := fmt.Sprintf(`
        -- getSenders
        SELECT `+senderColumns+`
        FROM `+"`app_senders`"+`
        WHERE 1=1
        AND (%d = 0 OR user_id = %d)
        ORDER BY sender_name COLLATE NOCASE ASC
        `, userId, userId)
    rows, err := r.Query(sql)
    if err != nil {
        return nil, err
    }
    senderEntities := make([]*entities.SenderEntity, 0, len(rows))
    for _, row := range rows {
        senderEntities = append(senderEntities, entities.SenderEntityFromPrimitives(row))
    }
    return senderEntities, nil
}

func (r *SendersReaderSqliteRepository) DoesSenderNameExistInUser(userId int, senderName string) (bool, error) {
    sql := `
        -- doesSenderNameExistInUser
        SELECT COUNT(*) AS n_rows
        FROM ` + "`app_senders`" + `
        WHERE 1=1
        AND user_id = ?
        AND sender_name = ? COLLATE NOCASE
        `
    rows, err := r.Query(sql, userId, senderName)
    if err != nil || len(rows) == 0 {
        return false, err
    }
    return application.ToInt(rows[0]["n_rows"]) > 0, nil
}

func (r *SendersReaderSqliteRepository) getFirstSenderEntity(rows []map[string]any, err error) (*entities.SenderEntity, error) {
    if err != nil || len(rows) == 0 {
        return nil, err
    }
    return entities.SenderEntityFromPrimitives(rows[0]), nil
}
