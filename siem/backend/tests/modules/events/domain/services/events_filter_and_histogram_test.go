package services

import (
    "strings"
    "testing"
    "time"

    "siem-backend/tests"

    "siem-backend/src/modules/events/domain/entities"
    "siem-backend/src/modules/events/domain/services"
)

func TestHostMetricSeverityByMonitorThresholds(t *testing.T) {
    eventSeverityResolver := services.EventSeverityResolverGetInstance()
    thresholds := func(value float64) map[string]any {
        return map[string]any{"value": value, "warning_threshold": 80.0, "critical_threshold": 95.0}
    }
    tests.AssertEquals(t, "debug", eventSeverityResolver.GetSeverity("host_metric", "x", 0, "", thresholds(40)))
    tests.AssertEquals(t, "warning", eventSeverityResolver.GetSeverity("host_metric", "x", 0, "", thresholds(80)))
    tests.AssertEquals(t, "critical", eventSeverityResolver.GetSeverity("host_metric", "x", 0, "", thresholds(99.5)))
    // sin umbrales en el monitor nunca alerta; metrica no disponible, warning
    tests.AssertEquals(t, "debug", eventSeverityResolver.GetSeverity("host_metric", "x", 0, "", map[string]any{"value": 100.0}))
    tests.AssertEquals(t, "warning", eventSeverityResolver.GetSeverity("host_metric", "x", 0, "", map[string]any{"is_unavailable": true}))
}

func TestFilterParsesIncludedAndExcludedValues(t *testing.T) {
    eventsFilterEntity := entities.EventsFilterEntityFromPrimitives(map[string]any{
        "severity": "error, critical,,error", "not_event_type": "docker_stats", "search": " boom ", "from": "2026-11-15 00:00:00",
    })
    tests.AssertEquals(t, []string{"error", "critical"}, eventsFilterEntity.GetIncludedValues("severity"))
    tests.AssertEquals(t, []string{"docker_stats"}, eventsFilterEntity.GetExcludedValues("event_type"))
    tests.AssertEmpty(t, eventsFilterEntity.GetIncludedValues("event_type"))
    tests.AssertEquals(t, "boom", eventsFilterEntity.GetSearch())
    tests.AssertEmpty(t, services.EventsSearchValidatorGetInstance().GetFilterErrors(eventsFilterEntity))
}

func TestFilterErrors(t *testing.T) {
    eventsFilterEntity := entities.EventsFilterEntityFromPrimitives(map[string]any{
        "severity": "error,fatal", "not_event_type": "bogus", "source_type": "tape",
        "from": "2026-11-16 00:00:00", "to": "15/11/2026", "search": strings.Repeat("x", 201),
    })
    errors := strings.Join(services.EventsSearchValidatorGetInstance().GetFilterErrors(eventsFilterEntity), "\n")
    for _, expectedError := range []string{
        "severity: must be one of", "event_type: must be one of", "source_type: must be one of",
        "to: must be a UTC date time", "search: max length is 200",
    } {
        tests.AssertStringContainsString(t, expectedError, errors)
    }

    wrongOrder := entities.EventsFilterEntityFromPrimitives(map[string]any{"from": "2026-11-16 00:00:00", "to": "2026-11-15 00:00:00"})
    tests.AssertEquals(t, []string{"from: must be before to"}, services.EventsSearchValidatorGetInstance().GetFilterErrors(wrongOrder))
}

func TestHistogramRangeAndBuckets(t *testing.T) {
    histogramBucketCalculator := services.HistogramBucketCalculatorGetInstance()
    now := time.Date(2026, 11, 15, 18, 0, 0, 0, time.UTC)

    // sin rango: las ultimas 24 h, en tramos de 30 min (48 barras)
    fromTime, toTime := histogramBucketCalculator.GetRange("", "", now)
    tests.AssertEquals(t, "2026-11-14 18:00:00", histogramBucketCalculator.GetFormatted(fromTime.Unix()))
    tests.AssertEquals(t, 1800, histogramBucketCalculator.GetBucketSeconds(fromTime, toTime))
    tests.AssertEquals(t, 49, len(histogramBucketCalculator.GetBucketStarts(fromTime, toTime, 1800)))

    // 15 minutos: tramos de 30 s; un rango de un año: semanas
    fromTime, toTime = histogramBucketCalculator.GetRange("2026-11-15 17:45:00", "2026-11-15 18:00:00", now)
    tests.AssertEquals(t, 30, histogramBucketCalculator.GetBucketSeconds(fromTime, toTime))
    fromTime, toTime = histogramBucketCalculator.GetRange("2025-11-15 00:00:00", "2026-11-15 00:00:00", now)
    tests.AssertEquals(t, 604800, histogramBucketCalculator.GetBucketSeconds(fromTime, toTime))

    // los tramos van alineados a su ancho, como los agrupa el SQL
    fromTime, toTime = histogramBucketCalculator.GetRange("2026-11-15 17:02:13", "2026-11-15 17:10:00", now)
    bucketStarts := histogramBucketCalculator.GetBucketStarts(fromTime, toTime, 60)
    tests.AssertEquals(t, "2026-11-15 17:02:00", histogramBucketCalculator.GetFormatted(bucketStarts[0]))
    tests.AssertEquals(t, "2026-11-15 17:10:00", histogramBucketCalculator.GetFormatted(bucketStarts[len(bucketStarts)-1]))
}
