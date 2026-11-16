// Package console registro de comandos de consola.
//
//  ./trafficker <command-alias> --clave=valor .. --clave=valor
package console

import (
    connectionscommands "trafficker/src/modules/connections/infrastructure/commands"
    "trafficker/src/console/application"
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
            Info:    "version del agente y fecha",
            Example: "",
        },
        {
            Alias: "connections-list",
            Service: func(input []string) application.InterfaceConsole {
                return connectionscommands.ListConnectionsCommandGetInstance(input)
            },
            Info:    "conexiones y escuchas de ahora con proceso y usuario; sin .env tambien pais y organizacion (whois RDAP). --filter: terminos AND, 'remote' = IP publica; --whois=0|1; --format=json. Sin .env es lo que hace ./trafficker sin argumentos",
            Example: `--filter="established,remote"`,
        },
    }
}
