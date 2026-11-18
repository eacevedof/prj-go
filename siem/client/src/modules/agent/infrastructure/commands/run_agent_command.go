package commands

import (
    "context"
    "errors"
    "fmt"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"

    "siem-client/src/modules/agentconfig/application/checkagentconfig"
    "siem-client/src/modules/agentconfig/application/syncremoteconfig"
    agentconfigexceptions "siem-client/src/modules/agentconfig/domain/exceptions"
    "siem-client/src/modules/collectors/application/collectfilelines"
    "siem-client/src/modules/collectors/application/runduecommands"
    "siem-client/src/modules/connections/application/collectconnectionchanges"
    "siem-client/src/modules/docker/application/collectcontainerlogs"
    "siem-client/src/modules/docker/application/collectcontainerstates"
    "siem-client/src/modules/docker/application/collectcontainerstats"
    "siem-client/src/modules/docker/application/collectdaemonevents"
    dockerexceptions "siem-client/src/modules/docker/domain/exceptions"
    "siem-client/src/modules/hostmetrics/application/collecthostmetrics"
    sharedcommands "siem-client/src/modules/shared/infrastructure/commands"
    sharedexceptions "siem-client/src/modules/shared/infrastructure/exceptions"
    "siem-client/src/modules/shared/infrastructure/repositories/configuration"
    "siem-client/src/modules/shipper/application/sendspooledevents"
    shipperexceptions "siem-client/src/modules/shipper/domain/exceptions"
)

const (
    dateTimeFormat = "2006-01-02 15:04:05"
    // cada cuanto se pide al backend la config remota (los monitores del panel)
    remoteConfigSyncInterval = 60 * time.Second
    stepRemoteMonitors       = "config-monitors"
)

// RunAgentCommand ./console agent-run   el agente: un bucle que hasta Ctrl+C (o SIGTERM del servicio) hace
//   - al arrancar y cada 60 s: baja del backend su config remota (los monitores que se crean en el panel).
//   - cada collect_interval_seconds: lineas nuevas de las carpetas vigiladas + comandos a los que les toca + metricas
//     del host a las que les toca (cpu, memoria, disco, gpu, red) + conexiones de red que abren o cierran (si el
//     panel tiene un monitor connections para este agente; solo Linux) + Docker
//     (estados, eventos del daemon, recursos si tocan, logs) -> spool.
//   - cada send_interval_seconds: spool -> backend (o el latido si no hay nada).
//
// Es el controller de todos los casos de uso: el unico que captura. Un fallo se imprime y se loguea, y el bucle sigue
// (el backend caido no para la recogida: el spool aguanta). Un mismo error no se repite en pantalla en cada vuelta.
// Si la configuracion no vale al arrancar, no arranca (exit 1).
type RunAgentCommand struct {
    sharedcommands.AbstractConsoleCommand
    lastErrorByStep map[string]string
}

func RunAgentCommandGetInstance(input []string) *RunAgentCommand {
    return &RunAgentCommand{
        AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input),
        lastErrorByStep:        map[string]string{},
    }
}

