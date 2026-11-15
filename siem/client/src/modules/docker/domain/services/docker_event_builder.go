package services

import (
    "fmt"
    "strings"
    "time"

    collectorsenums "siem-client/src/modules/collectors/domain/enums"
    "siem-client/src/modules/docker/domain/entities"
    "siem-client/src/modules/docker/domain/enums"
    "siem-client/src/modules/shared/application"
)

const (
    dateTimeFormat = "2006-01-02 15:04:05"
    // el backend guarda hasta 8000 caracteres por mensaje
    MESSAGE_MAX_LENGTH = 8000
    ATTRIBUTES_KEY     = "attributes"
    // eventos del daemon que no son de un contenedor (imagenes, redes, volumenes, el propio daemon)
    DAEMON_SOURCE_NAME = "docker-daemon"
)

// DockerEventBuilder monta los eventos de Docker con el contrato del backend: un mensaje legible para la tabla de
// eventos y los attributes con los datos sueltos. La gravedad no se manda: la decide el backend con los attributes.
// source_name = nombre del contenedor (cada contenedor es su propia fuente).
type DockerEventBuilder struct{}

func DockerEventBuilderGetInstance() *DockerEventBuilder {
    return &DockerEventBuilder{}
}

// GetContainerStateEvent "cont-web: running -> exited (exit 137, OOM killed) · image nginx".
func (b *DockerEventBuilder) GetContainerStateEvent(containerChangeEntity *entities.ContainerChangeEntity, occurredAt time.Time) map[string]any {
    containerEntity := containerChangeEntity.GetContainerEntity()
    attributes := map[string]any{
        enums.DockerAttributeEnum.CONTAINER_ID:  containerEntity.GetShortId(),
        enums.DockerAttributeEnum.IMAGE:         containerEntity.GetImage(),
        enums.DockerAttributeEnum.CHANGE_KIND:   containerChangeEntity.GetChangeKind(),
        enums.DockerAttributeEnum.STATUS:        containerEntity.GetStatus(),
        enums.DockerAttributeEnum.HEALTH:        containerEntity.GetHealth(),
        enums.DockerAttributeEnum.EXIT_CODE:     containerEntity.GetExitCode(),
        enums.DockerAttributeEnum.IS_OOM_KILLED: containerEntity.IsOomKilled(),
        enums.DockerAttributeEnum.RESTART_COUNT: containerEntity.GetRestartCount(),
    }
    if containerChangeEntity.GetPreviousStatus() != "" {
        attributes[enums.DockerAttributeEnum.PREVIOUS_STATUS] = containerChangeEntity.GetPreviousStatus()
        attributes[enums.DockerAttributeEnum.PREVIOUS_HEALTH] = containerChangeEntity.GetPreviousHealth()
    }

    exitCode := 0
    if containerEntity.IsExited() {
        exitCode = containerEntity.GetExitCode()
    }
    return map[string]any{
        collectorsenums.CollectorsKeyEnum.EVENT_TYPE:  enums.DockerEventTypeEnum.DOCKER_STATE,
        collectorsenums.CollectorsKeyEnum.SOURCE_NAME: containerEntity.GetName(),
        collectorsenums.CollectorsKeyEnum.MESSAGE:     b.getContainerStateMessage(containerChangeEntity),
        collectorsenums.CollectorsKeyEnum.EXIT_CODE:   exitCode,
        collectorsenums.CollectorsKeyEnum.OCCURRED_AT: occurredAt.UTC().Format(dateTimeFormat),
        ATTRIBUTES_KEY: attributes,
    }
}

func (b *DockerEventBuilder) getContainerStateMessage(containerChangeEntity *entities.ContainerChangeEntity) string {
    containerEntity := containerChangeEntity.GetContainerEntity()
    name := containerEntity.GetName()
    currentState := b.getStateLabel(containerEntity)
    image := " · image " + containerEntity.GetImage()

    switch containerChangeEntity.GetChangeKind() {
    case enums.ContainerChangeKindEnum.STATUS_CHANGED:
        return fmt.Sprintf("%s: %s -> %s%s", name, containerChangeEntity.GetPreviousStatus(), currentState, image)
    case enums.ContainerChangeKindEnum.RESTARTED:
        return fmt.Sprintf("%s: restarted (restart count %d), now %s%s", name, containerEntity.GetRestartCount(), currentState, image)
    case enums.ContainerChangeKindEnum.HEALTH_CHANGED:
        return fmt.Sprintf("%s: health %s -> %s%s", name, b.getOrNone(containerChangeEntity.GetPreviousHealth()), b.getOrNone(containerEntity.GetHealth()), image)
    case enums.ContainerChangeKindEnum.REMOVED:
        return fmt.Sprintf("%s: removed (was %s)%s", name, containerChangeEntity.GetPreviousStatus(), image)
    case enums.ContainerChangeKindEnum.APPEARED:
        return fmt.Sprintf("%s: new container, %s%s", name, currentState, image)
    }
    return fmt.Sprintf("%s: %s%s", name, currentState, image)
}

