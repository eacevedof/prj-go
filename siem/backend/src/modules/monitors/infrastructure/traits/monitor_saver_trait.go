package traits

import (
    "siem-backend/src/core/config"
    agentsrepositories "siem-backend/src/modules/agents/infrastructure/repositories"
    "siem-backend/src/modules/monitors/application/monitorinput"
    "siem-backend/src/modules/monitors/domain/enums"
    "siem-backend/src/modules/monitors/domain/exceptions"
    "siem-backend/src/modules/monitors/domain/services"
    "siem-backend/src/modules/monitors/infrastructure/repositories"
    sharedenums "siem-backend/src/modules/shared/domain/enums"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"
    sharedtraits "siem-backend/src/modules/shared/infrastructure/traits"
)

// MonitorSaverTrait las comprobaciones que comparten el alta y la edicion de un monitor: campos validos, agente del
// proyecto y nombre libre. Se embebe en CreateMonitorService y UpdateMonitorService.
type MonitorSaverTrait struct {
    sharedtraits.ErrorStackTrait

    MonitorFieldsValidator         *services.MonitorFieldsValidator
    MonitorsReaderSqliteRepository *repositories.MonitorsReaderSqliteRepository
    MonitorsWriterSqliteRepository *repositories.MonitorsWriterSqliteRepository
    agentsReaderSqliteRepository   *agentsrepositories.AgentsReaderSqliteRepository
}

func NewMonitorSaverTrait() MonitorSaverTrait {
    return MonitorSaverTrait{
        MonitorFieldsValidator:         services.MonitorFieldsValidatorGetInstance(),
        MonitorsReaderSqliteRepository: repositories.MonitorsReaderSqliteRepositoryGetInstance(),
        MonitorsWriterSqliteRepository: repositories.MonitorsWriterSqliteRepositoryGetInstance(),
        agentsReaderSqliteRepository:   agentsrepositories.AgentsReaderSqliteRepositoryGetInstance(),
    }
}

// FailIfWrongInput *FieldsException (400) con el detalle por campo. monitorType: el del alta o, al editar, el que ya
// tiene el monitor (no cambia).
func (t MonitorSaverTrait) FailIfWrongInput(monitorInput monitorinput.MonitorInput, monitorType string) error {
    errors := t.MonitorFieldsValidator.GetErrors(
        monitorInput.GetMonitorName(), monitorType, monitorInput.GetIntervalSeconds(), monitorInput.GetMonitorConfig(),
        monitorInput.GetWarningThreshold(), monitorInput.GetCriticalThreshold(),
    )
    if len(errors) == 0 {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        t.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
    )
}

// GetAgentIdOrFail 0 si el monitor es de todos los agentes (agent_uuid vacio); *FieldsException (400) si el agente no
// existe o es de otro proyecto.
func (t MonitorSaverTrait) GetAgentIdOrFail(agentUuid string, projectId int) (int, error) {
    if agentUuid == "" {
        return 0, nil
    }
    agentEntity, err := t.agentsReaderSqliteRepository.GetAgentByAgentUuid(agentUuid)
    if err != nil {
        return 0, err
    }
    if agentEntity == nil || agentEntity.GetProjectId() != projectId {
        return 0, sharedexceptions.FieldsExceptionBadRequestCustom(
            config.Tr(enums.MonitorsMessageEnum.AGENT_X_NOT_FOUND_IN_PROJECT, agentUuid),
        )
    }
    return agentEntity.GetId(), nil
}

// FailIfMonitorNameExists *MonitorsException (409). excludedMonitorId: el propio monitor al editar (0 en el alta).
func (t MonitorSaverTrait) FailIfMonitorNameExists(projectId int, monitorName string, excludedMonitorId int) error {
    doesMonitorNameExist, err := t.MonitorsReaderSqliteRepository.DoesMonitorNameExistInProject(projectId, monitorName, excludedMonitorId)
    if err != nil {
        return err
    }
    if doesMonitorNameExist {
        return exceptions.ConflictCustom(config.Tr(enums.MonitorsMessageEnum.MONITOR_NAME_X_ALREADY_EXISTS, monitorName))
    }
    return nil
}

// GetThresholdsOrNil folder y log_pattern no tienen valor numerico: sus umbrales se guardan siempre NULL.
func (t MonitorSaverTrait) GetThresholdsOrNil(monitorInput monitorinput.MonitorInput, monitorType string) (*float64, *float64) {
    if !t.MonitorFieldsValidator.IsMetricType(monitorType) {
        return nil, nil
    }
    return monitorInput.GetWarningThreshold(), monitorInput.GetCriticalThreshold()
}
