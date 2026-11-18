package collectconnectionchanges

import (
    "time"

    "siem-client/src/core/config"

    "siem-client/src/modules/shared/application"
    "siem-client/src/modules/shared/infrastructure/components/uuid"
    "siem-client/src/modules/shared/infrastructure/repositories/spool"
    "siem-client/src/modules/shared/infrastructure/repositories/state"

    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"
    agentconfigenums "siem-client/src/modules/agentconfig/domain/enums"
    agentconfigtraits "siem-client/src/modules/agentconfig/infrastructure/traits"

    collectorsservices "siem-client/src/modules/collectors/domain/services"

    "siem-client/src/modules/connections/domain/entities"
    "siem-client/src/modules/connections/domain/enums"
    "siem-client/src/modules/connections/domain/services"
    "siem-client/src/modules/connections/infrastructure/repositories"
)

// CollectConnectionChangesService una pasada del monitor connections que baja del panel del backend (uno por agente).
// Si le toca (interval_seconds): foto de los sockets con su proceso, filtro de la config del monitor (loopback,
// procesos y puertos ignorados), diff contra la foto anterior (state/connections.json) y SOLO los cambios al spool
// como net_connection. Sin monitor en la config no muestrea y borra el estado (si vuelve, empieza de cero). Si en este
// sistema no se puede muestrear (no es Linux, /proc no se puede leer) lo avisa UNA vez con un evento is_unavailable.
type CollectConnectionChangesService struct {
    agentconfigtraits.AgentConfigLoaderTrait

    uuidGenerator                 *uuid.UuidGenerator
    eventsBatcher                 *collectorsservices.EventsBatcher
    connectionEntitiesBuilder     *services.ConnectionEntitiesBuilder
    connectionsFilter             *services.ConnectionsFilter
    connectionChangesDetector     *services.ConnectionChangesDetector
    connectionEventBuilder        *services.ConnectionEventBuilder
    socketsReaderProcRepository   *repositories.SocketsReaderProcRepository
    processesReaderProcRepository *repositories.ProcessesReaderProcRepository
    usersReaderPasswdRepository   *repositories.UsersReaderPasswdRepository
    stateReaderJsonRepository     *state.StateReaderJsonRepository
    stateWriterJsonRepository     *state.StateWriterJsonRepository
    spoolWriterDiskRepository     *spool.SpoolWriterDiskRepository
}

func CollectConnectionChangesServiceGetInstance() *CollectConnectionChangesService {
    return &CollectConnectionChangesService{
        AgentConfigLoaderTrait:        agentconfigtraits.NewAgentConfigLoaderTrait(),
        uuidGenerator:                 uuid.GetInstance(),
        eventsBatcher:                 collectorsservices.EventsBatcherGetInstance(),
        connectionEntitiesBuilder:     services.ConnectionEntitiesBuilderGetInstance(),
        connectionsFilter:             services.ConnectionsFilterGetInstance(),
        connectionChangesDetector:     services.ConnectionChangesDetectorGetInstance(),
        connectionEventBuilder:        services.ConnectionEventBuilderGetInstance(),
        socketsReaderProcRepository:   repositories.SocketsReaderProcRepositoryGetInstance(),
        processesReaderProcRepository: repositories.ProcessesReaderProcRepositoryGetInstance(),
        usersReaderPasswdRepository:   repositories.UsersReaderPasswdRepositoryGetInstance(),
        stateReaderJsonRepository:     state.StateReaderJsonRepositoryGetInstance(),
        stateWriterJsonRepository:     state.StateWriterJsonRepositoryGetInstance(),
        spoolWriterDiskRepository:     spool.SpoolWriterDiskRepositoryGetInstance(),
    }
}