// getStateLabel "exited (exit 137, OOM killed)", "running (unhealthy)".
func (b *DockerEventBuilder) getStateLabel(containerEntity *entities.ContainerEntity) string {
    details := []string{}
    if containerEntity.IsExited() {
        details = append(details, fmt.Sprintf("exit %d", containerEntity.GetExitCode()))
    }
    if containerEntity.IsOomKilled() {
        details = append(details, "OOM killed")
    }
    if containerEntity.GetHealth() != "" {
        details = append(details, containerEntity.GetHealth())
    }
    if len(details) == 0 {
        return containerEntity.GetStatus()
    }
    return containerEntity.GetStatus() + " (" + strings.Join(details, ", ") + ")"
}

// GetContainerStatsEvent "cpu 12.3% · mem 256.0 MiB / 7.7 GiB (3.3%) · net rx 1.2 MiB tx 300.0 KiB".
func (b *DockerEventBuilder) GetContainerStatsEvent(containerEntity *entities.ContainerEntity, containerStatsEntity *entities.ContainerStatsEntity, occurredAt time.Time) map[string]any {
    attributes := map[string]any{
        enums.DockerAttributeEnum.CONTAINER_ID:       containerEntity.GetShortId(),
        enums.DockerAttributeEnum.MEMORY_USAGE_BYTES: containerStatsEntity.GetMemoryUsageBytes(),
        enums.DockerAttributeEnum.MEMORY_LIMIT_BYTES: containerStatsEntity.GetMemoryLimitBytes(),
        enums.DockerAttributeEnum.MEMORY_PERCENT:     containerStatsEntity.GetMemoryPercent(),
        enums.DockerAttributeEnum.NET_RX_BYTES:       containerStatsEntity.GetNetRxBytes(),
        enums.DockerAttributeEnum.NET_TX_BYTES:       containerStatsEntity.GetNetTxBytes(),
    }
    cpuLabel := "cpu n/a"
    if containerStatsEntity.HasCpuPercent() {
        attributes[enums.DockerAttributeEnum.CPU_PERCENT] = containerStatsEntity.GetCpuPercent()
        cpuLabel = fmt.Sprintf("cpu %.1f%%", containerStatsEntity.GetCpuPercent())
    }
    message := fmt.Sprintf("%s · mem %s / %s (%.1f%%) · net rx %s tx %s", cpuLabel,
        b.getBytesLabel(containerStatsEntity.GetMemoryUsageBytes()), b.getBytesLabel(containerStatsEntity.GetMemoryLimitBytes()),
        containerStatsEntity.GetMemoryPercent(),
        b.getBytesLabel(containerStatsEntity.GetNetRxBytes()), b.getBytesLabel(containerStatsEntity.GetNetTxBytes()))

    return map[string]any{
        collectorsenums.CollectorsKeyEnum.EVENT_TYPE:  enums.DockerEventTypeEnum.DOCKER_STATS,
        collectorsenums.CollectorsKeyEnum.SOURCE_NAME: containerEntity.GetName(),
        collectorsenums.CollectorsKeyEnum.MESSAGE:     message,
        collectorsenums.CollectorsKeyEnum.OCCURRED_AT: occurredAt.UTC().Format(dateTimeFormat),
        ATTRIBUTES_KEY: attributes,
    }
}

// GetContainerLogEvent una linea del log; file_path lleva el stream (stdout / stderr). Gravedad por palabras clave.
func (b *DockerEventBuilder) GetContainerLogEvent(containerEntity *entities.ContainerEntity, containerLogLineEntity *entities.ContainerLogLineEntity) map[string]any {
    return map[string]any{
        collectorsenums.CollectorsKeyEnum.EVENT_TYPE:  enums.DockerEventTypeEnum.DOCKER_LOG,
        collectorsenums.CollectorsKeyEnum.SOURCE_NAME: containerEntity.GetName(),
        collectorsenums.CollectorsKeyEnum.FILE_PATH:   containerLogLineEntity.GetStream(),
        collectorsenums.CollectorsKeyEnum.MESSAGE:     b.getTruncated(containerLogLineEntity.GetText()),
        collectorsenums.CollectorsKeyEnum.OCCURRED_AT: containerLogLineEntity.GetOccurredAt().UTC().Format(dateTimeFormat),
    }
}

