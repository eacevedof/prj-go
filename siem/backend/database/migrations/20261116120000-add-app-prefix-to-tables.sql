-- prefijo app_ en todas las tablas de la aplicacion (convencion de Eduardo). RENAME conserva datos e indices y SQLite
-- reescribe las FK de las otras tablas que apuntan a la renombrada. La tabla migrations NO se toca: es del runner
-- (la crea si no existe y anota aqui cada fichero; renombrarla desde una migracion le haria repetirlas todas).
-- Los indices conservan su nombre (idx_events_..., uk_agents_...): solo son nombres.

ALTER TABLE projects RENAME TO app_projects;
ALTER TABLE agents RENAME TO app_agents;
ALTER TABLE sources RENAME TO app_sources;
ALTER TABLE ingest_queue RENAME TO app_ingest_queue;
ALTER TABLE events RENAME TO app_events;
ALTER TABLE monitors RENAME TO app_monitors;
