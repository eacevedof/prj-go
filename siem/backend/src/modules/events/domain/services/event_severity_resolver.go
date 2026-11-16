package services

import (
    "regexp"
    "slices"
    "strings"

    "siem-backend/src/modules/shared/application"

    "siem-backend/src/modules/events/domain/enums"
)

// de mas grave a menos: gana el primero que case. \b: "terror" o "panicked_at" no cuentan como palabra suelta.
var (
    criticalPattern = regexp.MustCompile(`(?i)\b(critical|crit|fatal|panic|emerg|emergency|out of memory|segfault)\b`)
    errorPattern    = regexp.MustCompile(`(?i)\b(error|err|exception|failed|failure|denied|refused|unauthorized)\b`)
    warningPattern  = regexp.MustCompile(`(?i)\b(warning|warn|deprecated|timeout|timed out|retry|retrying)\b`)
    debugPattern    = regexp.MustCompile(`(?i)\b(debug|trace)\b`)
)

// EventSeverityResolver gravedad de un evento cuando el cliente no la declara (la regla vive solo aqui):
//   - un comando: su codigo de salida (distinto de 0 = error).
//   - una linea de log (de fichero o de un contenedor): la palabra clave mas grave que contenga; sin ninguna, info.
//   - Docker: los attributes del evento (OOM, codigo de salida, health, % de cpu y memoria, accion del daemon).
type EventSeverityResolver struct{}

func EventSeverityResolverGetInstance() *EventSeverityResolver {
    return &EventSeverityResolver{}
}

// GetSeverity attributes puede ser nil (los eventos que no son de Docker no lo traen).
func (r *EventSeverityResolver) GetSeverity(eventType string, message string, exitCode int, declaredSeverity string, attributes map[string]any) string {
    declaredSeverity = strings.ToLower(strings.TrimSpace(declaredSeverity))
    if slices.Contains(enums.EventSeverityEnum.VALUES, declaredSeverity) {
        return declaredSeverity
    }

    switch eventType {
    case enums.EventTypeEnum.COMMAND_RESULT:
        if exitCode != 0 {
            return enums.EventSeverityEnum.ERROR
        }
        return enums.EventSeverityEnum.INFO
    case enums.EventTypeEnum.DOCKER_STATE:
        return r.getDockerStateSeverity(attributes)
    case enums.EventTypeEnum.DOCKER_STATS:
        return r.getDockerStatsSeverity(attributes)
    case enums.EventTypeEnum.DOCKER_EVENT:
        return r.getDockerEventSeverity(attributes)
    case enums.EventTypeEnum.HOST_METRIC:
        return r.getHostMetricSeverity(attributes)
    }
    return r.getSeverityByKeywords(message)
}

// getHostMetricSeverity una muestra de un monitor del panel (cpu, memoria, disco, gpu, red) contra los umbrales de ESE
// monitor (los pone el backend en attributes al procesar el lote; sin umbral, esa regla no aplica). Metrica no
// disponible en el host, warning. Por debajo de los umbrales, debug: una muestra por minuto no infla la portada.
func (r *EventSeverityResolver) getHostMetricSeverity(attributes map[string]any) string {
    if application.ToInt(attributes[enums.MonitorAttributeEnum.IS_UNAVAILABLE]) == 1 {
        return enums.EventSeverityEnum.WARNING
    }
    value := r.getFloat(attributes, enums.MonitorAttributeEnum.VALUE)
    if criticalThreshold, hasThreshold := attributes[enums.MonitorAttributeEnum.CRITICAL_THRESHOLD].(float64); hasThreshold && value >= criticalThreshold {
        return enums.EventSeverityEnum.CRITICAL
    }
    if warningThreshold, hasThreshold := attributes[enums.MonitorAttributeEnum.WARNING_THRESHOLD].(float64); hasThreshold && value >= warningThreshold {
        return enums.EventSeverityEnum.WARNING
    }
    return enums.EventSeverityEnum.DEBUG
}

