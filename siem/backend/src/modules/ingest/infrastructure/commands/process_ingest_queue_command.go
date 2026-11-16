package commands

import (
    "errors"

    sharedcommands "siem-backend/src/modules/shared/infrastructure/commands"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"

    "siem-backend/src/modules/ingest/domain/enums"
    "siem-backend/src/modules/ingest/application/processingestqueue"
)

// ProcessIngestQueueCommand ./console ingest-queue-process --max_batches=100
// Lo mismo que hace el worker de la API, a mano o por cron (con APP_QUEUE_WORKER_INTERVAL_SECONDS=0).
type ProcessIngestQueueCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func ProcessIngestQueueCommandGetInstance(input []string) *ProcessIngestQueueCommand {
    return &ProcessIngestQueueCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *ProcessIngestQueueCommand) Run() {
    processIngestQueueResultDto, err := processingestqueue.ProcessIngestQueueServiceGetInstance().Invoke(
        processingestqueue.ProcessIngestQueueDtoFromPrimitives(map[string]any{
            "maxBatches": c.GetOptionValue(enums.IngestKeyEnum.MAX_BATCHES),
        }),
    )
    if err == nil {
        c.PrintResult(processIngestQueueResultDto.ToArray())
        return
    }

    var fieldsException *sharedexceptions.FieldsException
    if errors.As(err, &fieldsException) {
        c.ExitWithDomainError(fieldsException)
    }
    c.ExitWithUnexpectedError(err, "ProcessIngestQueueCommand.Run")
}
