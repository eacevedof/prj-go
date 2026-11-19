package repositories

import (
    sharedrepositories "notifier-backend/src/modules/shared/infrastructure/repositories"

    "notifier-backend/src/modules/apitokens/domain/entities"
)

type ApiTokensWriterSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func ApiTokensWriterSqliteRepositoryGetInstance() *ApiTokensWriterSqliteRepository {
    return &ApiTokensWriterSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// InsertApiToken en el alta updated_at = created_at.
func (r *ApiTokensWriterSqliteRepository) InsertApiToken(apiTokenEntity *entities.ApiTokenEntity) (int, error) {
    sql := `
        -- insertApiToken
        INSERT INTO ` + "`app_api_tokens`" + `
            (sender_id, token_uuid, token_name, api_token_hash, token_prefix, default_to, is_enabled, created_at, updated_at)
        VALUES
            (?, ?, ?, ?, ?, ?, 1, ?, ?)
        `
    _, err := r.Execute(sql,
        apiTokenEntity.GetSenderId(),
        apiTokenEntity.GetTokenUuid(),
        apiTokenEntity.GetTokenName(),
        apiTokenEntity.GetApiTokenHash(),
        apiTokenEntity.GetTokenPrefix(),
        apiTokenEntity.GetDefaultToJson(),
        apiTokenEntity.GetCreatedAt(),
        apiTokenEntity.GetCreatedAt(),
    )
    if err != nil {
        return 0, err
    }
    return r.GetLastInsertId(), nil
}

func (r *ApiTokensWriterSqliteRepository) UpdateLastUsedAtById(apiTokenId int, lastUsedAt string) (int, error) {
    sql := `
        -- updateLastUsedAtById
        UPDATE ` + "`app_api_tokens`" + `
        SET last_used_at = ?
        WHERE 1=1
        AND id = ?
        `
    return r.Execute(sql, lastUsedAt, apiTokenId)
}

func (r *ApiTokensWriterSqliteRepository) UpdateIsEnabledById(apiTokenId int, isEnabled bool, updatedAt string) (int, error) {
    sql := `
        -- updateIsEnabledById
        UPDATE ` + "`app_api_tokens`" + `
        SET is_enabled = ?, updated_at = ?
        WHERE 1=1
        AND id = ?
        `
    enabledFlag := 0
    if isEnabled {
        enabledFlag = 1
    }
    return r.Execute(sql, enabledFlag, updatedAt, apiTokenId)
}
