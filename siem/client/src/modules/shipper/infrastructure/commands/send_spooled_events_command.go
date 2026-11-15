package commands

import (
    "errors"

    agentconfigexceptions "siem-client/src/modules/agentconfig/domain/exceptions"
    sharedcommands "siem-client/src/modules/shared/infrastructure/commands"
    "siem-client/src/modules/shipper/application/sendspooledevents"
    "siem-client/src/modules/shipper/domain/exceptions"
)

// SendSpooledEventsCommand ./console events-send --max_batches=20   envia el spool al backend (lo que hace agent-run
// cada send_interval_seconds).
type SendSpooledEventsCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func SendSpooledEventsCommandGetInstance(input []string) *SendSpooledEventsCommand {
    return &SendSpooledEventsCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *SendSpooledEventsCommand) Run() {
    sendSpooledEventsResultDto, err := sendspooledevents.SendSpooledEventsServiceGetInstance().Invoke(
        sendspooledevents.SendSpooledEventsDtoFromPrimitives(map[string]any{"maxBatches": c.GetOptionValue("max_batches")}),
    )
    if err == nil {
        c.PrintResult(sendSpooledEventsResultDto.ToArray())
        return
    }

    var agentConfigException *agentconfigexceptions.AgentConfigException
    var shipperException *exceptions.ShipperException
    switch {
    case errors.As(err, &agentConfigException):
        c.ExitWithDomainError(agentConfigException)
    case errors.As(err, &shipperException):
        c.ExitWithDomainError(shipperException)
    }
    c.ExitWithUnexpectedError(err, "SendSpooledEventsCommand.Run")
}
