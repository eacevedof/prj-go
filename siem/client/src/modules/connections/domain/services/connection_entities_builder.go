package services

import (
    "cmp"
    "slices"
    "strconv"

    "siem-client/src/modules/connections/domain/entities"
    "siem-client/src/modules/connections/domain/enums"
)

// ConnectionEntitiesBuilder junta la foto de sockets con sus procesos y usuarios y decide la direccion de cada uno:
//   - listen: socket en escucha.
//   - inbound: su puerto local es uno que este host escucha con el mismo protocolo (alguien se ha conectado a el).
//   - outbound: el resto (un proceso de este host se ha conectado fuera, desde un puerto efimero).
//
// Fuera los sockets con inodo 0: son conexiones ya cerradas que el kernel aun guarda (TIME_WAIT) y no son de nadie.
type ConnectionEntitiesBuilder struct{}

func ConnectionEntitiesBuilderGetInstance() *ConnectionEntitiesBuilder {
    return &ConnectionEntitiesBuilder{}
}

// GetConnectionEntities ordenadas: escuchas, entrantes y salientes; dentro, por proceso y por extremo remoto.
func (b *ConnectionEntitiesBuilder) GetConnectionEntities(
    socketEntities []*entities.SocketEntity,
    processEntityByInode map[uint64]*entities.ProcessEntity,
    userNameByUid map[int]string,
) []*entities.ConnectionEntity {
    listeningPorts := b.getListeningPorts(socketEntities)
    connectionEntities := []*entities.ConnectionEntity{}
    for _, socketEntity := range socketEntities {
        if socketEntity.GetInode() == 0 {
            continue
        }
        processEntity := processEntityByInode[socketEntity.GetInode()]
        if processEntity == nil {
            processEntity = entities.UnknownProcessEntity()
        }
        connectionEntities = append(connectionEntities, entities.NewConnectionEntity(
            socketEntity,
            processEntity,
            b.getUserName(userNameByUid, socketEntity.GetUid()),
            b.getDirection(socketEntity, listeningPorts),
        ))
    }
    slices.SortStableFunc(connectionEntities, b.compare)
    return connectionEntities
}

// getListeningPorts "tcp|443" -> true. Un puerto escuchado en IPv4 o IPv6 cuenta para los dos.
func (b *ConnectionEntitiesBuilder) getListeningPorts(socketEntities []*entities.SocketEntity) map[string]bool {
    listeningPorts := map[string]bool{}
    for _, socketEntity := range socketEntities {
        if socketEntity.IsListening() {
            listeningPorts[socketEntity.GetProtocol()+"|"+strconv.Itoa(socketEntity.GetLocalPort())] = true
        }
    }
    return listeningPorts
}

func (b *ConnectionEntitiesBuilder) getDirection(socketEntity *entities.SocketEntity, listeningPorts map[string]bool) string {
    if socketEntity.IsListening() {
        return enums.DirectionEnum.LISTEN
    }
    if listeningPorts[socketEntity.GetProtocol()+"|"+strconv.Itoa(socketEntity.GetLocalPort())] {
        return enums.DirectionEnum.INBOUND
    }
    return enums.DirectionEnum.OUTBOUND
}

// getUserName el nombre de /etc/passwd; si no esta, el uid tal cual.
func (b *ConnectionEntitiesBuilder) getUserName(userNameByUid map[int]string, uid int) string {
    if userName, isFound := userNameByUid[uid]; isFound {
        return userName
    }
    return strconv.Itoa(uid)
}

var directionOrder = map[string]int{
    enums.DirectionEnum.LISTEN:   0,
    enums.DirectionEnum.INBOUND:  1,
    enums.DirectionEnum.OUTBOUND: 2,
}

func (b *ConnectionEntitiesBuilder) compare(first *entities.ConnectionEntity, second *entities.ConnectionEntity) int {
    return cmp.Or(
        cmp.Compare(directionOrder[first.GetDirection()], directionOrder[second.GetDirection()]),
        cmp.Compare(first.GetProcessEntity().GetName(), second.GetProcessEntity().GetName()),
        first.GetSocketEntity().GetRemoteAddr().Compare(second.GetSocketEntity().GetRemoteAddr()),
        cmp.Compare(first.GetSocketEntity().GetRemotePort(), second.GetSocketEntity().GetRemotePort()),
        cmp.Compare(first.GetSocketEntity().GetLocalPort(), second.GetSocketEntity().GetLocalPort()),
    )
}
