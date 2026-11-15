package commands

import (
    "errors"

    "siem-client/src/modules/agentconfig/application/checkagentconfig"
    "siem-client/src/modules/agentconfig/domain/exceptions"
    sharedcommands "siem-client/src/modules/shared/infrastructure/commands"
)

// CheckAgentConfigCommand ./console config-check   valida siem-client.json + endpoint y token del .env.
type CheckAgentConfigCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func CheckAgentConfigCommandGetInstance(input []string) *CheckAgentConfigCommand {
    return &CheckAgentConfigCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *CheckAgentConfigCommand) Run() {
    checkAgentConfigResultDto, err := checkagentconfig.CheckAgentConfigServiceGetInstance().Invoke(
        checkagentconfig.CheckAgentConfigDtoFromPrimitives(map[string]any{}),
    )
    if err == nil {
        c.PrintResult(checkAgentConfigResultDto.ToArray())
        return
    }

    var agentConfigException *exceptions.AgentConfigException
    if errors.As(err, &agentConfigException) {
        c.ExitWithDomainError(agentConfigException)
    }
    c.ExitWithUnexpectedError(err, "CheckAgentConfigCommand.Run")
}
