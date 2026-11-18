package monitorinput

import (
    "siem-backend/src/modules/shared/application"

    "siem-backend/src/modules/monitors/domain/enums"
    "siem-backend/src/modules/monitors/domain/entities"
)

// MonitorInput los campos editables de un monitor, comunes al alta y a la edicion: los embeben CreateMonitorDto y
// UpdateMonitorDto. Solo primitivos. is_enabled ausente = activo; interval_seconds ausente = 60 (connections, 30).
type MonitorInput struct {
    agentUuid         string
    monitorType       string
    monitorName       string
    isEnabled         bool
    intervalSeconds   int
    monitorConfig     map[string]any
    warningThreshold  *float64
    criticalThreshold *float64
}

func MonitorInputFromPrimitives(primitives map[string]any) MonitorInput {
    isEnabled := true
    if rawIsEnabled, isSet := primitives["isEnabled"]; isSet && rawIsEnabled != nil {
        isEnabled = application.ToInt(rawIsEnabled) == 1
    }
    monitorConfig, _ := primitives["config"].(map[string]any)
    if monitorConfig == nil {
        monitorConfig = map[string]any{}
    }
    monitorType := application.GetTrimmedString(primitives, "monitorType")
    defaultIntervalSeconds := enums.MonitorFieldsLimitEnum.INTERVAL_SECONDS_DEFAULT
    if monitorType == enums.MonitorTypeEnum.CONNECTIONS {
        defaultIntervalSeconds = enums.MonitorFieldsLimitEnum.CONNECTIONS_INTERVAL_SECONDS_DEFAULT
    }
    return MonitorInput{
        agentUuid:         application.GetTrimmedString(primitives, "agentUuid"),
        monitorType:       monitorType,
        monitorName:       application.GetTrimmedString(primitives, "monitorName"),
        isEnabled:         isEnabled,
        intervalSeconds:   application.GetIntOrDefault(primitives, "intervalSeconds", defaultIntervalSeconds),
        monitorConfig:     monitorConfig,
        warningThreshold:  entities.GetFloatOrNil(primitives["warningThreshold"]),
        criticalThreshold: entities.GetFloatOrNil(primitives["criticalThreshold"]),
    }
}

// GetAgentUuid "" = todos los agentes del proyecto.
func (i MonitorInput) GetAgentUuid() string {
    return i.agentUuid
}

func (i MonitorInput) GetMonitorType() string {
    return i.monitorType
}

func (i MonitorInput) GetMonitorName() string {
    return i.monitorName
}

func (i MonitorInput) IsEnabled() bool {
    return i.isEnabled
}

func (i MonitorInput) GetIntervalSeconds() int {
    return i.intervalSeconds
}

func (i MonitorInput) GetMonitorConfig() map[string]any {
    return i.monitorConfig
}

func (i MonitorInput) GetWarningThreshold() *float64 {
    return i.warningThreshold
}

func (i MonitorInput) GetCriticalThreshold() *float64 {
    return i.criticalThreshold
}
