package boot

// Si al arrancar se ha cargado un .env. Sin .env (y sin token del siem en el entorno real) trafficker va en modo
// local: no hay siem, lista por consola y resuelve el whois el mismo (ver EnvironmentReaderRawRepository.IsLocalMode).
var isEnvFileLoaded = false

func SetEnvFileLoaded() {
    isEnvFileLoaded = true
}

func IsEnvFileLoaded() bool {
    return isEnvFileLoaded
}
