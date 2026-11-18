package services

import (
    "maps"
    "net"
    "strconv"
    "strings"
    "time"

    "siem-client/src/core/config"

    "siem-client/src/modules/shared/application"

    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"

    collectorsenums "siem-client/src/modules/collectors/domain/enums"

    "siem-client/src/modules/connections/domain/enums"
    "siem-client/src/modules/connections/domain/entities"
)

// ConnectionEventBuilder monta el evento net_connection con el contrato del backend. source_name = el proceso (asi las
// facetas del visor dicen "que aplicacion"); attributes = los primitivos de la conexion + accion + monitor. La
// gravedad no se manda: la decide el backend (entrantes desde IP publica y escuchas nuevas, info; el resto, debug).
// Formato del mensaje: ConnectionEventMessageEnum.
type ConnectionEventBuilder struct{}

func ConnectionEventBuilderGetInstance() *ConnectionEventBuilder {
    return &ConnectionEventBuilder{}
}

func (b *ConnectionEventBuilder) GetChangeEvent(remoteMonitorEntity *agentconfigentities.RemoteMonitorEntity, connectionChangeEntity *entities.ConnectionChangeEntity, occurredAt time.Time) map[string]any {
    connectionPrimitives := connectionChangeEntity.GetConnectionPrimitives()
    attributes := b.getBaseAttributes(remoteMonitorEntity)
    maps.Copy(attributes, connectionPrimitives)
    attributes[enums.NetConnectionAttributeEnum.ACTION] = connectionChangeEntity.GetAction()
    if connectionChangeEntity.IsClose() {
        attributes[enums.NetConnectionAttributeEnum.DURATION_SECONDS] = connectionChangeEntity.GetDurationSeconds()
    }
    if connectionChangeEntity.IsInitial() {
        attributes[enums.NetConnectionAttributeEnum.IS_INITIAL] = true
    }

    sourceName := application.GetTrimmedString(connectionPrimitives, enums.ConnectionKeyEnum.PROCESS)
    if sourceName == "" {
        sourceName = enums.UNKNOWN_PROCESS
    }
    return b.getEvent(sourceName, b.getMessage(connectionChangeEntity), attributes, occurredAt)
}

// GetUnavailableEvent en este sistema no se puede muestrear (no es Linux, o /proc no se puede leer): sin conexion,
// con is_unavailable y el motivo en el mensaje. source_name = el monitor. Quien lo llama lo manda una sola vez.
func (b *ConnectionEventBuilder) GetUnavailableEvent(remoteMonitorEntity *agentconfigentities.RemoteMonitorEntity, reason string, occurredAt time.Time) map[string]any {
    attributes := b.getBaseAttributes(remoteMonitorEntity)
    attributes[enums.NetConnectionAttributeEnum.IS_UNAVAILABLE] = true
    return b.getEvent(
        remoteMonitorEntity.GetMonitorName(),
        config.Tr(enums.ConnectionEventMessageEnum.X_NOT_AVAILABLE_X, remoteMonitorEntity.GetMonitorType(), reason),
        attributes,
        occurredAt,
    )
}

func (b *ConnectionEventBuilder) getMessage(connectionChangeEntity *entities.ConnectionChangeEntity) string {
    connectionPrimitives := connectionChangeEntity.GetConnectionPrimitives()
    action := connectionChangeEntity.GetAction()
    protocol := application.GetTrimmedString(connectionPrimitives, enums.ConnectionKeyEnum.PROTOCOL)
    localEndpoint := b.getEndpoint(connectionPrimitives, enums.ConnectionKeyEnum.LOCAL_IP, enums.ConnectionKeyEnum.LOCAL_PORT)

    processLabel := b.getProcessLabel(connectionPrimitives)

    messageParts := []string{}
    if connectionChangeEntity.IsListen() {
        messageParts = append(messageParts, config.Tr(enums.ConnectionEventMessageEnum.LISTEN_X, action, protocol, processLabel, localEndpoint))
    } else {
        messageParts = append(messageParts, config.Tr(
            enums.ConnectionEventMessageEnum.CONNECTION_X,
            action,
            protocol,
            application.GetTrimmedString(connectionPrimitives, enums.ConnectionKeyEnum.DIRECTION),
            processLabel,
            localEndpoint,
            b.getEndpoint(connectionPrimitives, enums.ConnectionKeyEnum.REMOTE_IP, enums.ConnectionKeyEnum.REMOTE_PORT),
        ))
    }
    if connectionChangeEntity.IsClose() {
        messageParts = append(messageParts, config.Tr(enums.ConnectionEventMessageEnum.AFTER_X_SECONDS, connectionChangeEntity.GetDurationSeconds()))
    }
    if connectionChangeEntity.IsInitial() {
        messageParts = append(messageParts, enums.ConnectionEventMessageEnum.ALREADY_OPEN)
    }
    return strings.Join(messageParts, enums.ConnectionEventMessageEnum.MESSAGE_PART_SEPARATOR)
}

// getProcessLabel "firefox(4321)" o "?" si no se sabe de quien es.
func (b *ConnectionEventBuilder) getProcessLabel(connectionPrimitives map[string]any) string {
    processName := application.GetTrimmedString(connectionPrimitives, enums.ConnectionKeyEnum.PROCESS)
    pid := application.GetInt(connectionPrimitives, enums.ConnectionKeyEnum.PID)
    if processName == "" || pid == 0 {
        return enums.UNKNOWN_PROCESS
    }
    return config.Tr(enums.ConnectionFormatEnum.PROCESS_LABEL, processName, pid)
}

// getEndpoint "192.168.1.5:52780" / "[::1]:631"; "*" sin IP (el remoto de una escucha).
func (b *ConnectionEventBuilder) getEndpoint(connectionPrimitives map[string]any, ipKey string, portKey string) string {
    ip := application.GetTrimmedString(connectionPrimitives, ipKey)
    if ip == "" {
        return enums.ConnectionFormatEnum.ANY_ENDPOINT
    }
    return net.JoinHostPort(ip, strconv.Itoa(application.GetInt(connectionPrimitives, portKey)))
}

func (b *ConnectionEventBuilder) getBaseAttributes(remoteMonitorEntity *agentconfigentities.RemoteMonitorEntity) map[string]any {
    return map[string]any{
        enums.NetConnectionAttributeEnum.MONITOR_UUID: remoteMonitorEntity.GetMonitorUuid(),
        enums.NetConnectionAttributeEnum.MONITOR_TYPE: remoteMonitorEntity.GetMonitorType(),
    }
}

func (b *ConnectionEventBuilder) getEvent(sourceName string, message string, attributes map[string]any, occurredAt time.Time) map[string]any {
    return map[string]any{
        collectorsenums.CollectorsKeyEnum.EVENT_TYPE:  enums.NET_CONNECTION_EVENT_TYPE,
        collectorsenums.CollectorsKeyEnum.SOURCE_NAME: sourceName,
        collectorsenums.CollectorsKeyEnum.MESSAGE:     message,
        collectorsenums.CollectorsKeyEnum.OCCURRED_AT: occurredAt.UTC().Format(enums.EVENT_DATE_TIME_FORMAT),
        collectorsenums.CollectorsKeyEnum.ATTRIBUTES:  attributes,
    }
}
