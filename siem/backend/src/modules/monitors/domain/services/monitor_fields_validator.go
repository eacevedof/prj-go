package services

import (
    "path"
    "regexp"
    "slices"
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
    }
    return errors
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
