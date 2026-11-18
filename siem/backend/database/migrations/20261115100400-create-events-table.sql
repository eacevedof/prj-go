-- events: una linea de log o el resultado de un comando. project_id va desnormalizado (se deduce de agent_id) para
-- que el dashboard filtre por proyecto y fecha con un solo indice. occurred_at = hora del cliente (UTC);
-- received_at = hora en que el backend recibio el lote.

CREATE TABLE events
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    agent_id INTEGER NOT NULL REFERENCES agents (id) ON DELETE CASCADE,
    source_id INTEGER NOT NULL REFERENCES sources (id) ON DELETE CASCADE,
    event_type VARCHAR(20) NOT NULL CHECK (event_type IN ('log_line', 'command_result', 'docker_state', 'docker_stats',
        'docker_event', 'docker_log', 'host_metric', 'net_connection')),
    severity VARCHAR(10) NOT NULL DEFAULT 'info'
        CHECK (severity IN ('debug', 'info', 'warning', 'error', 'critical')),
    message TEXT NOT NULL,
    file_path VARCHAR(500) NULL,
    exit_code INTEGER NULL,
    duration_ms INTEGER NULL,
    -- datos estructurados del evento en JSON (Docker: estado, exit_code, oom, % de cpu y memoria...); NULL si no tiene
    attributes TEXT NULL,
    occurred_at TEXT NOT NULL,
    received_at TEXT NOT NULL
);

CREATE INDEX idx_events_project_id_occurred_at ON events (project_id, occurred_at);
CREATE INDEX idx_events_agent_id_occurred_at ON events (agent_id, occurred_at);
CREATE INDEX idx_events_source_id_occurred_at ON events (source_id, occurred_at);
CREATE INDEX idx_events_severity_occurred_at ON events (severity, occurred_at);