func (c *RunAgentCommand) Run() {
    // antes de validar: asi el primer arranque ya cuenta con los monitores del panel (si el backend responde)
    c.syncRemoteConfig()
    checkAgentConfigResultDto, err := checkagentconfig.CheckAgentConfigServiceGetInstance().Invoke(
        checkagentconfig.CheckAgentConfigDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        var agentConfigException *agentconfigexceptions.AgentConfigException
        if errors.As(err, &agentConfigException) {
            c.ExitWithDomainError(agentConfigException)
        }
        c.ExitWithUnexpectedError(err, "RunAgentCommand.Run")
    }

    environment := configuration.GetInstance()
    c.printLine(c.GetBlue(fmt.Sprintf(
        "%s v%s -> %s · %d watchers, %d commands, %d metrics, connections %s, docker %s, remote config %s · collect every %d s, send every %d s",
        environment.GetAppName(), environment.GetAppVersion(), checkAgentConfigResultDto.GetEndpointUrl(),
        checkAgentConfigResultDto.GetWatchersCount(), checkAgentConfigResultDto.GetCommandsCount(),
        checkAgentConfigResultDto.GetMetricMonitorsCount(), c.getOnOff(checkAgentConfigResultDto.IsConnectionsEnabled()),
        c.getOnOff(checkAgentConfigResultDto.IsDockerEnabled()),
        c.getOnOff(checkAgentConfigResultDto.IsRemoteConfigEnabled()),
        checkAgentConfigResultDto.GetCollectIntervalSeconds(), checkAgentConfigResultDto.GetSendIntervalSeconds(),
    )))
    c.printStepWarningIfChanged(stepRemoteMonitors, strings.Join(checkAgentConfigResultDto.GetRemoteConfigWarnings(), " · "))

    stopContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    collectInterval := time.Duration(checkAgentConfigResultDto.GetCollectIntervalSeconds()) * time.Second
    sendInterval := time.Duration(checkAgentConfigResultDto.GetSendIntervalSeconds()) * time.Second
    var nextCollectAt, nextSendAt time.Time
    nextSyncAt := time.Now().Add(remoteConfigSyncInterval)
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    for {
        now := time.Now()
        if !now.Before(nextSyncAt) {
            c.syncRemoteConfig()
            c.printRemoteMonitorsWarnings()
            nextSyncAt = now.Add(remoteConfigSyncInterval)
        }
        if !now.Before(nextCollectAt) {
            c.collectFileLines()
            c.runDueCommands()
            c.collectHostMetrics()
            c.collectConnectionChanges()
            c.collectDocker()
            nextCollectAt = now.Add(collectInterval)
        }
        if !now.Before(nextSendAt) {
            c.sendSpooledEvents()
            nextSendAt = now.Add(sendInterval)
        }

        select {
        case <-stopContext.Done():
            c.printLine("stop signal received: bye (pending events stay in the spool)")
            return
        case <-ticker.C:
        }
    }
}

// syncRemoteConfig baja del backend los monitores del panel (al arrancar y cada remoteConfigSyncInterval). Si falla,
// el agente sigue con la ultima config guardada y se avisa una vez.
func (c *RunAgentCommand) syncRemoteConfig() {
    const step = "config-sync"
    syncRemoteConfigResultDto, err := syncremoteconfig.SyncRemoteConfigServiceGetInstance().Invoke(
        syncremoteconfig.SyncRemoteConfigDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.printStepError(step, err)
        return
    }
    if syncRemoteConfigResultDto.IsDisabled() {
        return
    }
    if c.lastErrorByStep[step] != "" {
        c.printLine(c.GetGreen(step + ": remote config reachable again"))
    }
    c.lastErrorByStep[step] = ""
    if syncRemoteConfigResultDto.IsChanged() {
        c.printLine(fmt.Sprintf("%s: remote config updated, %d monitors", step, syncRemoteConfigResultDto.GetMonitors()))
    }
}

// printRemoteMonitorsWarnings monitores del panel que no se han podido aplicar (regex o path no validos, nombre que
// choca con un watcher local): se avisa cuando cambia la lista, no en cada pasada.
func (c *RunAgentCommand) printRemoteMonitorsWarnings() {
    checkAgentConfigResultDto, err := checkagentconfig.CheckAgentConfigServiceGetInstance().Invoke(
        checkagentconfig.CheckAgentConfigDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        return
    }
    c.printStepWarningIfChanged(stepRemoteMonitors, strings.Join(checkAgentConfigResultDto.GetRemoteConfigWarnings(), " · "))
}

func (c *RunAgentCommand) collectHostMetrics() {
    const step = "metrics-collect"
    collectHostMetricsResultDto, err := collecthostmetrics.CollectHostMetricsServiceGetInstance().Invoke(
        collecthostmetrics.CollectHostMetricsDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.printStepError(step, err)
        return
    }
    c.lastErrorByStep[step] = ""
    if collectHostMetricsResultDto.GetSampledMonitors() > 0 || collectHostMetricsResultDto.GetUnavailableMonitors() > 0 {
        c.printLine(fmt.Sprintf("%s: %d sampled, %d not available", step,
            collectHostMetricsResultDto.GetSampledMonitors(), collectHostMetricsResultDto.GetUnavailableMonitors()))
    }
}

