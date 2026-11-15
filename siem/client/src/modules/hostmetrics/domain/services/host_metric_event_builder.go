package services

import (
    "maps"
    "time"

    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"
    collectorsenums "siem-client/src/modules/collectors/domain/enums"
    "siem-client/src/modules/hostmetrics/domain/entities"
    "siem-client/src/modules/hostmetrics/domain/enums"
)

const dateTimeFormat = "2006-01-02 15:04:05"

// HostMetricEventBuilder monta el evento host_metric con el contrato del backend. source_name = nombre del monitor.
// La gravedad no se manda: la decide el backend comparando attributes.value con los umbrales del monitor (que no
// viajan al agente).
type HostMetricEventBuilder struct{}

func HostMetricEventBuilderGetInstance() *HostMetricEventBuilder {
    return &HostMetricEventBuilder{}
}

func (b *HostMetricEventBuilder) GetMetricEvent(remoteMonitorEntity *agentconfigentities.RemoteMonitorEntity, hostMetricSampleEntity *entities.HostMetricSampleEntity, occurredAt time.Time) map[string]any {
    attributes := b.getBaseAttributes(remoteMonitorEntity)
    maps.Copy(attributes, hostMetricSampleEntity.GetAttributes())
    attributes[enums.HostMetricAttributeEnum.VALUE] = hostMetricSampleEntity.GetValue()
    attributes[enums.HostMetricAttributeEnum.UNIT] = hostMetricSampleEntity.GetUnit()
    return b.getEvent(remoteMonitorEntity, hostMetricSampleEntity.GetMessage(), attributes, occurredAt)
}

// GetUnavailableEvent la metrica no se puede tomar en este host (sin nvidia-smi, SO no soportado, path que no
// existe): sin value, con is_unavailable y el motivo en el mensaje. Quien lo llama lo manda una sola vez.
func (b *HostMetricEventBuilder) GetUnavailableEvent(remoteMonitorEntity *agentconfigentities.RemoteMonitorEntity, reason string, occurredAt time.Time) map[string]any {
    attributes := b.getBaseAttributes(remoteMonitorEntity)
    attributes[enums.HostMetricAttributeEnum.IS_UNAVAILABLE] = true
    return b.getEvent(remoteMonitorEntity, remoteMonitorEntity.GetMonitorType()+" not available: "+reason, attributes, occurredAt)
}

func (b *HostMetricEventBuilder) getBaseAttributes(remoteMonitorEntity *agentconfigentities.RemoteMonitorEntity) map[string]any {
    return map[string]any{
        enums.HostMetricAttributeEnum.MONITOR_UUID: remoteMonitorEntity.GetMonitorUuid(),
        enums.HostMetricAttributeEnum.MONITOR_TYPE: remoteMonitorEntity.GetMonitorType(),
    }
}

func (b *HostMetricEventBuilder) getEvent(remoteMonitorEntity *agentconfigentities.RemoteMonitorEntity, message string, attributes map[string]any, occurredAt time.Time) map[string]any {
    return map[string]any{
        collectorsenums.CollectorsKeyEnum.EVENT_TYPE:  enums.HOST_METRIC_EVENT_TYPE,
        collectorsenums.CollectorsKeyEnum.SOURCE_NAME: remoteMonitorEntity.GetMonitorName(),
        collectorsenums.CollectorsKeyEnum.MESSAGE:     message,
        collectorsenums.CollectorsKeyEnum.OCCURRED_AT: occurredAt.UTC().Format(dateTimeFormat),
        collectorsenums.CollectorsKeyEnum.ATTRIBUTES:  attributes,
    }
}
