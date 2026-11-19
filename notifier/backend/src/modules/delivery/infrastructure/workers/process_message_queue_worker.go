package workers

import (
    "fmt"
    "time"

    "notifier-backend/src/modules/shared/infrastructure/repositories/configuration"
    "notifier-backend/src/modules/shared/infrastructure/traits"

    "notifier-backend/src/modules/delivery/application/processmessagequeue"
    "notifier-backend/src/modules/delivery/domain/enums"
)

// ProcessMessageQueueWorker el "controller" de la cola de un canal dentro de la API: lo dispara el reloj. Cada
// APP_QUEUE_WORKER_INTERVAL_SECONDS llama a ProcessMessageQueueService hasta dejar sin turno pendiente la cola de su
// canal. Hay uno por canal, cada uno en su goroutine: si Telegram cae, el email sigue. Como todo controller, es el
// unico que captura: un fallo se loguea y la siguiente pasada lo vuelve a intentar.
type ProcessMessageQueueWorker struct {
    traits.LogTrait
    channel         string
    intervalSeconds int
}

func ProcessMessageQueueWorkerGetInstance(channel string) *ProcessMessageQueueWorker {
    return &ProcessMessageQueueWorker{
        LogTrait:        traits.NewLogTrait(nil),
        channel:         channel,
        intervalSeconds: configuration.GetInstance().GetQueueWorkerIntervalSeconds(),
    }
}

// IsEnabled APP_QUEUE_WORKER_INTERVAL_SECONDS=0 apaga los workers (la cola se vacia con ./console message-queue-process).
func (w *ProcessMessageQueueWorker) IsEnabled() bool {
    return w.intervalSeconds > 0
}

// Run bucle sin fin: se lanza en su propia goroutine al arrancar la API.
func (w *ProcessMessageQueueWorker) Run() {
    ticker := time.NewTicker(time.Duration(w.intervalSeconds) * time.Second)
    defer ticker.Stop()
    for range ticker.C {
        w.emptyQueue()
    }
}

func (w *ProcessMessageQueueWorker) emptyQueue() {
    logTitle := "ProcessMessageQueueWorker.emptyQueue " + w.channel
    defer func() {
        if recovered := recover(); recovered != nil {
            w.LogException(fmt.Errorf("%v", recovered), logTitle)
        }
    }()

    for {
        processMessageQueueResultDto, err := processmessagequeue.ProcessMessageQueueServiceGetInstance().Invoke(
            processmessagequeue.ProcessMessageQueueDtoFromPrimitives(map[string]any{
                "channel":     w.channel,
                "maxMessages": enums.DeliveryLimitEnum.MAX_MESSAGES_PER_WORKER_ROUND,
            }),
        )
        if err != nil {
            w.LogException(err, logTitle)
            return
        }
        if processMessageQueueResultDto.HasFailures() {
            w.LogWarning(processMessageQueueResultDto.ToArray(), logTitle+" failures")
        }
        // menos mensajes que el maximo, o ninguno tratado (todos de perfiles deshabilitados): sin bucle, a la
        // siguiente pasada
        if processMessageQueueResultDto.GetDueMessages() < enums.DeliveryLimitEnum.MAX_MESSAGES_PER_WORKER_ROUND ||
            processMessageQueueResultDto.GetHandledMessages() == 0 {
            return
        }
    }
}
