package services

import (
    "regexp"
    "strings"

    "siem-backend/src/modules/monitors/domain/enums"
    "siem-backend/src/modules/shared/application"
)

const resourceKeySeparator = "|"

var windowsPathPattern = regexp.MustCompile(`^[A-Za-z]:[\\/]`)

// MonitorResourceKeyer la clave del recurso que vigila un monitor: dos monitores con la misma clave sobre el mismo
// agente vigilan lo mismo (repetido). cpu, memoria y gpu: el tipo (uno por agente) · disco: la ruta · red: la
// interfaz · carpeta: ruta + glob · log: ruta + glob + regex. Su gemelo en el front es MonitorDuplicateFinder.js:
// si cambia una regla aqui, cambia alli.
type MonitorResourceKeyer struct{}

func MonitorResourceKeyerGetInstance() *MonitorResourceKeyer {
    return &MonitorResourceKeyer{}
}

// GetResourceKey normalizedConfig: la de MonitorFieldsValidator.GetNormalizedConfig (con sus valores por defecto).
func (k *MonitorResourceKeyer) GetResourceKey(monitorType string, normalizedConfig map[string]any) string {
    keyParts := []string{monitorType}
    switch monitorType {
    case enums.MonitorTypeEnum.DISK:
        keyParts = append(keyParts, getComparablePath(normalizedConfig))
    case enums.MonitorTypeEnum.NETWORK:
        keyParts = append(keyParts, strings.ToLower(application.GetTrimmedString(normalizedConfig, enums.MonitorsKeyEnum.INTERFACE)))
    case enums.MonitorTypeEnum.FOLDER:
        keyParts = append(keyParts, getComparablePath(normalizedConfig), application.GetTrimmedString(normalizedConfig, enums.MonitorsKeyEnum.PATTERN))
    case enums.MonitorTypeEnum.LOG_PATTERN:
        keyParts = append(
            keyParts,
            getComparablePath(normalizedConfig),
            application.GetTrimmedString(normalizedConfig, enums.MonitorsKeyEnum.PATTERN),
            application.ToString(normalizedConfig[enums.MonitorsKeyEnum.MATCH]),
        )
    }
    return strings.Join(keyParts, resourceKeySeparator)
}

// getComparablePath sin separador final (salvo la raiz) y, en Windows, sin distinguir mayusculas: /var/log/ = /var/log
// y C:\Logs = c:\logs.
func getComparablePath(normalizedConfig map[string]any) string {
    comparablePath := application.GetTrimmedString(normalizedConfig, enums.MonitorsKeyEnum.PATH)
    isWindowsPath := windowsPathPattern.MatchString(comparablePath)
    // la raiz se queda como esta: "/" o "C:\"
    rootLength := 1
    if isWindowsPath {
        rootLength = 3
    }
    for len(comparablePath) > rootLength && strings.ContainsAny(comparablePath[len(comparablePath)-1:], `/\`) {
        comparablePath = comparablePath[:len(comparablePath)-1]
    }
    if isWindowsPath {
        return strings.ToLower(comparablePath)
    }
    return comparablePath
}
