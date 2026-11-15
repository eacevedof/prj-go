/**
 * Tipos de monitor (recurso a vigilar) del contrato de la API. Los de metricas dan un valor cada interval_seconds y
 * admiten umbrales; los de lineas (carpeta, patron de log) mandan lineas en cada pasada del agente.
 */
export const MonitorTypeEnum = Object.freeze({
    CPU: "cpu",
    MEMORY: "memory",
    DISK: "disk",
    GPU: "gpu",
    NETWORK: "network",
    FOLDER: "folder",
    LOG_PATTERN: "log_pattern",
});

export const MONITOR_TYPES = Object.freeze(Object.values(MonitorTypeEnum));

export const METRIC_MONITOR_TYPES = Object.freeze([
    MonitorTypeEnum.CPU, MonitorTypeEnum.MEMORY, MonitorTypeEnum.DISK, MonitorTypeEnum.GPU, MonitorTypeEnum.NETWORK,
]);

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
});

export const MONITOR_SEVERITIES = Object.freeze(["debug", "info", "warning", "error", "critical"]);

export const MonitorDefaultEnum = Object.freeze({
    INTERVAL_SECONDS: 60,
    DISK_PATH: "/",
    FILE_PATTERN: "*",
});
