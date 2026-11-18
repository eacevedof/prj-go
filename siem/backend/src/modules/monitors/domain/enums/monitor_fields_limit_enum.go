package enums

// MonitorFieldsLimitEnum limites y valores por defecto de los campos de un monitor (MonitorFieldsValidator). Si
// cambian, cambiar tambien los mensajes de MonitorsMessageEnum y los limites del front.
var MonitorFieldsLimitEnum = struct {
    INTERVAL_SECONDS_DEFAULT int
    INTERVAL_SECONDS_MIN     int
    INTERVAL_SECONDS_MAX     int
    // umbrales de los tipos que miden un %
    PERCENT_THRESHOLD_MIN float64
    PERCENT_THRESHOLD_MAX float64

    DEFAULT_DISK_PATH     string
    DEFAULT_FILES_PATTERN string
    // nombre de fichero cualquiera con el que se prueba que config.pattern es un glob valido
    GLOB_PROBE_FILE_NAME string

    // connections: muestreo por defecto (casi en tiempo real: a donde se envia informacion y que proceso lo hace) y
    // listas de procesos y puertos a ignorar
    CONNECTIONS_INTERVAL_SECONDS_DEFAULT int
    IGNORE_LIST_MAX_ITEMS                int
    PROCESS_NAME_MAX_LENGTH              int
    PORT_MIN                             int
    PORT_MAX                             int
    // separador de esas listas cuando llegan como texto ("sshd, nginx")
    LIST_SEPARATOR string
}{
    INTERVAL_SECONDS_DEFAULT: 60,
    INTERVAL_SECONDS_MIN:     10,
    INTERVAL_SECONDS_MAX:     86400,
    PERCENT_THRESHOLD_MIN:    0,
    PERCENT_THRESHOLD_MAX:    100,

    DEFAULT_DISK_PATH:     "/",
    DEFAULT_FILES_PATTERN: "*",
    GLOB_PROBE_FILE_NAME:  "probe",

    CONNECTIONS_INTERVAL_SECONDS_DEFAULT: 30,
    IGNORE_LIST_MAX_ITEMS:                50,
    PROCESS_NAME_MAX_LENGTH:              64,
    PORT_MIN:                             1,
    PORT_MAX:                             65535,
    LIST_SEPARATOR:                       ",",
}
