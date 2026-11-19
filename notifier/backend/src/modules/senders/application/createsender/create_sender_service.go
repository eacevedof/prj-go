package createsender

import (
    "notifier-backend/src/core/config"

    sharedenums "notifier-backend/src/modules/shared/domain/enums"
    "notifier-backend/src/modules/shared/infrastructure/components/datetimer"
    "notifier-backend/src/modules/shared/infrastructure/components/secretcipher"
    "notifier-backend/src/modules/shared/infrastructure/components/token"
    "notifier-backend/src/modules/shared/infrastructure/components/uuid"
    sharedexceptions "notifier-backend/src/modules/shared/infrastructure/exceptions"
    "notifier-backend/src/modules/shared/infrastructure/traits"

    usersentities "notifier-backend/src/modules/users/domain/entities"
    usersrepositories "notifier-backend/src/modules/users/infrastructure/repositories"

    apitokensentities "notifier-backend/src/modules/apitokens/domain/entities"
    apitokensenums "notifier-backend/src/modules/apitokens/domain/enums"
    apitokensservices "notifier-backend/src/modules/apitokens/domain/services"
    apitokensrepositories "notifier-backend/src/modules/apitokens/infrastructure/repositories"

    "notifier-backend/src/modules/senders/domain/entities"
    "notifier-backend/src/modules/senders/domain/enums"
    "notifier-backend/src/modules/senders/domain/exceptions"
    "notifier-backend/src/modules/senders/domain/services"
    "notifier-backend/src/modules/senders/infrastructure/repositories"
)

// CreateSenderService alta de un perfil de envio (canal email en la fase 1) para un usuario, y con el su primer token
// de envio: perfil y token en la misma transaccion. El secreto se guarda cifrado (SecretCipher) y no vuelve a salir.
type CreateSenderService struct {
    traits.ErrorStackTrait

    dateTimer                       *datetimer.DateTimer
    uuidGenerator                   *uuid.UuidGenerator
    tokenGenerator                  *token.TokenGenerator
    secretCipher                    *secretcipher.SecretCipher
    senderFieldsValidator           *services.SenderFieldsValidator
    senderConfigBuilder             *services.SenderConfigBuilder
    apiTokenFieldsValidator         *apitokensservices.ApiTokenFieldsValidator
    usersReaderSqliteRepository     *usersrepositories.UsersReaderSqliteRepository
    sendersReaderSqliteRepository   *repositories.SendersReaderSqliteRepository
    sendersWriterSqliteRepository   *repositories.SendersWriterSqliteRepository
    apiTokensReaderSqliteRepository *apitokensrepositories.ApiTokensReaderSqliteRepository
    apiTokensWriterSqliteRepository *apitokensrepositories.ApiTokensWriterSqliteRepository

    createSenderDto *CreateSenderDto
    senderFields    services.SenderFields
    userEntity      *usersentities.UserEntity
}

