package commands

import (
    sharedcommands "notifier-backend/src/modules/shared/infrastructure/commands"

    "notifier-backend/src/modules/apitokens/application/createapitoken"
    "notifier-backend/src/modules/apitokens/application/disableapitoken"
    "notifier-backend/src/modules/apitokens/domain/enums"
)

// CreateApiTokenCommand ./console token-create --sender_uuid=<uuid> --token_name=n8n [--default_to="a@x.com,b@x.com"]
// Imprime el api_token (solo esta vez).
type CreateApiTokenCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func CreateApiTokenCommandGetInstance(input []string) *CreateApiTokenCommand {
    return &CreateApiTokenCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *CreateApiTokenCommand) Run() {
    createApiTokenResultDto, err := createapitoken.CreateApiTokenServiceGetInstance().Invoke(
        createapitoken.CreateApiTokenDtoFromPrimitives(map[string]any{
            "senderUuid":   c.GetOptionValue(enums.ApiTokensKeyEnum.SENDER_UUID),
            "tokenName":    c.GetOptionValue(enums.ApiTokensKeyEnum.TOKEN_NAME),
            "defaultToCsv": c.GetOptionValue(enums.ApiTokensKeyEnum.DEFAULT_TO),
        }),
    )
    if err != nil {
        c.ExitWithError(err, "CreateApiTokenCommand.Run")
    }
    c.PrintResult(createApiTokenResultDto.ToArray())
}

// DisableApiTokenCommand ./console token-disable --token_uuid=<uuid>
type DisableApiTokenCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func DisableApiTokenCommandGetInstance(input []string) *DisableApiTokenCommand {
    return &DisableApiTokenCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *DisableApiTokenCommand) Run() {
    disableApiTokenResultDto, err := disableapitoken.DisableApiTokenServiceGetInstance().Invoke(
        disableapitoken.DisableApiTokenDtoFromPrimitives(map[string]any{
            "tokenUuid": c.GetOptionValue(enums.ApiTokensKeyEnum.TOKEN_UUID),
        }),
    )
    if err != nil {
        c.ExitWithError(err, "DisableApiTokenCommand.Run")
    }
    c.PrintResult(disableApiTokenResultDto.ToArray())
}
