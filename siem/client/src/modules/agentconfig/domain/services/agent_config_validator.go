package services

import (
    "net/url"
    "path"
    "path/filepath"
    "regexp"
    "slices"
    "strings"

    "siem-client/src/core/config"
    "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/agentconfig/domain/enums"
)

var namePattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,100}$`)

// AgentConfigValidator reglas de la configuracion del agente (siem-client.json + endpoint y token del .env).
// Devuelve la pila de errores (vacia = todo bien).
type AgentConfigValidator struct{}

func AgentConfigValidatorGetInstance() *AgentConfigValidator {
    return &AgentConfigValidator{}
}

// GetErrors la config local sola (sin config remota): tiene que vigilar algo.
func (v *AgentConfigValidator) GetErrors(agentConfigEntity *entities.AgentConfigEntity, endpointUrl string, apiToken string) []string {
    return v.GetErrorsWithRemoteConfig(agentConfigEntity, endpointUrl, apiToken, false)
}

// GetErrorsWithRemoteConfig valida la config LOCAL (los monitores remotos los filtra RemoteMonitorsMerger). Con la
// config remota activa el agente puede no tener nada local que vigilar: lo que vigila le llega del backend.
func (v *AgentConfigValidator) GetErrorsWithRemoteConfig(agentConfigEntity *entities.AgentConfigEntity, endpointUrl string, apiToken string, isRemoteConfigEnabled bool) []string {
    errors := v.getConnectionErrors(endpointUrl, apiToken)
    errors = append(errors, v.getIntervalsErrors(agentConfigEntity)...)
    isDockerEnabled := agentConfigEntity.GetDockerConfigEntity().IsEnabled()
    hasNothingToWatch := len(agentConfigEntity.GetWatcherEntities()) == 0 && len(agentConfigEntity.GetScheduledCommandEntities()) == 0 && !isDockerEnabled
    if hasNothingToWatch && !isRemoteConfigEnabled {
        errors = append(errors, enums.AgentConfigMessageEnum.NOTHING_TO_WATCH)
    }
    errors = append(errors, v.getWatchersErrors(agentConfigEntity.GetWatcherEntities())...)
    errors = append(errors, v.getCommandsErrors(agentConfigEntity.GetScheduledCommandEntities())...)
    if !isDockerEnabled {
        return errors
    }
    return append(errors, v.getDockerErrors(agentConfigEntity.GetDockerConfigEntity())...)
}

// getDockerErrors solo con el bloque activo: apagado, sus valores no importan.
func (v *AgentConfigValidator) getDockerErrors(dockerConfigEntity *entities.DockerConfigEntity) []string {
    errors := []string{}
    if !strings.HasPrefix(dockerConfigEntity.GetSocketPath(), "/") {
        errors = append(errors, enums.AgentConfigMessageEnum.DOCKER_SOCKET_PATH_MUST_BE_ABSOLUTE)
    }
    if interval := dockerConfigEntity.GetStatsIntervalSeconds(); interval < 10 || interval > 3600 {
        errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.X_MUST_BE_BETWEEN_X_AND_X, "docker.stats_interval_seconds", 10, 3600))
    }
    if maxLines := dockerConfigEntity.GetLogsMaxLinesPerContainer(); maxLines < 1 || maxLines > 5000 {
        errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.X_MUST_BE_BETWEEN_X_AND_X, "docker.logs_max_lines_per_container", 1, 5000))
    }
    for key, patterns := range map[string][]string{
        "logs_include": dockerConfigEntity.GetLogsIncludePatterns(),
        "logs_exclude": dockerConfigEntity.GetLogsExcludePatterns(),
    } {
        for _, pattern := range patterns {
            if _, err := path.Match(pattern, "probe"); err != nil {
                errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.DOCKER_X_PATTERN_X_IS_NOT_VALID, key, pattern))
            }
        }
    }
    return errors
}

func (v *AgentConfigValidator) getConnectionErrors(endpointUrl string, apiToken string) []string {
    errors := []string{}
    parsedUrl, err := url.Parse(endpointUrl)
    if err != nil || (parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https") || parsedUrl.Host == "" {
        errors = append(errors, enums.AgentConfigMessageEnum.ENDPOINT_URL_IS_NOT_VALID)
    }
    if apiToken == "" {
        errors = append(errors, enums.AgentConfigMessageEnum.API_TOKEN_IS_REQUIRED)
    }
    return errors
}

func (v *AgentConfigValidator) getIntervalsErrors(agentConfigEntity *entities.AgentConfigEntity) []string {
    errors := []string{}
    for _, rule := range []struct {
        key      string
        value    int
        minValue int
        maxValue int
    }{
        {enums.AgentConfigKeyEnum.COLLECT_INTERVAL_SECONDS, agentConfigEntity.GetCollectIntervalSeconds(), 1, 3600},
        {enums.AgentConfigKeyEnum.SEND_INTERVAL_SECONDS, agentConfigEntity.GetSendIntervalSeconds(), 1, 3600},
        {enums.AgentConfigKeyEnum.HEARTBEAT_INTERVAL_SECONDS, agentConfigEntity.GetHeartbeatIntervalSeconds(), 5, 3600},
        // el backend admite hasta 2000 eventos por lote
        {enums.AgentConfigKeyEnum.MAX_EVENTS_PER_BATCH, agentConfigEntity.GetMaxEventsPerBatch(), 1, 2000},
        {enums.AgentConfigKeyEnum.MAX_SPOOL_BATCHES, agentConfigEntity.GetMaxSpoolBatches(), 10, 100000},
    } {
        if rule.value < rule.minValue || rule.value > rule.maxValue {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.X_MUST_BE_BETWEEN_X_AND_X, rule.key, rule.minValue, rule.maxValue))
        }
    }
    return errors
}

func (v *AgentConfigValidator) getWatchersErrors(watcherEntities []*entities.WatcherEntity) []string {
    errors := []string{}
    names := []string{}
    for _, watcherEntity := range watcherEntities {
        name := watcherEntity.GetName()
        if !namePattern.MatchString(name) {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.X_NAME_IS_NOT_VALID, enums.AgentConfigKeyEnum.WATCHERS))
        } else if slices.Contains(names, name) {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.X_NAME_X_IS_DUPLICATED, enums.AgentConfigKeyEnum.WATCHERS, name))
        }
        names = append(names, name)

        if !filepath.IsAbs(watcherEntity.GetPath()) {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.WATCHER_X_PATH_MUST_BE_ABSOLUTE, name))
        }
        if _, err := filepath.Match(watcherEntity.GetPattern(), "probe"); err != nil {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.WATCHER_X_PATTERN_IS_NOT_VALID, name))
        }
        if !slices.Contains(enums.WatcherStartAtEnum.VALUES, watcherEntity.GetStartAt()) {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.WATCHER_X_START_AT_IS_NOT_VALID, name))
        }
        if !watcherEntity.IsMatchValid() {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.WATCHER_X_MATCH_IS_NOT_VALID, name))
        }
        if severity := watcherEntity.GetSeverity(); severity != "" && !slices.Contains(enums.LogSeverityEnum.VALUES, severity) {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.WATCHER_X_SEVERITY_IS_NOT_VALID, name))
        }
    }
    return errors
}

func (v *AgentConfigValidator) getCommandsErrors(scheduledCommandEntities []*entities.ScheduledCommandEntity) []string {
    errors := []string{}
    names := []string{}
    for _, scheduledCommandEntity := range scheduledCommandEntities {
        name := scheduledCommandEntity.GetName()
        if !namePattern.MatchString(name) {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.X_NAME_IS_NOT_VALID, enums.AgentConfigKeyEnum.COMMANDS))
        } else if slices.Contains(names, name) {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.X_NAME_X_IS_DUPLICATED, enums.AgentConfigKeyEnum.COMMANDS, name))
        }
        names = append(names, name)

        if scheduledCommandEntity.GetCommandLine() == "" {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.COMMAND_X_COMMAND_IS_REQUIRED, name))
        }
        if interval := scheduledCommandEntity.GetIntervalSeconds(); interval < 5 || interval > 86400 {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.X_MUST_BE_BETWEEN_X_AND_X, "command "+name+" interval_seconds", 5, 86400))
        }
        if timeout := scheduledCommandEntity.GetTimeoutSeconds(); timeout < 1 || timeout > 600 {
            errors = append(errors, config.Tr(enums.AgentConfigMessageEnum.X_MUST_BE_BETWEEN_X_AND_X, "command "+name+" timeout_seconds", 1, 600))
        }
    }
    return errors
}