// Invoke error: *AgentConfigException (400, 404) o uno inesperado (disco). No poder leer los sockets no es un error
// del servicio: es un evento is_unavailable (una vez).
func (s *CollectConnectionChangesService) Invoke(collectConnectionChangesDto *CollectConnectionChangesDto) (*CollectConnectionChangesResultDto, error) {
    agentConfigEntity, err := s.GetValidAgentConfigEntityOrFail()
    if err != nil {
        return nil, err
    }
    connectionsMonitorEntity := agentConfigEntity.GetConnectionsMonitorEntityOrNil()
    if connectionsMonitorEntity == nil {
        // el monitor se ha borrado (o nunca hubo): deja de muestrear y olvida la foto
        if err := s.stateWriterJsonRepository.DeleteState(s.GetStatePath(), enums.ConnectionsStateEnum.STATE_FILE); err != nil {
            return nil, err
        }
        return CollectConnectionChangesResultDtoFromPrimitives(map[string]any{"isDisabled": true}), nil
    }

    connectionsState, err := s.stateReaderJsonRepository.GetState(s.GetStatePath(), enums.ConnectionsStateEnum.STATE_FILE)
    if err != nil {
        return nil, err
    }
    // la foto de otro monitor (se borro y se creo otro) no vale: se empieza de cero
    if application.GetTrimmedString(connectionsState, enums.ConnectionsStateEnum.MONITOR_UUID) != connectionsMonitorEntity.GetMonitorUuid() {
        connectionsState = map[string]any{}
    }

    now := time.Now().UTC()
    lastRunUnix := int64(application.GetInt(connectionsState, enums.ConnectionsStateEnum.LAST_RUN_UNIX))
    isDue := collectConnectionChangesDto.IsForced() || lastRunUnix == 0 || now.Unix()-lastRunUnix >= int64(connectionsMonitorEntity.GetIntervalSeconds())
    if !isDue {
        return CollectConnectionChangesResultDtoFromPrimitives(map[string]any{}), nil
    }

    snapshot, unavailableReason := s.getSnapshot(connectionsMonitorEntity)
    if unavailableReason != "" {
        return s.getUnavailableResult(agentConfigEntity, connectionsMonitorEntity, connectionsState, unavailableReason, now)
    }

    // nil = no hay foto anterior: lo que este abierto sale como estado de partida
    previousConnections, _ := connectionsState[enums.ConnectionsStateEnum.CONNECTIONS].(map[string]any)
    connectionChangeEntities, nextConnections := s.connectionChangesDetector.GetChangesAndNextConnections(
        previousConnections, snapshot.allConnectionEntities, snapshot.filteredConnectionEntities, now.Unix(),
    )
    events := []map[string]any{}
    opened, closed := 0, 0
    for _, connectionChangeEntity := range connectionChangeEntities {
        events = append(events, s.connectionEventBuilder.GetChangeEvent(connectionsMonitorEntity, connectionChangeEntity, now))
        switch connectionChangeEntity.GetAction() {
        case enums.ConnectionActionEnum.OPENED, enums.ConnectionActionEnum.LISTEN_OPENED:
            opened++
        default:
            closed++
        }
    }

    spooledBatches, err := s.saveSpoolBatches(agentConfigEntity, events)
    if err != nil {
        return nil, err
    }
    err = s.stateWriterJsonRepository.SaveState(s.GetStatePath(), enums.ConnectionsStateEnum.STATE_FILE, map[string]any{
        enums.ConnectionsStateEnum.MONITOR_UUID:   connectionsMonitorEntity.GetMonitorUuid(),
        enums.ConnectionsStateEnum.LAST_RUN_UNIX:  now.Unix(),
        enums.ConnectionsStateEnum.IS_UNAVAILABLE: false,
        enums.ConnectionsStateEnum.CONNECTIONS:    nextConnections,
    })
    if err != nil {
        return nil, err
    }

    return CollectConnectionChangesResultDtoFromPrimitives(map[string]any{
        "isSampled":           true,
        "connections":         len(nextConnections),
        "opened":              opened,
        "closed":              closed,
        "unreadableProcesses": snapshot.unreadableProcesses,
        "spooledBatches":      spooledBatches,
        "results":             events,
    }), nil
}

// connectionsSnapshot la foto de ahora: todas las conexiones (para saber cuales siguen existiendo) y las que pasan el
// filtro del monitor (las que se vigilan).
type connectionsSnapshot struct {
    allConnectionEntities      []*entities.ConnectionEntity
    filteredConnectionEntities []*entities.ConnectionEntity
    unreadableProcesses        int
}

