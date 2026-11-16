package entities

import (
    "strconv"
    "strings"
)

// ConnectionEntity un socket con su proceso, su usuario y su direccion (listen / inbound / outbound).
type ConnectionEntity struct {
    socketEntity  *SocketEntity
    processEntity *ProcessEntity
    userName      string
    direction     string
}

func NewConnectionEntity(socketEntity *SocketEntity, processEntity *ProcessEntity, userName string, direction string) *ConnectionEntity {
    return &ConnectionEntity{
        socketEntity:  socketEntity,
        processEntity: processEntity,
        userName:      userName,
        direction:     direction,
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
    return e.socketEntity.GetKey() + "|" + strconv.Itoa(e.processEntity.GetPid())
}

// GetProcessLabel "firefox(4321)" o "?" si no se sabe.
func (e *ConnectionEntity) GetProcessLabel() string {
    if e.processEntity.IsUnknown() {
        return e.processEntity.GetName()
    }
    return e.processEntity.GetName() + "(" + strconv.Itoa(e.processEntity.GetPid()) + ")"
}

// GetSearchableText lo que mira el filtro de texto de connections-list, en minusculas.
func (e *ConnectionEntity) GetSearchableText() string {
    return strings.ToLower(strings.Join([]string{
        e.socketEntity.GetProtocol(),
        e.socketEntity.GetState(),
        e.direction,
        e.socketEntity.GetLocalEndpoint(),
        e.socketEntity.GetRemoteEndpoint(),
        e.GetProcessLabel(),
        e.processEntity.GetExe(),
        e.userName,
    }, " "))
}

// ToPrimitives lo que imprime connections-list --format=json y lo que viaja en los attributes del evento.
func (e *ConnectionEntity) ToPrimitives() map[string]any {
    return map[string]any{
        "protocol":    e.socketEntity.GetProtocol(),
        "ip_version":  e.socketEntity.GetIpVersion(),
        "direction":   e.direction,
        "state":       e.socketEntity.GetState(),
        "local_ip":    e.socketEntity.GetLocalIp(),
        "local_port":  e.socketEntity.GetLocalPort(),
        "remote_ip":   e.getRemoteIpOrEmpty(),
        "remote_port": e.socketEntity.GetRemotePort(),
        "pid":         e.processEntity.GetPid(),
        "process":     e.processEntity.GetName(),
        "exe":         e.processEntity.GetExe(),
        "user":        e.userName,
    }
}

func (e *ConnectionEntity) getRemoteIpOrEmpty() string {
    if e.socketEntity.IsListening() {
        return ""
    }
    return e.socketEntity.GetRemoteIp()
}
