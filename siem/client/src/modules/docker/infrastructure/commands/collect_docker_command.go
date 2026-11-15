package commands

import (
    "errors"

    agentconfigexceptions "siem-client/src/modules/agentconfig/domain/exceptions"
    "siem-client/src/modules/docker/application/collectcontainerlogs"
    "siem-client/src/modules/docker/application/collectcontainerstates"
    "siem-client/src/modules/docker/application/collectcontainerstats"
    "siem-client/src/modules/docker/application/collectdaemonevents"
    dockerexceptions "siem-client/src/modules/docker/domain/exceptions"
    sharedcommands "siem-client/src/modules/shared/infrastructure/commands"
)

// CollectDockerCommand ./console docker-collect [--force_stats=1]   una pasada de los cuatro colectores de Docker
// (estados, eventos del daemon, recursos si tocan, logs) al spool. Para probar el socket y ver que se manda.
type CollectDockerCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func CollectDockerCommandGetInstance(input []string) *CollectDockerCommand {
    return &CollectDockerCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *CollectDockerCommand) Run() {
    collectContainerStatesResultDto, err := collectcontainerstates.CollectContainerStatesServiceGetInstance().Invoke(
        collectcontainerstates.CollectContainerStatesDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.exitWithError(err)
    }
    collectDaemonEventsResultDto, err := collectdaemonevents.CollectDaemonEventsServiceGetInstance().Invoke(
        collectdaemonevents.CollectDaemonEventsDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.exitWithError(err)
    }
    collectContainerStatsResultDto, err := collectcontainerstats.CollectContainerStatsServiceGetInstance().Invoke(
        collectcontainerstats.CollectContainerStatsDtoFromPrimitives(map[string]any{"isForced": c.GetOptionValue("force_stats")}),
    )
    if err != nil {
        c.exitWithError(err)
    }
    collectContainerLogsResultDto, err := collectcontainerlogs.CollectContainerLogsServiceGetInstance().Invoke(
        collectcontainerlogs.CollectContainerLogsDtoFromPrimitives(map[string]any{}),
    )
    if err != nil {
        c.exitWithError(err)
    }

    c.PrintResult(map[string]any{
        "states": collectContainerStatesResultDto.ToArray(),
        "events": collectDaemonEventsResultDto.ToArray(),
        "stats":  collectContainerStatsResultDto.ToArray(),
        "logs":   collectContainerLogsResultDto.ToArray(),
    })
}

func (c *CollectDockerCommand) exitWithError(err error) {
    var agentConfigException *agentconfigexceptions.AgentConfigException
    if errors.As(err, &agentConfigException) {
        c.ExitWithDomainError(agentConfigException)
    }
    var dockerException *dockerexceptions.DockerException
    if errors.As(err, &dockerException) {
        c.ExitWithDomainError(dockerException)
    }
    c.ExitWithUnexpectedError(err, "CollectDockerCommand.Run")
}
