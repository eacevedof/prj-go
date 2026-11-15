-- ingest_queue: la cola. POST /api/v1/ingest-batches solo inserta aqui el lote crudo (payload = JSON de events) y
-- responde 202; el worker lo convierte en filas de events. (agent_id, batch_uuid) unico: un reenvio del cliente no
-- duplica. Un lote procesado queda como done con el payload vacio (lo borra events-purge pasadas 24 h).

CREATE TABLE ingest_queue
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    agent_id INTEGER NOT NULL REFERENCES agents (id) ON DELETE CASCADE,
    batch_uuid VARCHAR(60) NOT NULL,
    payload TEXT NOT NULL,
    events_count INTEGER NOT NULL DEFAULT 0,
    queue_status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (queue_status IN ('pending', 'done', 'failed')),
    attempts INTEGER NOT NULL DEFAULT 0,
    error_message VARCHAR(500) NULL,
    received_at TEXT NOT NULL,
    processed_at TEXT NULL,
    UNIQUE (agent_id, batch_uuid)
);

CREATE INDEX idx_ingest_queue_queue_status_id ON ingest_queue (queue_status, id);
