-- app_message_attempts: una fila por intento de entrega, con lo que respondio el proveedor (recortado, sin secretos).

CREATE TABLE app_message_attempts
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id INTEGER NOT NULL REFERENCES app_messages (id) ON DELETE CASCADE,
    attempt_number INTEGER NOT NULL,
    attempt_status VARCHAR(20) NOT NULL CHECK (attempt_status IN ('sent', 'temporary_error', 'permanent_error')),
    provider_response VARCHAR(1000) NULL,
    duration_ms INTEGER NOT NULL DEFAULT 0,
    attempted_at TEXT NOT NULL
);

CREATE INDEX idx_message_attempts_message_id ON app_message_attempts (message_id);
