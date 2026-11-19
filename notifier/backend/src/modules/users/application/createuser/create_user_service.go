package createuser

import (
    "notifier-backend/src/core/config"

    sharedenums "notifier-backend/src/modules/shared/domain/enums"
    "notifier-backend/src/modules/shared/infrastructure/components/datetimer"
    "notifier-backend/src/modules/shared/infrastructure/components/password"
    "notifier-backend/src/modules/shared/infrastructure/components/uuid"
    sharedexceptions "notifier-backend/src/modules/shared/infrastructure/exceptions"
    "notifier-backend/src/modules/shared/infrastructure/traits"

    "notifier-backend/src/modules/users/domain/entities"
    "notifier-backend/src/modules/users/domain/enums"
    "notifier-backend/src/modules/users/domain/exceptions"
    "notifier-backend/src/modules/users/domain/services"
    "notifier-backend/src/modules/users/infrastructure/repositories"
)

// CreateUserService alta de un usuario. En la fase 1 solo por consola (./console user-create): quien tiene la consola
// es el administrador del servidor. El primer root se crea asi y no puede haber dos (indice uk_users_root).
type CreateUserService struct {
    traits.ErrorStackTrait

    dateTimer                   *datetimer.DateTimer
    uuidGenerator               *uuid.UuidGenerator
    passwordHasher              *password.PasswordHasher
    userFieldsValidator         *services.UserFieldsValidator
    usersReaderSqliteRepository *repositories.UsersReaderSqliteRepository
    usersWriterSqliteRepository *repositories.UsersWriterSqliteRepository

    createUserDto *CreateUserDto
}

func CreateUserServiceGetInstance() *CreateUserService {
    return &CreateUserService{
        dateTimer:                   datetimer.GetInstance(),
        uuidGenerator:               uuid.GetInstance(),
        passwordHasher:              password.GetInstance(),
        userFieldsValidator:         services.UserFieldsValidatorGetInstance(),
        usersReaderSqliteRepository: repositories.UsersReaderSqliteRepositoryGetInstance(),
        usersWriterSqliteRepository: repositories.UsersWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *FieldsException (400), *UsersException (409) o uno inesperado.
func (s *CreateUserService) Invoke(createUserDto *CreateUserDto) (*CreateUserResultDto, error) {
    s.createUserDto = createUserDto

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }
    if err := s.failIfUserEmailOrRootExists(); err != nil {
        return nil, err
    }

    passwordHash, err := s.passwordHasher.GetPasswordHash(createUserDto.GetPassword())
    if err != nil {
        return nil, err
    }
    userId, err := s.usersWriterSqliteRepository.InsertUser(
        entities.UserEntityFromPrimitives(map[string]any{
            enums.UsersKeyEnum.USER_UUID:     s.uuidGenerator.GetUniqueIdWithPrefix(entities.USER_UUID_PREFIX),
            enums.UsersKeyEnum.USER_NAME:     createUserDto.GetUserName(),
            enums.UsersKeyEnum.USER_EMAIL:    createUserDto.GetUserEmail(),
            enums.UsersKeyEnum.PASSWORD_HASH: passwordHash,
            enums.UsersKeyEnum.USER_ROLE:     createUserDto.GetUserRole(),
            enums.UsersKeyEnum.CREATED_AT:    s.dateTimer.GetNowDtIntoTargetTz(),
        }),
    )
    if err != nil {
        return nil, err
    }

    userEntity, err := s.usersReaderSqliteRepository.GetUserById(userId)
    if err != nil {
        return nil, err
    }
    return CreateUserResultDtoFromPrimitives(map[string]any{"user": entities.ToPrimitivesOrEmpty(userEntity)}), nil
}

func (s *CreateUserService) failIfWrongInput() error {
    errors := s.userFieldsValidator.GetErrors(
        s.createUserDto.GetUserName(), s.createUserDto.GetUserEmail(), s.createUserDto.GetUserRole(), s.createUserDto.GetPassword(),
    )
    if len(errors) == 0 {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
    )
}

func (s *CreateUserService) failIfUserEmailOrRootExists() error {
    userEmail := s.createUserDto.GetUserEmail()
    doesUserEmailExist, err := s.usersReaderSqliteRepository.DoesUserEmailExist(userEmail)
    if err != nil {
        return err
    }
    if doesUserEmailExist {
        return exceptions.ConflictCustom(config.Tr(enums.UsersMessageEnum.USER_EMAIL_X_ALREADY_EXISTS, userEmail))
    }
    if s.createUserDto.GetUserRole() != enums.UserRoleEnum.ROOT {
        return nil
    }
    doesRootExist, err := s.usersReaderSqliteRepository.DoesRootExist()
    if err != nil {
        return err
    }
    if doesRootExist {
        return exceptions.ConflictCustom(config.Tr(enums.UsersMessageEnum.ROOT_ALREADY_EXISTS))
    }
    return nil
}
