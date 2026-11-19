package createapitoken

import (
    "notifier-backend/src/core/config"

    sharedenums "notifier-backend/src/modules/shared/domain/enums"
    "notifier-backend/src/modules/shared/infrastructure/components/datetimer"
    "notifier-backend/src/modules/shared/infrastructure/components/token"
    "notifier-backend/src/modules/shared/infrastructure/components/uuid"
    sharedexceptions "notifier-backend/src/modules/shared/infrastructure/exceptions"
    "notifier-backend/src/modules/shared/infrastructure/repositories/configuration"
    "notifier-backend/src/modules/shared/infrastructure/traits"

    sendersentities "notifier-backend/src/modules/senders/domain/entities"
    sendersrepositories "notifier-backend/src/modules/senders/infrastructure/repositories"

    "notifier-backend/src/modules/apitokens/domain/entities"
    "notifier-backend/src/modules/apitokens/domain/enums"
    "notifier-backend/src/modules/apitokens/domain/exceptions"
    "notifier-backend/src/modules/apitokens/domain/services"
    "notifier-backend/src/modules/apitokens/infrastructure/repositories"
)

// CreateApiTokenService un token de envio mas para un perfil (uno por app cliente). El token en claro solo se
// devuelve aqui; en BD queda su sha256 y el prefijo visible.
type CreateApiTokenService struct {
    traits.ErrorStackTrait

    dateTimer                       *datetimer.DateTimer
    uuidGenerator                   *uuid.UuidGenerator
    tokenGenerator                  *token.TokenGenerator
    apiTokenFieldsValidator         *services.ApiTokenFieldsValidator
    sendersReaderSqliteRepository   *sendersrepositories.SendersReaderSqliteRepository
    apiTokensReaderSqliteRepository *repositories.ApiTokensReaderSqliteRepository
    apiTokensWriterSqliteRepository *repositories.ApiTokensWriterSqliteRepository

    createApiTokenDto *CreateApiTokenDto
    defaultTo         []string
    senderEntity      *sendersentities.SenderEntity
}

func CreateApiTokenServiceGetInstance() *CreateApiTokenService {
    return &CreateApiTokenService{
        dateTimer:                       datetimer.GetInstance(),
        uuidGenerator:                   uuid.GetInstance(),
        tokenGenerator:                  token.GetInstance(),
        apiTokenFieldsValidator:         services.ApiTokenFieldsValidatorGetInstance(),
        sendersReaderSqliteRepository:   sendersrepositories.SendersReaderSqliteRepositoryGetInstance(),
        apiTokensReaderSqliteRepository: repositories.ApiTokensReaderSqliteRepositoryGetInstance(),
        apiTokensWriterSqliteRepository: repositories.ApiTokensWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *FieldsException (400), *ApiTokensException (404, 409) o uno inesperado.
func (s *CreateApiTokenService) Invoke(createApiTokenDto *CreateApiTokenDto) (*CreateApiTokenResultDto, error) {
    s.createApiTokenDto = createApiTokenDto
    s.defaultTo = s.apiTokenFieldsValidator.GetDefaultToFromCsv(createApiTokenDto.GetDefaultToCsv())

    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }
    if err := s.loadSenderOrFail(); err != nil {
        return nil, err
    }
    if err := s.failIfTokenNameExists(); err != nil {
        return nil, err
    }

    apiToken := enums.ApiTokenFormatEnum.PREFIX + s.tokenGenerator.GetNewToken()
    apiTokenId, err := s.apiTokensWriterSqliteRepository.InsertApiToken(
        entities.ApiTokenEntityFromPrimitives(map[string]any{
            enums.ApiTokensKeyEnum.SENDER_ID:      s.senderEntity.GetId(),
            enums.ApiTokensKeyEnum.TOKEN_UUID:     s.uuidGenerator.GetUniqueIdWithPrefix(entities.TOKEN_UUID_PREFIX),
            enums.ApiTokensKeyEnum.TOKEN_NAME:     createApiTokenDto.GetTokenName(),
            enums.ApiTokensKeyEnum.API_TOKEN_HASH: s.tokenGenerator.GetTokenHash(apiToken),
            enums.ApiTokensKeyEnum.TOKEN_PREFIX:   s.apiTokenFieldsValidator.GetVisiblePrefix(apiToken),
            enums.ApiTokensKeyEnum.DEFAULT_TO:     s.defaultTo,
            enums.ApiTokensKeyEnum.CREATED_AT:     s.dateTimer.GetNowDtIntoTargetTz(),
        }),
    )
    if err != nil {
        return nil, err
    }

    apiTokenEntity, err := s.apiTokensReaderSqliteRepository.GetApiTokenById(apiTokenId)
    if err != nil {
        return nil, err
    }
    return CreateApiTokenResultDtoFromPrimitives(map[string]any{
        "token":    entities.ToPrimitivesOrEmpty(apiTokenEntity),
        "apiToken": apiToken,
    }), nil
}

func (s *CreateApiTokenService) failIfWrongInput() error {
    errors := []string{}
    if s.createApiTokenDto.GetSenderUuid() == "" {
        errors = append(errors, enums.ApiTokensKeyEnum.SENDER_UUID+": "+config.Tr(enums.ApiTokensMessageEnum.FIELD_IS_REQUIRED))
    }
    errors = append(errors, s.apiTokenFieldsValidator.GetErrors(
        s.createApiTokenDto.GetTokenName(), s.defaultTo, configuration.GetInstance().GetMessageMaxRecipients(),
    )...)
    if len(errors) == 0 {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
    )
}

func (s *CreateApiTokenService) loadSenderOrFail() error {
    senderUuid := s.createApiTokenDto.GetSenderUuid()
    senderEntity, err := s.sendersReaderSqliteRepository.GetSenderBySenderUuid(senderUuid)
    if err != nil {
        return err
    }
    if senderEntity == nil {
        return exceptions.NotFoundCustom(config.Tr(enums.ApiTokensMessageEnum.SENDER_X_NOT_FOUND, senderUuid))
    }
    s.senderEntity = senderEntity
    return nil
}

func (s *CreateApiTokenService) failIfTokenNameExists() error {
    tokenName := s.createApiTokenDto.GetTokenName()
    doesTokenNameExist, err := s.apiTokensReaderSqliteRepository.DoesTokenNameExistInSender(s.senderEntity.GetId(), tokenName)
    if err != nil {
        return err
    }
    if doesTokenNameExist {
        return exceptions.ConflictCustom(config.Tr(enums.ApiTokensMessageEnum.TOKEN_NAME_X_ALREADY_EXISTS, tokenName))
    }
    return nil
}
