package enums

// ProjectHealthEnum el semaforo de un proyecto en el dashboard (ProjectHealthResolver). No se guarda: se calcula.
var ProjectHealthEnum = struct {
    OK       string
    WARNING  string
    CRITICAL string
    EMPTY    string
}{
    OK:       "ok",
    WARNING:  "warning",
    CRITICAL: "critical",
    EMPTY:    "empty",
}
