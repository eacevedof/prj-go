package entities

import (
    "strconv"
    "strings"

    "siem-client/src/core/config"

    "siem-client/src/modules/connections/domain/enums"
)

// ConnectionEntity un socket con su proceso, su usuario y su direccion (listen / inbound / outbound).
type ConnectionEntity struct {
    socketEntity  *SocketEntity
    processEntity *ProcessEntity
    userName      string
    direction     string
}

// ConnectionEntityFromSocketAndProcess la une ConnectionEntitiesBuilder, que es quien decide la direccion.
func ConnectionEntityFromSocketAndProcess(socketEntity *SocketEntity, processEntity *ProcessEntity, userName string, direction string) *ConnectionEntity {
    return &ConnectionEntity{
        socketEntity:  socketEntity,
        processEntity: processEntity,
        userName:      strings.TrimSpace(userName),
        direction:     strings.TrimSpace(direction),
    }
}

func (e *ConnectionEntity) GetSocketEntity() *SocketEntity {
    return e.socketEntity
}

func (e *ConnectionEntity) GetProcessEntity() *ProcessEntity {
    return e.processEntity
}

func (e *ConnectionEntity) GetUserName() string {
    return e.userName
}

func (e *ConnectionEntity) GetDirection() string {
    return e.direction
}

// GetKey identifica la conexion entre dos fotos: la del socket mas el pid (un socket heredado por otro proceso cuenta
// como otra conexion).
func (e *ConnectionEntity) GetKey() string {
    return e.socketEntity.GetKey() + enums.ConnectionFormatEnum.KEY_SEPARATOR + strconv.Itoa(e.processEntity.GetPid())
}

// GetProcessLabel "firefox(4321)" o "?" si no se sabe.
func (e *ConnectionEntity) GetProcessLabel() string {
    if e.processEntity.IsUnknown() {
        return e.processEntity.GetName()
    }
    return config.Tr(enums.ConnectionFormatEnum.PROCESS_LABEL, e.processEntity.GetName(), e.processEntity.GetPid())
}

// IsListening un socket en escucha (sus eventos son listen_opened / listen_closed).
func (e *ConnectionEntity) IsListening() bool {
    return e.socketEntity.IsListening()
}

// ToPrimitives lo que viaja en los attributes del evento net_connection y lo que se guarda en state/connections.json
// (con opened_at_unix) para el diff de la pasada siguiente.
func (e *ConnectionEntity) ToPrimitives() map[string]any {
    return map[string]any{
        enums.ConnectionKeyEnum.PROTOCOL:    e.socketEntity.GetProtocol(),
        enums.ConnectionKeyEnum.IP_VERSION:  e.socketEntity.GetIpVersion(),
        enums.ConnectionKeyEnum.DIRECTION:   e.direction,
        enums.ConnectionKeyEnum.STATE:       e.socketEntity.GetState(),
        enums.ConnectionKeyEnum.LOCAL_IP:    e.socketEntity.GetLocalIp(),
        enums.ConnectionKeyEnum.LOCAL_PORT:  e.socketEntity.GetLocalPort(),
        enums.ConnectionKeyEnum.REMOTE_IP:   e.getRemoteIpOrEmpty(),
        enums.ConnectionKeyEnum.REMOTE_PORT: e.socketEntity.GetRemotePort(),
        enums.ConnectionKeyEnum.PID:         e.processEntity.GetPid(),
        enums.ConnectionKeyEnum.PROCESS:     e.processEntity.GetName(),
        enums.ConnectionKeyEnum.EXE:         e.processEntity.GetExe(),
        enums.ConnectionKeyEnum.USER:        e.userName,
    }
}

func (e *ConnectionEntity) getRemoteIpOrEmpty() string {
    if e.socketEntity.IsListening() {
        return ""
    }
    return e.socketEntity.GetRemoteIp()
}