// collectConnectionChanges el monitor connections del panel: sin el no hace nada (y olvida la foto). Los procesos que
// no se han podido leer (sin root) se avisan cuando cambia su numero, no en cada pasada.
func (c *RunAgentCommand) collectConnectionChanges() {
    const step = "connections-collect"
    const stepUnreadable = "connections-processes"
    collectConnectionChangesResultDto, err := collectconnectionchanges.CollectConnectionChangesServiceGetInstance().Invoke(
        collectconnectionchanges.CollectConnectionChangesDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.printStepError(step, err)
        return
    }
    c.lastErrorByStep[step] = ""
    if collectConnectionChangesResultDto.IsUnavailable() {
        c.printLine(c.GetYellow(step + ": not available on this system (warned to the backend once)"))
    }
    if !collectConnectionChangesResultDto.IsSampled() {
        return
    }
    unreadableWarning := ""
    if unreadableProcesses := collectConnectionChangesResultDto.GetUnreadableProcesses(); unreadableProcesses > 0 {
        unreadableWarning = fmt.Sprintf("%d processes could not be read: their sockets show process ? (run as root)", unreadableProcesses)
    }
    c.printStepWarningIfChanged(stepUnreadable, unreadableWarning)
    if collectConnectionChangesResultDto.GetOpened() > 0 || collectConnectionChangesResultDto.GetClosed() > 0 {
        c.printLine(fmt.Sprintf("%s: %d opened, %d closed · %d watched", step,
            collectConnectionChangesResultDto.GetOpened(), collectConnectionChangesResultDto.GetClosed(),
            collectConnectionChangesResultDto.GetConnections()))
    }
}

func (c *RunAgentCommand) collectFileLines() {
    const step = "files-collect"
    collectFileLinesResultDto, err := collectfilelines.CollectFileLinesServiceGetInstance().Invoke(
        collectfilelines.CollectFileLinesDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.printStepError(step, err)
        return
    }

    problems := []string{}
    if missingFolders := collectFileLinesResultDto.GetMissingFolders(); len(missingFolders) > 0 {
        problems = append(problems, "missing folders: "+strings.Join(missingFolders, ", "))
    }
    if unreadableFiles := collectFileLinesResultDto.GetUnreadableFiles(); len(unreadableFiles) > 0 {
        problems = append(problems, "unreadable files: "+strings.Join(unreadableFiles, ", "))
    }
    c.printStepWarningIfChanged(step, strings.Join(problems, " · "))
    if collectFileLinesResultDto.GetLinesCollected() > 0 {
        c.printLine(fmt.Sprintf("%s: %d lines from %d files", step,
            collectFileLinesResultDto.GetLinesCollected(), collectFileLinesResultDto.GetFilesRead()))
    }
}

func (c *RunAgentCommand) runDueCommands() {
    const step = "commands-run"
    runDueCommandsResultDto, err := runduecommands.RunDueCommandsServiceGetInstance().Invoke(
        runduecommands.RunDueCommandsDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.printStepError(step, err)
        return
    }
    c.lastErrorByStep[step] = ""
    if runDueCommandsResultDto.GetRanCommands() > 0 {
        c.printLine(fmt.Sprintf("%s: %d ran, %d with exit code != 0", step,
            runDueCommandsResultDto.GetRanCommands(), runDueCommandsResultDto.GetFailedCommands()))
    }
}

