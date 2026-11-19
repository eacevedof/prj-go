-- app_senders: perfiles de envio. provider_code (ProviderEnum) fija el canal y los datos por defecto; config guarda lo
-- no secreto en JSON (from_email, smtp_host...) y secret_encrypted el secreto con AES-256-GCM (SecretCipher).

CREATE TABLE app_senders
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES app_users (id) ON DELETE CASCADE,
    sender_uuid VARCHAR(50) NOT NULL UNIQUE,
    sender_name VARCHAR(100) NOT NULL,
    channel VARCHAR(10) NOT NULL CHECK (channel IN ('email', 'whatsapp', 'telegram')),
    provider_code VARCHAR(20) NOT NULL
        CHECK (provider_code IN ('gmail', 'yahoo', 'outlook', 'resend', 'smtp_custom', 'telegram_bot', 'whatsapp_cloud')),
    config TEXT NOT NULL DEFAULT '{}',
    secret_encrypted TEXT NOT NULL DEFAULT '',
    is_enabled INTEGER NOT NULL DEFAULT 1 CHECK (is_enabled IN (0, 1)),
    daily_limit INTEGER NULL,
    minute_limit INTEGER NULL,
    last_used_at TEXT NULL,
    last_error_at TEXT NULL,
    last_error_message VARCHAR(500) NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE UNIQUE INDEX uk_senders_user_id_sender_name ON app_senders (user_id, sender_name COLLATE NOCASE);
