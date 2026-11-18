package entities

import (
    "regexp"
    "strings"

    "siem-client/src/modules/agentconfig/domain/enums"
)

const (
    DEFAULT_COLLECT_INTERVAL_SECONDS   = 5
    DEFAULT_SEND_INTERVAL_SECONDS      = 10
    DEFAULT_HEARTBEAT_INTERVAL_SECONDS = 30
    DEFAULT_MAX_EVENTS_PER_BATCH       = 500
    DEFAULT_MAX_SPOOL_BATCHES          = 5000
    DEFAULT_WATCHER_PATTERN            = "*"
    DEFAULT_COMMAND_INTERVAL_SECONDS   = 60
    DEFAULT_COMMAND_TIMEOUT_SECONDS    = 30
)

// AgentConfigPrimitives siem-client.json tal cual (con los defaults aplicados): es el formato del fichero y la
// salida de config-check. Un 0 / "" en el fichero = el valor por defecto.
type AgentConfigPrimitives struct {
    CollectIntervalSeconds   int                         `json:"collect_interval_seconds"`
    SendIntervalSeconds      int                         `json:"send_interval_seconds"`
    HeartbeatIntervalSeconds int                         `json:"heartbeat_interval_seconds"`
    MaxEventsPerBatch        int                         `json:"max_events_per_batch"`
    MaxSpoolBatches          int                         `json:"max_spool_batches"`
    Watchers                 []WatcherPrimitives         `json:"watchers"`
    Commands                 []ScheduledCommandPrimitives `json:"commands"`
    Docker                   DockerConfigPrimitives       `json:"docker"`
}

// AgentConfigEntity que vigila el agente y cada cuanto. Se relee en cada pasada: cambiar carpetas o comandos no
// exige reiniciar (los intervalos de collect/send si).
type AgentConfigEntity struct {
    collectIntervalSeconds   int
    sendIntervalSeconds      int
    heartbeatIntervalSeconds int
    maxEventsPerBatch        int
    maxSpoolBatches          int
    watcherEntities          []*WatcherEntity
    scheduledCommandEntities []*ScheduledCommandEntity
    dockerConfigEntity       *DockerConfigEntity
    // lo que llega del panel del backend (config remota): las metricas del host; folder y log_pattern ya van dentro de
    // watcherEntities. Los avisos son monitores remotos que no se han podido aplicar.
    metricMonitorEntities []*RemoteMonitorEntity
    remoteConfigWarnings  []string
    // el monitor connections del panel (uno por agente); nil = no se muestrean las conexiones de red
    connectionsMonitorEntity *RemoteMonitorEntity
}

// AddRemoteMonitors suma a la config local lo que salio de la config remota (RemoteMonitorsMerger ya resolvio los
// choques de nombre y descarto lo que no vale).
func (e *AgentConfigEntity) AddRemoteMonitors(watcherEntities []*WatcherEntity, metricMonitorEntities []*RemoteMonitorEntity, warnings []string) {
    e.watcherEntities = append(e.watcherEntities, watcherEntities...)
    e.metricMonitorEntities = append(e.metricMonitorEntities, metricMonitorEntities...)
    e.remoteConfigWarnings = append(e.remoteConfigWarnings, warnings...)
}

// GetMetricMonitorEntities los monitores cpu, memory, disk, gpu y network (modulo hostmetrics).
func (e *AgentConfigEntity) GetMetricMonitorEntities() []*RemoteMonitorEntity {
    return e.metricMonitorEntities
}

// SetConnectionsMonitorEntity lo pone RemoteMonitorsMerger (el primero que llega; uno de mas queda como aviso).
func (e *AgentConfigEntity) SetConnectionsMonitorEntity(connectionsMonitorEntity *RemoteMonitorEntity) {
    e.connectionsMonitorEntity = connectionsMonitorEntity
}

// GetConnectionsMonitorEntityOrNil nil si el panel no tiene un monitor connections para este agente.
func (e *AgentConfigEntity) GetConnectionsMonitorEntityOrNil() *RemoteMonitorEntity {
    return e.connectionsMonitorEntity
}

func (e *AgentConfigEntity) GetRemoteConfigWarnings() []string {
    return e.remoteConfigWarnings
}

