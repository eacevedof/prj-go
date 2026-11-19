package configuration

import (
    "os"
    "strconv"
    "strings"
    "sync"

    "notifier-backend/src/core/boot"
    "notifier-backend/src/modules/shared/domain/enums"
)

// version publicada en /api/v1: es la del CONTRATO (la misma que task-php). Al subirla, anotar el
// cambio en changelog.md
const (
    APP_VERSION        = "0.1.0"
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

// GetSecretsKey la clave con la que SecretCipher cifra los secretos de los perfiles de envio (AES-256-GCM). Si
// cambia, los secretos guardados dejan de poder descifrarse: hay que volver a meterlos.
func (e *EnvironmentReaderRawRepository) GetSecretsKey() string {
    return e.get(enums.AppEnvKeyEnum.APP_SECRETS_KEY)
}

// GetQueueWorkerIntervalSeconds cada cuanto los workers de la API vacian su cola de messages. 0 = workers apagados
// (se procesa con ./console message-queue-process).
func (e *EnvironmentReaderRawRepository) GetQueueWorkerIntervalSeconds() int {
    return e.getIntOrDefault(enums.AppEnvKeyEnum.APP_QUEUE_WORKER_INTERVAL_SECONDS, enums.AppEnvDefaultEnum.QUEUE_WORKER_INTERVAL_SECONDS)
}

// GetMessageMaxRecipients tope de destinatarios (to + cc + bcc) de un mensaje.
func (e *EnvironmentReaderRawRepository) GetMessageMaxRecipients() int {
    return e.getPositiveIntOrDefault(enums.AppEnvKeyEnum.APP_MESSAGE_MAX_RECIPIENTS, enums.AppEnvDefaultEnum.MESSAGE_MAX_RECIPIENTS)
}

// GetAttachmentAllowedHosts hosts (separados por comas) de los que el worker puede descargar adjuntos por url. Vacio =
// ningun adjunto por url se acepta.
func (e *EnvironmentReaderRawRepository) GetAttachmentAllowedHosts() []string {
    allowedHosts := []string{}
    for _, allowedHost := range strings.Split(e.get(enums.AppEnvKeyEnum.APP_ATTACHMENT_ALLOWED_HOSTS), ",") {
        if allowedHost = strings.ToLower(strings.TrimSpace(allowedHost)); allowedHost != "" {
            allowedHosts = append(allowedHosts, allowedHost)
        }
    }
    return allowedHosts
}

// GetAttachmentMaxBytes tope de tamano de cada adjunto (texto o descargado).
func (e *EnvironmentReaderRawRepository) GetAttachmentMaxBytes() int {
    return e.getPositiveIntOrDefault(enums.AppEnvKeyEnum.APP_ATTACHMENT_MAX_BYTES, enums.AppEnvDefaultEnum.ATTACHMENT_MAX_BYTES)
}

// GetAttachmentTimeoutSeconds tiempo maximo para descargar un adjunto por url.
func (e *EnvironmentReaderRawRepository) GetAttachmentTimeoutSeconds() int {
    return e.getPositiveIntOrDefault(enums.AppEnvKeyEnum.APP_ATTACHMENT_TIMEOUT_SECONDS, enums.AppEnvDefaultEnum.ATTACHMENT_TIMEOUT_SECONDS)
}

func (e *EnvironmentReaderRawRepository) getPositiveIntOrDefault(key string, defaultValue int) int {
    value := e.getIntOrDefault(key, defaultValue)
    if value == 0 {
        return defaultValue
    }
    return value
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
