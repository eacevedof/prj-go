/**
 * Tipos de monitor (recurso a vigilar) del contrato de la API. Los de metricas dan un valor cada interval_seconds y
 * admiten umbrales; los de lineas (carpeta, patron de log) mandan lineas en cada pasada del agente; connections muestrea
 * las conexiones de red cada interval_seconds y manda solo lo que abre o cierra (sin umbrales, solo Linux).
 */
export const MonitorTypeEnum = Object.freeze({
    CPU: "cpu",
    MEMORY: "memory",
    DISK: "disk",
    GPU: "gpu",
    NETWORK: "network",
    FOLDER: "folder",
    LOG_PATTERN: "log_pattern",
    CONNECTIONS: "connections",
});

export const MONITOR_TYPES = Object.freeze(Object.values(MonitorTypeEnum));

export const METRIC_MONITOR_TYPES = Object.freeze([
    MonitorTypeEnum.CPU, MonitorTypeEnum.MEMORY, MonitorTypeEnum.DISK, MonitorTypeEnum.GPU, MonitorTypeEnum.NETWORK,
]);

// los que usan interval_seconds (las metricas y el muestreo de conexiones); los de lineas van en cada pasada
export const INTERVAL_MONITOR_TYPES = Object.freeze([...METRIC_MONITOR_TYPES, MonitorTypeEnum.CONNECTIONS]);

// el valor de estos es un porcentaje (0-100); el de network, Mbit/s
export const PERCENT_MONITOR_TYPES = Object.freeze([
    MonitorTypeEnum.CPU, MonitorTypeEnum.MEMORY, MonitorTypeEnum.DISK, MonitorTypeEnum.GPU,
]);

export const MONITOR_TYPE_LABELS = Object.freeze({
    [MonitorTypeEnum.CPU]: "CPU (% de uso)",
    [MonitorTypeEnum.MEMORY]: "Memoria (% usada)",
    [MonitorTypeEnum.DISK]: "Disco (% usado de una ruta)",
    [MonitorTypeEnum.GPU]: "GPU (% de uso, nvidia-smi)",
    [MonitorTypeEnum.NETWORK]: "Red (entrada/salida, Mbit/s)",
    [MonitorTypeEnum.FOLDER]: "Carpeta (todas las líneas nuevas)",
    [MonitorTypeEnum.LOG_PATTERN]: "Patrón de log (solo las líneas que casan)",
    [MonitorTypeEnum.CONNECTIONS]: "Conexiones de red (abre/cierra por proceso, solo Linux)",
});

export const MONITOR_SEVERITIES = Object.freeze(["debug", "info", "warning", "error", "critical"]);

export const MonitorDefaultEnum = Object.freeze({
    INTERVAL_SECONDS: 60,
    // casi en tiempo real: a donde se envia informacion y que proceso lo hace
    CONNECTIONS_INTERVAL_SECONDS: 30,
    DISK_PATH: "/",
    FILE_PATTERN: "*",
});

// listas de connections (procesos y puertos a ignorar): las mismas reglas que MonitorFieldsValidator del backend
export const MonitorListLimitEnum = Object.freeze({
    SEPARATOR: ",",
    MAX_ITEMS: 50,
    PROCESS_NAME_MAX_LENGTH: 64,
    PORT_MIN: 1,
    PORT_MAX: 65535,
});
