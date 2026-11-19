package authenticateuser

import (
    "notifier-backend/src/core/config"

    "notifier-backend/src/modules/shared/infrastructure/components/token"

    "notifier-backend/src/modules/users/domain/enums"
    "notifier-backend/src/modules/users/domain/exceptions"
    "notifier-backend/src/modules/users/infrastructure/repositories"
)

// AuthenticateUserService el Bearer del panel -> el usuario y su rol. Mismo 401 para token ausente, desconocido o
// usuario deshabilitado: no se da pista de cual fallo.
type AuthenticateUserService struct {
    tokenGenerator              *token.TokenGenerator
    usersReaderSqliteRepository *repositories.UsersReaderSqliteRepository
}

func AuthenticateUserServiceGetInstance() *AuthenticateUserService {
    return &AuthenticateUserService{
        tokenGenerator:              token.GetInstance(),
        usersReaderSqliteRepository: repositories.UsersReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *UsersException (401) o uno inesperado.
func (s *AuthenticateUserService) Invoke(authenticateUserDto *AuthenticateUserDto) (*AuthenticateUserResultDto, error) {
    unauthorizedException := exceptions.UnauthorizedCustom(config.Tr(enums.UsersMessageEnum.MISSING_OR_INVALID_API_TOKEN))
    if authenticateUserDto.GetApiToken() == "" {
        return nil, unauthorizedException
    }

    userEntity, err := s.usersReaderSqliteRepository.GetUserByApiTokenHash(s.tokenGenerator.GetTokenHash(authenticateUserDto.GetApiToken()))
    if err != nil {
        return nil, err
    }
    if userEntity == nil || !userEntity.IsEnabled() {
        return nil, unauthorizedException
    }
    return AuthenticateUserResultDtoFromPrimitives(map[string]any{
        "authUserId":   userEntity.GetId(),
        "authUserUuid": userEntity.GetUserUuid(),
        "authUserRole": userEntity.GetUserRole(),
    }), nil
}
