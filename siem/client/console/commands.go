// Package console registro de comandos de consola.
//
//  ./siem-client <command-alias> --clave=valor .. --clave=valor
package console

import (
    "siem-client/src/console/application"

    agentconfigcommands "siem-client/src/modules/agentconfig/infrastructure/commands"

    collectorscommands "siem-client/src/modules/collectors/infrastructure/commands"

    hostmetricscommands "siem-client/src/modules/hostmetrics/infrastructure/commands"

    connectionscommands "siem-client/src/modules/connections/infrastructure/commands"

    dockercommands "siem-client/src/modules/docker/infrastructure/commands"

    shippercommands "siem-client/src/modules/shipper/infrastructure/commands"

    agentcommands "siem-client/src/modules/agent/infrastructure/commands"
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
            Alias: "version",
            Service: func(input []string) application.InterfaceConsole {
                return application.VersionCommandGetInstance(input)
            },
            Info:    "version del cliente y fecha",
            Example: "",
        },
        {
            Alias: "config-check",
            Service: func(input []string) application.InterfaceConsole {
                return agentconfigcommands.CheckAgentConfigCommandGetInstance(input)
            },
            Info:    "valida la config (siem-client.json + monitores del panel ya bajados) y el endpoint/token del .env",
            Example: "",
        },

        {
            Alias: "config-sync",
            Service: func(input []string) application.InterfaceConsole {
                return agentconfigcommands.SyncRemoteConfigCommandGetInstance(input)
            },
            Info:    "baja del backend los monitores de este agente (lo que se crea en el panel) a state/remote-config.json",
            Example: "",
        },

        // el agente: lo que se deja corriendo como servicio
        {
            Alias: "agent-run",
            Service: func(input []string) application.InterfaceConsole {
                return agentcommands.RunAgentCommandGetInstance(input)
            },
            Info:    "bucle del agente: sincroniza la config, recoge (ficheros, comandos, metricas, conexiones, docker) y envia hasta Ctrl+C",
            Example: "",
        },

        // los pasos del bucle, sueltos (pruebas, cron)
        {
            Alias: "files-collect",
            Service: func(input []string) application.InterfaceConsole {
                return collectorscommands.CollectFileLinesCommandGetInstance(input)
            },
            Info:    "una pasada por las carpetas vigiladas: lineas nuevas al spool",
            Example: "",
        },
        {
            Alias: "commands-run",
            Service: func(input []string) application.InterfaceConsole {
                return collectorscommands.RunDueCommandsCommandGetInstance(input)
            },
            Info:    "lanza los comandos a los que les toca (--force=1: todos) y deja el resultado en el spool",
            Example: "--force=1",
        },
        {
            Alias: "metrics-collect",
            Service: func(input []string) application.InterfaceConsole {
                return hostmetricscommands.CollectHostMetricsCommandGetInstance(input)
            },
            Info:    "muestrea las metricas del host a las que les toca (--force=1: todas) y deja el resultado en el spool",
            Example: "--force=1",
        },
        {
            Alias: "connections-collect",
            Service: func(input []string) application.InterfaceConsole {
                return connectionscommands.CollectConnectionChangesCommandGetInstance(input)
            },
            Info:    "una pasada del monitor connections (--force=1: ya): conexiones que abren o cierran al spool (solo Linux)",
            Example: "--force=1",
        },
        {
            Alias: "docker-collect",
            Service: func(input []string) application.InterfaceConsole {
                return dockercommands.CollectDockerCommandGetInstance(input)
            },
            Info:    "una pasada por Docker (estados, eventos del daemon, recursos si tocan, logs) al spool",
            Example: "--force_stats=1",
        },
        {
            Alias: "events-send",
            Service: func(input []string) application.InterfaceConsole {
                return shippercommands.SendSpooledEventsCommandGetInstance(input)
            },
            Info:    "envia el spool al backend (o el latido si esta vacio)",
            Example: "--max_batches=20",
        },
    }
}