// collectDocker los cuatro colectores de Docker; con el bloque docker apagado no hacen nada. Si Docker no responde se
// avisa una vez (como el backend caido) y el resto del agente sigue.
func (c *RunAgentCommand) collectDocker() {
    const step = "docker-collect"
    collectContainerStatesResultDto, err := collectcontainerstates.CollectContainerStatesServiceGetInstance().Invoke(
        collectcontainerstates.CollectContainerStatesDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.printStepError(step, err)
        return
    }
    if collectContainerStatesResultDto.IsDisabled() {
        return
    }
    collectDaemonEventsResultDto, err := collectdaemonevents.CollectDaemonEventsServiceGetInstance().Invoke(
        collectdaemonevents.CollectDaemonEventsDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.printStepError(step, err)
        return
    }
    collectContainerStatsResultDto, err := collectcontainerstats.CollectContainerStatsServiceGetInstance().Invoke(
        collectcontainerstats.CollectContainerStatsDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.printStepError(step, err)
        return
    }
    collectContainerLogsResultDto, err := collectcontainerlogs.CollectContainerLogsServiceGetInstance().Invoke(
        collectcontainerlogs.CollectContainerLogsDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.printStepError(step, err)
        return
    }

    if c.lastErrorByStep[step] != "" {
        c.printLine(c.GetGreen(step + ": docker reachable again"))
    }
    c.lastErrorByStep[step] = ""
    changes := collectContainerStatesResultDto.GetChanges()
    daemonEvents := collectDaemonEventsResultDto.GetDaemonEvents()
    sampledContainers := collectContainerStatsResultDto.GetSampledContainers()
    logLines := collectContainerLogsResultDto.GetLogLines()
    if changes+daemonEvents+sampledContainers+logLines == 0 {
        return
    }
    c.printLine(fmt.Sprintf("%s: %d containers (%d running) · %d changes · %d daemon events · %d stats · %d log lines", step,
        collectContainerStatesResultDto.GetContainers(), collectContainerStatesResultDto.GetRunningContainers(),
        changes, daemonEvents, sampledContainers, logLines))
}

func (c *RunAgentCommand) sendSpooledEvents() {
    const step = "events-send"
    sendSpooledEventsResultDto, err := sendspooledevents.SendSpooledEventsServiceGetInstance().Invoke(
        sendspooledevents.SendSpooledEventsDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.printStepError(step, err)
        return
    }
    if c.lastErrorByStep[step] != "" {
        c.printLine(c.GetGreen(step + ": backend reachable again"))
    }
    c.lastErrorByStep[step] = ""

    if sendSpooledEventsResultDto.GetSentBatches() > 0 || sendSpooledEventsResultDto.GetRejectedBatches() > 0 || sendSpooledEventsResultDto.GetDroppedBatches() > 0 {
        c.printLine(fmt.Sprintf("%s: %d events in %d batches sent · %d rejected · %d dropped · %d pending", step,
            sendSpooledEventsResultDto.GetSentEvents(), sendSpooledEventsResultDto.GetSentBatches(),
            sendSpooledEventsResultDto.GetRejectedBatches(), sendSpooledEventsResultDto.GetDroppedBatches(),
            sendSpooledEventsResultDto.GetPendingBatches()))
    }
}

// printStepError errores de dominio en rojo tal cual; inesperados, ademas al log con traza. Solo si cambia.
func (c *RunAgentCommand) printStepError(step string, err error) {
    var agentConfigException *agentconfigexceptions.AgentConfigException
    var shipperException *shipperexceptions.ShipperException
    var fieldsException *sharedexceptions.FieldsException
    var dockerException *dockerexceptions.DockerException
    isDomainError := errors.As(err, &agentConfigException) || errors.As(err, &shipperException) || errors.As(err, &fieldsException) ||
        errors.As(err, &dockerException)
    if !isDomainError {
        c.LogException(err, "RunAgentCommand."+step)
    }
    if c.lastErrorByStep[step] == err.Error() {
        return
    }
    c.lastErrorByStep[step] = err.Error()
    c.LogWarning(err.Error(), "RunAgentCommand."+step)
    c.printLine(c.GetRed(step + ": " + err.Error()))
}

func (c *RunAgentCommand) printStepWarningIfChanged(step string, warning string) {
    if c.lastErrorByStep[step] == warning {
        return
    }
    c.lastErrorByStep[step] = warning
    if warning != "" {
        c.LogWarning(warning, "RunAgentCommand."+step)
        c.printLine(c.GetYellow(step + ": " + warning))
    }
}

func (c *RunAgentCommand) printLine(line string) {
    fmt.Printf("[%s] %s\n", time.Now().UTC().Format(dateTimeFormat), line)
}

func (c *RunAgentCommand) getOnOff(isOn bool) string {
    if isOn {
        return "on"
    }
    return "off"
}
