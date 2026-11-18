package services

import (
    "path"
    "regexp"
    "slices"
    "strconv"
    "strings"

    "siem-backend/src/core/config"

    "siem-backend/src/modules/shared/application"

    eventsenums "siem-backend/src/modules/events/domain/enums"

    "siem-backend/src/modules/monitors/domain/enums"
)

var (
    monitorNamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,100}$`)
    interfacePattern   = regexp.MustCompile(`^[A-Za-z0-9._:-]{0,50}$`)
    // /var/log (Linux, macOS) o C:\logs (Windows)
    absolutePathPattern = regexp.MustCompile(`^(/|[A-Za-z]:[\\/])`)
    // nombre de un proceso (comm de /proc/<pid>/comm): sin espacios ni comas (la coma separa la lista en el panel)
    processNamePattern = regexp.MustCompile(`^[^\s,]+$`)
)

// MonitorFieldsValidator reglas de un monitor: nombre, tipo, intervalo, umbrales (solo en los tipos con valor
// numerico) y la config propia de cada tipo. Devuelve la pila de errores (vacia = todo bien).
type MonitorFieldsValidator struct{}

func MonitorFieldsValidatorGetInstance() *MonitorFieldsValidator {
    return &MonitorFieldsValidator{}
}

func (v *MonitorFieldsValidator) GetErrors(monitorName string, monitorType string, intervalSeconds int, monitorConfig map[string]any, warningThreshold *float64, criticalThreshold *float64) []string {
    errors := []string{}
    if !monitorNamePattern.MatchString(monitorName) {
        errors = append(errors, enums.MonitorsMessageEnum.MONITOR_NAME_IS_NOT_VALID)
    }
    if !slices.Contains(enums.MonitorTypeEnum.VALUES, monitorType) {
        return append(errors, enums.MonitorsMessageEnum.MONITOR_TYPE_IS_NOT_VALID)
    }
    if intervalSeconds < enums.MonitorFieldsLimitEnum.INTERVAL_SECONDS_MIN || intervalSeconds > enums.MonitorFieldsLimitEnum.INTERVAL_SECONDS_MAX {
        errors = append(errors, enums.MonitorsMessageEnum.INTERVAL_SECONDS_MUST_BE_10_TO_86400)
    }
    errors = append(errors, v.getThresholdsErrors(monitorType, warningThreshold, criticalThreshold)...)
    return append(errors, v.getConfigErrors(monitorType, v.GetNormalizedConfig(monitorType, monitorConfig))...)
}

// GetNormalizedConfig solo las claves de ese tipo, con sus valores por defecto: lo que se guarda y lo que recibe el
// agente (una clave de mas en la peticion se descarta).
func (v *MonitorFieldsValidator) GetNormalizedConfig(monitorType string, monitorConfig map[string]any) map[string]any {
    switch monitorType {
    case enums.MonitorTypeEnum.DISK:
        return map[string]any{
            enums.MonitorsKeyEnum.PATH: v.getStringOrDefault(monitorConfig, enums.MonitorsKeyEnum.PATH, enums.MonitorFieldsLimitEnum.DEFAULT_DISK_PATH),
        }
    case enums.MonitorTypeEnum.NETWORK:
        return map[string]any{
            enums.MonitorsKeyEnum.INTERFACE: application.GetTrimmedString(monitorConfig, enums.MonitorsKeyEnum.INTERFACE),
        }
    case enums.MonitorTypeEnum.FOLDER, enums.MonitorTypeEnum.LOG_PATTERN:
        normalizedConfig := map[string]any{
            enums.MonitorsKeyEnum.PATH:         application.GetTrimmedString(monitorConfig, enums.MonitorsKeyEnum.PATH),
            enums.MonitorsKeyEnum.PATTERN:      v.getStringOrDefault(monitorConfig, enums.MonitorsKeyEnum.PATTERN, enums.MonitorFieldsLimitEnum.DEFAULT_FILES_PATTERN),
            enums.MonitorsKeyEnum.IS_RECURSIVE: application.GetInt(monitorConfig, enums.MonitorsKeyEnum.IS_RECURSIVE) == 1,
        }
        if monitorType == enums.MonitorTypeEnum.LOG_PATTERN {
            // el regex no se recorta: un espacio al principio o al final puede ser parte del patron
            normalizedConfig[enums.MonitorsKeyEnum.MATCH] = application.ToString(monitorConfig[enums.MonitorsKeyEnum.MATCH])
            normalizedConfig[enums.MonitorsKeyEnum.SEVERITY] = strings.ToLower(application.GetTrimmedString(monitorConfig, enums.MonitorsKeyEnum.SEVERITY))
        }
        return normalizedConfig
    case enums.MonitorTypeEnum.CONNECTIONS:
        return map[string]any{
            enums.MonitorsKeyEnum.IS_LOOPBACK_INCLUDED: application.GetInt(monitorConfig, enums.MonitorsKeyEnum.IS_LOOPBACK_INCLUDED) == 1,
            enums.MonitorsKeyEnum.IGNORE_PROCESSES:     v.getNormalizedProcessNames(monitorConfig[enums.MonitorsKeyEnum.IGNORE_PROCESSES]),
            enums.MonitorsKeyEnum.IGNORE_PORTS:         v.getNormalizedPorts(monitorConfig[enums.MonitorsKeyEnum.IGNORE_PORTS]),
        }
    }
    return map[string]any{}
}

// IsMetricType true si el tipo manda un valor numerico (admite umbrales).
func (v *MonitorFieldsValidator) IsMetricType(monitorType string) bool {
    return slices.Contains(enums.MonitorTypeEnum.METRIC_VALUES, monitorType)
}

func (v *MonitorFieldsValidator) getThresholdsErrors(monitorType string, warningThreshold *float64, criticalThreshold *float64) []string {
    errors := []string{}
    if !v.IsMetricType(monitorType) {
        return errors
    }
    isPercent := slices.Contains(enums.MonitorTypeEnum.PERCENT_VALUES, monitorType)
    for _, threshold := range []struct {
        key   string
        value *float64
    }{
        {enums.MonitorsKeyEnum.WARNING_THRESHOLD, warningThreshold},
        {enums.MonitorsKeyEnum.CRITICAL_THRESHOLD, criticalThreshold},
    } {
        if threshold.value == nil {
            continue
        }
        if isPercent && (*threshold.value < enums.MonitorFieldsLimitEnum.PERCENT_THRESHOLD_MIN || *threshold.value > enums.MonitorFieldsLimitEnum.PERCENT_THRESHOLD_MAX) {
            errors = append(errors, config.Tr(enums.MonitorsMessageEnum.X_MUST_BE_BETWEEN_0_AND_100, threshold.key))
        }
        if !isPercent && *threshold.value <= 0 {
            errors = append(errors, config.Tr(enums.MonitorsMessageEnum.X_MUST_BE_GREATER_THAN_0, threshold.key))
        }
    }
    if warningThreshold != nil && criticalThreshold != nil && *warningThreshold >= *criticalThreshold {
        errors = append(errors, enums.MonitorsMessageEnum.WARNING_MUST_BE_LOWER_THAN_CRITICAL)
    }
    return errors
}

func (v *MonitorFieldsValidator) getConfigErrors(monitorType string, normalizedConfig map[string]any) []string {
    errors := []string{}
    switch monitorType {
    case enums.MonitorTypeEnum.DISK:
        if !absolutePathPattern.MatchString(application.GetTrimmedString(normalizedConfig, enums.MonitorsKeyEnum.PATH)) {
            errors = append(errors, config.Tr(enums.MonitorsMessageEnum.CONFIG_X_MUST_BE_AN_ABSOLUTE_PATH, enums.MonitorsKeyEnum.PATH))
        }
    case enums.MonitorTypeEnum.NETWORK:
        if !interfacePattern.MatchString(application.GetTrimmedString(normalizedConfig, enums.MonitorsKeyEnum.INTERFACE)) {
            errors = append(errors, enums.MonitorsMessageEnum.CONFIG_INTERFACE_IS_NOT_VALID)
        }
    case enums.MonitorTypeEnum.FOLDER, enums.MonitorTypeEnum.LOG_PATTERN:
        if !absolutePathPattern.MatchString(application.GetTrimmedString(normalizedConfig, enums.MonitorsKeyEnum.PATH)) {
            errors = append(errors, config.Tr(enums.MonitorsMessageEnum.CONFIG_X_MUST_BE_AN_ABSOLUTE_PATH, enums.MonitorsKeyEnum.PATH))
        }
        if _, err := path.Match(application.GetTrimmedString(normalizedConfig, enums.MonitorsKeyEnum.PATTERN), enums.MonitorFieldsLimitEnum.GLOB_PROBE_FILE_NAME); err != nil {
            errors = append(errors, enums.MonitorsMessageEnum.CONFIG_PATTERN_IS_NOT_A_VALID_GLOB)
        }
        if monitorType == enums.MonitorTypeEnum.LOG_PATTERN {
            errors = append(errors, v.getLogPatternErrors(normalizedConfig)...)
        }
    case enums.MonitorTypeEnum.CONNECTIONS:
        errors = append(errors, v.getConnectionsErrors(normalizedConfig)...)
    }
    return errors
}

// getConnectionsErrors las listas ya vienen normalizadas (sin vacios ni repetidos): aqui solo tamano y forma.
func (v *MonitorFieldsValidator) getConnectionsErrors(normalizedConfig map[string]any) []string {
    errors := []string{}
    processNames, _ := normalizedConfig[enums.MonitorsKeyEnum.IGNORE_PROCESSES].([]string)
    isProcessNameInvalid := slices.ContainsFunc(processNames, func(processName string) bool {
        return len(processName) > enums.MonitorFieldsLimitEnum.PROCESS_NAME_MAX_LENGTH || !processNamePattern.MatchString(processName)
    })
    if len(processNames) > enums.MonitorFieldsLimitEnum.IGNORE_LIST_MAX_ITEMS || isProcessNameInvalid {
        errors = append(errors, enums.MonitorsMessageEnum.CONFIG_IGNORE_PROCESSES_IS_NOT_VALID)
    }
    ports, _ := normalizedConfig[enums.MonitorsKeyEnum.IGNORE_PORTS].([]int)
    isPortInvalid := slices.ContainsFunc(ports, func(port int) bool {
        return port < enums.MonitorFieldsLimitEnum.PORT_MIN || port > enums.MonitorFieldsLimitEnum.PORT_MAX
    })
    if len(ports) > enums.MonitorFieldsLimitEnum.IGNORE_LIST_MAX_ITEMS || isPortInvalid {
        errors = append(errors, enums.MonitorsMessageEnum.CONFIG_IGNORE_PORTS_IS_NOT_VALID)
    }
    return errors
}

// getNormalizedProcessNames de una lista JSON o de un texto separado por comas: recortados, sin vacios ni repetidos y
// en el orden en que llegan (el nombre de un proceso distingue mayusculas).
func (v *MonitorFieldsValidator) getNormalizedProcessNames(rawList any) []string {
    processNames := []string{}
    for _, rawItem := range v.getListItems(rawList) {
        processName := strings.TrimSpace(application.ToString(rawItem))
        if processName != "" && !slices.Contains(processNames, processName) {
            processNames = append(processNames, processName)
        }
    }
    return processNames
}

// getNormalizedPorts de una lista JSON o de un texto separado por comas: sin repetidos y ordenados. Lo que no es un
// entero pasa a 0, que no es un puerto valido: asi el error lo ve getConnectionsErrors en vez de perderse.
func (v *MonitorFieldsValidator) getNormalizedPorts(rawList any) []int {
    ports := []int{}
    for _, rawItem := range v.getListItems(rawList) {
        portText := strings.TrimSpace(application.ToString(rawItem))
        if portText == "" {
            continue
        }
        port, err := strconv.Atoi(portText)
        if err != nil {
            port = 0
        }
        if !slices.Contains(ports, port) {
            ports = append(ports, port)
        }
    }
    slices.Sort(ports)
    return ports
}

// getListItems []any (JSON), []string o un texto con comas; cualquier otra cosa, lista vacia.
func (v *MonitorFieldsValidator) getListItems(rawList any) []any {
    switch typedList := rawList.(type) {
    case []any:
        return typedList
    case []string:
        items := make([]any, 0, len(typedList))
        for _, item := range typedList {
            items = append(items, item)
        }
        return items
    case []int:
        items := make([]any, 0, len(typedList))
        for _, item := range typedList {
            items = append(items, item)
        }
        return items
    case string:
        items := []any{}
        for _, item := range strings.Split(typedList, enums.MonitorFieldsLimitEnum.LIST_SEPARATOR) {
            items = append(items, item)
        }
        return items
    }
    return []any{}
}

// getLogPatternErrors el regex lo compila el agente con el paquete regexp de Go: se valida con el mismo.
func (v *MonitorFieldsValidator) getLogPatternErrors(normalizedConfig map[string]any) []string {
    errors := []string{}
    match := application.ToString(normalizedConfig[enums.MonitorsKeyEnum.MATCH])
    if strings.TrimSpace(match) == "" {
        errors = append(errors, enums.MonitorsMessageEnum.CONFIG_MATCH_IS_REQUIRED)
    } else if _, err := regexp.Compile(match); err != nil {
        errors = append(errors, enums.MonitorsMessageEnum.CONFIG_MATCH_IS_NOT_A_VALID_REGEX)
    }
    severity := application.GetTrimmedString(normalizedConfig, enums.MonitorsKeyEnum.SEVERITY)
    if severity != "" && !slices.Contains(eventsenums.EventSeverityEnum.VALUES, severity) {
        errors = append(errors, enums.MonitorsMessageEnum.CONFIG_SEVERITY_IS_NOT_VALID)
    }
    return errors
}

func (v *MonitorFieldsValidator) getStringOrDefault(primitives map[string]any, key string, defaultValue string) string {
    if value := application.GetTrimmedString(primitives, key); value != "" {
        return value
    }
    return defaultValue
}
