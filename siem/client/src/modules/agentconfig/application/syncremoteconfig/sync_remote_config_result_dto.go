package syncremoteconfig

import "siem-client/src/modules/shared/application"

// SyncRemoteConfigPrimitives lo que imprime config-sync (claves en este orden).
type SyncRemoteConfigPrimitives struct {
    IsDisabled     bool   `json:"is_disabled"`
    AgentConfigUrl string `json:"agent_config_url"`
    ConfigHash     string `json:"config_hash"`
    IsChanged      bool   `json:"is_changed"`
    Monitors       int    `json:"monitors"`
}

type SyncRemoteConfigResultDto struct {
    isDisabled     bool
    agentConfigUrl string
    configHash     string
    // la config que ha llegado es distinta de la guardada (y se ha guardado)
    isChanged bool
    monitors  int
}

func SyncRemoteConfigResultDtoFromPrimitives(primitives map[string]any) *SyncRemoteConfigResultDto {
    return &SyncRemoteConfigResultDto{
        isDisabled:     application.GetInt(primitives, "isDisabled") == 1,
        agentConfigUrl: application.GetTrimmedString(primitives, "agentConfigUrl"),
        configHash:     application.GetTrimmedString(primitives, "configHash"),
        isChanged:      application.GetInt(primitives, "isChanged") == 1,
        monitors:       application.GetInt(primitives, "monitors"),
    }
}

func (d *SyncRemoteConfigResultDto) IsDisabled() bool {
    return d.isDisabled
}

func (d *SyncRemoteConfigResultDto) IsChanged() bool {
    return d.isChanged
}

func (d *SyncRemoteConfigResultDto) GetMonitors() int {
    return d.monitors
}

func (d *SyncRemoteConfigResultDto) ToArray() SyncRemoteConfigPrimitives {
    return SyncRemoteConfigPrimitives{
        IsDisabled:     d.isDisabled,
        AgentConfigUrl: d.agentConfigUrl,
        ConfigHash:     d.configHash,
        IsChanged:      d.isChanged,
        Monitors:       d.monitors,
    }
}
