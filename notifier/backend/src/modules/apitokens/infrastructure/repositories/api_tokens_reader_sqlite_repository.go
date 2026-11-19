package repositories

import (
    "fmt"

    "notifier-backend/src/modules/shared/application"
    sharedrepositories "notifier-backend/src/modules/shared/infrastructure/repositories"

    "notifier-backend/src/modules/apitokens/domain/entities"
)

const apiTokenColumns = `id, sender_id, token_uuid, token_name, token_prefix, default_to, is_enabled, last_used_at, created_at`

type ApiTokensReaderSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func ApiTokensReaderSqliteRepositoryGetInstance() *ApiTokensReaderSqliteRepository {
    return &ApiTokensReaderSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

func (r *ApiTokensReaderSqliteRepository) GetApiTokenById(apiTokenId int) (*entities.ApiTokenEntity, error) {
    sql := fmt.Sprintf(`
        -- getApiTokenById
        SELECT `+apiTokenColumns+`
        FROM `+"`app_api_tokens`"+`
        WHERE 1=1
        AND id = %d
        `, apiTokenId)
    return r.getFirstApiTokenEntity(r.Query(sql))
}

func (r *ApiTokensReaderSqliteRepository) GetApiTokenByTokenUuid(tokenUuid string) (*entities.ApiTokenEntity, error) {
    sql := `
        -- getApiTokenByTokenUuid
        SELECT ` + apiTokenColumns + `
        FROM ` + "`app_api_tokens`" + `
        WHERE 1=1
        AND token_uuid = ?
        `
    return r.getFirstApiTokenEntity(r.Query(sql, tokenUuid))
}

func (r *ApiTokensReaderSqliteRepository) DoesTokenNameExistInSender(senderId int, tokenName string) (bool, error) {
    sql := `
        -- doesTokenNameExistInSender
        SELECT COUNT(*) AS n_rows
        FROM ` + "`app_api_tokens`" + `
        WHERE 1=1
        AND sender_id = ?
        AND token_name = ? COLLATE NOCASE
        `
    rows, err := r.Query(sql, senderId, tokenName)
    if err != nil || len(rows) == 0 {
        return false, err
    }
    return application.ToInt(rows[0]["n_rows"]) > 0, nil
}

// GetAuthDataByApiTokenHash el token de envio con lo que hace falta para encolar: su perfil (canal, proveedor) y si
// el token, el perfil y el dueno estan habilitados. nil si el hash no existe.
func (r *ApiTokensReaderSqliteRepository) GetAuthDataByApiTokenHash(apiTokenHash string) (map[string]any, error) {
    sql := `
        -- getAuthDataByApiTokenHash
        SELECT
            t.id, t.sender_id, t.token_uuid, t.default_to, t.is_enabled,
            s.channel, s.provider_code, s.user_id, s.is_enabled AS sender_enabled,
            u.is_enabled AS user_enabled
        FROM ` + "`app_api_tokens`" + ` t
        INNER JOIN ` + "`app_senders`" + ` s ON s.id = t.sender_id
        INNER JOIN ` + "`app_users`" + ` u ON u.id = s.user_id
        WHERE 1=1
        AND t.api_token_hash = ?
        `
    rows, err := r.Query(sql, apiTokenHash)
    if err != nil || len(rows) == 0 {
        return nil, err
    }
    return rows[0], nil
}

func (r *ApiTokensReaderSqliteRepository) getFirstApiTokenEntity(rows []map[string]any, err error) (*entities.ApiTokenEntity, error) {
    if err != nil || len(rows) == 0 {
        return nil, err
    }
    return entities.ApiTokenEntityFromPrimitives(rows[0]), nil
}
