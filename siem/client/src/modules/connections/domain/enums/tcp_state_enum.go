package enums

// TcpStateEnum los estados de include/net/tcp_states.h tal como salen en la columna "st" de /proc/net/tcp (hex).
var TcpStateEnum = struct {
    ESTABLISHED string
    LISTEN      string
    // UDP no tiene estados: "07" (CLOSE) en un socket udp es un socket abierto sin conectar, "01" uno conectado
    UNCONN string
}{
    ESTABLISHED: "ESTABLISHED",
    LISTEN:      "LISTEN",
    UNCONN:      "UNCONN",
}

// TcpStateByHex "01" -> ESTABLISHED. Un valor que no esta aqui se muestra tal cual (hex).
var TcpStateByHex = map[string]string{
    "01": "ESTABLISHED",
    "02": "SYN_SENT",
    "03": "SYN_RECV",
    "04": "FIN_WAIT1",
    "05": "FIN_WAIT2",
    "06": "TIME_WAIT",
    "07": "CLOSE",
    "08": "CLOSE_WAIT",
    "09": "LAST_ACK",
    "0A": "LISTEN",
    "0B": "CLOSING",
    "0C": "NEW_SYN_RECV",
}
