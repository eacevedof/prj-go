package services

import "notifier-backend/src/modules/delivery/domain/enums"

// RetryScheduler cuanto espera un mensaje tras un intento fallido y cuando se deja de intentar.
type RetryScheduler struct{}

func RetrySchedulerGetInstance() *RetryScheduler {
    return &RetryScheduler{}
}

// IsLastAttempt attemptsDone incluye el que acaba de fallar.
func (s *RetryScheduler) IsLastAttempt(attemptsDone int, maxAttempts int) bool {
    return attemptsDone >= maxAttempts
}

// GetDelaySeconds la espera tras el intento attemptsDone (1 = el primero). Mas alla de la lista, la ultima.
func (s *RetryScheduler) GetDelaySeconds(attemptsDone int) int {
    retryDelays := enums.DeliveryLimitEnum.RETRY_DELAYS_SECONDS
    if attemptsDone < 1 {
        return retryDelays[0]
    }
    if attemptsDone > len(retryDelays) {
        return retryDelays[len(retryDelays)-1]
    }
    return retryDelays[attemptsDone-1]
}