func CreateSenderServiceGetInstance() *CreateSenderService {
    return &CreateSenderService{
        dateTimer:                       datetimer.GetInstance(),
        uuidGenerator:                   uuid.GetInstance(),
        tokenGenerator:                  token.GetInstance(),
        secretCipher:                    secretcipher.GetInstance(),
        senderFieldsValidator:           services.SenderFieldsValidatorGetInstance(),
        senderConfigBuilder:             services.SenderConfigBuilderGetInstance(),
        apiTokenFieldsValidator:         apitokensservices.ApiTokenFieldsValidatorGetInstance(),
        usersReaderSqliteRepository:     usersrepositories.UsersReaderSqliteRepositoryGetInstance(),
        sendersReaderSqliteRepository:   repositories.SendersReaderSqliteRepositoryGetInstance(),
        sendersWriterSqliteRepository:   repositories.SendersWriterSqliteRepositoryGetInstance(),
        apiTokensReaderSqliteRepository: apitokensrepositories.ApiTokensReaderSqliteRepositoryGetInstance(),
        apiTokensWriterSqliteRepository: apitokensrepositories.ApiTokensWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *FieldsException (400), *SendersException (404, 409) o uno inesperado.
func (s *CreateSenderService) Invoke(createSenderDto *CreateSenderDto) (*CreateSenderResultDto, error) {
    s.createSenderDto = createSenderDto

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }
    if err := s.loadUserOrFail(); err != nil {
        return nil, err
    }
    if err := s.failIfSenderNameExists(); err != nil {
        return nil, err
    }

    secretEncrypted, err := s.getSecretEncrypted()
    if err != nil {
        return nil, err
    }
    apiToken := apitokensenums.ApiTokenFormatEnum.PREFIX + s.tokenGenerator.GetNewToken()
    senderId, apiTokenId := 0, 0
    err = s.sendersWriterSqliteRepository.RunInSharedTransaction(
        func() error {
            var err error
            if senderId, err = s.insertSender(secretEncrypted); err != nil {
                return err
            }
            apiTokenId, err = s.insertFirstApiToken(senderId, apiToken)
            return err
        },
        &s.apiTokensWriterSqliteRepository.AbstractSqliteRepository,
    )
    if err != nil {
        return nil, err
    }

    senderEntity, err := s.sendersReaderSqliteRepository.GetSenderById(senderId)
    if err != nil {
        return nil, err
    }
    apiTokenEntity, err := s.apiTokensReaderSqliteRepository.GetApiTokenById(apiTokenId)
    if err != nil {
        return nil, err
    }
    return CreateSenderResultDtoFromPrimitives(map[string]any{
        "sender":   entities.ToPrimitivesOrEmpty(senderEntity),
        "token":    apitokensentities.ToPrimitivesOrEmpty(apiTokenEntity),
        "apiToken": apiToken,
    }), nil
}

func (s *CreateSenderService) failIfWrongInput() error {
    createSenderDto := s.createSenderDto
    errors := []string{}
    if createSenderDto.GetUserUuid() == "" {
        errors = append(errors, enums.SendersKeyEnum.USER_UUID+": "+config.Tr(enums.SendersMessageEnum.FIELD_IS_REQUIRED))
    }
    providerErrors := s.senderFieldsValidator.GetProviderErrors(createSenderDto.GetProviderCode())
    errors = append(errors, providerErrors...)
    if len(providerErrors) == 0 {
        s.senderFields = s.senderConfigBuilder.GetWithProviderDefaults(services.SenderFields{
            SenderName:   createSenderDto.GetSenderName(),
            ProviderCode: createSenderDto.GetProviderCode(),
            FromEmail:    createSenderDto.GetFromEmail(),
            FromName:     createSenderDto.GetFromName(),
            SmtpHost:     createSenderDto.GetSmtpHost(),
            SmtpPort:     createSenderDto.GetSmtpPort(),
            SmtpSecurity: createSenderDto.GetSmtpSecurity(),
            Username:     createSenderDto.GetUsername(),
            Secret:       createSenderDto.GetSecret(),
            DailyLimit:   createSenderDto.GetDailyLimit(),
            MinuteLimit:  createSenderDto.GetMinuteLimit(),
        })
        errors = append(errors, s.senderFieldsValidator.GetEmailSenderErrors(s.senderFields)...)
    }
    errors = append(errors, s.apiTokenFieldsValidator.GetErrors(s.getTokenName(), nil, 0)...)
    if len(errors) == 0 {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
    )
}

func (s *CreateSenderService) loadUserOrFail() error {
    userUuid := s.createSenderDto.GetUserUuid()
    userEntity, err := s.usersReaderSqliteRepository.GetUserByUserUuid(userUuid)
    if err != nil {
        return err
    }
    if userEntity == nil {
        return exceptions.NotFoundCustom(config.Tr(enums.SendersMessageEnum.USER_X_NOT_FOUND, userUuid))
    }
    s.userEntity = userEntity
    return nil
}

func (s *CreateSenderService) failIfSenderNameExists() error {
    senderName := s.senderFields.SenderName
    doesSenderNameExist, err := s.sendersReaderSqliteRepository.DoesSenderNameExistInUser(s.userEntity.GetId(), senderName)
    if err != nil {
        return err
    }
    if doesSenderNameExist {
        return exceptions.ConflictCustom(config.Tr(enums.SendersMessageEnum.SENDER_NAME_X_ALREADY_EXISTS, senderName))
    }
    return nil
}

// getSecretEncrypted sin secreto (un smtp_custom sin autenticacion) se guarda vacio.
func (s *CreateSenderService) getSecretEncrypted() (string, error) {
    if s.senderFields.Secret == "" {
        return "", nil
    }
    return s.secretCipher.GetEncrypted(s.senderFields.Secret)
}

func (s *CreateSenderService) insertSender(secretEncrypted string) (int, error) {
    providerCode := s.senderFields.ProviderCode
    return s.sendersWriterSqliteRepository.InsertSender(
        entities.SenderEntityFromPrimitives(map[string]any{
            enums.SendersKeyEnum.USER_ID:          s.userEntity.GetId(),
            enums.SendersKeyEnum.SENDER_UUID:      s.uuidGenerator.GetUniqueIdWithPrefix(entities.SENDER_UUID_PREFIX),
            enums.SendersKeyEnum.SENDER_NAME:      s.senderFields.SenderName,
            enums.SendersKeyEnum.CHANNEL:          enums.ProviderEnum.DEFINITIONS[providerCode].Channel,
            enums.SendersKeyEnum.PROVIDER_CODE:    providerCode,
            enums.SendersKeyEnum.CONFIG:           s.senderConfigBuilder.GetEmailConfig(s.senderFields),
            enums.SendersKeyEnum.SECRET_ENCRYPTED: secretEncrypted,
            enums.SendersKeyEnum.DAILY_LIMIT:      s.senderFields.DailyLimit,
            enums.SendersKeyEnum.MINUTE_LIMIT:     s.senderFields.MinuteLimit,
            enums.SendersKeyEnum.CREATED_AT:       s.dateTimer.GetNowDtIntoTargetTz(),
        }),
    )
}

func (s *CreateSenderService) insertFirstApiToken(senderId int, apiToken string) (int, error) {
    return s.apiTokensWriterSqliteRepository.InsertApiToken(
        apitokensentities.ApiTokenEntityFromPrimitives(map[string]any{
            apitokensenums.ApiTokensKeyEnum.SENDER_ID:      senderId,
            apitokensenums.ApiTokensKeyEnum.TOKEN_UUID:     s.uuidGenerator.GetUniqueIdWithPrefix(apitokensentities.TOKEN_UUID_PREFIX),
            apitokensenums.ApiTokensKeyEnum.TOKEN_NAME:     s.getTokenName(),
            apitokensenums.ApiTokensKeyEnum.API_TOKEN_HASH: s.tokenGenerator.GetTokenHash(apiToken),
            apitokensenums.ApiTokensKeyEnum.TOKEN_PREFIX:   s.apiTokenFieldsValidator.GetVisiblePrefix(apiToken),
            apitokensenums.ApiTokensKeyEnum.CREATED_AT:     s.dateTimer.GetNowDtIntoTargetTz(),
        }),
    )
}

func (s *CreateSenderService) getTokenName() string {
    if tokenName := s.createSenderDto.GetTokenName(); tokenName != "" {
        return tokenName
    }
    return apitokensenums.ApiTokenFormatEnum.DEFAULT_NAME
}
