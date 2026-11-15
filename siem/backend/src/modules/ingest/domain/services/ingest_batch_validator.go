package services

import (
    "regexp"

    "siem-backend/src/modules/ingest/domain/enums"
)

const MAX_EVENTS_PER_BATCH = 2000

var batchUuidPattern = regexp.MustCompile(`^[A-Za-z0-9-]{8,60}$`)

// IngestBatchValidator reglas de un lote que llega del siem-client. Devuelve la pila de errores (vacia = todo bien).
type IngestBatchValidator struct{}

func IngestBatchValidatorGetInstance() *IngestBatchValidator {
    return &IngestBatchValidator{}
}

// GetErrors eventsCount < 0 = "events" no era un array.
func (v *IngestBatchValidator) GetErrors(batchUuid string, eventsCount int) []string {
    errors := []string{}
    if !batchUuidPattern.MatchString(batchUuid) {
        errors = append(errors, enums.IngestMessageEnum.BATCH_UUID_IS_NOT_VALID)
    }
    if eventsCount < 0 {
        errors = append(errors, enums.IngestMessageEnum.EVENTS_MUST_BE_AN_ARRAY)
    }
    if eventsCount > MAX_EVENTS_PER_BATCH {
        errors = append(errors, enums.IngestMessageEnum.EVENTS_MAX_PER_BATCH_IS_2000)
    }
    return errors
}
