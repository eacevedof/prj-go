-- app_api_tokens: tokens de envio. Uno por app cliente: revocar uno no rompe las demas que usan el mismo perfil.
-- Solo se guarda el sha256; token_prefix es lo que se ve en el panel para reconocerlo.

CREATE TABLE app_api_tokens
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL REFERENCES app_senders (id) ON DELETE CASCADE,
    token_uuid VARCHAR(50) NOT NULL UNIQUE,
    token_name VARCHAR(100) NOT NULL,
    api_token_hash CHAR(64) NOT NULL UNIQUE,
    token_prefix VARCHAR(12) NOT NULL,
    default_to TEXT NULL,
    is_enabled INTEGER NOT NULL DEFAULT 1 CHECK (is_enabled IN (0, 1)),
    last_used_at TEXT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE UNIQUE INDEX uk_api_tokens_sender_id_token_name ON app_api_tokens (sender_id, token_name COLLATE NOCASE);
