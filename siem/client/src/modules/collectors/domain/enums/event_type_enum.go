package enums

// EventTypeEnum los dos tipos de evento que entiende el backend (EventTypeEnum de siem-backend: mismos valores).
var EventTypeEnum = struct {
    LOG_LINE       string
    COMMAND_RESULT string
}{
    LOG_LINE:       "log_line",
    COMMAND_RESULT: "command_result",
}
