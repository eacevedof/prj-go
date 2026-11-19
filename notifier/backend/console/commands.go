// Package console registro de comandos de consola.
//
//  ./console <command-alias> --clave=valor .. --clave=valor
//  Ejemplo:
//      ./console help
package console

import (
    "notifier-backend/src/console/application"
    apitokenscommands "notifier-backend/src/modules/apitokens/infrastructure/commands"
    deliverycommands "notifier-backend/src/modules/delivery/infrastructure/commands"
    devopscommands "notifier-backend/src/modules/devops/infrastructure/commands"
    senderscommands "notifier-backend/src/modules/senders/infrastructure/commands"
    userscommands "notifier-backend/src/modules/users/infrastructure/commands"
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

        // usuarios, perfiles de envio y tokens (fase 1: por consola; el panel llega en la fase 2)
        {
            Alias: "user-create",
            Service: func(input []string) application.InterfaceConsole {
                return userscommands.CreateUserCommandGetInstance(input)
            },
            Info:    "alta de un usuario (root solo uno; sin --user_role, owner); imprime su user_uuid",
            Example: `--user_name="Ana" --user_email=ana@example.com --user_role=owner --password="una-clave-larga"`,
        },
        {
            Alias: "sender-create",
            Service: func(input []string) application.InterfaceConsole {
                return senderscommands.CreateSenderCommandGetInstance(input)
            },
            Info:    "alta de un perfil de envio (gmail, smtp_custom) y su primer token; imprime el api_token (solo esta vez)",
            Example: `--user_uuid=<uuid> --sender_name=avisos --provider_code=gmail --from_email=yo@gmail.com --from_name="Avisos" --secret="app password" --token_name=n8n`,
        },
        {
            Alias: "sender-list",
            Service: func(input []string) application.InterfaceConsole {
                return senderscommands.ListSendersCommandGetInstance(input)
            },
            Info:    "perfiles de envio (de un usuario o todos), sin secretos",
            Example: "--user_uuid=<uuid>",
        },
        {
            Alias: "token-create",
            Service: func(input []string) application.InterfaceConsole {
                return apitokenscommands.CreateApiTokenCommandGetInstance(input)
            },
            Info:    "otro token de envio para un perfil (uno por app); imprime el api_token (solo esta vez)",
            Example: `--sender_uuid=<uuid> --token_name=web --default_to="a@example.com,b@example.com"`,
        },
        {
            Alias: "token-disable",
            Service: func(input []string) application.InterfaceConsole {
                return apitokenscommands.DisableApiTokenCommandGetInstance(input)
            },
            Info:    "revoca un token de envio (da 401 desde ese momento)",
            Example: "--token_uuid=<uuid>",
        },

        // mantenimiento
        {
            Alias: "message-queue-process",
            Service: func(input []string) application.InterfaceConsole {
                return deliverycommands.ProcessMessageQueueCommandGetInstance(input)
            },
            Info:    "entrega los mensajes pendientes de un canal (lo mismo que el worker de la API)",
            Example: "--channel=email --max_messages=50",
        },
    }
}
