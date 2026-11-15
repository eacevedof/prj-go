package workers

import (
    "fmt"
    "time"

    "siem-backend/src/modules/ingest/application/processingestqueue"
    "siem-backend/src/modules/shared/infrastructure/repositories/configuration"
    "siem-backend/src/modules/shared/infrastructure/traits"
)

const maxBatchesPerRound = 50

// ProcessIngestQueueWorker el "controller" de la cola dentro de la API: no entra por HTTP ni por consola, lo dispara
// el reloj. Cada APP_QUEUE_WORKER_INTERVAL_SECONDS llama a ProcessIngestQueueService hasta dejar la cola vacia.
// Como todo controller, es el unico que captura: un fallo se loguea y la siguiente pasada lo vuelve a intentar.
type ProcessIngestQueueWorker struct {
    traits.LogTrait
    intervalSeconds int
}

func ProcessIngestQueueWorkerGetInstance() *ProcessIngestQueueWorker {
    return &ProcessIngestQueueWorker{
        LogTrait:        traits.NewLogTrait(nil),
        intervalSeconds: configuration.GetInstance().GetQueueWorkerIntervalSeconds(),
    }
}

// IsEnabled APP_QUEUE_WORKER_INTERVAL_SECONDS=0 apaga el worker (la cola se vacia con ./console ingest-queue-process).
func (w *ProcessIngestQueueWorker) IsEnabled() bool {
    return w.intervalSeconds > 0
}

// Run bucle sin fin: se lanza en su propia goroutine al arrancar la API.
func (w *ProcessIngestQueueWorker) Run() {
    ticker := time.NewTicker(time.Duration(w.intervalSeconds) * time.Second)
    defer ticker.Stop()
    for range ticker.C {
        w.processUntilQueueIsEmpty()
    }
}

func (w *ProcessIngestQueueWorker) processUntilQueueIsEmpty() {
    const logTitle = "ProcessIngestQueueWorker.processUntilQueueIsEmpty"
    defer func() {
        if recovered := recover(); recovered != nil {
            w.LogException(fmt.Errorf("%v", recovered), logTitle)
        }
    }()

    for {
        processIngestQueueResultDto, err := processingestqueue.ProcessIngestQueueServiceGetInstance().Invoke(
            processingestqueue.ProcessIngestQueueDtoFromPrimitives(map[string]any{"maxBatches": maxBatchesPerRound}),
        )
        if err != nil {
            w.LogException(err, logTitle)
            return
        }
        if processIngestQueueResultDto.GetFailedBatches() > 0 {
            w.LogWarning(processIngestQueueResultDto.ToArray(), logTitle+" failed batches")
        }
        // menos lotes que el maximo, o alguno fallido (seguira pending): la cola queda para la siguiente pasada
        if processIngestQueueResultDto.GetProcessedBatches() < maxBatchesPerRound {
            return
        }
    }
}