func AgentConfigEntityFromPrimitives(agentConfigPrimitives AgentConfigPrimitives) *AgentConfigEntity {
    watcherEntities := make([]*WatcherEntity, 0, len(agentConfigPrimitives.Watchers))
    for _, watcherPrimitives := range agentConfigPrimitives.Watchers {
        watcherEntities = append(watcherEntities, WatcherEntityFromPrimitives(watcherPrimitives))
    }
    scheduledCommandEntities := make([]*ScheduledCommandEntity, 0, len(agentConfigPrimitives.Commands))
    for _, scheduledCommandPrimitives := range agentConfigPrimitives.Commands {
        scheduledCommandEntities = append(scheduledCommandEntities, ScheduledCommandEntityFromPrimitives(scheduledCommandPrimitives))
    }
    return &AgentConfigEntity{
        collectIntervalSeconds:   getIntOrDefault(agentConfigPrimitives.CollectIntervalSeconds, DEFAULT_COLLECT_INTERVAL_SECONDS),
        sendIntervalSeconds:      getIntOrDefault(agentConfigPrimitives.SendIntervalSeconds, DEFAULT_SEND_INTERVAL_SECONDS),
        heartbeatIntervalSeconds: getIntOrDefault(agentConfigPrimitives.HeartbeatIntervalSeconds, DEFAULT_HEARTBEAT_INTERVAL_SECONDS),
        maxEventsPerBatch:        getIntOrDefault(agentConfigPrimitives.MaxEventsPerBatch, DEFAULT_MAX_EVENTS_PER_BATCH),
        maxSpoolBatches:          getIntOrDefault(agentConfigPrimitives.MaxSpoolBatches, DEFAULT_MAX_SPOOL_BATCHES),
        watcherEntities:          watcherEntities,
        scheduledCommandEntities: scheduledCommandEntities,
        dockerConfigEntity:       DockerConfigEntityFromPrimitives(agentConfigPrimitives.Docker),
    }
}

func (e *AgentConfigEntity) ToPrimitives() AgentConfigPrimitives {
    watchers := make([]WatcherPrimitives, 0, len(e.watcherEntities))
    for _, watcherEntity := range e.watcherEntities {
        watchers = append(watchers, watcherEntity.ToPrimitives())
    }
    commands := make([]ScheduledCommandPrimitives, 0, len(e.scheduledCommandEntities))
    for _, scheduledCommandEntity := range e.scheduledCommandEntities {
        commands = append(commands, scheduledCommandEntity.ToPrimitives())
    }
    return AgentConfigPrimitives{
        CollectIntervalSeconds:   e.collectIntervalSeconds,
        SendIntervalSeconds:      e.sendIntervalSeconds,
        HeartbeatIntervalSeconds: e.heartbeatIntervalSeconds,
        MaxEventsPerBatch:        e.maxEventsPerBatch,
        MaxSpoolBatches:          e.maxSpoolBatches,
        Watchers:                 watchers,
        Commands:                 commands,
        Docker:                   e.dockerConfigEntity.ToPrimitives(),
    }
}

func (e *AgentConfigEntity) GetCollectIntervalSeconds() int {
    return e.collectIntervalSeconds
}

func (e *AgentConfigEntity) GetSendIntervalSeconds() int {
    return e.sendIntervalSeconds
}

func (e *AgentConfigEntity) GetHeartbeatIntervalSeconds() int {
    return e.heartbeatIntervalSeconds
}

func (e *AgentConfigEntity) GetMaxEventsPerBatch() int {
    return e.maxEventsPerBatch
}

func (e *AgentConfigEntity) GetMaxSpoolBatches() int {
    return e.maxSpoolBatches
}

func (e *AgentConfigEntity) GetWatcherEntities() []*WatcherEntity {
    return e.watcherEntities
}

func (e *AgentConfigEntity) GetScheduledCommandEntities() []*ScheduledCommandEntity {
    return e.scheduledCommandEntities
}

// WatcherPrimitives una carpeta vigilada: los ficheros de path que casan con pattern (glob sobre el nombre).
// match (opcional): regex; solo se envian las lineas que casan. severity (opcional): gravedad que se declara en cada
// linea (si no, la deduce el backend). monitor_uuid / monitor_type: solo en los watchers que salen de un monitor del
// panel (folder, log_pattern); en siem-client.json no se escriben.
type WatcherPrimitives struct {
    Name        string `json:"name"`
    Path        string `json:"path"`
    Pattern     string `json:"pattern"`
    IsRecursive bool   `json:"is_recursive"`
    StartAt     string `json:"start_at"`
    Match       string `json:"match,omitempty"`
    Severity    string `json:"severity,omitempty"`
    MonitorUuid string `json:"monitor_uuid,omitempty"`
    MonitorType string `json:"monitor_type,omitempty"`
}

type WatcherEntity struct {
    name        string
    path        string
    pattern     string
    isRecursive bool
    startAt     string
    match       string
    // nil = sin filtro, o match no es un regex valido (lo avisa el validador / el merger de monitores)
    matchRegexp *regexp.Regexp
    severity    string
    monitorUuid string
    monitorType string
}

