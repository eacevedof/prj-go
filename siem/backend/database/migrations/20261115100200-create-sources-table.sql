-- sources: cada cosa que vigila un agente (un watcher de carpeta, un comando programado o, por contenedor Docker, su
-- estado (docker), sus recursos (docker_stats) y sus logs (docker_log)). Nace sola con su primer
-- evento. last_* es el ultimo estado conocido: el front lo pinta sin recorrer events.

CREATE TABLE sources
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    agent_id INTEGER NOT NULL REFERENCES agents (id) ON DELETE CASCADE,
    source_type VARCHAR(20) NOT NULL CHECK (source_type IN ('file', 'command', 'docker', 'docker_stats', 'docker_log', 'metric', 'network')),
    source_name VARCHAR(255) NOT NULL,
    events_count INTEGER NOT NULL DEFAULT 0,
    last_event_at TEXT NULL,
    last_severity VARCHAR(10) NULL,
    last_exit_code INTEGER NULL,
    last_message VARCHAR(500) NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (agent_id, source_type, source_name)
);
