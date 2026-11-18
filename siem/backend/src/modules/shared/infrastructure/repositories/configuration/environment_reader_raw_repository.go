package configuration

import (
    "os"
    "strconv"
    "strings"
    "sync"

    "siem-backend/src/core/boot"
    "siem-backend/src/modules/shared/domain/enums"
)

// version publicada en /api/v1: es la del CONTRATO (la misma que task-php). Al subirla, anotar el
// cambio en changelog.md
const (
    APP_VERSION        = "0.11.0"
    APP_VERSION_UPDATE = "2026-11-18"
)

// EnvironmentReaderRawRepository fuente UNICA de lectura de entorno: un getter semantico por cada clave del .env
// que la app consume (claves en AppEnvKeyEnum). El .env lo carga AppbootLoadenv() (core/config) en el entorno del
// proceso antes del primer GetInstance(); las variables REALES del contenedor/OS ganan sobre el fichero.
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

func (e *EnvironmentReaderRawRepository) GetBaseUrl() string {
    return strings.TrimRight(e.get(enums.AppEnvKeyEnum.APP_BASE_URL), "/")
}

func (e *EnvironmentReaderRawRepository) GetEmailTo() string {
    return e.get(enums.AppEnvKeyEnum.APP_EMAIL_TO)
}

func (e *EnvironmentReaderRawRepository) GetEmailFrom() string {
    return e.get(enums.AppEnvKeyEnum.APP_EMAIL_FROM)
}

// GetSendmailPath el sendmail_path del php.ini de task-php: el comando al que FuncMailer entrega el email.
func (e *EnvironmentReaderRawRepository) GetSendmailPath() string {
    return e.get(enums.AppEnvKeyEnum.APP_SENDMAIL_PATH)
}

func (e *EnvironmentReaderRawRepository) GetLogPath() string {
    return e.get(enums.AppEnvKeyEnum.APP_LOG_PATH)
}

func (e *EnvironmentReaderRawRepository) GetAppName() string {
    return e.get(enums.AppEnvKeyEnum.APP_NAME)
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

func (e *EnvironmentReaderRawRepository) IsDev() bool {
    return e.GetEnvironment() == boot.ENV.DEVELOP
}

func (e *EnvironmentReaderRawRepository) IsTest() bool {
    return e.GetEnvironment() == boot.ENV.TEST
}

// IsDebug filter_var(FILTER_VALIDATE_BOOLEAN): "1", "true", "on", "yes".
func (e *EnvironmentReaderRawRepository) IsDebug() bool {
    rawDebugValue := strings.ToLower(e.get(enums.AppEnvKeyEnum.APP_DEBUG))
    if rawDebugValue == "on" || rawDebugValue == "yes" {
        return true
    }
    isDebug, _ := strconv.ParseBool(rawDebugValue)
    return isDebug
}

func (e *EnvironmentReaderRawRepository) GetColor() string {
    return boot.ENV.COLORS[e.GetEnvironment()]
}

func (e *EnvironmentReaderRawRepository) GetSqlitePath() string {
    return e.get(enums.AppEnvKeyEnum.APP_SQLITE_PATH)
}

func (e *EnvironmentReaderRawRepository) GetDefaultLang() string {
    return e.get(enums.AppEnvKeyEnum.APP_DEFAULT_LANG)
}

func (e *EnvironmentReaderRawRepository) GetDomain() string {
    return e.get(enums.AppEnvKeyEnum.APP_DOMAIN)
}

func (e *EnvironmentReaderRawRepository) get(key string) string {
    return strings.TrimSpace(os.Getenv(key))
}

// GetAppStartedAt arranque del proceso en UTC ("Y-m-d H:i:s"): el dashboard calcula con el el uptime.
func (e *EnvironmentReaderRawRepository) GetAppStartedAt() string {
    return boot.GetAppStartedAt()
}

// GetDashboardToken el Bearer que exige el front (dashboard). Vacio = el dashboard no deja entrar a nadie.
func (e *EnvironmentReaderRawRepository) GetDashboardToken() string {
    return e.get(enums.AppEnvKeyEnum.APP_DASHBOARD_TOKEN)
}

// GetQueueWorkerIntervalSeconds cada cuanto el worker de la API vacia ingest_queue. 0 = worker apagado (se procesa
// con ./console ingest-queue-process).
func (e *EnvironmentReaderRawRepository) GetQueueWorkerIntervalSeconds() int {
    return e.getIntOrDefault(enums.AppEnvKeyEnum.APP_QUEUE_WORKER_INTERVAL_SECONDS, enums.AppEnvDefaultEnum.QUEUE_WORKER_INTERVAL_SECONDS)
}

// GetEventsRetentionDays dias que se conservan los eventos (./console events-purge sin --retention_days).
func (e *EnvironmentReaderRawRepository) GetEventsRetentionDays() int {
    return e.getIntOrDefault(enums.AppEnvKeyEnum.APP_EVENTS_RETENTION_DAYS, enums.AppEnvDefaultEnum.EVENTS_RETENTION_DAYS)
}

func (e *EnvironmentReaderRawRepository) getIntOrDefault(key string, defaultValue int) int {
    rawValue := e.get(key)
    if rawValue == "" {
        return defaultValue
    }
    value, err := strconv.Atoi(rawValue)
    if err != nil || value < 0 {
        return defaultValue
    }
    return value
}