func WatcherEntityFromPrimitives(watcherPrimitives WatcherPrimitives) *WatcherEntity {
    match := trim(watcherPrimitives.Match)
    var matchRegexp *regexp.Regexp
    if match != "" {
        matchRegexp, _ = regexp.Compile(match)
    }
    return &WatcherEntity{
        name:        trim(watcherPrimitives.Name),
        path:        trim(watcherPrimitives.Path),
        pattern:     getStringOrDefault(watcherPrimitives.Pattern, DEFAULT_WATCHER_PATTERN),
        isRecursive: watcherPrimitives.IsRecursive,
        startAt:     getStringOrDefault(watcherPrimitives.StartAt, enums.WatcherStartAtEnum.END),
        match:       match,
        matchRegexp: matchRegexp,
        severity:    strings.ToLower(trim(watcherPrimitives.Severity)),
        monitorUuid: trim(watcherPrimitives.MonitorUuid),
        monitorType: trim(watcherPrimitives.MonitorType),
    }
}

func (e *WatcherEntity) ToPrimitives() WatcherPrimitives {
    return WatcherPrimitives{
        Name: e.name, Path: e.path, Pattern: e.pattern, IsRecursive: e.isRecursive, StartAt: e.startAt,
        Match: e.match, Severity: e.severity, MonitorUuid: e.monitorUuid, MonitorType: e.monitorType,
    }
}

func (e *WatcherEntity) GetMatch() string {
    return e.match
}

// IsMatchValid sin match, o un regex que compila.
func (e *WatcherEntity) IsMatchValid() bool {
    return e.match == "" || e.matchRegexp != nil
}

// IsLineMatching sin match pasan todas las lineas.
func (e *WatcherEntity) IsLineMatching(line string) bool {
    return e.matchRegexp == nil || e.matchRegexp.MatchString(line)
}

func (e *WatcherEntity) GetSeverity() string {
    return e.severity
}

func (e *WatcherEntity) GetMonitorUuid() string {
    return e.monitorUuid
}

func (e *WatcherEntity) GetMonitorType() string {
    return e.monitorType
}

func (e *WatcherEntity) IsFromMonitor() bool {
    return e.monitorUuid != ""
}

func (e *WatcherEntity) GetName() string {
    return e.name
}

func (e *WatcherEntity) GetPath() string {
    return e.path
}

func (e *WatcherEntity) GetPattern() string {
    return e.pattern
}

func (e *WatcherEntity) IsRecursive() bool {
    return e.isRecursive
}

func (e *WatcherEntity) GetStartAt() string {
    return e.startAt
}

func (e *WatcherEntity) IsStartAtEnd() bool {
    return e.startAt == enums.WatcherStartAtEnum.END
}

// ScheduledCommandPrimitives un comando que se lanza cada interval_seconds y se corta a los timeout_seconds.
type ScheduledCommandPrimitives struct {
    Name            string `json:"name"`
    Command         string `json:"command"`
    IntervalSeconds int    `json:"interval_seconds"`
    TimeoutSeconds  int    `json:"timeout_seconds"`
}

type ScheduledCommandEntity struct {
    name            string
    commandLine     string
    intervalSeconds int
    timeoutSeconds  int
}

func ScheduledCommandEntityFromPrimitives(scheduledCommandPrimitives ScheduledCommandPrimitives) *ScheduledCommandEntity {
    return &ScheduledCommandEntity{
        name:            trim(scheduledCommandPrimitives.Name),
        commandLine:     trim(scheduledCommandPrimitives.Command),
        intervalSeconds: getIntOrDefault(scheduledCommandPrimitives.IntervalSeconds, DEFAULT_COMMAND_INTERVAL_SECONDS),
        timeoutSeconds:  getIntOrDefault(scheduledCommandPrimitives.TimeoutSeconds, DEFAULT_COMMAND_TIMEOUT_SECONDS),
    }
}

func (e *ScheduledCommandEntity) ToPrimitives() ScheduledCommandPrimitives {
    return ScheduledCommandPrimitives{
        Name:            e.name,
        Command:         e.commandLine,
        IntervalSeconds: e.intervalSeconds,
        TimeoutSeconds:  e.timeoutSeconds,
    }
}

func (e *ScheduledCommandEntity) GetName() string {
    return e.name
}

func (e *ScheduledCommandEntity) GetCommandLine() string {
    return e.commandLine
}

func (e *ScheduledCommandEntity) GetIntervalSeconds() int {
    return e.intervalSeconds
}

func (e *ScheduledCommandEntity) GetTimeoutSeconds() int {
    return e.timeoutSeconds
}

func (e *AgentConfigEntity) GetDockerConfigEntity() *DockerConfigEntity {
    return e.dockerConfigEntity
}
