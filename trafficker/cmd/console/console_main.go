package main

import (
    "errors"
    "fmt"
    "os"
    "runtime/debug"
    "strings"

    "trafficker/src/console/application"
    "trafficker/src/core/boot"
    "trafficker/src/core/config"
    "trafficker/src/modules/shared/infrastructure/components/uuid"
    "trafficker/src/modules/shared/infrastructure/repositories/configuration"
)

// ConsoleMain ./console <alias> --clave=valor ... -> el command del alias.
type ConsoleMain struct {
    argv           []string
    commandConfigs []application.CommandConfig
}

func ConsoleMainGetInstance(argv []string, commandConfigs []application.CommandConfig) *ConsoleMain {
    boot.SetSapiNameCli()
    boot.EnableConsoleColors()
    return &ConsoleMain{argv: argv, commandConfigs: commandConfigs}
}

func (m *ConsoleMain) Exec() {
    defer func() {
        if recovered := recover(); recovered != nil {
            err := fmt.Errorf("%v", recovered)
            debugError(err)
            onError(err)
            os.Exit(2)
        }
    }()

    if err := m.invokeConsoleServiceWithArgs(); err != nil {
        debugError(err)
        onError(err)
        os.Exit(2)
    }
}

func (m *ConsoleMain) invokeConsoleServiceWithArgs() error {
    alias := ""
    if len(m.argv) > 1 {
        alias = strings.TrimSpace(m.argv[1])
    }
    if alias == "" {
        alias = m.getDefaultAlias()
    }

    for _, commandConfig := range m.commandConfigs {
        if commandConfig.Alias != alias {
            continue
        }
        // a partir de la pos 2 en adelante son parametros de input
        var input []string
        if len(m.argv) > 2 {
            input = m.argv[2:]
        }
        commandConfig.Service(input).Run()
        return nil
    }
    return errors.New("no command found for alias " + alias)
}

// getDefaultAlias sin argumentos: en modo local (sin .env) la lista de conexiones, como el localhost-traffic de Rust;
// con .env, la ayuda.
func (m *ConsoleMain) getDefaultAlias() string {
    if configuration.GetInstance().IsLocalMode() {
        return "connections-list"
    }
    return "help"
}

func onError(err error) {
    errorUuid := uuid.GetInstance().GetUniqueIdWithPrefix("")
    config.Lgerr(err.Error(), "console-exception "+errorUuid, "error")
    config.Lgerr(string(debug.Stack()), "console-exception "+errorUuid+" TRACE", "error")
}

func debugError(err error) {
    environmentReaderRawRepository := configuration.GetInstance()
    if environmentReaderRawRepository.IsProduction() || !environmentReaderRawRepository.IsDebug() {
        fmt.Fprintln(os.Stderr, err.Error())
        return
    }
    fmt.Printf("Exception: %s\nresponse: 500\n", err.Error())
}
