package boot

import "time"

// Equivalente a php_sapi_name(): "server" para la API y "cli" para la consola. Lo fija el main de cada
// binario al arrancar (cmd/console llama a SetSapiNameCli) y es de solo lectura despues.
const (
    SAPI_SERVER = "server"
    SAPI_CLI    = "cli"
)

var sapiName = SAPI_SERVER

func SetSapiNameCli() {
    sapiName = SAPI_CLI
}

func GetSapiName() string {
    return sapiName
}

func IsCli() bool {
    return sapiName == SAPI_CLI
}

// arranque del proceso (UTC): lo lee EnvironmentReaderRawRepository.GetAppStartedAt para el uptime del dashboard
var appStartedAt = time.Now().UTC().Format("2006-01-02 15:04:05")

func GetAppStartedAt() string {
    return appStartedAt
}
