package authenticateapitoken

import (
    "notifier-backend/src/core/config"

    "notifier-backend/src/modules/shared/application"
    "notifier-backend/src/modules/shared/infrastructure/components/token"

    "notifier-backend/src/modules/apitokens/domain/enums"
    "notifier-backend/src/modules/apitokens/domain/exceptions"
    "notifier-backend/src/modules/apitokens/infrastructure/repositories"
)

// AuthenticateApiTokenService un token de envio -> su perfil. Mismo 401 para todo fallo (sin token, desconocido, o el
// token, el perfil o su dueno deshabilitados): no se da pista de cual.
type AuthenticateApiTokenService struct {
    tokenGenerator                  *token.TokenGenerator
    apiTokensReaderSqliteRepository *repositories.ApiTokensReaderSqliteRepository
}

func AuthenticateApiTokenServiceGetInstance() *AuthenticateApiTokenService {
    return &AuthenticateApiTokenService{
        tokenGenerator:                  token.GetInstance(),
        apiTokensReaderSqliteRepository: repositories.ApiTokensReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *ApiTokensException (401) o uno inesperado.
func (s *AuthenticateApiTokenService) Invoke(authenticateApiTokenDto *AuthenticateApiTokenDto) (*AuthenticateApiTokenResultDto, error) {
    unauthorizedException := exceptions.UnauthorizedCustom(config.Tr(enums.ApiTokensMessageEnum.INVALID_OR_MISSING_API_TOKEN))
    if authenticateApiTokenDto.GetApiToken() == "" {
        return nil, unauthorizedException
    }

    authData, err := s.apiTokensReaderSqliteRepository.GetAuthDataByApiTokenHash(
        s.tokenGenerator.GetTokenHash(authenticateApiTokenDto.GetApiToken()),
    )
    if err != nil {
        return nil, err
    }
    if !s.isEverythingEnabled(authData) {
        return nil, unauthorizedException
    }

    keys := enums.ApiTokensKeyEnum
    return AuthenticateApiTokenResultDtoFromPrimitives(map[string]any{
        "authApiTokenId": application.GetInt(authData, keys.ID),
        "authSenderId":   application.GetInt(authData, keys.SENDER_ID),
        "authChannel":    application.GetTrimmedString(authData, keys.CHANNEL),
        "authDefaultTo":  application.GetTrimmedString(authData, keys.DEFAULT_TO),
    }), nil
}

func (s *AuthenticateApiTokenService) isEverythingEnabled(authData map[string]any) bool {
    keys := enums.ApiTokensKeyEnum
    return authData != nil &&
        application.GetInt(authData, keys.IS_ENABLED) == 1 &&
        application.GetInt(authData, keys.SENDER_ENABLED) == 1 &&
        application.GetInt(authData, keys.USER_ENABLED) == 1
}
