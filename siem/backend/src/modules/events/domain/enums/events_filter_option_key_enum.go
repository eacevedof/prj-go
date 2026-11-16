package enums

// EventsFilterOptionKeyEnum columnas de las filas que lee EventsFilterOptionsReaderSqliteRepository y claves de
// GET /api/v1/events/filter-options (las mismas: uuid y nombre de proyectos, agentes y monitores).
var EventsFilterOptionKeyEnum = struct {
    PROJECT_UUID string
    PROJECT_NAME string
    AGENT_UUID   string
    AGENT_NAME   string
    MONITOR_UUID string
    MONITOR_NAME string
    MONITOR_TYPE string
}{
    PROJECT_UUID: "project_uuid",
    PROJECT_NAME: "project_name",
    AGENT_UUID:   "agent_uuid",
    AGENT_NAME:   "agent_name",
    MONITOR_UUID: "monitor_uuid",
    MONITOR_NAME: "monitor_name",
    MONITOR_TYPE: "monitor_type",
}
