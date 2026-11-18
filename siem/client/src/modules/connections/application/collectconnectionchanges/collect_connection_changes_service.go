package collectconnectionchanges

import (
    "time"

    "siem-client/src/core/config"

    "siem-client/src/modules/shared/application"
    "siem-client/src/modules/shared/infrastructure/components/uuid"
    "siem-client/src/modules/shared/infrastructure/repositories/spool"
    "siem-client/src/modules/shared/infrastructure/repositories/state"

    agentconfigenums "siem-client/src/modules/agentconfig/domain/enums"
    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"
    agentconfigtraits "siem-client/src/modules/agentconfig/infrastructure/traits"

    collectorsservices "siem-client/src/modules/collectors/domain/services"

    "siem-client/src/modules/connections/domain/enums"
    "siem-client/src/modules/connections/domain/entities"
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
    stateReaderJsonRepository     *state.StateReaderJsonRepository
    stateWriterJsonRepository     *state.StateWriterJsonRepository
    spoolWriterDiskRepository     *spool.SpoolWriterDiskRepository
    eventsBatcher                 *collectorsservices.EventsBatcher
    socketsReaderProcRepository   *repositories.SocketsReaderProcRepository
    processesReaderProcRepository *repositories.ProcessesReaderProcRepository
    usersReaderPasswdRepository   *repositories.UsersReaderPasswdRepository
    connectionEntitiesBuilder     *services.ConnectionEntitiesBuilder
    connectionsFilter             *services.ConnectionsFilter
    connectionChangesDetector     *services.ConnectionChangesDetector
    connectionEventBuilder        *services.ConnectionEventBuilder

    collectConnectionChangesDto *CollectConnectionChangesDto
    agentConfigEntity           *agentconfigentities.AgentConfigEntity
    connectionsMonitorEntity    *agentconfigentities.RemoteMonitorEntity
    // state/connections.json (vacio si no hay, o si es de otro monitor)
    connectionsState map[string]any
    now              time.Time
    // la foto de ahora: todas (para saber cuales siguen existiendo) y las que pasan el filtro del monitor
    allConnectionEntities      []*entities.ConnectionEntity
    filteredConnectionEntities []*entities.ConnectionEntity
    unreadableProcesses        int
}

func CollectConnectionChangesServiceGetInstance() *CollectConnectionChangesService {
    return &CollectConnectionChangesService{
        AgentConfigLoaderTrait:        agentconfigtraits.NewAgentConfigLoaderTrait(),
        uuidGenerator:                 uuid.GetInstance(),
        stateReaderJsonRepository:     state.StateReaderJsonRepositoryGetInstance(),
        stateWriterJsonRepository:     state.StateWriterJsonRepositoryGetInstance(),
        spoolWriterDiskRepository:     spool.SpoolWriterDiskRepositoryGetInstance(),
        eventsBatcher:                 collectorsservices.EventsBatcherGetInstance(),
        socketsReaderProcRepository:   repositories.SocketsReaderProcRepositoryGetInstance(),
        processesReaderProcRepository: repositories.ProcessesReaderProcRepositoryGetInstance(),
        usersReaderPasswdRepository:   repositories.UsersReaderPasswdRepositoryGetInstance(),
        connectionEntitiesBuilder:     services.ConnectionEntitiesBuilderGetInstance(),
        connectionsFilter:             services.ConnectionsFilterGetInstance(),
        connectionChangesDetector:     services.ConnectionChangesDetectorGetInstance(),
        connectionEventBuilder:        services.ConnectionEventBuilderGetInstance(),
    }
}

// Invoke error: *AgentConfigException (400, 404) o uno inesperado (disco). No poder leer los sockets no es un error
// del servicio: es un evento is_unavailable (una vez).
func (s *CollectConnectionChangesService) Invoke(collectConnectionChangesDto *CollectConnectionChangesDto) (*CollectConnectionChangesResultDto, error) {
    s.collectConnectionChangesDto = collectConnectionChangesDto
    agentConfigEntity, err := s.GetValidAgentConfigEntityOrFail()
    if err != nil {
        return nil, err
    }
    s.agentConfigEntity = agentConfigEntity
    s.connectionsMonitorEntity = agentConfigEntity.GetConnectionsMonitorEntityOrNil()
    if s.connectionsMonitorEntity == nil {
        // el monitor se ha borrado (o nunca hubo): deja de muestrear y olvida la foto
        if err := s.stateWriterJsonRepository.DeleteState(s.GetStatePath(), enums.ConnectionsStateEnum.STATE_FILE); err != nil {
            return nil, err
        }
        return CollectConnectionChangesResultDtoFromPrimitives(map[string]any{"isDisabled": true}), nil
    }

    if err := s.loadConnectionsState(); err != nil {
        return nil, err
    }
    s.now = time.Now().UTC()
    if !s.isDue() {
        return CollectConnectionChangesResultDtoFromPrimitives(map[string]any{}), nil
    }

    if !s.socketsReaderProcRepository.IsSupported() {
        return s.getUnavailableResult(config.Tr(enums.ConnectionsMessageEnum.OS_X_IS_NOT_SUPPORTED, s.socketsReaderProcRepository.GetOsName()))
    }
    if err := s.tryToLoadConnectionEntities(); err != nil {
        return s.getUnavailableResult(config.Tr(enums.ConnectionsMessageEnum.CANNOT_READ_SOCKETS_X, err.Error()))
    }

    // nil = no hay foto anterior: lo que este abierto sale como estado de partida
    previousConnections, _ := s.connectionsState[enums.ConnectionsStateEnum.CONNECTIONS].(map[string]any)
    connectionChangeEntities, nextConnections := s.connectionChangesDetector.GetChangesAndNextConnections(
        previousConnections, s.allConnectionEntities, s.filteredConnectionEntities, s.now.Unix(),
    )
    events := []map[string]any{}
    opened, closed := 0, 0
    for _, connectionChangeEntity := range connectionChangeEntities {
        events = append(events, s.connectionEventBuilder.GetChangeEvent(s.connectionsMonitorEntity, connectionChangeEntity, s.now))
        if connectionChangeEntity.IsClose() {
            closed++
            continue
        }
        opened++
    }

    spooledBatches, err := s.saveSpoolBatches(events)
    if err != nil {
        return nil, err
    }
    err = s.stateWriterJsonRepository.SaveState(s.GetStatePath(), enums.ConnectionsStateEnum.STATE_FILE, map[string]any{
        enums.ConnectionsStateEnum.MONITOR_UUID:   s.connectionsMonitorEntity.GetMonitorUuid(),
        enums.ConnectionsStateEnum.LAST_RUN_UNIX:  s.now.Unix(),
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
        "unreadableProcesses": s.unreadableProcesses,
        "spooledBatches":      spooledBatches,
        "results":             events,
    }), nil
}

