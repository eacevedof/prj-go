package commands

import (
    sharedcommands "notifier-backend/src/modules/shared/infrastructure/commands"

    "notifier-backend/src/modules/senders/application/createsender"
    "notifier-backend/src/modules/senders/application/listsenders"
    "notifier-backend/src/modules/senders/domain/enums"
)

// CreateSenderCommand ./console sender-create --user_uuid=<uuid> --sender_name=avisos --provider_code=gmail
// --from_email=yo@gmail.com --from_name="Avisos" --secret="app password" [--token_name=n8n]
// smtp_custom ademas: --smtp_host --smtp_port --smtp_security=starttls|ssl|none [--username]. Imprime el api_token.
type CreateSenderCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func CreateSenderCommandGetInstance(input []string) *CreateSenderCommand {
    return &CreateSenderCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *CreateSenderCommand) Run() {
    keys := enums.SendersKeyEnum
    createSenderResultDto, err := createsender.CreateSenderServiceGetInstance().Invoke(
        createsender.CreateSenderDtoFromPrimitives(map[string]any{
            "userUuid":     c.GetOptionValue(keys.USER_UUID),
            "senderName":   c.GetOptionValue(keys.SENDER_NAME),
            "providerCode": c.GetOptionValue(keys.PROVIDER_CODE),
            "fromEmail":    c.GetOptionValue(keys.FROM_EMAIL),
            "fromName":     c.GetOptionValue(keys.FROM_NAME),
            "smtpHost":     c.GetOptionValue(keys.SMTP_HOST),
            "smtpPort":     c.GetOptionValue(keys.SMTP_PORT),
            "smtpSecurity": c.GetOptionValue(keys.SMTP_SECURITY),
            "username":     c.GetOptionValue(keys.USERNAME),
            "secret":       c.GetOptionValue(keys.SECRET),
            "dailyLimit":   c.GetOptionValue(keys.DAILY_LIMIT),
            "minuteLimit":  c.GetOptionValue(keys.MINUTE_LIMIT),
            "tokenName":    c.GetOptionValue(keys.TOKEN_NAME),
        }),
    )
    if err != nil {
        c.ExitWithError(err, "CreateSenderCommand.Run")
    }
    c.PrintResult(createSenderResultDto.ToArray())
}

// ListSendersCommand ./console sender-list [--user_uuid=<uuid>]
type ListSendersCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func ListSendersCommandGetInstance(input []string) *ListSendersCommand {
    return &ListSendersCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *ListSendersCommand) Run() {
    listSendersResultDto, err := listsenders.ListSendersServiceGetInstance().Invoke(
        listsenders.ListSendersDtoFromPrimitives(map[string]any{"userUuid": c.GetOptionValue(enums.SendersKeyEnum.USER_UUID)}),
    )
    if err != nil {
        c.ExitWithError(err, "ListSendersCommand.Run")
    }
    c.PrintResult(listSendersResultDto.ToArray())
}
