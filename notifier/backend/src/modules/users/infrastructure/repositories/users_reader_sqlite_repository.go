package repositories

import (
    "fmt"

    "notifier-backend/src/modules/shared/application"
    sharedrepositories "notifier-backend/src/modules/shared/infrastructure/repositories"

    "notifier-backend/src/modules/users/domain/entities"
    "notifier-backend/src/modules/users/domain/enums"
)

const userColumns = `id, user_uuid, user_name, user_email, user_role, is_enabled, created_at, updated_at`

type UsersReaderSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func UsersReaderSqliteRepositoryGetInstance() *UsersReaderSqliteRepository {
    return &UsersReaderSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

func (r *UsersReaderSqliteRepository) GetUserById(userId int) (*entities.UserEntity, error) {
    sql := fmt.Sprintf(`
        -- getUserById
        SELECT `+userColumns+`
        FROM `+"`app_users`"+`
        WHERE 1=1
        AND id = %d
        `, userId)
    return r.getFirstUserEntity(r.Query(sql))
}

func (r *UsersReaderSqliteRepository) GetUserByUserUuid(userUuid string) (*entities.UserEntity, error) {
    sql := `
        -- getUserByUserUuid
        SELECT ` + userColumns + `
        FROM ` + "`app_users`" + `
        WHERE 1=1
        AND user_uuid = ?
        `
    return r.getFirstUserEntity(r.Query(sql, userUuid))
}

// GetUserWithPasswordHashByUserEmail el unico que lee password_hash: el log-in.
func (r *UsersReaderSqliteRepository) GetUserWithPasswordHashByUserEmail(userEmail string) (*entities.UserEntity, error) {
    sql := `
        -- getUserWithPasswordHashByUserEmail
        SELECT ` + userColumns + `, password_hash
        FROM ` + "`app_users`" + `
        WHERE 1=1
        AND user_email = ? COLLATE NOCASE
        `
    return r.getFirstUserEntity(r.Query(sql, userEmail))
}

func (r *UsersReaderSqliteRepository) GetUserByApiTokenHash(apiTokenHash string) (*entities.UserEntity, error) {
    sql := `
        -- getUserByApiTokenHash
        SELECT ` + userColumns + `
        FROM ` + "`app_users`" + `
        WHERE 1=1
        AND api_token_hash = ?
        `
    return r.getFirstUserEntity(r.Query(sql, apiTokenHash))
}

func (r *UsersReaderSqliteRepository) DoesUserEmailExist(userEmail string) (bool, error) {
    sql := `
        -- doesUserEmailExist
        SELECT COUNT(*) AS n_rows
        FROM ` + "`app_users`" + `
        WHERE 1=1
        AND user_email = ? COLLATE NOCASE
        `
    return r.hasRows(r.Query(sql, userEmail))
}

func (r *UsersReaderSqliteRepository) DoesRootExist() (bool, error) {
    sql := `
        -- doesRootExist
        SELECT COUNT(*) AS n_rows
        FROM ` + "`app_users`" + `
        WHERE 1=1
        AND user_role = ?
        `
    return r.hasRows(r.Query(sql, enums.UserRoleEnum.ROOT))
}

func (r *UsersReaderSqliteRepository) hasRows(rows []map[string]any, err error) (bool, error) {
    if err != nil || len(rows) == 0 {
        return false, err
    }
    return application.ToInt(rows[0]["n_rows"]) > 0, nil
}

func (r *UsersReaderSqliteRepository) getFirstUserEntity(rows []map[string]any, err error) (*entities.UserEntity, error) {
    if err != nil || len(rows) == 0 {
        return nil, err
    }
    return entities.UserEntityFromPrimitives(rows[0]), nil
}
