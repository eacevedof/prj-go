package getagentconfig

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"

    "siem-backend/src/modules/monitors/infrastructure/repositories"
)

// GetAgentConfigService lo que tiene que vigilar un agente: los monitores activos de su proyecto que son de todos los
// agentes o suyos. El agente lo pide cada minuto con su token: crear, editar o borrar un monitor en el panel se
// aplica sin tocar el servidor.
type GetAgentConfigService struct {
    monitorsReaderSqliteRepository *repositories.MonitorsReaderSqliteRepository
}

func GetAgentConfigServiceGetInstance() *GetAgentConfigService {
    return &GetAgentConfigService{
        monitorsReaderSqliteRepository: repositories.MonitorsReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: solo uno inesperado (el 401 del token lo da antes AuthenticateAgent en el controller).
func (s *GetAgentConfigService) Invoke(getAgentConfigDto *GetAgentConfigDto) (*GetAgentConfigResultDto, error) {
    monitorEntities, err := s.monitorsReaderSqliteRepository.GetEnabledMonitorsForAgent(
        getAgentConfigDto.GetAuthProjectId(), getAgentConfigDto.GetAuthAgentId(),
    )
    if err != nil {
        return nil, err
    }

    monitors := make([]AgentMonitorPrimitives, 0, len(monitorEntities))
    for _, monitorEntity := range monitorEntities {
        monitors = append(monitors, AgentMonitorPrimitives{
            MonitorUuid:     monitorEntity.GetMonitorUuid(),
            MonitorType:     monitorEntity.GetMonitorType(),
            MonitorName:     monitorEntity.GetMonitorName(),
            IntervalSeconds: monitorEntity.GetIntervalSeconds(),
            Config:          monitorEntity.GetConfig(),
        })
    }
    return GetAgentConfigResultDtoFromPrimitives(map[string]any{
        "configHash": s.getConfigHash(monitors),
        "monitors":   monitors,
    }), nil
}

// getConfigHash json.Marshal ordena las claves de los maps: la misma lista da siempre el mismo hash.
func (s *GetAgentConfigService) getConfigHash(monitors []AgentMonitorPrimitives) string {
    monitorsJson, err := json.Marshal(monitors)
    if err != nil {
        return ""
    }
    hash := sha256.Sum256(monitorsJson)
    return hex.EncodeToString(hash[:])
}
