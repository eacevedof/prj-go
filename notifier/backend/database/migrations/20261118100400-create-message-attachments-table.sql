-- app_message_attachments: adjuntos de un mensaje. text = el contenido va aqui; url = lo descarga el worker al
-- enviar (solo https y de APP_ATTACHMENT_ALLOWED_HOSTS).

CREATE TABLE app_message_attachments
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id INTEGER NOT NULL REFERENCES app_messages (id) ON DELETE CASCADE,
    attachment_type VARCHAR(5) NOT NULL CHECK (attachment_type IN ('text', 'url')),
    file_name VARCHAR(255) NOT NULL,
    content TEXT NULL,
    source_url VARCHAR(1000) NULL,
    mime_type VARCHAR(100) NULL,
    size_bytes INTEGER NULL
);

CREATE INDEX idx_message_attachments_message_id ON app_message_attachments (message_id);
