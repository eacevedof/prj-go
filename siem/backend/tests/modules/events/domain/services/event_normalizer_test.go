package services

import (
    "strings"
    "testing"

    "siem-backend/tests"

    "siem-backend/src/modules/events/domain/enums"
    "siem-backend/src/modules/events/domain/services"
)

const receivedAt = "2026-11-15 10:00:00"

func TestUnknownEventTypeOrMissingSourceIsDiscarded(t *testing.T) {
    eventNormalizer := services.EventNormalizerGetInstance()
    tests.AssertNil(t, nilIfNil(eventNormalizer.GetEventEntityOrNil(map[string]any{"event_type": "bogus", "source_name": "x", "message": "y"}, 1, 1, receivedAt)))
    tests.AssertNil(t, nilIfNil(eventNormalizer.GetEventEntityOrNil(map[string]any{"event_type": "log_line", "message": "y"}, 1, 1, receivedAt)))
    tests.AssertNil(t, nilIfNil(eventNormalizer.GetEventEntityOrNil(map[string]any{"event_type": "log_line", "source_name": "x", "message": "  "}, 1, 1, receivedAt)))
}

func TestCommandWithoutOutputGetsPlaceholder(t *testing.T) {
    eventEntity := services.EventNormalizerGetInstance().GetEventEntityOrNil(
        map[string]any{"event_type": "command_result", "source_name": "disk", "exit_code": float64(1)}, 1, 1, receivedAt,
    )
    tests.AssertEquals(t, services.EMPTY_COMMAND_OUTPUT, eventEntity.GetMessage())
    tests.AssertEquals(t, "command", eventEntity.GetSourceType())
    tests.AssertEquals(t, "error", eventEntity.GetSeverity())
}

func TestOccurredAtFallsBackToReceivedAt(t *testing.T) {
    eventNormalizer := services.EventNormalizerGetInstance()
    for occurredAt, expectedOccurredAt := range map[string]string{
        "2026-11-15 09:59:00": "2026-11-15 09:59:00",
        "not a date":          receivedAt,
        "2026-11-15 12:00:00": receivedAt, // del futuro: reloj del cliente descolocado
    } {
        eventEntity := eventNormalizer.GetEventEntityOrNil(
            map[string]any{"event_type": "log_line", "source_name": "app", "message": "hi", "occurred_at": occurredAt}, 1, 1, receivedAt,
        )
        tests.AssertEquals(t, expectedOccurredAt, eventEntity.GetOccurredAt(), occurredAt)
    }
}

func TestLongMessageIsTruncated(t *testing.T) {
    eventEntity := services.EventNormalizerGetInstance().GetEventEntityOrNil(
        map[string]any{"event_type": "log_line", "source_name": "app", "message": strings.Repeat("ñ", enums.EventsLimitEnum.MESSAGE_MAX_LENGTH+10)}, 1, 1, receivedAt,
    )
    tests.AssertEquals(t, enums.EventsLimitEnum.MESSAGE_MAX_LENGTH, len([]rune(eventEntity.GetMessage())))
}

func TestDockerEventsGetTheirSourceTypeAndAttributes(t *testing.T) {
    eventNormalizer := services.EventNormalizerGetInstance()
    for eventType, expectedSourceType := range map[string]string{
        "docker_state": "docker", "docker_event": "docker", "docker_stats": "docker_stats", "docker_log": "docker_log",
    } {
        eventEntity := eventNormalizer.GetEventEntityOrNil(map[string]any{
            "event_type": eventType, "source_name": "cont-web", "message": "m",
            "attributes": map[string]any{"status": "running"},
        }, 1, 1, receivedAt)
        tests.AssertEquals(t, expectedSourceType, eventEntity.GetSourceType(), eventType)
        tests.AssertEquals(t, `{"status":"running"}`, eventEntity.GetAttributes(), eventType)
    }
}

func TestAttributesThatAreNotAnObjectOrTooBigAreIgnored(t *testing.T) {
    eventNormalizer := services.EventNormalizerGetInstance()
    notAnObject := eventNormalizer.GetEventEntityOrNil(map[string]any{
        "event_type": "docker_state", "source_name": "c", "message": "m", "attributes": "status=running",
    }, 1, 1, receivedAt)
    tests.AssertEquals(t, "", notAnObject.GetAttributes())

    tooBig := eventNormalizer.GetEventEntityOrNil(map[string]any{
        "event_type": "docker_state", "source_name": "c", "message": "m",
        "attributes": map[string]any{"blob": strings.Repeat("x", enums.EventsLimitEnum.ATTRIBUTES_MAX_BYTES)},
    }, 1, 1, receivedAt)
    tests.AssertEquals(t, "", tooBig.GetAttributes())
}

// nilIfNil un *EventEntity nil dentro de un any no es nil: se normaliza para AssertNil.
func nilIfNil[T any](value *T) any {
    if value == nil {
        return nil
    }
    return value
}
