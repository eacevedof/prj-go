// API de notifier-backend (esqueleto de task-go). go run ./cmd/api -addr=:8080
// En un proceso: la API (/api/v1/...), los estaticos de public/ y un worker por canal que vacia su cola de messages.
package main

import (
    "flag"
    "log"
    "net/http"
    "time"

    "notifier-backend/src/core/config"
    "notifier-backend/src/core/routes"
    "notifier-backend/src/modules/delivery/application/processmessagequeue"
    deliveryworkers "notifier-backend/src/modules/delivery/infrastructure/workers"
)

func main() {
    listenAddress := flag.String("addr", ":8080", "direccion de escucha")
    flag.Parse()

    if err := config.AppbootLoadenv(); err != nil {
        log.Fatal(err)
    }

    enabledWorkers := 0
    for _, channel := range processmessagequeue.ChannelsWithWorker {
        processMessageQueueWorker := deliveryworkers.ProcessMessageQueueWorkerGetInstance(channel)
        if processMessageQueueWorker.IsEnabled() {
            go processMessageQueueWorker.Run()
            enabledWorkers++
        }
    }

    server := &http.Server{
        Addr:              *listenAddress,
        Handler:           PublicFilesServerGetInstance(IndexMainGetInstance(routes.GetRoutes())),
        ReadHeaderTimeout: 10 * time.Second,
    }
    log.Printf("notifier-backend api listening on %s (queue workers: %d)", *listenAddress, enabledWorkers)
    log.Fatal(server.ListenAndServe())
}
