package commands

import (
    "errors"

    agentconfigexceptions "siem-client/src/modules/agentconfig/domain/exceptions"
    "siem-client/src/modules/collectors/application/runduecommands"
    sharedcommands "siem-client/src/modules/shared/infrastructure/commands"
)

// RunDueCommandsCommand ./console commands-run [--force=1]   lanza los comandos a los que les toca (con --force=1,
// todos) y deja el resultado en el spool.
type RunDueCommandsCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func RunDueCommandsCommandGetInstance(input []string) *RunDueCommandsCommand {
    return &RunDueCommandsCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *RunDueCommandsCommand) Run() {
    runDueCommandsResultDto, err := runduecommands.RunDueCommandsServiceGetInstance().Invoke(
        runduecommands.RunDueCommandsDtoFromPrimitives(map[string]any{"isForced": c.GetOptionValue("force")}),
    )
    if err == nil {
        c.PrintResult(runDueCommandsResultDto.ToArray())
        return
    }

    var agentConfigException *agentconfigexceptions.AgentConfigException
    if errors.As(err, &agentConfigException) {
        c.ExitWithDomainError(agentConfigException)
    }
    c.ExitWithUnexpectedError(err, "RunDueCommandsCommand.Run")
}
