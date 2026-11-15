-- projects: lo que se monitoriza (una app, un servidor, un cliente). Agrupa agentes.

CREATE TABLE projects
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_uuid VARCHAR(50) NOT NULL UNIQUE,
    project_name VARCHAR(100) NOT NULL,
    project_description VARCHAR(500) NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uk_projects_project_name ON projects (project_name COLLATE NOCASE);
