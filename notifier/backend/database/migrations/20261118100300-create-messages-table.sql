-- app_messages: la cola y el historico a la vez (un mensaje = una unidad de trabajo). channel y sender_id van
-- desnormalizados (salen del token) para que el worker de cada canal filtre sin joins.

CREATE TABLE app_messages
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_uuid VARCHAR(50) NOT NULL UNIQUE,
    api_token_id INTEGER NOT NULL REFERENCES app_api_tokens (id) ON DELETE CASCADE,
    sender_id INTEGER NOT NULL REFERENCES app_senders (id) ON DELETE CASCADE,
    channel VARCHAR(10) NOT NULL CHECK (channel IN ('email', 'whatsapp', 'telegram')),
    recipients TEXT NOT NULL,
    subject VARCHAR(255) NULL,
    body TEXT NOT NULL DEFAULT '',
    body_format VARCHAR(5) NOT NULL DEFAULT 'text' CHECK (body_format IN ('html', 'text')),
    message_status VARCHAR(10) NOT NULL DEFAULT 'pending'
        CHECK (message_status IN ('pending', 'sending', 'sent', 'failed', 'cancelled')),
    attempts INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 5,
    next_attempt_at TEXT NOT NULL,
    error_message VARCHAR(500) NULL,
    provider_message_id VARCHAR(255) NULL,
    idempotency_key VARCHAR(100) NULL,
    request_ip VARCHAR(45) NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    sent_at TEXT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE (api_token_id, idempotency_key)
);

CREATE INDEX idx_messages_worker ON app_messages (channel, message_status, next_attempt_at, id);
CREATE INDEX idx_messages_sender_id_created_at ON app_messages (sender_id, created_at);
CREATE INDEX idx_messages_message_status_created_at ON app_messages (message_status, created_at);
