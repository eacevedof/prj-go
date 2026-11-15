package commands

import (
    "errors"

    agentconfigexceptions "siem-client/src/modules/agentconfig/domain/exceptions"
    "siem-client/src/modules/collectors/application/collectfilelines"
    sharedcommands "siem-client/src/modules/shared/infrastructure/commands"
)

// CollectFileLinesCommand ./console files-collect   una pasada por las carpetas vigiladas (lo que hace agent-run cada
// collect_interval_seconds). Deja los eventos en el spool: los envia events-send.
type CollectFileLinesCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func CollectFileLinesCommandGetInstance(input []string) *CollectFileLinesCommand {
    return &CollectFileLinesCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *CollectFileLinesCommand) Run() {
    collectFileLinesResultDto, err := collectfilelines.CollectFileLinesServiceGetInstance().Invoke(
        collectfilelines.CollectFileLinesDtoFromPrimitives(map[string]any{}),
    )
    if err == nil {
        c.PrintResult(collectFileLinesResultDto.ToArray())
        return
    }

    var agentConfigException *agentconfigexceptions.AgentConfigException
    if errors.As(err, &agentConfigException) {
        c.ExitWithDomainError(agentConfigException)
    }
    c.ExitWithUnexpectedError(err, "CollectFileLinesCommand.Run")
}
