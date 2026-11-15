package entities

import "path"

const (
    DEFAULT_DOCKER_SOCKET_PATH                  = "/var/run/docker.sock"
    DEFAULT_DOCKER_STATS_INTERVAL_SECONDS       = 60
    DEFAULT_DOCKER_LOGS_MAX_LINES_PER_CONTAINER = 200
    DEFAULT_DOCKER_LOGS_INCLUDE_PATTERN         = "*"
)

// DockerConfigPrimitives el bloque "docker" de siem-client.json. Sin bloque o con is_enabled false no se vigila
// Docker. Con el bloque activo, estados y health van siempre; recursos, logs y eventos del daemon se pueden apagar
// uno a uno (un is_*_enabled ausente = true).
type DockerConfigPrimitives struct {
    IsEnabled                 bool     `json:"is_enabled"`
    SocketPath                string   `json:"socket_path"`
    IsStatsEnabled            *bool    `json:"is_stats_enabled"`
    StatsIntervalSeconds      int      `json:"stats_interval_seconds"`
    IsEventsEnabled           *bool    `json:"is_events_enabled"`
    IsLogsEnabled             *bool    `json:"is_logs_enabled"`
    LogsMaxLinesPerContainer  int      `json:"logs_max_lines_per_container"`
    LogsIncludePatterns       []string `json:"logs_include"`
    LogsExcludePatterns       []string `json:"logs_exclude"`
}

// DockerConfigEntity que se vigila de Docker y cada cuanto. Los patrones de logs son globs sobre el nombre del
// contenedor (cont-web-*): un contenedor manda sus logs si casa con algun include y con ningun exclude.
type DockerConfigEntity struct {
    isEnabled                bool
    socketPath               string
    isStatsEnabled           bool
    statsIntervalSeconds     int
    isEventsEnabled          bool
    isLogsEnabled            bool
    logsMaxLinesPerContainer int
    logsIncludePatterns      []string
    logsExcludePatterns      []string
}

func DockerConfigEntityFromPrimitives(dockerConfigPrimitives DockerConfigPrimitives) *DockerConfigEntity {
    logsIncludePatterns := getTrimmedStrings(dockerConfigPrimitives.LogsIncludePatterns)
    if len(logsIncludePatterns) == 0 {
        logsIncludePatterns = []string{DEFAULT_DOCKER_LOGS_INCLUDE_PATTERN}
    }
    return &DockerConfigEntity{
        isEnabled:                dockerConfigPrimitives.IsEnabled,
        socketPath:               getStringOrDefault(dockerConfigPrimitives.SocketPath, DEFAULT_DOCKER_SOCKET_PATH),
        isStatsEnabled:           getBoolOrTrue(dockerConfigPrimitives.IsStatsEnabled),
        statsIntervalSeconds:     getIntOrDefault(dockerConfigPrimitives.StatsIntervalSeconds, DEFAULT_DOCKER_STATS_INTERVAL_SECONDS),
        isEventsEnabled:          getBoolOrTrue(dockerConfigPrimitives.IsEventsEnabled),
        isLogsEnabled:            getBoolOrTrue(dockerConfigPrimitives.IsLogsEnabled),
        logsMaxLinesPerContainer: getIntOrDefault(dockerConfigPrimitives.LogsMaxLinesPerContainer, DEFAULT_DOCKER_LOGS_MAX_LINES_PER_CONTAINER),
        logsIncludePatterns:      logsIncludePatterns,
        logsExcludePatterns:      getTrimmedStrings(dockerConfigPrimitives.LogsExcludePatterns),
    }
}

func (e *DockerConfigEntity) ToPrimitives() DockerConfigPrimitives {
    return DockerConfigPrimitives{
        IsEnabled:                e.isEnabled,
        SocketPath:               e.socketPath,
        IsStatsEnabled:           &e.isStatsEnabled,
        StatsIntervalSeconds:     e.statsIntervalSeconds,
        IsEventsEnabled:          &e.isEventsEnabled,
        IsLogsEnabled:            &e.isLogsEnabled,
        LogsMaxLinesPerContainer: e.logsMaxLinesPerContainer,
        LogsIncludePatterns:      e.logsIncludePatterns,
        LogsExcludePatterns:      e.logsExcludePatterns,
    }
}

func (e *DockerConfigEntity) IsEnabled() bool {
    return e.isEnabled
}

func (e *DockerConfigEntity) GetSocketPath() string {
    return e.socketPath
}

func (e *DockerConfigEntity) IsStatsEnabled() bool {
    return e.isEnabled && e.isStatsEnabled
}

func (e *DockerConfigEntity) GetStatsIntervalSeconds() int {
    return e.statsIntervalSeconds
}

func (e *DockerConfigEntity) IsEventsEnabled() bool {
    return e.isEnabled && e.isEventsEnabled
}

func (e *DockerConfigEntity) IsLogsEnabled() bool {
    return e.isEnabled && e.isLogsEnabled
}

func (e *DockerConfigEntity) GetLogsMaxLinesPerContainer() int {
    return e.logsMaxLinesPerContainer
}

func (e *DockerConfigEntity) GetLogsIncludePatterns() []string {
    return e.logsIncludePatterns
}

func (e *DockerConfigEntity) GetLogsExcludePatterns() []string {
    return e.logsExcludePatterns
}

// IsContainerLogged casa con algun include y con ningun exclude (un patron no valido no casa: lo avisa el validador).
func (e *DockerConfigEntity) IsContainerLogged(containerName string) bool {
    return isMatchingAnyPattern(containerName, e.logsIncludePatterns) && !isMatchingAnyPattern(containerName, e.logsExcludePatterns)
}

func isMatchingAnyPattern(name string, patterns []string) bool {
    for _, pattern := range patterns {
        if isMatch, err := path.Match(pattern, name); err == nil && isMatch {
            return true
        }
    }
    return false
}

func getBoolOrTrue(value *bool) bool {
    return value == nil || *value
}

func getTrimmedStrings(values []string) []string {
    trimmedValues := []string{}
    for _, value := range values {
        if trimmedValue := trim(value); trimmedValue != "" {
            trimmedValues = append(trimmedValues, trimmedValue)
        }
    }
    return trimmedValues
}
