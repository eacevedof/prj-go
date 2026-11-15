package collecthostmetrics

import (
    "time"

    "siem-client/src/core/config"
    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"
    agentconfigenums "siem-client/src/modules/agentconfig/domain/enums"
    agentconfigtraits "siem-client/src/modules/agentconfig/infrastructure/traits"
    collectorsservices "siem-client/src/modules/collectors/domain/services"
    "siem-client/src/modules/hostmetrics/domain/entities"
    "siem-client/src/modules/hostmetrics/domain/enums"
    "siem-client/src/modules/hostmetrics/domain/services"
    "siem-client/src/modules/hostmetrics/infrastructure/repositories"
    "siem-client/src/modules/shared/application"
    "siem-client/src/modules/shared/infrastructure/components/uuid"
    "siem-client/src/modules/shared/infrastructure/repositories/spool"
    "siem-client/src/modules/shared/infrastructure/repositories/state"
)

const DEFAULT_DISK_PATH = "/"

// CollectHostMetricsService una pasada por los monitores de metricas del host (cpu, memory, disk, gpu, network) que
// bajan del panel del backend: a los que les toca (interval_seconds) se les toma la muestra y pasa al spool como
// host_metric. Por monitor se guarda en state/host-metrics.json cuando se muestreo, los contadores de la muestra
// anterior (cpu y red salen de una diferencia: su primera muestra no manda evento) y si ya se aviso de que no esta
// disponible (sin nvidia-smi, SO no soportado, path que no existe): ese aviso se manda UNA vez hasta que vuelva.
type CollectHostMetricsService struct {
    agentconfigtraits.AgentConfigLoaderTrait

    uuidGenerator                    *uuid.UuidGenerator
    eventsBatcher                    *collectorsservices.EventsBatcher
    cpuUsageCalculator               *services.CpuUsageCalculator
    memoryUsageCalculator            *services.MemoryUsageCalculator
    networkRateCalculator            *services.NetworkRateCalculator
    gpuStatsParser                   *services.GpuStatsParser
    hostMetricSampleBuilder          *services.HostMetricSampleBuilder
    hostMetricEventBuilder           *services.HostMetricEventBuilder
    procFilesReaderDiskRepository    *repositories.ProcFilesReaderDiskRepository
    diskUsageReaderSyscallRepository *repositories.DiskUsageReaderSyscallRepository
    gpuStatsReaderShellRepository    *repositories.GpuStatsReaderShellRepository
    stateReaderJsonRepository        *state.StateReaderJsonRepository
    stateWriterJsonRepository        *state.StateWriterJsonRepository
    spoolWriterDiskRepository        *spool.SpoolWriterDiskRepository
}

func CollectHostMetricsServiceGetInstance() *CollectHostMetricsService {
    return &CollectHostMetricsService{
        AgentConfigLoaderTrait:           agentconfigtraits.NewAgentConfigLoaderTrait(),
        uuidGenerator:                    uuid.GetInstance(),
        eventsBatcher:                    collectorsservices.EventsBatcherGetInstance(),
        cpuUsageCalculator:               services.CpuUsageCalculatorGetInstance(),
        memoryUsageCalculator:            services.MemoryUsageCalculatorGetInstance(),
        networkRateCalculator:            services.NetworkRateCalculatorGetInstance(),
        gpuStatsParser:                   services.GpuStatsParserGetInstance(),
        hostMetricSampleBuilder:          services.HostMetricSampleBuilderGetInstance(),
        hostMetricEventBuilder:           services.HostMetricEventBuilderGetInstance(),
        procFilesReaderDiskRepository:    repositories.ProcFilesReaderDiskRepositoryGetInstance(),
        diskUsageReaderSyscallRepository: repositories.DiskUsageReaderSyscallRepositoryGetInstance(),
        gpuStatsReaderShellRepository:    repositories.GpuStatsReaderShellRepositoryGetInstance(),
        stateReaderJsonRepository:        state.StateReaderJsonRepositoryGetInstance(),
        stateWriterJsonRepository:        state.StateWriterJsonRepositoryGetInstance(),
        spoolWriterDiskRepository:        spool.SpoolWriterDiskRepositoryGetInstance(),
    }
}

