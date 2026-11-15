package entities

import (
    "strings"
    "testing"

    "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/agentconfig/domain/services"
    "siem-client/tests"
)

func TestDockerDefaultsWithEveryFeatureOn(t *testing.T) {
    dockerConfigEntity := entities.DockerConfigEntityFromPrimitives(entities.DockerConfigPrimitives{IsEnabled: true})
    tests.AssertEquals(t, "/var/run/docker.sock", dockerConfigEntity.GetSocketPath())
    tests.AssertEquals(t, true, dockerConfigEntity.IsStatsEnabled())
    tests.AssertEquals(t, true, dockerConfigEntity.IsEventsEnabled())
    tests.AssertEquals(t, true, dockerConfigEntity.IsLogsEnabled())
    tests.AssertEquals(t, entities.DEFAULT_DOCKER_STATS_INTERVAL_SECONDS, dockerConfigEntity.GetStatsIntervalSeconds())
    tests.AssertEquals(t, []string{"*"}, dockerConfigEntity.GetLogsIncludePatterns())
}

func TestDockerDisabledTurnsEveryFeatureOff(t *testing.T) {
    dockerConfigEntity := entities.DockerConfigEntityFromPrimitives(entities.DockerConfigPrimitives{})
    tests.AssertEquals(t, false, dockerConfigEntity.IsStatsEnabled())
    tests.AssertEquals(t, false, dockerConfigEntity.IsLogsEnabled())
}

func TestDockerFeatureCanBeTurnedOffAlone(t *testing.T) {
    isOff := false
    dockerConfigEntity := entities.DockerConfigEntityFromPrimitives(entities.DockerConfigPrimitives{IsEnabled: true, IsLogsEnabled: &isOff})
    tests.AssertEquals(t, false, dockerConfigEntity.IsLogsEnabled())
    tests.AssertEquals(t, true, dockerConfigEntity.IsStatsEnabled())
}

func TestContainerLogsIncludeAndExclude(t *testing.T) {
    dockerConfigEntity := entities.DockerConfigEntityFromPrimitives(entities.DockerConfigPrimitives{
        IsEnabled: true, LogsIncludePatterns: []string{"cont-*"}, LogsExcludePatterns: []string{"cont-siem-*"},
    })
    tests.AssertEquals(t, true, dockerConfigEntity.IsContainerLogged("cont-web"))
    tests.AssertEquals(t, false, dockerConfigEntity.IsContainerLogged("cont-siem-backend-app"))
    tests.AssertEquals(t, false, dockerConfigEntity.IsContainerLogged("redis"))
}

func TestDockerOnlyConfigIsValidAndBadValuesAreReported(t *testing.T) {
    validator := services.AgentConfigValidatorGetInstance()
    dockerOnly := entities.AgentConfigEntityFromPrimitives(entities.AgentConfigPrimitives{Docker: entities.DockerConfigPrimitives{IsEnabled: true}})
    tests.AssertEmpty(t, validator.GetErrors(dockerOnly, "http://localhost:9400/api/v1/ingest-batches", "token"))

    badDocker := entities.AgentConfigEntityFromPrimitives(entities.AgentConfigPrimitives{Docker: entities.DockerConfigPrimitives{
        IsEnabled: true, SocketPath: "docker.sock", StatsIntervalSeconds: 5, LogsMaxLinesPerContainer: 9000, LogsExcludePatterns: []string{"["},
    }})
    errors := strings.Join(validator.GetErrors(badDocker, "http://localhost:9400/api/v1/ingest-batches", "token"), "\n")
    for _, expectedError := range []string{
        "docker.socket_path: must be an absolute path",
        "docker.stats_interval_seconds: must be between 10 and 3600",
        "docker.logs_max_lines_per_container: must be between 1 and 5000",
        "docker.logs_exclude: [ is not a valid glob",
    } {
        tests.AssertStringContainsString(t, expectedError, errors)
    }
}
