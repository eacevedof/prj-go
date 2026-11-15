package commands

import (
    "errors"

    agentconfigexceptions "siem-client/src/modules/agentconfig/domain/exceptions"
    "siem-client/src/modules/hostmetrics/application/collecthostmetrics"
    sharedcommands "siem-client/src/modules/shared/infrastructure/commands"
)

// CollectHostMetricsCommand ./console metrics-collect [--force=1]   muestrea los monitores de metricas del host (cpu,
// memory, disk, gpu, network) a los que les toca (con --force=1, todos) y deja el resultado en el spool. Los
// monitores salen de la config remota guardada (./console config-sync la baja).
type CollectHostMetricsCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func CollectHostMetricsCommandGetInstance(input []string) *CollectHostMetricsCommand {
    return &CollectHostMetricsCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *CollectHostMetricsCommand) Run() {
    collectHostMetricsResultDto, err := collecthostmetrics.CollectHostMetricsServiceGetInstance().Invoke(
        collecthostmetrics.CollectHostMetricsDtoFromPrimitives(map[string]any{"isForced": c.GetOptionValue("force")}),
    )
    if err == nil {
        c.PrintResult(collectHostMetricsResultDto.ToArray())
        return
    }

    var agentConfigException *agentconfigexceptions.AgentConfigException
    if errors.As(err, &agentConfigException) {
        c.ExitWithDomainError(agentConfigException)
    }
    c.ExitWithUnexpectedError(err, "CollectHostMetricsCommand.Run")
}
