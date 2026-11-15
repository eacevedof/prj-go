package traits

import (
    "siem-client/src/core/config"
    "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/agentconfig/domain/enums"
    "siem-client/src/modules/agentconfig/domain/exceptions"
    "siem-client/src/modules/agentconfig/domain/services"
    "siem-client/src/modules/agentconfig/infrastructure/repositories"
    "siem-client/src/modules/shared/infrastructure/repositories/configuration"
    "siem-client/src/modules/shared/infrastructure/repositories/state"
    sharedtraits "siem-client/src/modules/shared/infrastructure/traits"
)

// AgentConfigLoaderTrait lo que hacen todos los casos de uso del agente antes de nada: montar la configuracion.
//   - local: siem-client.json (carpetas, comandos, docker). Con la config remota activa es OPCIONAL: sin fichero, vacia.
//   - remota (APP_SIEM_REMOTE_CONFIG=1, defecto): los monitores del panel del backend, guardados por SyncRemoteConfig en
//     state/remote-config.json. Se suman a la local (folder / log_pattern como watchers; metricas para hostmetrics).
//
// Se embebe en el servicio (como ErrorStackTrait) para que cada pasada relea la configuracion: cambiar algo en el
// fichero o en el panel no exige reiniciar el agente.
type AgentConfigLoaderTrait struct {
    sharedtraits.ErrorStackTrait
    environmentReaderRawRepository  *configuration.EnvironmentReaderRawRepository
    agentConfigReaderJsonRepository *repositories.AgentConfigReaderJsonRepository
    stateReaderJsonRepository       *state.StateReaderJsonRepository
}

func NewAgentConfigLoaderTrait() AgentConfigLoaderTrait {
    return AgentConfigLoaderTrait{
        environmentReaderRawRepository:  configuration.GetInstance(),
        agentConfigReaderJsonRepository: repositories.AgentConfigReaderJsonRepositoryGetInstance(),
        stateReaderJsonRepository:       state.StateReaderJsonRepositoryGetInstance(),
    }
}

// GetValidAgentConfigEntityOrFail error: *AgentConfigException (404 sin fichero y sin config remota, 400 fichero o
// valores no validos). Un monitor remoto que no vale no es un error: se ignora y queda en GetRemoteConfigWarnings.
func (t AgentConfigLoaderTrait) GetValidAgentConfigEntityOrFail() (*entities.AgentConfigEntity, error) {
    isRemoteConfigEnabled := t.environmentReaderRawRepository.IsRemoteConfigEnabled()
    agentConfigEntity, err := t.getLocalAgentConfigEntityOrFail(isRemoteConfigEnabled)
    if err != nil {
        return nil, err
    }

    errors := services.AgentConfigValidatorGetInstance().GetErrorsWithRemoteConfig(
        agentConfigEntity,
        t.environmentReaderRawRepository.GetSiemEndpointUrl(),
        t.environmentReaderRawRepository.GetSiemApiToken(),
        isRemoteConfigEnabled,
    )
    if len(errors) > 0 {
        return nil, exceptions.BadRequestCustom(
            t.GetErrorStackMessage(config.Tr(enums.AgentConfigMessageEnum.CONFIG_HAS_ERRORS), errors),
        )
    }

    if isRemoteConfigEnabled {
        services.RemoteMonitorsMergerGetInstance().Merge(agentConfigEntity, t.getSavedRemoteMonitorEntities())
    }
    return agentConfigEntity, nil
}

// getLocalAgentConfigEntityOrFail sin fichero: con la config remota, una config vacia (todo por defecto); sin ella, 404.
func (t AgentConfigLoaderTrait) getLocalAgentConfigEntityOrFail(isRemoteConfigEnabled bool) (*entities.AgentConfigEntity, error) {
    pathConfigFile := t.environmentReaderRawRepository.GetConfigPath()
    if !t.agentConfigReaderJsonRepository.DoesConfigFileExist(pathConfigFile) {
        if isRemoteConfigEnabled {
            return entities.AgentConfigEntityFromPrimitives(entities.AgentConfigPrimitives{}), nil
        }
        return nil, exceptions.NotFoundCustom(config.Tr(enums.AgentConfigMessageEnum.CONFIG_FILE_X_NOT_FOUND, pathConfigFile))
    }

    agentConfigEntity, err := t.agentConfigReaderJsonRepository.GetAgentConfigEntity(pathConfigFile)
    if err != nil {
        return nil, exceptions.BadRequestCustom(
            config.Tr(enums.AgentConfigMessageEnum.CONFIG_FILE_X_IS_NOT_VALID_X, pathConfigFile, err.Error()),
        )
    }
    return agentConfigEntity, nil
}

// getSavedRemoteMonitorEntities la ultima config remota guardada. Sin fichero todavia (o ilegible): ninguno; el
// agente arranca igual y espera a la siguiente sincronizacion.
func (t AgentConfigLoaderTrait) getSavedRemoteMonitorEntities() []*entities.RemoteMonitorEntity {
    remoteConfig, err := t.stateReaderJsonRepository.GetState(t.environmentReaderRawRepository.GetStatePath(), enums.RemoteConfigKeyEnum.STATE_FILE)
    if err != nil {
        return []*entities.RemoteMonitorEntity{}
    }
    return services.RemoteConfigParserGetInstance().GetRemoteMonitorEntities(remoteConfig)
}

func (t AgentConfigLoaderTrait) GetConfigPath() string {
    return t.environmentReaderRawRepository.GetConfigPath()
}

func (t AgentConfigLoaderTrait) IsRemoteConfigEnabled() bool {
    return t.environmentReaderRawRepository.IsRemoteConfigEnabled()
}

func (t AgentConfigLoaderTrait) GetSiemEndpointUrl() string {
    return t.environmentReaderRawRepository.GetSiemEndpointUrl()
}

func (t AgentConfigLoaderTrait) GetSiemApiToken() string {
    return t.environmentReaderRawRepository.GetSiemApiToken()
}

func (t AgentConfigLoaderTrait) GetSpoolPath() string {
    return t.environmentReaderRawRepository.GetSpoolPath()
}

func (t AgentConfigLoaderTrait) GetStatePath() string {
    return t.environmentReaderRawRepository.GetStatePath()
}
