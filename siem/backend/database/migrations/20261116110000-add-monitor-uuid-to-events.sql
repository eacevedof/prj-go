-- events.monitor_uuid: el monitor que genero el evento (host_metric, log_pattern), copiado de attributes.monitor_uuid al
-- insertar. Columna propia con indice para filtrar el visor por monitor: json_extract sobre attributes recorria la
-- tabla entera (facetas de 5 s a 9 s). Sin FK: un monitor borrado deja sus eventos con su uuid. Se rellena la de los
-- eventos que ya habia.

ALTER TABLE events ADD COLUMN monitor_uuid VARCHAR(50) NULL;

UPDATE events
SET monitor_uuid = json_extract(attributes, '$.monitor_uuid')
WHERE attributes IS NOT NULL
AND json_extract(attributes, '$.monitor_uuid') IS NOT NULL;

CREATE INDEX idx_events_monitor_uuid_occurred_at ON events (monitor_uuid, occurred_at);