// Invoke error: *AgentConfigException (400, 404) o uno inesperado (disco). Una metrica que no se puede tomar no es un
// error del servicio: es un evento is_unavailable (una vez).
func (s *CollectHostMetricsService) Invoke(collectHostMetricsDto *CollectHostMetricsDto) (*CollectHostMetricsResultDto, error) {
    agentConfigEntity, err := s.GetValidAgentConfigEntityOrFail()
    if err != nil {
        return nil, err
    }
    metricMonitorEntities := agentConfigEntity.GetMetricMonitorEntities()
    hostMetricsState, err := s.stateReaderJsonRepository.GetState(s.GetStatePath(), enums.HostMetricsStateEnum.STATE_FILE)
    if err != nil {
        return nil, err
    }
    monitorStatesByUuid, _ := hostMetricsState[enums.HostMetricsStateEnum.MONITORS].(map[string]any)
    if len(metricMonitorEntities) == 0 && len(monitorStatesByUuid) == 0 {
        return CollectHostMetricsResultDtoFromPrimitives(map[string]any{}), nil
    }

    // solo los monitores que siguen en la config: uno que se borra en el panel deja de figurar en el estado
    nextMonitorStatesByUuid := map[string]any{}
    events := []map[string]any{}
    sampledMonitors, unavailableMonitors := 0, 0
    now := time.Now().UTC()
    for _, remoteMonitorEntity := range metricMonitorEntities {
        monitorState, _ := monitorStatesByUuid[remoteMonitorEntity.GetMonitorUuid()].(map[string]any)
        lastRunUnix := int64(application.GetInt(monitorState, enums.HostMetricsStateEnum.LAST_RUN_UNIX))
        isDue := collectHostMetricsDto.IsForced() || lastRunUnix == 0 || now.Unix()-lastRunUnix >= int64(remoteMonitorEntity.GetIntervalSeconds())
        if !isDue {
            nextMonitorStatesByUuid[remoteMonitorEntity.GetMonitorUuid()] = monitorState
            continue
        }

        hostMetricSampleEntity, nextMonitorState, unavailableReason := s.getSample(remoteMonitorEntity, monitorState, now)
        nextMonitorState[enums.HostMetricsStateEnum.LAST_RUN_UNIX] = now.Unix()
        nextMonitorState[enums.HostMetricsStateEnum.IS_UNAVAILABLE] = unavailableReason != ""
        nextMonitorStatesByUuid[remoteMonitorEntity.GetMonitorUuid()] = nextMonitorState

        switch {
        case unavailableReason != "":
            if monitorState[enums.HostMetricsStateEnum.IS_UNAVAILABLE] != true {
                events = append(events, s.hostMetricEventBuilder.GetUnavailableEvent(remoteMonitorEntity, unavailableReason, now))
                unavailableMonitors++
            }
        case hostMetricSampleEntity.HasValue():
            events = append(events, s.hostMetricEventBuilder.GetMetricEvent(remoteMonitorEntity, hostMetricSampleEntity, now))
            sampledMonitors++
        }
    }

    spooledBatches := 0
    for _, batchEvents := range s.eventsBatcher.GetBatches(events, agentConfigEntity.GetMaxEventsPerBatch()) {
        if err := s.spoolWriterDiskRepository.SaveSpoolBatch(s.GetSpoolPath(), s.uuidGenerator.GetUniqueIdWithPrefix("b-"), batchEvents); err != nil {
            return nil, err
        }
        spooledBatches++
    }
    err = s.stateWriterJsonRepository.SaveState(s.GetStatePath(), enums.HostMetricsStateEnum.STATE_FILE, map[string]any{
        enums.HostMetricsStateEnum.MONITORS: nextMonitorStatesByUuid,
    })
    if err != nil {
        return nil, err
    }

    return CollectHostMetricsResultDtoFromPrimitives(map[string]any{
        "metricMonitors":      len(metricMonitorEntities),
        "sampledMonitors":     sampledMonitors,
        "unavailableMonitors": unavailableMonitors,
        "spooledBatches":      spooledBatches,
        "results":             events,
    }), nil
}

