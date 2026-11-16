package commands

import (
    "errors"

    sharedcommands "siem-backend/src/modules/shared/infrastructure/commands"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"

    "siem-backend/src/modules/agents/domain/enums"
    "siem-backend/src/modules/agents/application/createagent"
    "siem-backend/src/modules/agents/domain/exceptions"
)

// CreateAgentCommand ./console agent-create --project_uuid=<uuid> --agent_name=web-01 --heartbeat_interval_seconds=60
// Imprime el agente con su api_token en claro: es la unica vez que se ve (va al .env del siem-client).
type CreateAgentCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func CreateAgentCommandGetInstance(input []string) *CreateAgentCommand {
    return &CreateAgentCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *CreateAgentCommand) Run() {
    createAgentResultDto, err := createagent.CreateAgentServiceGetInstance().Invoke(
        createagent.CreateAgentDtoFromPrimitives(map[string]any{
            "projectUuid":              c.GetOptionValue(enums.AgentsKeyEnum.PROJECT_UUID),
            "agentName":                c.GetOptionValue(enums.AgentsKeyEnum.AGENT_NAME),
            "heartbeatIntervalSeconds": c.GetOptionValue(enums.AgentsKeyEnum.HEARTBEAT_INTERVAL_SECONDS),
        }),
    )
    if err == nil {
        c.PrintResult(createAgentResultDto.ToArray())
        return
    }

    var fieldsException *sharedexceptions.FieldsException
    var agentsException *exceptions.AgentsException
    switch {
    case errors.As(err, &fieldsException):
        c.ExitWithDomainError(fieldsException)
    case errors.As(err, &agentsException):
        c.ExitWithDomainError(agentsException)
    default:
        c.ExitWithUnexpectedError(err, "CreateAgentCommand.Run")
    }
}
