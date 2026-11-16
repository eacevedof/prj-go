// Consola de trafficker (mismo esqueleto que task-go). go run ./cmd/console help
package main

import (
    "fmt"
    "os"

    "trafficker/console"
    "trafficker/src/core/config"
)

func main() {
    if err := config.AppbootLoadenv(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(2)
    }
    ConsoleMainGetInstance(os.Args, console.GetCommands()).Exec()
}
