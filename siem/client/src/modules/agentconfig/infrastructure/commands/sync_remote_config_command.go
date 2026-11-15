package commands

import (
    "errors"

    "siem-client/src/modules/agentconfig/application/syncremoteconfig"
    "siem-client/src/modules/agentconfig/domain/exceptions"
    sharedcommands "siem-client/src/modules/shared/infrastructure/commands"
)

// SyncRemoteConfigCommand ./console config-sync   baja del backend los monitores de este agente (lo que se crea en el
// panel) y los guarda en state/remote-config.json. El bucle agent-run lo hace solo cada 60 s.
type SyncRemoteConfigCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func SyncRemoteConfigCommandGetInstance(input []string) *SyncRemoteConfigCommand {
    return &SyncRemoteConfigCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *SyncRemoteConfigCommand) Run() {
    syncRemoteConfigResultDto, err := syncremoteconfig.SyncRemoteConfigServiceGetInstance().Invoke(
        syncremoteconfig.SyncRemoteConfigDtoFromPrimitives(map[string]any{}),
    )
    if err == nil {
        c.PrintResult(syncRemoteConfigResultDto.ToArray())
        return
    }

    var agentConfigException *exceptions.AgentConfigException
    if errors.As(err, &agentConfigException) {
        c.ExitWithDomainError(agentConfigException)
    }
    c.ExitWithUnexpectedError(err, "SyncRemoteConfigCommand.Run")
}