// GetDaemonEvent un evento de GET /events tal cual: "container die cont-web (exitCode=137) · image nginx".
func (b *DockerEventBuilder) GetDaemonEvent(rawDaemonEvent map[string]any) map[string]any {
    eventType := application.GetTrimmedString(rawDaemonEvent, "Type")
    action := application.GetTrimmedString(rawDaemonEvent, "Action")
    actor, _ := rawDaemonEvent["Actor"].(map[string]any)
    actorAttributes, _ := actor["Attributes"].(map[string]any)
    actorId := application.GetTrimmedString(actor, "ID")
    name := application.GetTrimmedString(actorAttributes, "name")
    image := application.GetTrimmedString(actorAttributes, "image")

    sourceName := DAEMON_SOURCE_NAME
    if eventType == "container" && name != "" {
        sourceName = name
    }
    attributes := map[string]any{
        enums.DockerAttributeEnum.TYPE:     eventType,
        enums.DockerAttributeEnum.ACTION:   action,
        enums.DockerAttributeEnum.ACTOR_ID: b.getShortId(actorId),
        enums.DockerAttributeEnum.NAME:     name,
        enums.DockerAttributeEnum.IMAGE:    image,
    }
    details := []string{}
    exitCode := 0
    if rawExitCode := application.GetTrimmedString(actorAttributes, "exitCode"); rawExitCode != "" {
        exitCode = application.ToInt(rawExitCode)
        attributes[enums.DockerAttributeEnum.EXIT_CODE] = exitCode
        details = append(details, "exitCode="+rawExitCode)
    }
    if signal := application.GetTrimmedString(actorAttributes, "signal"); signal != "" {
        attributes[enums.DockerAttributeEnum.SIGNAL] = signal
        details = append(details, "signal="+signal)
    }

    message := strings.TrimSpace(fmt.Sprintf("%s %s %s", eventType, action, b.getOr(name, b.getShortId(actorId))))
    if len(details) > 0 {
        message += " (" + strings.Join(details, ", ") + ")"
    }
    if image != "" {
        message += " · image " + image
    }

    occurredAt := time.Unix(int64(application.ToInt(rawDaemonEvent["time"])), 0)
    return map[string]any{
        collectorsenums.CollectorsKeyEnum.EVENT_TYPE:  enums.DockerEventTypeEnum.DOCKER_EVENT,
        collectorsenums.CollectorsKeyEnum.SOURCE_NAME: sourceName,
        collectorsenums.CollectorsKeyEnum.MESSAGE:     b.getTruncated(message),
        collectorsenums.CollectorsKeyEnum.EXIT_CODE:   exitCode,
        collectorsenums.CollectorsKeyEnum.OCCURRED_AT: occurredAt.UTC().Format(dateTimeFormat),
        ATTRIBUTES_KEY: attributes,
    }
}

// getBytesLabel en binario, como docker stats: 1.5 KiB, 256.0 MiB, 7.7 GiB.
func (b *DockerEventBuilder) getBytesLabel(bytes int64) string {
    units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
    value := float64(bytes)
    unitIndex := 0
    for value >= 1024 && unitIndex < len(units)-1 {
        value /= 1024
        unitIndex++
    }
    if unitIndex == 0 {
        return fmt.Sprintf("%d B", bytes)
    }
    return fmt.Sprintf("%.1f %s", value, units[unitIndex])
}

func (b *DockerEventBuilder) getShortId(id string) string {
    if len(id) <= entities.SHORT_ID_LENGTH {
        return id
    }
    return id[:entities.SHORT_ID_LENGTH]
}

func (b *DockerEventBuilder) getOrNone(value string) string {
    return b.getOr(value, "none")
}

func (b *DockerEventBuilder) getOr(value string, fallback string) string {
    if value == "" {
        return fallback
    }
    return value
}

func (b *DockerEventBuilder) getTruncated(text string) string {
    runes := []rune(text)
    if len(runes) <= MESSAGE_MAX_LENGTH {
        return text
    }
    return string(runes[:MESSAGE_MAX_LENGTH])
}
