package configuration

import (
    "os"
    "strconv"
    "strings"
    "sync"

    "trafficker/src/core/boot"
    "trafficker/src/modules/shared/domain/enums"
)

// Version del agente: la muestran `version` y `help`, viaja al backend del siem en cada envio (agent.client_version)
// y da nombre a los binarios de make dist. Al subirla, anotar el cambio en changelog.md.
const (
    APP_VERSION        = "0.1.0"
    APP_VERSION_UPDATE = "2026-11-16"
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

// GetSiemEndpointUrl la url completa del endpoint de ingesta del backend del siem (…/api/v1/ingest-batches).
func (e *EnvironmentReaderRawRepository) GetSiemEndpointUrl() string {
    return e.get(enums.AppEnvKeyEnum.APP_SIEM_ENDPOINT_URL)
}

// GetSiemApiToken el api_token del agente. Nunca se imprime ni se loguea.
func (e *EnvironmentReaderRawRepository) GetSiemApiToken() string {
    return e.get(enums.AppEnvKeyEnum.APP_SIEM_API_TOKEN)
}

// IsLocalMode sin .env y sin token del siem en el entorno real: no hay a donde enviar. trafficker solo lista por
// consola (sin argumentos = connections-list) y resuelve el whois de las IP remotas el mismo, con su cache.
func (e *EnvironmentReaderRawRepository) IsLocalMode() bool {
    return !boot.IsEnvFileLoaded() && e.GetSiemApiToken() == ""
}

// GetCollectIntervalSeconds cada cuanto se hace la foto de los sockets (minimo 1).
func (e *EnvironmentReaderRawRepository) GetCollectIntervalSeconds() int {
    return e.getPositiveInt(enums.AppEnvKeyEnum.APP_COLLECT_INTERVAL_SECONDS, enums.AppEnvDefaultEnum.APP_COLLECT_INTERVAL_SECONDS)
}

// GetSendIntervalSeconds cada cuanto se envia el spool al backend (minimo 1).
func (e *EnvironmentReaderRawRepository) GetSendIntervalSeconds() int {
    return e.getPositiveInt(enums.AppEnvKeyEnum.APP_SEND_INTERVAL_SECONDS, enums.AppEnvDefaultEnum.APP_SEND_INTERVAL_SECONDS)
}

// IsLoopbackIncluded "1"/"true"/"on"/"yes" = tambien las conexiones de 127.0.0.0/8 y ::1.
func (e *EnvironmentReaderRawRepository) IsLoopbackIncluded() bool {
    return e.isTrue(enums.AppEnvKeyEnum.APP_INCLUDE_LOOPBACK)
}

// GetIgnoredProcesses nombres de proceso (comm) a ignorar, en minusculas: "chrome, sshd" -> [chrome sshd].
func (e *EnvironmentReaderRawRepository) GetIgnoredProcesses() []string {
    return e.getList(enums.AppEnvKeyEnum.APP_IGNORE_PROCESSES)
}

// GetIgnoredPorts puertos (locales o remotos) a ignorar: "22, 53" -> [22 53]. Lo que no es un puerto se descarta.
func (e *EnvironmentReaderRawRepository) GetIgnoredPorts() []int {
    ports := []int{}
    for _, rawPort := range e.getList(enums.AppEnvKeyEnum.APP_IGNORE_PORTS) {
        port, err := strconv.Atoi(rawPort)
        if err != nil || port < 1 || port > 65535 {
            continue
        }
        ports = append(ports, port)
    }
    return ports
}

func (e *EnvironmentReaderRawRepository) GetSpoolPath() string {
    return e.get(enums.AppEnvKeyEnum.APP_SPOOL_PATH)
}

func (e *EnvironmentReaderRawRepository) GetStatePath() string {
    return e.get(enums.AppEnvKeyEnum.APP_STATE_PATH)
}

// GetProcPath raiz de /proc (en los tests, una carpeta con ficheros de ejemplo).
func (e *EnvironmentReaderRawRepository) GetProcPath() string {
    return e.get(enums.AppEnvKeyEnum.APP_PROC_PATH)
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
    return e.isTrue(enums.AppEnvKeyEnum.APP_DEBUG)
}

func (e *EnvironmentReaderRawRepository) isTrue(key string) bool {
    value := strings.ToLower(e.get(key))
    if value == "on" || value == "yes" {
        return true
    }
    isTrue, _ := strconv.ParseBool(value)
    return isTrue
}

func (e *EnvironmentReaderRawRepository) getPositiveInt(key string, defaultValue string) int {
    value, err := strconv.Atoi(e.get(key))
    if err != nil || value < 1 {
        value, _ = strconv.Atoi(defaultValue)
    }
    return value
}

func (e *EnvironmentReaderRawRepository) getList(key string) []string {
    items := []string{}
    for _, item := range strings.Split(e.get(key), ",") {
        item = strings.ToLower(strings.TrimSpace(item))
        if item != "" {
            items = append(items, item)
        }
    }
    return items
}

func (e *EnvironmentReaderRawRepository) get(key string) string {
    return strings.TrimSpace(os.Getenv(key))
}
