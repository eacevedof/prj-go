package entities

import (
    "slices"

    "siem-client/src/modules/agentconfig/domain/enums"
    "siem-client/src/modules/shared/application"
)

const DEFAULT_MONITOR_INTERVAL_SECONDS = 60

// RemoteMonitorEntity un monitor (recurso a vigilar) creado en el panel del backend, tal como llega en
// GET /api/v1/agent-config. Los umbrales no viajan: la gravedad la decide el backend. monitorName es el source_name
// de sus eventos.
type RemoteMonitorEntity struct {
    monitorUuid     string
    monitorType     string
    monitorName     string
    intervalSeconds int
    config          map[string]any
}

// RemoteMonitorEntityFromPrimitives el JSON ya decodificado (numeros como float64).
func RemoteMonitorEntityFromPrimitives(primitives map[string]any) *RemoteMonitorEntity {
    config, _ := primitives[enums.RemoteConfigKeyEnum.CONFIG].(map[string]any)
    if config == nil {
        config = map[string]any{}
    }
    intervalSeconds := application.GetInt(primitives, enums.RemoteConfigKeyEnum.INTERVAL_SECONDS)
    if intervalSeconds <= 0 {
        intervalSeconds = DEFAULT_MONITOR_INTERVAL_SECONDS
    }
    return &RemoteMonitorEntity{
        monitorUuid:     application.GetTrimmedString(primitives, enums.RemoteConfigKeyEnum.MONITOR_UUID),
        monitorType:     application.GetTrimmedString(primitives, enums.RemoteConfigKeyEnum.MONITOR_TYPE),
        monitorName:     application.GetTrimmedString(primitives, enums.RemoteConfigKeyEnum.MONITOR_NAME),
        intervalSeconds: intervalSeconds,
        config:          config,
    }
}

func (e *RemoteMonitorEntity) ToPrimitives() map[string]any {
    return map[string]any{
        enums.RemoteConfigKeyEnum.MONITOR_UUID:     e.monitorUuid,
        enums.RemoteConfigKeyEnum.MONITOR_TYPE:     e.monitorType,
        enums.RemoteConfigKeyEnum.MONITOR_NAME:     e.monitorName,
        enums.RemoteConfigKeyEnum.INTERVAL_SECONDS: e.intervalSeconds,
        enums.RemoteConfigKeyEnum.CONFIG:           e.config,
    }
}

func (e *RemoteMonitorEntity) GetMonitorUuid() string {
    return e.monitorUuid
}

func (e *RemoteMonitorEntity) GetMonitorType() string {
    return e.monitorType
}

func (e *RemoteMonitorEntity) GetMonitorName() string {
    return e.monitorName
}

func (e *RemoteMonitorEntity) GetIntervalSeconds() int {
    return e.intervalSeconds
}

// IsMetric cpu, memory, disk, gpu, network: los muestrea el modulo hostmetrics.
func (e *RemoteMonitorEntity) IsMetric() bool {
    return slices.Contains(enums.MonitorTypeEnum.METRIC_VALUES, e.monitorType)
}

// IsLines folder, log_pattern: se vigilan como watchers.
func (e *RemoteMonitorEntity) IsLines() bool {
    return slices.Contains(enums.MonitorTypeEnum.LINES_VALUES, e.monitorType)
}

func (e *RemoteMonitorEntity) GetConfigString(key string) string {
    return application.GetTrimmedString(e.config, key)
}

func (e *RemoteMonitorEntity) GetConfigBool(key string) bool {
    isTrue, _ := e.config[key].(bool)
    return isTrue
}
