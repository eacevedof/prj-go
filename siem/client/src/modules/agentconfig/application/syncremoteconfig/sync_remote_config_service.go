package syncremoteconfig

import (
    "net/http"
    "strconv"
    "time"

    "siem-client/src/core/config"
    "siem-client/src/modules/agentconfig/domain/enums"
    "siem-client/src/modules/agentconfig/domain/exceptions"
    "siem-client/src/modules/agentconfig/domain/services"
    "siem-client/src/modules/agentconfig/infrastructure/repositories"
    "siem-client/src/modules/shared/application"
    "siem-client/src/modules/shared/infrastructure/repositories/configuration"
    "siem-client/src/modules/shared/infrastructure/repositories/state"
)

const dateTimeFormat = "2006-01-02 15:04:05"

// SyncRemoteConfigService baja del backend lo que este agente tiene que vigilar (los monitores que se crean en el
// panel) y lo guarda en state/remote-config.json; de ahi lo lee AgentConfigLoaderTrait en cada pasada. Solo se guarda
// una respuesta 200 con la forma esperada: si el backend no responde o contesta otra cosa, se conserva la ultima
// config guardada (el agente sigue vigilando lo mismo) y se devuelve el error para avisarlo.
type SyncRemoteConfigService struct {
    environmentReaderRawRepository   *configuration.EnvironmentReaderRawRepository
    remoteConfigReaderHttpRepository *repositories.RemoteConfigReaderHttpRepository
    remoteConfigParser               *services.RemoteConfigParser
    stateReaderJsonRepository        *state.StateReaderJsonRepository
    stateWriterJsonRepository        *state.StateWriterJsonRepository
}

func SyncRemoteConfigServiceGetInstance() *SyncRemoteConfigService {
    return &SyncRemoteConfigService{
        environmentReaderRawRepository:   configuration.GetInstance(),
        remoteConfigReaderHttpRepository: repositories.RemoteConfigReaderHttpRepositoryGetInstance(),
        remoteConfigParser:               services.RemoteConfigParserGetInstance(),
        stateReaderJsonRepository:        state.StateReaderJsonRepositoryGetInstance(),
        stateWriterJsonRepository:        state.StateWriterJsonRepositoryGetInstance(),
    }
}

// Invoke error: *AgentConfigException (401 token rechazado; 503 sin respuesta o respuesta inesperada) o uno
// inesperado (disco). Con APP_SIEM_REMOTE_CONFIG=0 no hace nada (isDisabled).
func (s *SyncRemoteConfigService) Invoke(syncRemoteConfigDto *SyncRemoteConfigDto) (*SyncRemoteConfigResultDto, error) {
    if !s.environmentReaderRawRepository.IsRemoteConfigEnabled() {
        return SyncRemoteConfigResultDtoFromPrimitives(map[string]any{"isDisabled": 1}), nil
    }

    agentConfigUrl := s.remoteConfigParser.GetAgentConfigUrl(s.environmentReaderRawRepository.GetSiemEndpointUrl())
    statusCode, responseBody, err := s.remoteConfigReaderHttpRepository.GetRemoteConfigResponse(
        agentConfigUrl, s.environmentReaderRawRepository.GetSiemApiToken(),
    )
    if err != nil {
        return nil, exceptions.ServiceUnavailableCustom(config.Tr(enums.AgentConfigMessageEnum.REMOTE_CONFIG_UNAVAILABLE_X, err.Error()))
    }
    if statusCode == http.StatusUnauthorized {
        return nil, exceptions.UnauthorizedCustom(enums.AgentConfigMessageEnum.REMOTE_CONFIG_REJECTED_API_TOKEN)
    }
    if statusCode != http.StatusOK {
        return nil, exceptions.ServiceUnavailableCustom(config.Tr(enums.AgentConfigMessageEnum.REMOTE_CONFIG_UNEXPECTED_STATUS_X, strconv.Itoa(statusCode)))
    }
    remoteConfig, err := s.remoteConfigParser.GetRemoteConfigOrFail(responseBody)
    if err != nil {
        return nil, exceptions.ServiceUnavailableCustom(config.Tr(enums.AgentConfigMessageEnum.REMOTE_CONFIG_IS_NOT_VALID_X, err.Error()))
    }

    statePath := s.environmentReaderRawRepository.GetStatePath()
    savedRemoteConfig, err := s.stateReaderJsonRepository.GetState(statePath, enums.RemoteConfigKeyEnum.STATE_FILE)
    if err != nil {
        // un fichero de estado corrupto no debe impedir guardar la config buena que acaba de llegar
        savedRemoteConfig = map[string]any{}
    }
    configHash := application.GetTrimmedString(remoteConfig, enums.RemoteConfigKeyEnum.CONFIG_HASH)
    _, hasSavedMonitors := savedRemoteConfig[enums.RemoteConfigKeyEnum.MONITORS]
    isChanged := !hasSavedMonitors || configHash == "" || configHash != application.GetTrimmedString(savedRemoteConfig, enums.RemoteConfigKeyEnum.CONFIG_HASH)
    if isChanged {
        remoteConfig[enums.RemoteConfigKeyEnum.SYNCED_AT] = time.Now().UTC().Format(dateTimeFormat)
        if err := s.stateWriterJsonRepository.SaveState(statePath, enums.RemoteConfigKeyEnum.STATE_FILE, remoteConfig); err != nil {
            return nil, err
        }
    }

    isChangedAsInt := 0
    if isChanged {
        isChangedAsInt = 1
    }
    return SyncRemoteConfigResultDtoFromPrimitives(map[string]any{
        "agentConfigUrl": agentConfigUrl,
        "configHash":     configHash,
        "isChanged":      isChangedAsInt,
        "monitors":       len(s.remoteConfigParser.GetRemoteMonitorEntities(remoteConfig)),
    }), nil
}
