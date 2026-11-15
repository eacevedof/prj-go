package services

import (
    "strings"
    "testing"

    "siem-client/src/modules/agentconfig/domain/entities"
    "siem-client/src/modules/agentconfig/domain/services"
    "siem-client/tests"
)

const (
    endpointUrl = "http://localhost:9400/api/v1/ingest-batches"
    apiToken    = "token"
)

func TestValidConfigWithDefaults(t *testing.T) {
    agentConfigEntity := entities.AgentConfigEntityFromPrimitives(entities.AgentConfigPrimitives{
        Watchers: []entities.WatcherPrimitives{{Name: "app-logs", Path: getAbsolutePath()}},
        Commands: []entities.ScheduledCommandPrimitives{{Name: "disk", Command: "df -h"}},
    })
    tests.AssertEmpty(t, services.AgentConfigValidatorGetInstance().GetErrors(agentConfigEntity, endpointUrl, apiToken))
    tests.AssertEquals(t, entities.DEFAULT_COLLECT_INTERVAL_SECONDS, agentConfigEntity.GetCollectIntervalSeconds())
    tests.AssertEquals(t, "*", agentConfigEntity.GetWatcherEntities()[0].GetPattern())
    tests.AssertEquals(t, "end", agentConfigEntity.GetWatcherEntities()[0].GetStartAt())
    tests.AssertEquals(t, entities.DEFAULT_COMMAND_TIMEOUT_SECONDS, agentConfigEntity.GetScheduledCommandEntities()[0].GetTimeoutSeconds())
}

func TestInvalidConfigReportsEveryError(t *testing.T) {
    agentConfigEntity := entities.AgentConfigEntityFromPrimitives(entities.AgentConfigPrimitives{
        MaxEventsPerBatch: 5000,
        Watchers: []entities.WatcherPrimitives{
            {Name: "logs", Path: "relative/path", Pattern: "[", StartAt: "middle"},
            {Name: "logs", Path: getAbsolutePath()},
        },
        Commands: []entities.ScheduledCommandPrimitives{{Name: "bad name!", IntervalSeconds: 1}},
    })
    errors := strings.Join(services.AgentConfigValidatorGetInstance().GetErrors(agentConfigEntity, "ftp://x", ""), "\n")

    for _, expectedError := range []string{
        "APP_SIEM_ENDPOINT_URL",
        "APP_SIEM_API_TOKEN",
        "max_events_per_batch: must be between 1 and 2000",
        "path must be an absolute",
        "pattern is not a valid glob",
        "start_at must be end or beginning",
        "watchers.name: logs is duplicated",
        "commands.name: required",
        "command is required",
        "interval_seconds: must be between 5 and 86400",
    } {
        tests.AssertStringContainsString(t, expectedError, errors)
    }
}

func TestNothingToWatch(t *testing.T) {
    errors := services.AgentConfigValidatorGetInstance().GetErrors(
        entities.AgentConfigEntityFromPrimitives(entities.AgentConfigPrimitives{}), endpointUrl, apiToken,
    )
    tests.AssertEquals(t, []string{"watchers, commands or docker: at least one is required"}, errors)
}

func getAbsolutePath() string {
    return tests.GetTempAbsolutePath()
}
