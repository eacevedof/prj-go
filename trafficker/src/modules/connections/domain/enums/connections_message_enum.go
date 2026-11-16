package enums

// ConnectionsMessageEnum mensajes del modulo (en ingles, como el resto de la consola).
var ConnectionsMessageEnum = struct {
    PROC_NET_NOT_READABLE         string
    PROCESSES_PARTIALLY_READABLE  string
}{
    // {0} ruta de /proc/net · {1} error
    PROC_NET_NOT_READABLE: "cannot read the sockets in {0} (trafficker 0.1.0 only reads Linux /proc/net): {1}",
    // {0} procesos cuyos sockets no se han podido leer
    PROCESSES_PARTIALLY_READABLE: "{0} processes could not be read: their sockets show process ?. Run as root (in a container, also cap_add SYS_PTRACE and apparmor:unconfined)",
}
