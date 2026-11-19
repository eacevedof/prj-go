package commands

import (
    sharedcommands "notifier-backend/src/modules/shared/infrastructure/commands"

    "notifier-backend/src/modules/users/application/createuser"
    "notifier-backend/src/modules/users/domain/enums"
)

// CreateUserCommand ./console user-create --user_name="Ana" --user_email=ana@x.com --user_role=owner --password=...
// Sin --user_role, owner.
type CreateUserCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func CreateUserCommandGetInstance(input []string) *CreateUserCommand {
    return &CreateUserCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *CreateUserCommand) Run() {
    userRole := c.GetOptionValue(enums.UsersKeyEnum.USER_ROLE)
    if userRole == "" {
        userRole = enums.UserRoleEnum.OWNER
    }
    createUserResultDto, err := createuser.CreateUserServiceGetInstance().Invoke(
        createuser.CreateUserDtoFromPrimitives(map[string]any{
            "userName":  c.GetOptionValue(enums.UsersKeyEnum.USER_NAME),
            "userEmail": c.GetOptionValue(enums.UsersKeyEnum.USER_EMAIL),
            "userRole":  userRole,
            "password":  c.GetOptionValue(enums.UsersKeyEnum.PASSWORD),
        }),
    )
    if err != nil {
        c.ExitWithError(err, "CreateUserCommand.Run")
    }
    c.PrintResult(createUserResultDto.ToArray())
}
