-- indice por fecha a secas: el visor sin filtros (la pantalla de entrada: ultimas 24 h) solo acota por occurred_at y,
-- sin el, la lista, el histograma y cada faceta recorrian la tabla entera (SCAN events). Los otros indices llevan
-- delante project_id / agent_id / source_id / severity / monitor_uuid y no sirven para un rango solo.

CREATE INDEX idx_events_occurred_at ON events (occurred_at);