// loadConnectionsState la foto de otro monitor (se borro y se creo otro) no vale: se empieza de cero.
func (s *CollectConnectionChangesService) loadConnectionsState() error {
    connectionsState, err := s.stateReaderJsonRepository.GetState(s.GetStatePath(), enums.ConnectionsStateEnum.STATE_FILE)
    if err != nil {
        return err
    }
    if application.GetTrimmedString(connectionsState, enums.ConnectionsStateEnum.MONITOR_UUID) != s.connectionsMonitorEntity.GetMonitorUuid() {
        connectionsState = map[string]any{}
    }
    s.connectionsState = connectionsState
    return nil
}

// isDue le toca por interval_seconds, es la primera pasada o se ha forzado (connections-collect --force=1).
func (s *CollectConnectionChangesService) isDue() bool {
    lastRunUnix := int64(application.GetInt(s.connectionsState, enums.ConnectionsStateEnum.LAST_RUN_UNIX))
    return s.collectConnectionChangesDto.IsForced() || lastRunUnix == 0 ||
        s.now.Unix()-lastRunUnix >= int64(s.connectionsMonitorEntity.GetIntervalSeconds())
}

// tryToLoadConnectionEntities la foto de ahora, con el filtro de la config del monitor. Un error = no se puede
// muestrear en este sistema (quien llama lo convierte en el aviso is_unavailable).
func (s *CollectConnectionChangesService) tryToLoadConnectionEntities() error {
    socketEntities, err := s.socketsReaderProcRepository.GetAllSocketEntities(enums.ProcPathEnum.PROC_PATH)
    if err != nil {
        return err
    }
    processEntityByInode, unreadableProcesses, err := s.processesReaderProcRepository.GetProcessEntitiesByInode(enums.ProcPathEnum.PROC_PATH)
    if err != nil {
        return err
    }
    userNameByUid, err := s.usersReaderPasswdRepository.GetUserNameByUid(enums.ProcPathEnum.PASSWD_PATH)
    if err != nil {
        return err
    }

    s.allConnectionEntities = s.connectionEntitiesBuilder.GetConnectionEntities(socketEntities, processEntityByInode, userNameByUid)
    s.filteredConnectionEntities = s.connectionsFilter.GetFilteredConnectionEntities(
        s.allConnectionEntities,
        s.connectionsMonitorEntity.GetConfigBool(agentconfigenums.RemoteConfigKeyEnum.IS_LOOPBACK_INCLUDED),
        s.connectionsMonitorEntity.GetConfigStrings(agentconfigenums.RemoteConfigKeyEnum.IGNORE_PROCESSES),
        s.connectionsMonitorEntity.GetConfigInts(agentconfigenums.RemoteConfigKeyEnum.IGNORE_PORTS),
    )
    s.unreadableProcesses = unreadableProcesses
    return nil
}

// getUnavailableResult el aviso va una sola vez por monitor: el estado recuerda que ya se mando. La foto anterior se
// descarta (si vuelve a poder leer, empieza de cero).
func (s *CollectConnectionChangesService) getUnavailableResult(unavailableReason string) (*CollectConnectionChangesResultDto, error) {
    isAlreadyWarned := s.connectionsState[enums.ConnectionsStateEnum.IS_UNAVAILABLE] == true
    events := []map[string]any{}
    if !isAlreadyWarned {
        events = append(events, s.connectionEventBuilder.GetUnavailableEvent(s.connectionsMonitorEntity, unavailableReason, s.now))
    }
    spooledBatches, err := s.saveSpoolBatches(events)
    if err != nil {
        return nil, err
    }
    err = s.stateWriterJsonRepository.SaveState(s.GetStatePath(), enums.ConnectionsStateEnum.STATE_FILE, map[string]any{
        enums.ConnectionsStateEnum.MONITOR_UUID:   s.connectionsMonitorEntity.GetMonitorUuid(),
        enums.ConnectionsStateEnum.LAST_RUN_UNIX:  s.now.Unix(),
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

// saveSpoolBatches los eventos al spool en lotes de max_events_per_batch; devuelve cuantos lotes.
func (s *CollectConnectionChangesService) saveSpoolBatches(events []map[string]any) (int, error) {
    spooledBatches := 0
    for _, batchEvents := range s.eventsBatcher.GetBatches(events, s.agentConfigEntity.GetMaxEventsPerBatch()) {
        if err := s.spoolWriterDiskRepository.SaveSpoolBatch(s.GetSpoolPath(), s.uuidGenerator.GetUniqueIdWithPrefix("b-"), batchEvents); err != nil {
            return 0, err
        }
        spooledBatches++
    }
    return spooledBatches, nil
}
