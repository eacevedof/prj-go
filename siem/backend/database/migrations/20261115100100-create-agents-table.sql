-- agents: una instalacion del siem-client (un host) dentro de un proyecto. Se autentica con su api_token (aqui solo
-- el sha256). host_* y client_version los informa el propio cliente en cada envio; last_seen_at es su latido.

CREATE TABLE agents
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    agent_uuid VARCHAR(50) NOT NULL UNIQUE,
    agent_name VARCHAR(100) NOT NULL,
    api_token_hash CHAR(64) NOT NULL UNIQUE,
    heartbeat_interval_seconds INTEGER NOT NULL DEFAULT 60,
    is_enabled INTEGER NOT NULL DEFAULT 1 CHECK (is_enabled IN (0, 1)),
    host_name VARCHAR(150) NULL,
    host_os VARCHAR(50) NULL,
    host_ip VARCHAR(45) NULL,
    client_version VARCHAR(20) NULL,
    last_seen_at TEXT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uk_agents_project_id_agent_name ON agents (project_id, agent_name COLLATE NOCASE);