// getSample la muestra, lo que hay que guardar para la siguiente (contadores) y, si la metrica no se puede tomar, el
// motivo (entonces la muestra no cuenta).
func (s *CollectHostMetricsService) getSample(remoteMonitorEntity *agentconfigentities.RemoteMonitorEntity, monitorState map[string]any, now time.Time) (*entities.HostMetricSampleEntity, map[string]any, string) {
    switch remoteMonitorEntity.GetMonitorType() {
    case agentconfigenums.MonitorTypeEnum.DISK:
        return s.getDiskSample(remoteMonitorEntity)
    case agentconfigenums.MonitorTypeEnum.GPU:
        return s.getGpuSample()
    }

    // cpu, memory y network salen de /proc
    if !s.procFilesReaderDiskRepository.IsSupported() {
        return nil, map[string]any{}, config.Tr(enums.HostMetricsMessageEnum.OS_X_IS_NOT_SUPPORTED, s.procFilesReaderDiskRepository.GetOsName())
    }
    switch remoteMonitorEntity.GetMonitorType() {
    case agentconfigenums.MonitorTypeEnum.CPU:
        return s.getCpuSample(monitorState)
    case agentconfigenums.MonitorTypeEnum.MEMORY:
        return s.getMemorySample()
    }
    return s.getNetworkSample(remoteMonitorEntity, monitorState, now)
}

func (s *CollectHostMetricsService) getCpuSample(monitorState map[string]any) (*entities.HostMetricSampleEntity, map[string]any, string) {
    procStatContent, err := s.procFilesReaderDiskRepository.GetProcFileContent(repositories.PROC_STAT)
    if err != nil {
        return nil, map[string]any{}, config.Tr(enums.HostMetricsMessageEnum.CANNOT_READ_X, repositories.PROC_STAT)
    }
    total, idle, err := s.cpuUsageCalculator.GetCpuTimes(procStatContent)
    if err != nil {
        return nil, map[string]any{}, config.Tr(enums.HostMetricsMessageEnum.X_HAS_AN_UNEXPECTED_FORMAT, repositories.PROC_STAT)
    }
    nextMonitorState := map[string]any{
        enums.HostMetricsStateEnum.CPU_TOTAL: total,
        enums.HostMetricsStateEnum.CPU_IDLE:  idle,
    }

    usagePercent := s.cpuUsageCalculator.GetUsagePercent(
        getFloat(monitorState, enums.HostMetricsStateEnum.CPU_TOTAL), getFloat(monitorState, enums.HostMetricsStateEnum.CPU_IDLE), total, idle,
    )
    if usagePercent < 0 {
        return entities.HostMetricSampleEntityWithoutValue(), nextMonitorState, ""
    }
    loadAverages, hasLoadAverages := [3]float64{}, false
    if procLoadavgContent, err := s.procFilesReaderDiskRepository.GetProcFileContent(repositories.PROC_LOADAVG); err == nil {
        loadAverages, hasLoadAverages = s.cpuUsageCalculator.GetLoadAverages(procLoadavgContent)
    }
    return s.hostMetricSampleBuilder.GetCpuSample(usagePercent, loadAverages, hasLoadAverages), nextMonitorState, ""
}

func (s *CollectHostMetricsService) getMemorySample() (*entities.HostMetricSampleEntity, map[string]any, string) {
    procMeminfoContent, err := s.procFilesReaderDiskRepository.GetProcFileContent(repositories.PROC_MEMINFO)
    if err != nil {
        return nil, map[string]any{}, config.Tr(enums.HostMetricsMessageEnum.CANNOT_READ_X, repositories.PROC_MEMINFO)
    }
    totalBytes, usedBytes, err := s.memoryUsageCalculator.GetMemoryBytes(procMeminfoContent)
    if err != nil {
        return nil, map[string]any{}, config.Tr(enums.HostMetricsMessageEnum.X_HAS_AN_UNEXPECTED_FORMAT, repositories.PROC_MEMINFO)
    }
    usedPercent := s.memoryUsageCalculator.GetUsedPercent(totalBytes, usedBytes)
    return s.hostMetricSampleBuilder.GetMemorySample(usedPercent, usedBytes, totalBytes), map[string]any{}, ""
}

