package repositories

import (
    "siem-client/src/modules/shared/infrastructure/components/jsonposter"
    "siem-client/src/modules/shared/infrastructure/repositories/configuration"
)

// RemoteConfigReaderHttpRepository el backend como fuente de la configuracion: GET /api/v1/agent-config devuelve los
// monitores (recursos a vigilar) que el panel ha dado a este agente.
type RemoteConfigReaderHttpRepository struct {
    jsonPoster                     *jsonposter.JsonPoster
    environmentReaderRawRepository *configuration.EnvironmentReaderRawRepository
}

func RemoteConfigReaderHttpRepositoryGetInstance() *RemoteConfigReaderHttpRepository {
    return &RemoteConfigReaderHttpRepository{
        jsonPoster:                     jsonposter.GetInstance(),
        environmentReaderRawRepository: configuration.GetInstance(),
    }
}

// GetRemoteConfigResponse codigo HTTP y body. error = no hubo respuesta (red, DNS, timeout).
func (r *RemoteConfigReaderHttpRepository) GetRemoteConfigResponse(agentConfigUrl string, apiToken string) (int, string, error) {
    userAgent := r.environmentReaderRawRepository.GetAppName() + "/" + r.environmentReaderRawRepository.GetAppVersion()
    return r.jsonPoster.GetJson(agentConfigUrl, apiToken, userAgent)
}
