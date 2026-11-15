package commands

import (
    "bytes"
    "encoding/json"
    "fmt"
    "os"
    "strings"

    "siem-client/src/core/config"
    "siem-client/src/modules/shared/domain/enums"
    "siem-client/src/modules/shared/infrastructure/exceptions"
    "siem-client/src/modules/shared/infrastructure/traits"
)

const (
    EXIT_CODE_DOMAIN_ERROR         = 1
    EXIT_CODE_INFRASTRUCTURE_ERROR = 2
)

// AbstractConsoleCommand el controller de la consola (vault: CLI Commands as Controllers): lee --clave=valor, imprime
// el ResultDto en JSON por stdout (sin color: se puede pasar a jq) y los errores en rojo por stderr.
// 0 ok · 1 error de dominio (config mal, backend rechaza) · 2 error de infraestructura (se loguea).
type AbstractConsoleCommand struct {
    traits.CliColorTrait
    traits.LogTrait
    // los argumentos que van tras el alias del comando
    input []string
}

func NewAbstractConsoleCommand(input []string) AbstractConsoleCommand {
    return AbstractConsoleCommand{LogTrait: traits.NewLogTrait(), input: input}
}

// GetOptionValue --name="app logs" -> "app logs". Si no se pasa, "" (el servicio decide si es obligatorio).
func (c *AbstractConsoleCommand) GetOptionValue(optionName string) string {
    optionPrefix := "--" + optionName + "="
    for _, argument := range c.input {
        if strings.HasPrefix(argument, optionPrefix) {
            return strings.Trim(argument[len(optionPrefix):], `"'`)
        }
    }
    return ""
}

// PrintResult JSON con 4 espacios sin escapar unicode ni barras.
func (c *AbstractConsoleCommand) PrintResult(result any) {
    var jsonBuffer bytes.Buffer
    jsonEncoder := json.NewEncoder(&jsonBuffer)
    jsonEncoder.SetEscapeHTML(false)
    jsonEncoder.SetIndent("", "    ")
    _ = jsonEncoder.Encode(result)
    fmt.Print(jsonBuffer.String())
}

// PrintDomainError el error de dominio sin salir: lo usa el bucle del agente, que tiene que seguir vivo.
func (c *AbstractConsoleCommand) PrintDomainError(ex exceptions.InterfaceException) {
    c.printErrors(ex.GetCode(), []string{ex.Error()})
}

// PrintUnexpectedError como ExitWithUnexpectedError pero sin salir (bucle del agente): se loguea y se sigue.
func (c *AbstractConsoleCommand) PrintUnexpectedError(err error, logTitle string) {
    c.LogException(err, logTitle)
    c.printErrors(enums.ExceptionEnum.CODE_INTERNAL_SERVER_ERROR, []string{logTitle + ": " + err.Error()})
}

func (c *AbstractConsoleCommand) ExitWithDomainError(ex exceptions.InterfaceException) {
    c.PrintDomainError(ex)
    os.Exit(EXIT_CODE_DOMAIN_ERROR)
}

// ExitWithUnexpectedError se loguea con el metodo que fallo; quien tiene la consola si ve el mensaje real.
func (c *AbstractConsoleCommand) ExitWithUnexpectedError(err error, logTitle string) {
    c.LogException(err, logTitle)
    c.LogDebugCli(err.Error(), logTitle)
    c.printErrors(
        enums.ExceptionEnum.CODE_INTERNAL_SERVER_ERROR,
        []string{config.Tr(enums.SharedMessageEnum.AN_UNEXPECTED_ERROR_HAS_OCCURRED)},
    )
    os.Exit(EXIT_CODE_INFRASTRUCTURE_ERROR)
}

func (c *AbstractConsoleCommand) printErrors(code int, errors []string) {
    fmt.Fprintln(os.Stderr, c.GetRed(fmt.Sprintf("Error (%d):", code)))
    for _, errorMessage := range errors {
        fmt.Fprintln(os.Stderr, c.GetRed("  - "+errorMessage))
    }
}
