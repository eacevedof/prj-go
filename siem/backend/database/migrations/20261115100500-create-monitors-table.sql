-- monitors: los recursos a vigilar de un proyecto, creados desde el panel. El agente los descarga
-- (GET /api/v1/agent-config) y manda sus eventos con source_name = monitor_name y attributes.monitor_uuid.
-- agent_id NULL = lo aplican todos los agentes del proyecto. config = JSON propio del tipo (path, interface, pattern,
-- match...). Los umbrales (NULL = sin umbral) solo los usa el backend para decidir la gravedad de host_metric.

CREATE TABLE monitors
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    monitor_uuid VARCHAR(50) NOT NULL UNIQUE,
    project_id INTEGER NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    agent_id INTEGER NULL REFERENCES agents (id) ON DELETE CASCADE,
    monitor_type VARCHAR(20) NOT NULL
        CHECK (monitor_type IN ('cpu', 'memory', 'disk', 'gpu', 'network', 'folder', 'log_pattern', 'connections')),
    monitor_name VARCHAR(100) NOT NULL,
    is_enabled INTEGER NOT NULL DEFAULT 1 CHECK (is_enabled IN (0, 1)),
    interval_seconds INTEGER NOT NULL DEFAULT 60,
    config TEXT NOT NULL DEFAULT '{}',
    warning_threshold REAL NULL,
    critical_threshold REAL NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uk_monitors_project_id_monitor_name ON monitors (project_id, monitor_name COLLATE NOCASE);
CREATE INDEX idx_monitors_agent_id ON monitors (agent_id);
