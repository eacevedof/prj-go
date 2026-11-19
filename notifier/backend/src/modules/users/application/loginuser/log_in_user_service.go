package loginuser

import (
    "notifier-backend/src/core/config"

    "notifier-backend/src/modules/shared/infrastructure/components/datetimer"
    "notifier-backend/src/modules/shared/infrastructure/components/password"
    "notifier-backend/src/modules/shared/infrastructure/components/token"

    "notifier-backend/src/modules/users/domain/entities"
    "notifier-backend/src/modules/users/domain/enums"
    "notifier-backend/src/modules/users/domain/exceptions"
    "notifier-backend/src/modules/users/infrastructure/repositories"
)

// LogInUserService log-in del panel con email + password: devuelve un api_token NUEVO (el anterior deja de valer).
// Usuario inexistente, deshabilitado o password erronea responden lo mismo (401, mismo mensaje y mismo tiempo).
type LogInUserService struct {
    dateTimer                   *datetimer.DateTimer
    passwordHasher              *password.PasswordHasher
    tokenGenerator              *token.TokenGenerator
    usersReaderSqliteRepository *repositories.UsersReaderSqliteRepository
    usersWriterSqliteRepository *repositories.UsersWriterSqliteRepository

    logInUserDto *LogInUserDto
}

func LogInUserServiceGetInstance() *LogInUserService {
    return &LogInUserService{
        dateTimer:                   datetimer.GetInstance(),
        passwordHasher:              password.GetInstance(),
        tokenGenerator:              token.GetInstance(),
        usersReaderSqliteRepository: repositories.UsersReaderSqliteRepositoryGetInstance(),
        usersWriterSqliteRepository: repositories.UsersWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *UsersException (401) o uno inesperado.
func (s *LogInUserService) Invoke(logInUserDto *LogInUserDto) (*LogInUserResultDto, error) {
    s.logInUserDto = logInUserDto
    if logInUserDto.GetUserEmail() == "" || logInUserDto.GetPassword() == "" {
        return nil, s.getInvalidCredentialsException()
    }

    userEntity, err := s.tryToGetUserByCredentialsOrFail()
    if err != nil {
        return nil, err
    }

    now := s.dateTimer.GetNowDtIntoTargetTz()
    if err := s.updatePasswordHashIfNeedsRehash(userEntity, now); err != nil {
        return nil, err
    }
    apiToken := s.tokenGenerator.GetNewToken()
    if _, err := s.usersWriterSqliteRepository.UpdateApiTokenHashById(userEntity.GetId(), s.tokenGenerator.GetTokenHash(apiToken), now); err != nil {
        return nil, err
    }

    loggedInUserEntity, err := s.usersReaderSqliteRepository.GetUserById(userEntity.GetId())
    if err != nil {
        return nil, err
    }
    return LogInUserResultDtoFromPrimitives(map[string]any{
        "user":     entities.ToPrimitivesOrEmpty(loggedInUserEntity),
        "apiToken": apiToken,
    }), nil
}

// tryToGetUserByCredentialsOrFail el usuario inexistente tambien pasa por la verificacion (hash vacio): mismo tiempo de
// respuesta que una password erronea.
func (s *LogInUserService) tryToGetUserByCredentialsOrFail() (*entities.UserEntity, error) {
    userEntity, err := s.usersReaderSqliteRepository.GetUserWithPasswordHashByUserEmail(s.logInUserDto.GetUserEmail())
    if err != nil {
        return nil, err
    }
    passwordHash := ""
    if userEntity != nil {
        passwordHash = userEntity.GetPasswordHash()
    }
    isPasswordValid := s.passwordHasher.IsPasswordValid(s.logInUserDto.GetPassword(), passwordHash)
    if userEntity == nil || !isPasswordValid || !userEntity.IsEnabled() {
        return nil, s.getInvalidCredentialsException()
    }
    return userEntity, nil
}

func (s *LogInUserService) updatePasswordHashIfNeedsRehash(userEntity *entities.UserEntity, now string) error {
    if !s.passwordHasher.DoesPasswordHashNeedRehash(userEntity.GetPasswordHash()) {
        return nil
    }
    passwordHash, err := s.passwordHasher.GetPasswordHash(s.logInUserDto.GetPassword())
    if err != nil {
        return err
    }
    _, err = s.usersWriterSqliteRepository.UpdatePasswordHashById(userEntity.GetId(), passwordHash, now)
    return err
}

func (s *LogInUserService) getInvalidCredentialsException() error {
    return exceptions.UnauthorizedCustom(config.Tr(enums.UsersMessageEnum.INVALID_EMAIL_OR_PASSWORD))
}
