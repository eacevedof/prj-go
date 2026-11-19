package disableapitoken

import (
    "notifier-backend/src/core/config"

    "notifier-backend/src/modules/shared/infrastructure/components/datetimer"

    "notifier-backend/src/modules/apitokens/domain/entities"
    "notifier-backend/src/modules/apitokens/domain/enums"
    "notifier-backend/src/modules/apitokens/domain/exceptions"
    "notifier-backend/src/modules/apitokens/infrastructure/repositories"
)

// DisableApiTokenService revoca un token de envio: desde ese momento da 401. Los mensajes ya encolados con el se
// siguen entregando.
type DisableApiTokenService struct {
    dateTimer                       *datetimer.DateTimer
    apiTokensReaderSqliteRepository *repositories.ApiTokensReaderSqliteRepository
    apiTokensWriterSqliteRepository *repositories.ApiTokensWriterSqliteRepository
}

func DisableApiTokenServiceGetInstance() *DisableApiTokenService {
    return &DisableApiTokenService{
        dateTimer:                       datetimer.GetInstance(),
        apiTokensReaderSqliteRepository: repositories.ApiTokensReaderSqliteRepositoryGetInstance(),
        apiTokensWriterSqliteRepository: repositories.ApiTokensWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *ApiTokensException (404) o uno inesperado.
func (s *DisableApiTokenService) Invoke(disableApiTokenDto *DisableApiTokenDto) (*DisableApiTokenResultDto, error) {
    tokenUuid := disableApiTokenDto.GetTokenUuid()
    apiTokenEntity, err := s.apiTokensReaderSqliteRepository.GetApiTokenByTokenUuid(tokenUuid)
    if err != nil {
        return nil, err
    }
    if apiTokenEntity == nil {
        return nil, exceptions.NotFoundCustom(config.Tr(enums.ApiTokensMessageEnum.TOKEN_X_NOT_FOUND, tokenUuid))
    }
    if _, err := s.apiTokensWriterSqliteRepository.UpdateIsEnabledById(apiTokenEntity.GetId(), false, s.dateTimer.GetNowDtIntoTargetTz()); err != nil {
        return nil, err
    }

    disabledApiTokenEntity, err := s.apiTokensReaderSqliteRepository.GetApiTokenById(apiTokenEntity.GetId())
    if err != nil {
        return nil, err
    }
    return DisableApiTokenResultDtoFromPrimitives(map[string]any{"token": entities.ToPrimitivesOrEmpty(disabledApiTokenEntity)}), nil
}
