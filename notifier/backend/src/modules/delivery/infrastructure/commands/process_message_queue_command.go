package commands

import (
    sharedcommands "notifier-backend/src/modules/shared/infrastructure/commands"

    sendersenums "notifier-backend/src/modules/senders/domain/enums"

    "notifier-backend/src/modules/delivery/application/processmessagequeue"
)

// ProcessMessageQueueCommand ./console message-queue-process [--channel=email] [--max_messages=50]
// Lo mismo que hace el worker de la API, a mano o por cron (con APP_QUEUE_WORKER_INTERVAL_SECONDS=0).
type ProcessMessageQueueCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func ProcessMessageQueueCommandGetInstance(input []string) *ProcessMessageQueueCommand {
    return &ProcessMessageQueueCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *ProcessMessageQueueCommand) Run() {
    channel := c.GetOptionValue("channel")
    if channel == "" {
        channel = sendersenums.ChannelEnum.EMAIL
    }
    processMessageQueueResultDto, err := processmessagequeue.ProcessMessageQueueServiceGetInstance().Invoke(
        processmessagequeue.ProcessMessageQueueDtoFromPrimitives(map[string]any{
            "channel":     channel,
            "maxMessages": c.GetOptionValue("max_messages"),
        }),
    )
    if err != nil {
        c.ExitWithError(err, "ProcessMessageQueueCommand.Run")
    }
    c.PrintResult(processMessageQueueResultDto.ToArray())
}
