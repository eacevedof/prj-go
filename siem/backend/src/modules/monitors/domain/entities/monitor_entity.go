package entities

import (
    "encoding/json"

    "siem-backend/src/modules/monitors/domain/enums"
    "siem-backend/src/modules/shared/application"
)

const MONITOR_UUID_PREFIX = "monitor-"

// MonitorPrimitives un monitor tal como lo devuelve la API (claves en este orden). agent_uuid "" = todos los agentes
// del proyecto. Un umbral null = sin umbral.
type MonitorPrimitives struct {
    MonitorUuid       string         `json:"monitor_uuid"`
    ProjectUuid       string         `json:"project_uuid"`
    AgentUuid         string         `json:"agent_uuid"`
    AgentName         string         `json:"agent_name"`
    MonitorType       string         `json:"monitor_type"`
    MonitorName       string         `json:"monitor_name"`
    IsEnabled         bool           `json:"is_enabled"`
    IntervalSeconds   int            `json:"interval_seconds"`
    Config            map[string]any `json:"config"`
    WarningThreshold  *float64       `json:"warning_threshold"`
    CriticalThreshold *float64       `json:"critical_threshold"`
    CreatedAt         string         `json:"created_at"`
    UpdatedAt         string         `json:"updated_at"`
}

// MonitorEntity un recurso a vigilar de un proyecto. agentId 0 = lo aplican todos sus agentes.
type MonitorEntity struct {
    id                int
    monitorUuid       string
    projectId         int
    projectUuid       string
    agentId           int
    agentUuid         string
    agentName         string
    monitorType       string
    monitorName       string
    isEnabled         bool
    intervalSeconds   int
    config            map[string]any
    warningThreshold  *float64
    criticalThreshold *float64
    createdAt         string
    updatedAt         string
}

// MonitorEntityFromPrimitives de una fila de BD (config como texto JSON, umbrales NULL o numero) o de los primitivos
// de un servicio (config como map, umbrales *float64).
func MonitorEntityFromPrimitives(primitives map[string]any) *MonitorEntity {
    return &MonitorEntity{
        id:                application.GetInt(primitives, enums.MonitorsKeyEnum.ID),
        monitorUuid:       application.GetTrimmedString(primitives, enums.MonitorsKeyEnum.MONITOR_UUID),
        projectId:         application.GetInt(primitives, enums.MonitorsKeyEnum.PROJECT_ID),
        projectUuid:       application.GetTrimmedString(primitives, enums.MonitorsKeyEnum.PROJECT_UUID),
        agentId:           application.GetInt(primitives, enums.MonitorsKeyEnum.AGENT_ID),
        agentUuid:         application.GetTrimmedString(primitives, enums.MonitorsKeyEnum.AGENT_UUID),
        agentName:         application.GetTrimmedString(primitives, enums.MonitorsKeyEnum.AGENT_NAME),
        monitorType:       application.GetTrimmedString(primitives, enums.MonitorsKeyEnum.MONITOR_TYPE),
        monitorName:       application.GetTrimmedString(primitives, enums.MonitorsKeyEnum.MONITOR_NAME),
        isEnabled:         application.GetInt(primitives, enums.MonitorsKeyEnum.IS_ENABLED) == 1,
        intervalSeconds:   application.GetInt(primitives, enums.MonitorsKeyEnum.INTERVAL_SECONDS),
        config:            getConfig(primitives[enums.MonitorsKeyEnum.CONFIG]),
        warningThreshold:  GetFloatOrNil(primitives[enums.MonitorsKeyEnum.WARNING_THRESHOLD]),
        criticalThreshold: GetFloatOrNil(primitives[enums.MonitorsKeyEnum.CRITICAL_THRESHOLD]),
        createdAt:         application.GetTrimmedString(primitives, enums.MonitorsKeyEnum.CREATED_AT),
        updatedAt:         application.GetTrimmedString(primitives, enums.MonitorsKeyEnum.UPDATED_AT),
    }
}

func (e *MonitorEntity) ToPrimitives() MonitorPrimitives {
    return MonitorPrimitives{
        MonitorUuid:       e.monitorUuid,
        ProjectUuid:       e.projectUuid,
        AgentUuid:         e.agentUuid,
        AgentName:         e.agentName,
        MonitorType:       e.monitorType,
        MonitorName:       e.monitorName,
        IsEnabled:         e.isEnabled,
        IntervalSeconds:   e.intervalSeconds,
        Config:            e.config,
        WarningThreshold:  e.warningThreshold,
        CriticalThreshold: e.criticalThreshold,
        CreatedAt:         e.createdAt,
        UpdatedAt:         e.updatedAt,
    }
}

func (e *MonitorEntity) GetId() int {
    return e.id
}

func (e *MonitorEntity) GetMonitorUuid() string {
    return e.monitorUuid
}

func (e *MonitorEntity) GetProjectId() int {
    return e.projectId
}

func (e *MonitorEntity) GetAgentId() int {
    return e.agentId
}

func (e *MonitorEntity) GetAgentName() string {
    return e.agentName
}

func (e *MonitorEntity) GetMonitorType() string {
    return e.monitorType
}

func (e *MonitorEntity) GetMonitorName() string {
    return e.monitorName
}

func (e *MonitorEntity) IsEnabled() bool {
    return e.isEnabled
}

func (e *MonitorEntity) GetIntervalSeconds() int {
    return e.intervalSeconds
}

func (e *MonitorEntity) GetConfig() map[string]any {
    return e.config
}

// GetConfigJson el texto que se guarda en monitors.config.
func (e *MonitorEntity) GetConfigJson() string {
    configJson, err := json.Marshal(e.config)
    if err != nil {
        return "{}"
    }
    return string(configJson)
}

func (e *MonitorEntity) GetWarningThreshold() *float64 {
    return e.warningThreshold
}

func (e *MonitorEntity) GetCriticalThreshold() *float64 {
    return e.criticalThreshold
}

func (e *MonitorEntity) GetCreatedAt() string {
    return e.createdAt
}

func (e *MonitorEntity) GetUpdatedAt() string {
    return e.updatedAt
}

// getConfig un map tal cual; un texto JSON, decodificado; otra cosa, config vacia.
func getConfig(rawConfig any) map[string]any {
    switch typedConfig := rawConfig.(type) {
    case map[string]any:
        return typedConfig
    case string:
        return getConfigFromJson([]byte(typedConfig))
    case []byte:
        return getConfigFromJson(typedConfig)
    }
    return map[string]any{}
}

func getConfigFromJson(configJson []byte) map[string]any {
    config := map[string]any{}
    if json.Unmarshal(configJson, &config) != nil || config == nil {
        return map[string]any{}
    }
    return config
}

// GetFloatOrNil un umbral: numero (JSON o BD) -> su valor; nil, "" o algo que no es un numero -> nil (sin umbral).
func GetFloatOrNil(rawValue any) *float64 {
    var value float64
    switch typedValue := rawValue.(type) {
    case *float64:
        return typedValue
    case float64:
        value = typedValue
    case int:
        value = float64(typedValue)
    case int64:
        value = float64(typedValue)
    case json.Number:
        parsedValue, err := typedValue.Float64()
        if err != nil {
            return nil
        }
        value = parsedValue
    default:
        return nil
    }
    return &value
}

func ToPrimitivesOrEmpty(monitorEntity *MonitorEntity) MonitorPrimitives {
    if monitorEntity == nil {
        return MonitorPrimitives{Config: map[string]any{}}
    }
    return monitorEntity.ToPrimitives()
}
