package updatemonitor

import (
    "siem-backend/src/core/config"

    "siem-backend/src/modules/shared/infrastructure/components/datetimer"

    "siem-backend/src/modules/monitors/domain/enums"
    "siem-backend/src/modules/monitors/domain/entities"
    "siem-backend/src/modules/monitors/infrastructure/traits"
    "siem-backend/src/modules/monitors/domain/exceptions"
)

// UpdateMonitorService edicion de un monitor desde el panel: nombre, ambito (agente), activo, intervalo, config y
// umbrales. Cambiar el nombre abre una fuente nueva en cada agente (source_name = monitor_name); los eventos antiguos
// se quedan en la anterior.
type UpdateMonitorService struct {
    traits.MonitorSaverTrait

    dateTimer *datetimer.DateTimer
}

func UpdateMonitorServiceGetInstance() *UpdateMonitorService {
    return &UpdateMonitorService{
        MonitorSaverTrait: traits.NewMonitorSaverTrait(),
        dateTimer:         datetimer.GetInstance(),
    }
}

// Invoke error: *MonitorsException (404, 409), *FieldsException (400) o uno inesperado.
func (s *UpdateMonitorService) Invoke(updateMonitorDto *UpdateMonitorDto) (*UpdateMonitorResultDto, error) {
    monitorEntity, err := s.MonitorsReaderSqliteRepository.GetMonitorByMonitorUuid(updateMonitorDto.GetMonitorUuid())
    if err != nil {
        return nil, err
    }
    if monitorEntity == nil {
        return nil, exceptions.NotFoundCustom(config.Tr(enums.MonitorsMessageEnum.MONITOR_X_NOT_FOUND, updateMonitorDto.GetMonitorUuid()))
    }

    monitorType := monitorEntity.GetMonitorType()
    if err := s.FailIfWrongInput(updateMonitorDto.MonitorInput, monitorType); err != nil {
        return nil, err
    }
    agentId, err := s.GetAgentIdOrFail(updateMonitorDto.GetAgentUuid(), monitorEntity.GetProjectId())
    if err != nil {
        return nil, err
    }
    if err := s.FailIfMonitorNameExists(monitorEntity.GetProjectId(), updateMonitorDto.GetMonitorName(), monitorEntity.GetId()); err != nil {
        return nil, err
    }
    normalizedConfig := s.MonitorFieldsValidator.GetNormalizedConfig(monitorType, updateMonitorDto.GetMonitorConfig())
    if err := s.FailIfMonitorResourceExists(monitorEntity.GetProjectId(), agentId, monitorType, normalizedConfig, monitorEntity.GetId()); err != nil {
        return nil, err
    }

    warningThreshold, criticalThreshold := s.GetThresholdsOrNil(updateMonitorDto.MonitorInput, monitorType)
    _, err = s.MonitorsWriterSqliteRepository.UpdateMonitorById(entities.MonitorEntityFromPrimitives(map[string]any{
        enums.MonitorsKeyEnum.ID:                 monitorEntity.GetId(),
        enums.MonitorsKeyEnum.AGENT_ID:           agentId,
        enums.MonitorsKeyEnum.MONITOR_NAME:       updateMonitorDto.GetMonitorName(),
        enums.MonitorsKeyEnum.IS_ENABLED:         updateMonitorDto.IsEnabled(),
        enums.MonitorsKeyEnum.INTERVAL_SECONDS:   updateMonitorDto.GetIntervalSeconds(),
        enums.MonitorsKeyEnum.CONFIG:             normalizedConfig,
        enums.MonitorsKeyEnum.WARNING_THRESHOLD:  warningThreshold,
        enums.MonitorsKeyEnum.CRITICAL_THRESHOLD: criticalThreshold,
        enums.MonitorsKeyEnum.UPDATED_AT:         s.dateTimer.GetNowDtIntoTargetTz(),
    }))
    if err != nil {
        return nil, err
    }

    updatedMonitorEntity, err := s.MonitorsReaderSqliteRepository.GetMonitorById(monitorEntity.GetId())
    if err != nil {
        return nil, err
    }
    return UpdateMonitorResultDtoFromPrimitives(map[string]any{
        "monitor": entities.ToPrimitivesOrEmpty(updatedMonitorEntity),
    }), nil
}
