package enums

// HostMetricEventTypeEnum el tipo de evento de una metrica del host (EventTypeEnum de siem-backend).
const HOST_METRIC_EVENT_TYPE = "host_metric"

// HostMetricAttributeEnum claves de attributes de un evento host_metric. monitor_uuid, value e is_unavailable las
// usa el backend para decidir la gravedad (con los umbrales del monitor); el resto son datos propios de cada tipo.
var HostMetricAttributeEnum = struct {
    MONITOR_UUID   string
    MONITOR_TYPE   string
    VALUE          string
    UNIT           string
    IS_UNAVAILABLE string

    // cpu
    LOAD_1M  string
    LOAD_5M  string
    LOAD_15M string

    // memory, disk
    USED_BYTES  string
    TOTAL_BYTES string
    PATH        string

    // network
    INTERFACE string
    RX_MBIT   string
    TX_MBIT   string

    // gpu
    GPU_NAME           string
    GPU_MEMORY_PERCENT string
    GPU_TEMPERATURE    string
    GPU_COUNT          string
}{
    MONITOR_UUID:   "monitor_uuid",
    MONITOR_TYPE:   "monitor_type",
    VALUE:          "value",
    UNIT:           "unit",
    IS_UNAVAILABLE: "is_unavailable",

    LOAD_1M:  "load_1m",
    LOAD_5M:  "load_5m",
    LOAD_15M: "load_15m",

    USED_BYTES:  "used_bytes",
    TOTAL_BYTES: "total_bytes",
    PATH:        "path",

    INTERFACE: "interface",
    RX_MBIT:   "rx_mbit",
    TX_MBIT:   "tx_mbit",

    GPU_NAME:           "gpu_name",
    GPU_MEMORY_PERCENT: "gpu_memory_percent",
    GPU_TEMPERATURE:    "gpu_temperature",
    GPU_COUNT:          "gpu_count",
}

// MetricUnitEnum la unidad de value.
var MetricUnitEnum = struct {
    PERCENT         string
    MBIT_PER_SECOND string
}{
    PERCENT:         "percent",
    MBIT_PER_SECOND: "mbit_per_second",
}

// HostMetricsStateEnum state/host-metrics.json: por monitor, cuando se muestreo, si ya se aviso de que no esta
// disponible y los contadores de la muestra anterior (cpu y red se calculan por diferencia).
var HostMetricsStateEnum = struct {
    STATE_FILE string

    MONITORS       string
    LAST_RUN_UNIX  string
    IS_UNAVAILABLE string
    CPU_TOTAL      string
    CPU_IDLE       string
    NET_RX_BYTES   string
    NET_TX_BYTES   string
    NET_AT_MS      string
}{
    STATE_FILE: "host-metrics.json",

    MONITORS:       "monitors",
    LAST_RUN_UNIX:  "last_run_unix",
    IS_UNAVAILABLE: "is_unavailable",
    CPU_TOTAL:      "cpu_total",
    CPU_IDLE:       "cpu_idle",
    NET_RX_BYTES:   "net_rx_bytes",
    NET_TX_BYTES:   "net_tx_bytes",
    NET_AT_MS:      "net_at_ms",
}

// HostMetricsMessageEnum motivos por los que una metrica no esta disponible (van en el mensaje del evento).
var HostMetricsMessageEnum = struct {
    OS_X_IS_NOT_SUPPORTED     string
    CANNOT_READ_X             string
    X_HAS_AN_UNEXPECTED_FORMAT string
    PATH_X_NOT_AVAILABLE_X    string
    INTERFACE_X_NOT_FOUND     string
    NVIDIA_SMI_NOT_AVAILABLE_X string
}{
    OS_X_IS_NOT_SUPPORTED:      "host metrics are not supported on {0} (linux only)",
    CANNOT_READ_X:              "cannot read {0}",
    X_HAS_AN_UNEXPECTED_FORMAT: "{0} has an unexpected format",
    PATH_X_NOT_AVAILABLE_X:     "path {0} not available: {1}",
    INTERFACE_X_NOT_FOUND:      "network interface {0} not found",
    NVIDIA_SMI_NOT_AVAILABLE_X: "nvidia-smi not available: {0}",
}