// getDiskSample % como df: usados / (usados + disponibles para un usuario normal).
func (s *CollectHostMetricsService) getDiskSample(remoteMonitorEntity *agentconfigentities.RemoteMonitorEntity) (*entities.HostMetricSampleEntity, map[string]any, string) {
    path := remoteMonitorEntity.GetConfigString(agentconfigenums.RemoteConfigKeyEnum.PATH)
    if path == "" {
        path = DEFAULT_DISK_PATH
    }
    totalBytes, usedBytes, availableBytes, err := s.diskUsageReaderSyscallRepository.GetDiskUsageBytes(path)
    if err != nil {
        return nil, map[string]any{}, config.Tr(enums.HostMetricsMessageEnum.PATH_X_NOT_AVAILABLE_X, path, err.Error())
    }
    if usedBytes+availableBytes <= 0 {
        return nil, map[string]any{}, config.Tr(enums.HostMetricsMessageEnum.PATH_X_NOT_AVAILABLE_X, path, "empty filesystem")
    }
    usedPercent := s.memoryUsageCalculator.GetUsedPercent(usedBytes+availableBytes, usedBytes)
    return s.hostMetricSampleBuilder.GetDiskSample(path, usedPercent, usedBytes, totalBytes), map[string]any{}, ""
}

func (s *CollectHostMetricsService) getNetworkSample(remoteMonitorEntity *agentconfigentities.RemoteMonitorEntity, monitorState map[string]any, now time.Time) (*entities.HostMetricSampleEntity, map[string]any, string) {
    procNetDevContent, err := s.procFilesReaderDiskRepository.GetProcFileContent(repositories.PROC_NET_DEV)
    if err != nil {
        return nil, map[string]any{}, config.Tr(enums.HostMetricsMessageEnum.CANNOT_READ_X, repositories.PROC_NET_DEV)
    }
    interfaceName := remoteMonitorEntity.GetConfigString(agentconfigenums.RemoteConfigKeyEnum.INTERFACE)
    rxBytes, txBytes, isFound := s.networkRateCalculator.GetInterfaceBytes(procNetDevContent, interfaceName)
    if !isFound {
        interfaceLabel := interfaceName
        if interfaceLabel == "" {
            interfaceLabel = "(any)"
        }
        return nil, map[string]any{}, config.Tr(enums.HostMetricsMessageEnum.INTERFACE_X_NOT_FOUND, interfaceLabel)
    }
    nowMs := float64(now.UnixMilli())
    nextMonitorState := map[string]any{
        enums.HostMetricsStateEnum.NET_RX_BYTES: rxBytes,
        enums.HostMetricsStateEnum.NET_TX_BYTES: txBytes,
        enums.HostMetricsStateEnum.NET_AT_MS:    nowMs,
    }

    previousAtMs := getFloat(monitorState, enums.HostMetricsStateEnum.NET_AT_MS)
    if previousAtMs <= 0 {
        return entities.HostMetricSampleEntityWithoutValue(), nextMonitorState, ""
    }
    elapsedMs := nowMs - previousAtMs
    rxMbit := s.networkRateCalculator.GetMbitPerSecond(getFloat(monitorState, enums.HostMetricsStateEnum.NET_RX_BYTES), rxBytes, elapsedMs)
    txMbit := s.networkRateCalculator.GetMbitPerSecond(getFloat(monitorState, enums.HostMetricsStateEnum.NET_TX_BYTES), txBytes, elapsedMs)
    if rxMbit < 0 || txMbit < 0 {
        return entities.HostMetricSampleEntityWithoutValue(), nextMonitorState, ""
    }
    return s.hostMetricSampleBuilder.GetNetworkSample(interfaceName, rxMbit, txMbit), nextMonitorState, ""
}

func (s *CollectHostMetricsService) getGpuSample() (*entities.HostMetricSampleEntity, map[string]any, string) {
    nvidiaSmiOutput, err := s.gpuStatsReaderShellRepository.GetNvidiaSmiOutput()
    if err != nil {
        return nil, map[string]any{}, config.Tr(enums.HostMetricsMessageEnum.NVIDIA_SMI_NOT_AVAILABLE_X, err.Error())
    }
    gpuStats := s.gpuStatsParser.GetGpuStats(nvidiaSmiOutput)
    busiestGpuStat, hasGpu := s.gpuStatsParser.GetBusiestGpuStat(gpuStats)
    if !hasGpu {
        return nil, map[string]any{}, config.Tr(enums.HostMetricsMessageEnum.NVIDIA_SMI_NOT_AVAILABLE_X, "no gpu in its output")
    }
    return s.hostMetricSampleBuilder.GetGpuSample(busiestGpuStat, len(gpuStats)), map[string]any{}, ""
}

// getFloat los contadores del JSON de estado vuelven como float64 (jiffies y bytes acumulados no caben en int32).
func getFloat(primitives map[string]any, key string) float64 {
    value, _ := primitives[key].(float64)
    return value
}
