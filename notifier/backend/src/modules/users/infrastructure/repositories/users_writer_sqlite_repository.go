package repositories

import (
    sharedrepositories "notifier-backend/src/modules/shared/infrastructure/repositories"

    "notifier-backend/src/modules/users/domain/entities"
)

type UsersWriterSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func UsersWriterSqliteRepositoryGetInstance() *UsersWriterSqliteRepository {
    return &UsersWriterSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// InsertUser en el alta updated_at = created_at.
func (r *UsersWriterSqliteRepository) InsertUser(userEntity *entities.UserEntity) (int, error) {
    sql := `
        -- insertUser
        INSERT INTO ` + "`app_users`" + `
            (user_uuid, user_name, user_email, password_hash, user_role, is_enabled, created_at, updated_at)
        VALUES
            (?, ?, ?, ?, ?, 1, ?, ?)
        `
    _, err := r.Execute(sql,
        userEntity.GetUserUuid(),
        userEntity.GetUserName(),
        userEntity.GetUserEmail(),
        userEntity.GetPasswordHash(),
        userEntity.GetUserRole(),
        userEntity.GetCreatedAt(),
        userEntity.GetCreatedAt(),
    )
    if err != nil {
        return 0, err
    }
    return r.GetLastInsertId(), nil
}

func (r *UsersWriterSqliteRepository) UpdateApiTokenHashById(userId int, apiTokenHash string, updatedAt string) (int, error) {
    sql := `
        -- updateApiTokenHashById
        UPDATE ` + "`app_users`" + `
        SET api_token_hash = ?, updated_at = ?
        WHERE 1=1
        AND id = ?
        `
    return r.Execute(sql, apiTokenHash, updatedAt, userId)
}

func (r *UsersWriterSqliteRepository) UpdatePasswordHashById(userId int, passwordHash string, updatedAt string) (int, error) {
    sql := `
        -- updatePasswordHashById
        UPDATE ` + "`app_users`" + `
        SET password_hash = ?, updated_at = ?
        WHERE 1=1
        AND id = ?
        `
    return r.Execute(sql, passwordHash, updatedAt, userId)
}
