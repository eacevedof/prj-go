package entities

import (
    "notifier-backend/src/modules/shared/application"

    "notifier-backend/src/modules/users/domain/enums"
)

const USER_UUID_PREFIX = "user-"

// UserPrimitives representacion publica: sin id, password_hash ni api_token_hash.
type UserPrimitives struct {
    UserUuid  string `json:"user_uuid"`
    UserName  string `json:"user_name"`
    UserEmail string `json:"user_email"`
    UserRole  string `json:"user_role"`
    IsEnabled bool   `json:"is_enabled"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

type UserEntity struct {
    id           int
    userUuid     string
    userName     string
    userEmail    string
    passwordHash string
    apiTokenHash string
    userRole     string
    isEnabled    bool
    createdAt    string
    updatedAt    string
}

// UserEntityFromPrimitives desde una fila de app_users (columnas en snake_case).
func UserEntityFromPrimitives(primitives map[string]any) *UserEntity {
    return &UserEntity{
        id:           application.GetInt(primitives, enums.UsersKeyEnum.ID),
        userUuid:     application.GetTrimmedString(primitives, enums.UsersKeyEnum.USER_UUID),
        userName:     application.GetTrimmedString(primitives, enums.UsersKeyEnum.USER_NAME),
        userEmail:    application.GetTrimmedString(primitives, enums.UsersKeyEnum.USER_EMAIL),
        passwordHash: application.GetTrimmedString(primitives, enums.UsersKeyEnum.PASSWORD_HASH),
        apiTokenHash: application.GetTrimmedString(primitives, enums.UsersKeyEnum.API_TOKEN_HASH),
        userRole:     application.GetTrimmedString(primitives, enums.UsersKeyEnum.USER_ROLE),
        isEnabled:    application.GetInt(primitives, enums.UsersKeyEnum.IS_ENABLED) == 1,
        createdAt:    application.GetTrimmedString(primitives, enums.UsersKeyEnum.CREATED_AT),
        updatedAt:    application.GetTrimmedString(primitives, enums.UsersKeyEnum.UPDATED_AT),
    }
}

func (e *UserEntity) ToPrimitives() UserPrimitives {
    return UserPrimitives{
        UserUuid:  e.userUuid,
        UserName:  e.userName,
        UserEmail: e.userEmail,
        UserRole:  e.userRole,
        IsEnabled: e.isEnabled,
        CreatedAt: e.createdAt,
        UpdatedAt: e.updatedAt,
    }
}

// ToPrimitivesOrEmpty una lectura que no encuentra nada da el usuario vacio.
func ToPrimitivesOrEmpty(userEntity *UserEntity) UserPrimitives {
    if userEntity == nil {
        return UserPrimitives{}
    }
    return userEntity.ToPrimitives()
}

func (e *UserEntity) GetId() int {
    return e.id
}

func (e *UserEntity) GetUserUuid() string {
    return e.userUuid
}

func (e *UserEntity) GetUserName() string {
    return e.userName
}

func (e *UserEntity) GetUserEmail() string {
    return e.userEmail
}

func (e *UserEntity) GetPasswordHash() string {
    return e.passwordHash
}

func (e *UserEntity) GetApiTokenHash() string {
    return e.apiTokenHash
}

func (e *UserEntity) GetUserRole() string {
    return e.userRole
}

func (e *UserEntity) IsEnabled() bool {
    return e.isEnabled
}

func (e *UserEntity) GetCreatedAt() string {
    return e.createdAt
}

func (e *UserEntity) GetUpdatedAt() string {
    return e.updatedAt
}
