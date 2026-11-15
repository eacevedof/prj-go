package enums

// EventSeverityEnum gravedad de un evento, de menos a mas. Mismos valores que el CHECK de events.severity.
var EventSeverityEnum = struct {
    DEBUG    string
    INFO     string
    WARNING  string
    ERROR    string
    CRITICAL string
    VALUES   []string
}{
    DEBUG:    "debug",
    INFO:     "info",
    WARNING:  "warning",
    ERROR:    "error",
    CRITICAL: "critical",
    VALUES:   []string{"debug", "info", "warning", "error", "critical"},
}