// getSnapshot la foto, o el motivo por el que no se puede tomar en este sistema.
func (s *CollectConnectionChangesService) getSnapshot(connectionsMonitorEntity *agentconfigentities.RemoteMonitorEntity) (*connectionsSnapshot, string) {
    if !s.socketsReaderProcRepository.IsSupported() {
        return nil, config.Tr(enums.ConnectionsMessageEnum.OS_X_IS_NOT_SUPPORTED, s.socketsReaderProcRepository.GetOsName())
    }
    socketEntities, err := s.socketsReaderProcRepository.GetAllSocketEntities(repositories.PROC_PATH)
    if err != nil {
        return nil, config.Tr(enums.ConnectionsMessageEnum.CANNOT_READ_SOCKETS_X, err.Error())
    }
    processEntityByInode, unreadableProcesses, err := s.processesReaderProcRepository.GetProcessEntitiesByInode(repositories.PROC_PATH)
    if err != nil {
        return nil, config.Tr(enums.ConnectionsMessageEnum.CANNOT_READ_SOCKETS_X, err.Error())
    }

    allConnectionEntities := s.connectionEntitiesBuilder.GetConnectionEntities(
        socketEntities,
        processEntityByInode,
        s.usersReaderPasswdRepository.GetUserNameByUid(repositories.PASSWD_PATH),
    )
    filteredConnectionEntities := s.connectionsFilter.GetFilteredConnectionEntities(
        allConnectionEntities,
        connectionsMonitorEntity.GetConfigBool(agentconfigenums.RemoteConfigKeyEnum.IS_LOOPBACK_INCLUDED),
        connectionsMonitorEntity.GetConfigStrings(agentconfigenums.RemoteConfigKeyEnum.IGNORE_PROCESSES),
        connectionsMonitorEntity.GetConfigInts(agentconfigenums.RemoteConfigKeyEnum.IGNORE_PORTS),
    )
    return &connectionsSnapshot{
        allConnectionEntities:      allConnectionEntities,
        filteredConnectionEntities: filteredConnectionEntities,
        unreadableProcesses:        unreadableProcesses,
    }, ""
}

// getUnavailableResult el aviso va una sola vez por monitor: el estado recuerda que ya se mando. La foto anterior se
// descarta (si vuelve a poder leer, empieza de cero).
func (s *CollectConnectionChangesService) getUnavailableResult(
    agentConfigEntity *agentconfigentities.AgentConfigEntity,
    connectionsMonitorEntity *agentconfigentities.RemoteMonitorEntity,
    connectionsState map[string]any,
    unavailableReason string,
    now time.Time,
) (*CollectConnectionChangesResultDto, error) {
    isAlreadyWarned := connectionsState[enums.ConnectionsStateEnum.IS_UNAVAILABLE] == true
    events := []map[string]any{}
    if !isAlreadyWarned {
        events = append(events, s.connectionEventBuilder.GetUnavailableEvent(connectionsMonitorEntity, unavailableReason, now))
    }
    spooledBatches, err := s.saveSpoolBatches(agentConfigEntity, events)
    if err != nil {
        return nil, err
    }
    err = s.stateWriterJsonRepository.SaveState(s.GetStatePath(), enums.ConnectionsStateEnum.STATE_FILE, map[string]any{
        enums.ConnectionsStateEnum.MONITOR_UUID:   connectionsMonitorEntity.GetMonitorUuid(),
        enums.ConnectionsStateEnum.LAST_RUN_UNIX:  now.Unix(),
        enums.ConnectionsStateEnum.IS_UNAVAILABLE: true,
    })
    if err != nil {
        return nil, err
    }
    return CollectConnectionChangesResultDtoFromPrimitives(map[string]any{
        "isUnavailable":  !isAlreadyWarned,
        "spooledBatches": spooledBatches,
        "results":        events,
    }), nil
}

func (s *CollectConnectionChangesService) saveSpoolBatches(agentConfigEntity *agentconfigentities.AgentConfigEntity, events []map[string]any) (int, error) {
    spooledBatches := 0
    for _, batchEvents := range s.eventsBatcher.GetBatches(events, agentConfigEntity.GetMaxEventsPerBatch()) {
        if err := s.spoolWriterDiskRepository.SaveSpoolBatch(s.GetSpoolPath(), s.uuidGenerator.GetUniqueIdWithPrefix("b-"), batchEvents); err != nil {
            return 0, err
        }
        spooledBatches++
    }
    return spooledBatches, nil
}
