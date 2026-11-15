package enums

// AgentStatusEnum estado de un agente segun su ultimo latido (AgentStatusResolver). No se guarda: se calcula al leer.
var AgentStatusEnum = struct {
    ONLINE     string
    OFFLINE    string
    NEVER_SEEN string
    DISABLED   string
}{
    ONLINE:     "online",
    OFFLINE:    "offline",
    NEVER_SEEN: "never_seen",
    DISABLED:   "disabled",
}
