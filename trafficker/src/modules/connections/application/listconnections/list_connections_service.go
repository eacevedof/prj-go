package listconnections

import (
    "path/filepath"

    "trafficker/src/core/config"
    "trafficker/src/modules/connections/domain/enums"
    "trafficker/src/modules/connections/domain/exceptions"
    "trafficker/src/modules/connections/domain/services"
    "trafficker/src/modules/connections/infrastructure/repositories"
    "trafficker/src/modules/shared/infrastructure/repositories/configuration"
)

// ListConnectionsService la foto de ahora: cada socket con su proceso, usuario y direccion, con los filtros del .env
// (loopback, procesos y puertos ignorados) y los terminos del filtro. No envia nada ni toca el estado.
type ListConnectionsService struct {
    environmentReaderRawRepository *configuration.EnvironmentReaderRawRepository
    socketsReaderProcRepository    *repositories.SocketsReaderProcRepository
    processesReaderProcRepository  *repositories.ProcessesReaderProcRepository
    usersReaderPasswdRepository    *repositories.UsersReaderPasswdRepository
    connectionEntitiesBuilder      *services.ConnectionEntitiesBuilder
    connectionsFilter              *services.ConnectionsFilter
}

func ListConnectionsServiceGetInstance() *ListConnectionsService {
    return &ListConnectionsService{
        environmentReaderRawRepository: configuration.GetInstance(),
        socketsReaderProcRepository:    repositories.SocketsReaderProcRepositoryGetInstance(),
        processesReaderProcRepository:  repositories.ProcessesReaderProcRepositoryGetInstance(),
        usersReaderPasswdRepository:    repositories.UsersReaderPasswdRepositoryGetInstance(),
        connectionEntitiesBuilder:      services.ConnectionEntitiesBuilderGetInstance(),
        connectionsFilter:              services.ConnectionsFilterGetInstance(),
    }
}

// Invoke error: *ConnectionsException (503: no se puede leer /proc/net, p. ej. fuera de Linux).
func (s *ListConnectionsService) Invoke(listConnectionsDto *ListConnectionsDto) (*ListConnectionsResultDto, error) {
    procPath := s.environmentReaderRawRepository.GetProcPath()
    socketEntities, err := s.socketsReaderProcRepository.GetAllSocketEntities(procPath)
    if err != nil {
        return nil, exceptions.ServiceUnavailableCustom(
            config.Tr(enums.ConnectionsMessageEnum.PROC_NET_NOT_READABLE, filepath.Join(procPath, "net"), err.Error()),
        )
    }
    processEntityByInode, unreadableProcesses, err := s.processesReaderProcRepository.GetProcessEntitiesByInode(procPath)
    if err != nil {
        return nil, exceptions.ServiceUnavailableCustom(
            config.Tr(enums.ConnectionsMessageEnum.PROC_NET_NOT_READABLE, procPath, err.Error()),
        )
    }

    connectionEntities := s.connectionEntitiesBuilder.GetConnectionEntities(
        socketEntities,
        processEntityByInode,
        s.usersReaderPasswdRepository.GetUserNameByUid(repositories.PASSWD_PATH),
    )
    filteredConnectionEntities := s.connectionsFilter.GetFilteredConnectionEntities(
        connectionEntities,
        s.environmentReaderRawRepository.IsLoopbackIncluded(),
        s.environmentReaderRawRepository.GetIgnoredProcesses(),
        s.environmentReaderRawRepository.GetIgnoredPorts(),
        s.connectionsFilter.GetFilterTerms(listConnectionsDto.GetFilter()),
    )
    return NewListConnectionsResultDto(filteredConnectionEntities, unreadableProcesses), nil
}
