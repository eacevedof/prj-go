package enums

// EventsFilterFieldEnum los campos por los que se puede filtrar (incluir o, con el prefijo not_, excluir) en
// GET /api/v1/events, /events/histogram y /events/facets. Los valores son los parametros de la query. monitor_uuid es la
// columna events.monitor_uuid (attributes.monitor_uuid copiado al insertar: eventos host_metric y log_pattern).
var EventsFilterFieldEnum = struct {
    PROJECT_UUID string
    AGENT_UUID   string
    SEVERITY     string
    EVENT_TYPE   string
    SOURCE_TYPE  string
    SOURCE_NAME  string
    MONITOR_UUID string

    VALUES     []string
    NOT_PREFIX string
    // valores por campo en una peticion: de sobra para el panel y acota el IN (...)
    MAX_VALUES_PER_FIELD int
}{
    PROJECT_UUID: "project_uuid",
    AGENT_UUID:   "agent_uuid",
    SEVERITY:     "severity",
    EVENT_TYPE:   "event_type",
    SOURCE_TYPE:  "source_type",
    SOURCE_NAME:  "source_name",
    MONITOR_UUID: "monitor_uuid",

    VALUES:               []string{"project_uuid", "agent_uuid", "severity", "event_type", "source_type", "source_name", "monitor_uuid"},
    NOT_PREFIX:           "not_",
    MAX_VALUES_PER_FIELD: 20,
}

// EventsFacetFieldEnum los campos del panel lateral (top de valores con su recuento). agent_name y project_name se
// muestran por nombre; el front filtra por su uuid.
var EventsFacetFieldEnum = struct {
    VALUES []string
    TOP    int
}{
    VALUES: []string{"severity", "event_type", "source_type", "source_name", "agent_name", "project_name"},
    TOP:    10,
}
