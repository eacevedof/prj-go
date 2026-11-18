package enums

// ProcNetFileEnum los ficheros de /proc/net que se leen y el protocolo / version de IP de cada uno.
var ProcNetFileEnum = struct {
    TCP  string
    TCP6 string
    UDP  string
    UDP6 string
}{
    TCP:  "tcp",
    TCP6: "tcp6",
    UDP:  "udp",
    UDP6: "udp6",
}

// IpVersionEnum version de IP de un socket (la dice el fichero de /proc/net del que sale) y bytes de cada direccion.
var IpVersionEnum = struct {
    V4       int
    V6       int
    V4_BYTES int
    V6_BYTES int
}{
    V4:       4,
    V6:       6,
    V4_BYTES: 4,
    V6_BYTES: 16,
}

// ProcNetColumnEnum columnas de una linea de /proc/net/{tcp,tcp6,udp,udp6} (separadas por espacios):
//
//   sl  local_address rem_address   st tx_queue:rx_queue tr:tm->when retrnsmt   uid  timeout inode
//    0: 0100007F:0277 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 22512 1 ...
//
// Las IP van en hex, en palabras de 32 bits en el orden de bytes del host (little-endian en amd64/arm64); los puertos,
// en hex normal.
var ProcNetColumnEnum = struct {
    MIN_COLUMNS       int
    SLOT_SUFFIX       string
    LOCAL_ADDRESS     int
    REMOTE_ADDRESS    int
    STATE             int
    UID               int
    INODE             int
    ENDPOINT_SEPARATOR string
    HEX_BASE          int
    PORT_BITS         int
    DECIMAL_BASE      int
    INODE_BITS        int
    ADDRESS_WORD_BYTES int
}{
    MIN_COLUMNS:        10,
    SLOT_SUFFIX:        ":",
    LOCAL_ADDRESS:      1,
    REMOTE_ADDRESS:     2,
    STATE:              3,
    UID:                7,
    INODE:              9,
    ENDPOINT_SEPARATOR: ":",
    HEX_BASE:           16,
    PORT_BITS:          16,
    DECIMAL_BASE:       10,
    INODE_BITS:         64,
    ADDRESS_WORD_BYTES: 4,
}

// ProcPathEnum rutas que se leen de /proc (sockets y de que proceso es cada uno) y de /etc/passwd (nombre de usuario).
var ProcPathEnum = struct {
    PROC_PATH   string
    NET_DIR     string
    FD_DIR      string
    COMM_FILE   string
    EXE_LINK    string
    PASSWD_PATH string
    // un fd que es un socket enlaza a "socket:[<inodo>]"
    SOCKET_LINK_PREFIX string
    SOCKET_LINK_SUFFIX string
    // exe de un binario borrado o sustituido mientras corre (una actualizacion)
    DELETED_EXE_SUFFIX string
}{
    PROC_PATH:          "/proc",
    NET_DIR:            "net",
    FD_DIR:             "fd",
    COMM_FILE:          "comm",
    EXE_LINK:           "exe",
    PASSWD_PATH:        "/etc/passwd",
    SOCKET_LINK_PREFIX: "socket:[",
    SOCKET_LINK_SUFFIX: "]",
    DELETED_EXE_SUFFIX: " (deleted)",
}

// PasswdColumnEnum columnas de una linea de /etc/passwd: nombre:x:uid:gid:...
var PasswdColumnEnum = struct {
    SEPARATOR   string
    MIN_COLUMNS int
    USER_NAME   int
    UID         int
}{
    SEPARATOR:   ":",
    MIN_COLUMNS: 3,
    USER_NAME:   0,
    UID:         2,
}
