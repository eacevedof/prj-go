package commands

import (
    "errors"

    "siem-backend/src/modules/projects/application/createproject"
    "siem-backend/src/modules/projects/domain/enums"
    "siem-backend/src/modules/projects/domain/exceptions"
    sharedcommands "siem-backend/src/modules/shared/infrastructure/commands"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"
)

// CreateProjectCommand ./console project-create --project_name="Track Alert" --project_description="prod"
type CreateProjectCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func CreateProjectCommandGetInstance(input []string) *CreateProjectCommand {
    return &CreateProjectCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *CreateProjectCommand) Run() {
    createProjectResultDto, err := createproject.CreateProjectServiceGetInstance().Invoke(
        createproject.CreateProjectDtoFromPrimitives(map[string]any{
            "projectName":        c.GetOptionValue(enums.ProjectsKeyEnum.PROJECT_NAME),
            "projectDescription": c.GetOptionValue(enums.ProjectsKeyEnum.PROJECT_DESCRIPTION),
        }),
    )
    if err == nil {
        c.PrintResult(createProjectResultDto.ToArray())
        return
    }

    var fieldsException *sharedexceptions.FieldsException
    var projectsException *exceptions.ProjectsException
    switch {
    case errors.As(err, &fieldsException):
        c.ExitWithDomainError(fieldsException)
    case errors.As(err, &projectsException):
        c.ExitWithDomainError(projectsException)
    default:
        c.ExitWithUnexpectedError(err, "CreateProjectCommand.Run")
    }
}
