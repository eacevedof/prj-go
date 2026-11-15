package services

import (
    "testing"

    "siem-client/src/modules/collectors/domain/services"
    "siem-client/tests"
)

func TestEventsAreSplitInBatches(t *testing.T) {
    events := make([]map[string]any, 7)
    batches := services.EventsBatcherGetInstance().GetBatches(events, 3)
    tests.AssertEquals(t, 3, len(batches))
    tests.AssertEquals(t, 1, len(batches[2]))
    tests.AssertEmpty(t, services.EventsBatcherGetInstance().GetBatches([]map[string]any{}, 3))
}
