package services

import (
    "maps"
    "net"
    "strconv"
    "strings"
    "time"

    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"

    collectorsenums "siem-client/src/modules/collectors/domain/enums"

    "siem-client/src/modules/connections/domain/entities"
    "siem-client/src/modules/connections/domain/enums"
    "siem-client/src/modules/shared/application"
)

const (
    dateTimeFormat = "2006-01-02 15:04:05"
    // un socket en escucha no tiene extremo remoto
    anyRemoteEndpoint = "*"
)

// ConnectionEventBuilder monta el evento net_connection con el contrato del backend. source_name = el proceso (asi las
// facetas del visor dicen "que aplicacion"); attributes = los primitivos de la conexion + accion + monitor. La
// gravedad no se manda: la decide el backend (entrantes desde IP publica y escuchas nuevas, info; el resto, debug).
//
//   opened tcp outbound firefox(4321) 192.168.1.5:52780 -> 142.250.1.1:443
//   closed tcp outbound firefox(4321) 192.168.1.5:52780 -> 142.250.1.1:443 after 12 s
//   listen_opened tcp sshd(812) 0.0.0.0:22
type ConnectionEventBuilder struct{}

func ConnectionEventBuilderGetInstance() *ConnectionEventBuilder {
    return &ConnectionEventBuilder{}
}

func (b *ConnectionEventBuilder) GetChangeEvent(remoteMonitorEntity *agentconfigentities.RemoteMonitorEntity, connectionChangeEntity *entities.ConnectionChangeEntity, occurredAt time.Time) map[string]any {
    connectionPrimitives := connectionChangeEntity.GetConnectionPrimitives()
    attributes := b.getBaseAttributes(remoteMonitorEntity)
    maps.Copy(attributes, connectionPrimitives)
    attributes[enums.NetConnectionAttributeEnum.ACTION] = connectionChangeEntity.GetAction()
    if b.isClose(connectionChangeEntity.GetAction()) {
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
    return b.getEvent(remoteMonitorEntity.GetMonitorName(), remoteMonitorEntity.GetMonitorType()+" not available: "+reason, attributes, occurredAt)
}

func (b *ConnectionEventBuilder) getMessage(connectionChangeEntity *entities.ConnectionChangeEntity) string {
    connectionPrimitives := connectionChangeEntity.GetConnectionPrimitives()
    action := connectionChangeEntity.GetAction()
    localEndpoint := b.getEndpoint(connectionPrimitives, enums.ConnectionKeyEnum.LOCAL_IP, enums.ConnectionKeyEnum.LOCAL_PORT)
    messageParts := []string{action, application.GetTrimmedString(connectionPrimitives, enums.ConnectionKeyEnum.PROTOCOL)}
    if action == enums.ConnectionActionEnum.LISTEN_OPENED || action == enums.ConnectionActionEnum.LISTEN_CLOSED {
        messageParts = append(messageParts, b.getProcessLabel(connectionPrimitives), localEndpoint)
    } else {
        messageParts = append(
            messageParts,
            application.GetTrimmedString(connectionPrimitives, enums.ConnectionKeyEnum.DIRECTION),
            b.getProcessLabel(connectionPrimitives),
            localEndpoint,
            "->",
            b.getEndpoint(connectionPrimitives, enums.ConnectionKeyEnum.REMOTE_IP, enums.ConnectionKeyEnum.REMOTE_PORT),
        )
    }
    if b.isClose(action) {
        messageParts = append(messageParts, "after", strconv.Itoa(connectionChangeEntity.GetDurationSeconds()), "s")
    }
    if connectionChangeEntity.IsInitial() {
        messageParts = append(messageParts, "(already open)")
    }
    return strings.Join(messageParts, " ")
}

// getProcessLabel "firefox(4321)" o "?" si no se sabe de quien es.
func (b *ConnectionEventBuilder) getProcessLabel(connectionPrimitives map[string]any) string {
    processName := application.GetTrimmedString(connectionPrimitives, enums.ConnectionKeyEnum.PROCESS)
    pid := application.GetInt(connectionPrimitives, enums.ConnectionKeyEnum.PID)
    if processName == "" || pid == 0 {
        return enums.UNKNOWN_PROCESS
    }
    return processName + "(" + strconv.Itoa(pid) + ")"
}

// getEndpoint "192.168.1.5:52780" / "[::1]:631"; "*" sin IP (el remoto de una escucha).
func (b *ConnectionEventBuilder) getEndpoint(connectionPrimitives map[string]any, ipKey string, portKey string) string {
    ip := application.GetTrimmedString(connectionPrimitives, ipKey)
    if ip == "" {
        return anyRemoteEndpoint
    }
    return net.JoinHostPort(ip, strconv.Itoa(application.GetInt(connectionPrimitives, portKey)))
}

func (b *ConnectionEventBuilder) isClose(action string) bool {
    return action == enums.ConnectionActionEnum.CLOSED || action == enums.ConnectionActionEnum.LISTEN_CLOSED
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
        collectorsenums.CollectorsKeyEnum.OCCURRED_AT: occurredAt.UTC().Format(dateTimeFormat),
        collectorsenums.CollectorsKeyEnum.ATTRIBUTES:  attributes,
    }
}
