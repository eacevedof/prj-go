-- agents.client_path: ruta absoluta del ejecutable del siem-client en su host (os.Executable), la informa el propio
-- cliente en cada envio, como host_*. Sirve para decir donde desinstalarlo al borrar su proyecto o el agente.

ALTER TABLE agents ADD COLUMN client_path VARCHAR(500) NULL;
