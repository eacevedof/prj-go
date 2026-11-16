package commands

import (
    "errors"

    sharedcommands "siem-backend/src/modules/shared/infrastructure/commands"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"

    "siem-backend/src/modules/events/domain/enums"
    "siem-backend/src/modules/events/application/purgeevents"
)

// PurgeEventsCommand ./console events-purge --retention_days=30   (sin la opcion: APP_EVENTS_RETENTION_DAYS)
type PurgeEventsCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func PurgeEventsCommandGetInstance(input []string) *PurgeEventsCommand {
    return &PurgeEventsCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *PurgeEventsCommand) Run() {
    purgeEventsResultDto, err := purgeevents.PurgeEventsServiceGetInstance().Invoke(
        purgeevents.PurgeEventsDtoFromPrimitives(map[string]any{
            "retentionDays": c.GetOptionValue(enums.EventsKeyEnum.RETENTION_DAYS),
        }),
    )
    if err == nil {
        c.PrintResult(purgeEventsResultDto.ToArray())
        return
    }

    var fieldsException *sharedexceptions.FieldsException
    if errors.As(err, &fieldsException) {
        c.ExitWithDomainError(fieldsException)
    }
    c.ExitWithUnexpectedError(err, "PurgeEventsCommand.Run")
}
