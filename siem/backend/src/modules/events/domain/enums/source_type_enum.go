package enums

// SourceTypeEnum de donde sale un evento: un watcher de ficheros, un comando o un contenedor Docker (su estado, sus
// recursos o su log), una metrica o, en network, un proceso con conexiones de red. Mismos valores que el CHECK de sources.source_type.
var SourceTypeEnum = struct {
    FILE         string
    COMMAND      string
    DOCKER       string
    DOCKER_STATS string
    DOCKER_LOG   string
    METRIC       string
    NETWORK      string
    VALUES       []string
}{
    FILE:         "file",
    COMMAND:      "command",
    DOCKER:       "docker",
    DOCKER_STATS: "docker_stats",
    DOCKER_LOG:   "docker_log",
    METRIC:       "metric",
    NETWORK:      "network",
    VALUES:       []string{"file", "command", "docker", "docker_stats", "docker_log", "metric", "network"},
}
