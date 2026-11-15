package configuration

import (
    "os"
    "strconv"
    "strings"
    "sync"

    "siem-client/src/core/boot"
    "siem-client/src/modules/shared/domain/enums"
)

// Version del cliente: la muestran `version` y `help`, viaja al backend en cada envio (agent.client_version) y da
// nombre a los binarios de make dist. Al subirla, anotar el cambio en changelog.md.
const (
    APP_VERSION        = "0.3.0"
    APP_VERSION_UPDATE = "2026-11-15"
)

// EnvironmentReaderRawRepository fuente UNICA de lectura de entorno: un getter semantico por cada clave del .env
// que la app consume (claves en AppEnvKeyEnum).
type EnvironmentReaderRawRepository struct{}

var (
    instance     *EnvironmentReaderRawRepository
    instanceOnce sync.Once
)

func GetInstance() *EnvironmentReaderRawRepository {
    instanceOnce.Do(func() {
        instance = &EnvironmentReaderRawRepository{}
    })
    return instance
}

func (e *EnvironmentReaderRawRepository) GetAppVersion() string {
    return APP_VERSION
}

func (e *EnvironmentReaderRawRepository) GetAppVersionUpdate() string {
    return APP_VERSION_UPDATE
}

func (e *EnvironmentReaderRawRepository) GetAppName() string {
    return e.get(enums.AppEnvKeyEnum.APP_NAME)
}

func (e *EnvironmentReaderRawRepository) GetDefaultLang() string {
    return e.get(enums.AppEnvKeyEnum.APP_DEFAULT_LANG)
}

func (e *EnvironmentReaderRawRepository) GetLogPath() string {
    return e.get(enums.AppEnvKeyEnum.APP_LOG_PATH)
}

// GetSiemEndpointUrl la url completa del endpoint de ingesta del backend (…/api/v1/ingest-batches).
func (e *EnvironmentReaderRawRepository) GetSiemEndpointUrl() string {
    return e.get(enums.AppEnvKeyEnum.APP_SIEM_ENDPOINT_URL)
}

// GetSiemApiToken el api_token del agente. Nunca se imprime ni se loguea.
func (e *EnvironmentReaderRawRepository) GetSiemApiToken() string {
    return e.get(enums.AppEnvKeyEnum.APP_SIEM_API_TOKEN)
}

// IsRemoteConfigEnabled true (defecto) = el agente baja del backend lo que tiene que vigilar (monitores del panel) y
// siem-client.json pasa a ser opcional. "0" / "false" / "off" / "no" = solo el fichero local.
func (e *EnvironmentReaderRawRepository) IsRemoteConfigEnabled() bool {
    remoteConfig := strings.ToLower(e.get(enums.AppEnvKeyEnum.APP_SIEM_REMOTE_CONFIG))
    return remoteConfig != "0" && remoteConfig != "false" && remoteConfig != "off" && remoteConfig != "no"
}

// GetConfigPath el siem-client.json (opcional con la config remota): carpetas, comandos y docker locales.
func (e *EnvironmentReaderRawRepository) GetConfigPath() string {
    return e.get(enums.AppEnvKeyEnum.APP_CONFIG_PATH)
}

func (e *EnvironmentReaderRawRepository) GetSpoolPath() string {
    return e.get(enums.AppEnvKeyEnum.APP_SPOOL_PATH)
}

func (e *EnvironmentReaderRawRepository) GetStatePath() string {
    return e.get(enums.AppEnvKeyEnum.APP_STATE_PATH)
}

func (e *EnvironmentReaderRawRepository) GetEnvironment() string {
    return e.get(enums.AppEnvKeyEnum.APP_ENV)
}

func (e *EnvironmentReaderRawRepository) IsProduction() bool {
    return e.GetEnvironment() == boot.ENV.PRODUCTION
}

func (e *EnvironmentReaderRawRepository) IsLocal() bool {
    return e.GetEnvironment() == boot.ENV.LOCAL
}

// IsDebug filter_var(FILTER_VALIDATE_BOOLEAN): "1", "true", "on", "yes".
func (e *EnvironmentReaderRawRepository) IsDebug() bool {
    debug := strings.ToLower(e.get(enums.AppEnvKeyEnum.APP_DEBUG))
    if debug == "on" || debug == "yes" {
        return true
    }
    isDebug, _ := strconv.ParseBool(debug)
    return isDebug
}

func (e *EnvironmentReaderRawRepository) get(key string) string {
    return strings.TrimSpace(os.Getenv(key))
}
