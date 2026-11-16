package enums

// DashboardKeyEnum claves JSON de las respuestas del dashboard.
var DashboardKeyEnum = struct {
    SERVER   string
    PROJECTS string
    PROJECT  string
    AGENTS   string
    SOURCES  string

    // server
    APP_NAME                  string
    APP_VERSION               string
    APP_ENV                   string
    NOW                       string
    STARTED_AT                string
    UPTIME_SECONDS            string
    DB_SIZE_BYTES             string
    EVENTS_TOTAL              string
    QUEUE_PENDING_BATCHES     string
    QUEUE_PENDING_EVENTS      string
    QUEUE_FAILED_BATCHES      string
    QUEUE_OLDEST_PENDING_AT   string
    QUEUE_WORKER_INTERVAL     string
    EVENTS_RETENTION_DAYS     string

    // project
    PROJECT_ID          string
    PROJECT_UUID        string
    PROJECT_NAME        string
    PROJECT_DESCRIPTION string
    PROJECT_HEALTH      string
    AGENTS_TOTAL        string
    AGENTS_ONLINE       string
    EVENTS_24H          string
    WARNINGS_24H        string
    ERRORS_24H          string
    CRITICALS_24H       string
    LAST_EVENT_AT       string

    // agent
    AGENT_ID     string
    AGENT_STATUS string
}{
    SERVER:   "server",
    PROJECTS: "projects",
    PROJECT:  "project",
    AGENTS:   "agents",
    SOURCES:  "sources",

    APP_NAME:                "app_name",
    APP_VERSION:             "app_version",
    APP_ENV:                 "app_env",
    NOW:                     "now",
    STARTED_AT:              "started_at",
    UPTIME_SECONDS:          "uptime_seconds",
    DB_SIZE_BYTES:           "db_size_bytes",
    EVENTS_TOTAL:            "events_total",
    QUEUE_PENDING_BATCHES:   "queue_pending_batches",
    QUEUE_PENDING_EVENTS:    "queue_pending_events",
    QUEUE_FAILED_BATCHES:    "queue_failed_batches",
    QUEUE_OLDEST_PENDING_AT: "queue_oldest_pending_at",
    QUEUE_WORKER_INTERVAL:   "queue_worker_interval_seconds",
    EVENTS_RETENTION_DAYS:   "events_retention_days",

    PROJECT_ID:          "project_id",
    PROJECT_UUID:        "project_uuid",
    PROJECT_NAME:        "project_name",
    PROJECT_DESCRIPTION: "project_description",
    PROJECT_HEALTH:      "project_health",
    AGENTS_TOTAL:        "agents_total",
    AGENTS_ONLINE:       "agents_online",
    EVENTS_24H:          "events_24h",
    WARNINGS_24H:        "warnings_24h",
    ERRORS_24H:          "errors_24h",
    CRITICALS_24H:       "criticals_24h",
    LAST_EVENT_AT:       "last_event_at",

    AGENT_ID:     "agent_id",
    AGENT_STATUS: "agent_status",
}
