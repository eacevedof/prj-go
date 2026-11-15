package services

import (
    "strconv"
    "strings"

    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/collectors/domain/entities"
    "siem-client/src/modules/collectors/domain/enums"
)

const (
    // el backend guarda hasta 8000 caracteres por mensaje: mandar mas es gastar red
    MESSAGE_MAX_LENGTH   = 8000
    EMPTY_COMMAND_OUTPUT = "(no output)"
)

// CollectedEventBuilder monta los eventos con el contrato del backend. La gravedad no se manda: la deduce el backend
// (palabras clave en una linea, codigo de salida en un comando), asi la regla vive en un solo sitio.
type CollectedEventBuilder struct{}

func CollectedEventBuilderGetInstance() *CollectedEventBuilder {
    return &CollectedEventBuilder{}
}

func (b *CollectedEventBuilder) GetLogLineEvent(sourceName string, filePath string, line string, occurredAt string) map[string]any {
    return map[string]any{
        enums.CollectorsKeyEnum.EVENT_TYPE:  enums.EventTypeEnum.LOG_LINE,
        enums.CollectorsKeyEnum.SOURCE_NAME: sourceName,
        enums.CollectorsKeyEnum.FILE_PATH:   filePath,
        enums.CollectorsKeyEnum.MESSAGE:     b.getTruncated(line),
        enums.CollectorsKeyEnum.OCCURRED_AT: occurredAt,
    }
}

// GetCommandResultEvent un comando cortado por tiempo sale con exit_code -1 y el aviso delante de su salida.
func (b *CollectedEventBuilder) GetCommandResultEvent(sourceName string, commandResultEntity *entities.CommandResultEntity, timeoutSeconds int, occurredAt string) map[string]any {
    message := strings.TrimSpace(commandResultEntity.GetOutput())
    if message == "" {
        message = EMPTY_COMMAND_OUTPUT
    }
    if commandResultEntity.IsTimedOut() {
        message = "[timeout after " + strconv.Itoa(timeoutSeconds) + " s] " + message
    }
    return map[string]any{
        enums.CollectorsKeyEnum.EVENT_TYPE:  enums.EventTypeEnum.COMMAND_RESULT,
        enums.CollectorsKeyEnum.SOURCE_NAME: sourceName,
        enums.CollectorsKeyEnum.MESSAGE:     b.getTruncated(message),
        enums.CollectorsKeyEnum.EXIT_CODE:   commandResultEntity.GetExitCode(),
        enums.CollectorsKeyEnum.DURATION_MS: commandResultEntity.GetDurationMs(),
        enums.CollectorsKeyEnum.OCCURRED_AT: occurredAt,
    }
}

func (b *CollectedEventBuilder) getTruncated(text string) string {
    runes := []rune(text)
    if len(runes) <= MESSAGE_MAX_LENGTH {
        return text
    }
    return string(runes[:MESSAGE_MAX_LENGTH])
}

// GetWatcherLogLineEvent una linea de un watcher. Si el watcher declara severity, viaja en el evento (el backend
// respeta la declarada); si sale de un monitor del panel (folder / log_pattern), attributes lleva su uuid y su tipo.
func (b *CollectedEventBuilder) GetWatcherLogLineEvent(watcherEntity *agentconfigentities.WatcherEntity, filePath string, line string, occurredAt string) map[string]any {
    event := b.GetLogLineEvent(watcherEntity.GetName(), filePath, line, occurredAt)
    if watcherEntity.GetSeverity() != "" {
        event[enums.CollectorsKeyEnum.SEVERITY] = watcherEntity.GetSeverity()
    }
    if watcherEntity.IsFromMonitor() {
        event[enums.CollectorsKeyEnum.ATTRIBUTES] = map[string]any{
            enums.CollectorsKeyEnum.MONITOR_UUID: watcherEntity.GetMonitorUuid(),
            enums.CollectorsKeyEnum.MONITOR_TYPE: watcherEntity.GetMonitorType(),
        }
    }
    return event
}
