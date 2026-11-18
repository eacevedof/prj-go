package commands

import (
    "errors"

    sharedcommands "siem-client/src/modules/shared/infrastructure/commands"

    agentconfigexceptions "siem-client/src/modules/agentconfig/domain/exceptions"

    "siem-client/src/modules/connections/application/collectconnectionchanges"
)

// CollectConnectionChangesCommand ./console connections-collect [--force=1]   una pasada del monitor connections del
// panel (si le toca; con --force=1, ya): deja en el spool las conexiones que han abierto o cerrado desde la pasada
// anterior. No es un listado: para ver las conexiones en consola esta trafficker.
type CollectConnectionChangesCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func CollectConnectionChangesCommandGetInstance(input []string) *CollectConnectionChangesCommand {
    return &CollectConnectionChangesCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *CollectConnectionChangesCommand) Run() {
    collectConnectionChangesResultDto, err := collectconnectionchanges.CollectConnectionChangesServiceGetInstance().Invoke(
        collectconnectionchanges.CollectConnectionChangesDtoFromPrimitives(map[string]any{"isForced": c.GetOptionValue("force")}),
    )
    if err == nil {
        c.PrintResult(collectConnectionChangesResultDto.ToArray())
        return
    }

    var agentConfigException *agentconfigexceptions.AgentConfigException
    if errors.As(err, &agentConfigException) {
        c.ExitWithDomainError(agentConfigException)
    }
    c.ExitWithUnexpectedError(err, "CollectConnectionChangesCommand.Run")
}
