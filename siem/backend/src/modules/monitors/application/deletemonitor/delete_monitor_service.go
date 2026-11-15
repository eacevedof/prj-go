package deletemonitor

import (
    "siem-backend/src/core/config"
    "siem-backend/src/modules/monitors/domain/enums"
    "siem-backend/src/modules/monitors/domain/exceptions"
    "siem-backend/src/modules/monitors/infrastructure/repositories"
)

// DeleteMonitorService baja de un monitor desde el panel. El agente deja de vigilarlo en su siguiente sincronizacion;
// los eventos que ya mando se quedan (se van con la retencion).
type DeleteMonitorService struct {
    monitorsReaderSqliteRepository *repositories.MonitorsReaderSqliteRepository
    monitorsWriterSqliteRepository *repositories.MonitorsWriterSqliteRepository
}

func DeleteMonitorServiceGetInstance() *DeleteMonitorService {
    return &DeleteMonitorService{
        monitorsReaderSqliteRepository: repositories.MonitorsReaderSqliteRepositoryGetInstance(),
        monitorsWriterSqliteRepository: repositories.MonitorsWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *MonitorsException (404) o uno inesperado.
func (s *DeleteMonitorService) Invoke(deleteMonitorDto *DeleteMonitorDto) (*DeleteMonitorResultDto, error) {
    monitorEntity, err := s.monitorsReaderSqliteRepository.GetMonitorByMonitorUuid(deleteMonitorDto.GetMonitorUuid())
    if err != nil {
        return nil, err
    }
    if monitorEntity == nil {
        return nil, exceptions.NotFoundCustom(config.Tr(enums.MonitorsMessageEnum.MONITOR_X_NOT_FOUND, deleteMonitorDto.GetMonitorUuid()))
    }
    if _, err := s.monitorsWriterSqliteRepository.DeleteMonitorById(monitorEntity.GetId()); err != nil {
        return nil, err
    }
    return DeleteMonitorResultDtoFromPrimitives(map[string]any{"monitorUuid": monitorEntity.GetMonitorUuid()}), nil
}