// getDockerStateSeverity un contenedor parado con exit 0 es warning (se paro, pero sin fallo); con otro codigo, error;
// muerto o matado por falta de memoria, critical.
func (r *EventSeverityResolver) getDockerStateSeverity(attributes map[string]any) string {
    status := application.GetTrimmedString(attributes, enums.DockerAttributeEnum.STATUS)
    switch {
    case application.ToInt(attributes[enums.DockerAttributeEnum.IS_OOM_KILLED]) == 1 || status == enums.DockerAttributeValueEnum.STATUS_DEAD:
        return enums.EventSeverityEnum.CRITICAL
    case status == enums.DockerAttributeValueEnum.STATUS_EXITED && application.GetInt(attributes, enums.DockerAttributeEnum.EXIT_CODE) != 0:
        return enums.EventSeverityEnum.ERROR
    case application.GetTrimmedString(attributes, enums.DockerAttributeEnum.HEALTH) == enums.DockerAttributeValueEnum.HEALTH_UNHEALTHY:
        return enums.EventSeverityEnum.ERROR
    case status == enums.DockerAttributeValueEnum.STATUS_EXITED || status == enums.DockerAttributeValueEnum.STATUS_RESTARTING:
        return enums.EventSeverityEnum.WARNING
    case application.GetTrimmedString(attributes, enums.DockerAttributeEnum.CHANGE_KIND) == enums.DockerAttributeValueEnum.CHANGE_RESTARTED:
        return enums.EventSeverityEnum.WARNING
    }
    return enums.EventSeverityEnum.INFO
}

// getDockerStatsSeverity por debajo de los umbrales es debug: una muestra cada minuto por contenedor no debe inflar
// los contadores de la portada.
func (r *EventSeverityResolver) getDockerStatsSeverity(attributes map[string]any) string {
    memoryPercent := r.getFloat(attributes, enums.DockerAttributeEnum.MEMORY_PERCENT)
    cpuPercent := r.getFloat(attributes, enums.DockerAttributeEnum.CPU_PERCENT)
    switch {
    case memoryPercent >= enums.DockerThresholdEnum.MEMORY_ERROR_PERCENT:
        return enums.EventSeverityEnum.ERROR
    case memoryPercent >= enums.DockerThresholdEnum.MEMORY_WARNING_PERCENT:
        return enums.EventSeverityEnum.WARNING
    case cpuPercent >= enums.DockerThresholdEnum.CPU_WARNING_PERCENT:
        return enums.EventSeverityEnum.WARNING
    }
    return enums.EventSeverityEnum.DEBUG
}

// getDockerEventSeverity eventos del daemon: oom, critical; die con codigo distinto de 0 o health unhealthy, error;
// kill, warning; el resto (start, create, pull...) info.
func (r *EventSeverityResolver) getDockerEventSeverity(attributes map[string]any) string {
    action := application.GetTrimmedString(attributes, enums.DockerAttributeEnum.ACTION)
    switch {
    case action == enums.DockerAttributeValueEnum.ACTION_OOM:
        return enums.EventSeverityEnum.CRITICAL
    case action == enums.DockerAttributeValueEnum.ACTION_DIE && application.GetInt(attributes, enums.DockerAttributeEnum.EXIT_CODE) != 0:
        return enums.EventSeverityEnum.ERROR
    case strings.HasPrefix(action, enums.DockerAttributeValueEnum.ACTION_HEALTH_STATUS_PREFIX) && strings.Contains(action, enums.DockerAttributeValueEnum.HEALTH_UNHEALTHY):
        return enums.EventSeverityEnum.ERROR
    case action == enums.DockerAttributeValueEnum.ACTION_KILL:
        return enums.EventSeverityEnum.WARNING
    }
    return enums.EventSeverityEnum.INFO
}

func (r *EventSeverityResolver) getSeverityByKeywords(message string) string {
    switch {
    case criticalPattern.MatchString(message):
        return enums.EventSeverityEnum.CRITICAL
    case errorPattern.MatchString(message):
        return enums.EventSeverityEnum.ERROR
    case warningPattern.MatchString(message):
        return enums.EventSeverityEnum.WARNING
    case debugPattern.MatchString(message):
        return enums.EventSeverityEnum.DEBUG
    }
    return enums.EventSeverityEnum.INFO
}

// getFloat los numeros del JSON llegan como float64; un valor ausente o no numerico cuenta como 0.
func (r *EventSeverityResolver) getFloat(primitives map[string]any, key string) float64 {
    switch value := primitives[key].(type) {
    case float64:
        return value
    case int:
        return float64(value)
    }
    return 0
}
