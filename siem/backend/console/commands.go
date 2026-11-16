// Package console registro de comandos de consola.
//
//  ./console <command-alias> --clave=valor .. --clave=valor
//  Ejemplo:
//      ./console help
package console

import (
    "siem-backend/src/console/application"
    agentscommands "siem-backend/src/modules/agents/infrastructure/commands"
    devopscommands "siem-backend/src/modules/devops/infrastructure/commands"
    eventscommands "siem-backend/src/modules/events/infrastructure/commands"
    ingestcommands "siem-backend/src/modules/ingest/infrastructure/commands"
    projectscommands "siem-backend/src/modules/projects/infrastructure/commands"
)

// GetCommands alias => command, que hace (una linea) y los argumentos de un ejemplo. El orden es el del help.
func GetCommands() []application.CommandConfig {
    return []application.CommandConfig{
        {
            Alias: "help",
            Service: func(input []string) application.InterfaceConsole {
                return application.HelpCommandGetInstance(input, GetCommands())
            },
            Info:    "lista los comandos con un ejemplo de cada uno",
            Example: "",
        },

        {
            Alias: "run-migrations",
            Service: func(input []string) application.InterfaceConsole {
                return devopscommands.RunMigrationsCommandGetInstance(input)
            },
            Info:    "ejecuta los .sql pendientes de database/migrations (tabla migrations)",
            Example: "",
        },

        // alta de lo que se monitoriza: un proyecto y, dentro, un agente por host con siem-client
        {
            Alias: "project-create",
            Service: func(input []string) application.InterfaceConsole {
                return projectscommands.CreateProjectCommandGetInstance(input)
            },
            Info:    "alta de un proyecto; imprime su project_uuid",
            Example: `--project_name="Mi app" --project_description="produccion"`,
        },
        {
            Alias: "agent-create",
            Service: func(input []string) application.InterfaceConsole {
                return agentscommands.CreateAgentCommandGetInstance(input)
            },
            Info:    "alta de un agente (un host con siem-client); imprime su api_token (solo esta vez)",
            Example: "--project_uuid=<uuid> --agent_name=web-01 --heartbeat_interval_seconds=60",
        },

        // mantenimiento
        {
            Alias: "ingest-queue-process",
            Service: func(input []string) application.InterfaceConsole {
                return ingestcommands.ProcessIngestQueueCommandGetInstance(input)
            },
            Info:    "pasa a events los lotes pendientes de ingest_queue (lo mismo que el worker de la API)",
            Example: "--max_batches=100",
        },
        {
            Alias: "events-purge",
            Service: func(input []string) application.InterfaceConsole {
                return eventscommands.PurgeEventsCommandGetInstance(input)
            },
            Info:    "borra los eventos con mas de N dias (por defecto APP_EVENTS_RETENTION_DAYS) y los lotes ya procesados",
            Example: "--retention_days=30",
        },
    }
}
