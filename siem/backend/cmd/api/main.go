// API de siem-backend (esqueleto de task-go). go run ./cmd/api -addr=:8080
// Tres cosas en un proceso: la API (/api/v1/...), el front (public/) y el worker que vacia ingest_queue.
package main

import (
    "flag"
    "log"
    "net/http"
    "time"

    "siem-backend/src/core/config"
    "siem-backend/src/core/routes"
    ingestworkers "siem-backend/src/modules/ingest/infrastructure/workers"
)

func main() {
    listenAddress := flag.String("addr", ":8080", "direccion de escucha")
    flag.Parse()

    if err := config.AppbootLoadenv(); err != nil {
        log.Fatal(err)
    }

    processIngestQueueWorker := ingestworkers.ProcessIngestQueueWorkerGetInstance()
    if processIngestQueueWorker.IsEnabled() {
        go processIngestQueueWorker.Run()
    }

    server := &http.Server{
        Addr:              *listenAddress,
        Handler:           PublicFilesServerGetInstance(IndexMainGetInstance(routes.GetRoutes())),
        ReadHeaderTimeout: 10 * time.Second,
    }
    log.Printf("siem-backend api listening on %s (queue worker enabled: %t)", *listenAddress, processIngestQueueWorker.IsEnabled())
    log.Fatal(server.ListenAndServe())
}
