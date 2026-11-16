package enums

// AgentHostOsEnum los host_os que admite el panel: el runtime.GOOS que informa el siem-client en cada lote. Cualquier
// otro valor se guarda vacio (AgentHostDataSanitizer). Mismos valores que ClientOsEnum del front.
var AgentHostOsEnum = struct {
    LINUX   string
    WINDOWS string
    MACOS   string
}{
    LINUX:   "linux",
    WINDOWS: "windows",
    MACOS:   "darwin",
}
