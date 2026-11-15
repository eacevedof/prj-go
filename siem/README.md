# siem

SIEM básico: un **agente** (`client/`) por host vigila carpetas de logs y lanza comandos cada X tiempo, y envía lo
recogido a un **backend** (`backend/`) que lo encola, lo guarda en SQLite y lo enseña en un dashboard web.

```
client (por host)                          backend (uno)
  watchers ─┐                                POST /api/v1/ingest-batches ──► ingest_queue (202 enseguida)
  commands ─┴─► spool/ ──► events-send ──────►                                   │ worker cada 2 s
              (aguanta si el backend cae)                                        ▼
                                             projects · agents · sources · events ◄── dashboard (/)
```

## 1 - backend

### 1.1 Arranque en local (Docker)

1. `cd backend && make start` — crea `.env` y `docker/.env`, compila, pasa las migraciones y levanta
   `http://localhost:9400`.
2. En `backend/.env` poner un `APP_DASHBOARD_TOKEN` (p. ej. `openssl rand -hex 32`) y `make restart`. Es el token que
   pide el dashboard.
3. Abrir `http://localhost:9400` y entrar con el token del dashboard. Desde el panel se crea todo:
   1. **Proyecto** (portada → "nuevo proyecto").
   2. **Agente** (detalle del proyecto → "nuevo agente"): enseña su `api_token` **una sola vez**, con el `.env` listo
      para copiar.
   3. **Monitores** (detalle del proyecto → "monitores"): lo que vigilan los agentes de ese proyecto (cpu, memoria,
      disco, gpu, red, carpetas, patrones de log), con sus umbrales. Los agentes los descargan solos en < 1 min.
   4. **Descargas**: los binarios del agente con las instrucciones de instalación (antes, `cd client && make dist`).
4. Lo mismo por consola: `make console c='project-create --project_name="Mi app"'` y
   `make console c='agent-create --project_uuid=<uuid> --agent_name=web-01'`.

### 1.2 Comandos

| # | Comando | Qué hace |
|---|---|---|
| 1 | `make console c="help"` | lista los comandos de consola |
| 2 | `make console c="ingest-queue-process"` | vacía la cola a mano (lo hace solo el worker de la API) |
| 3 | `make console c="events-purge --retention_days=30"` | retención: borra eventos viejos y lotes ya procesados (cron diario) |
| 4 | `make test` · `make vet` · `make logs` · `make migrate` | lo habitual |

### 1.3 API

| # | Ruta | Auth | Para qué |
|---|---|---|---|
| 1 | `POST /api/v1/ingest-batches` | agente | el cliente envía un lote (202); `events: []` = latido |
| 2 | `GET /api/v1/agent-config` | agente | los monitores que tiene que vigilar ese agente: `{config_hash, monitors}` |
| 3 | `GET /api/v1/dashboard` | dashboard | servidor + semáforo por proyecto |
| 4 | `POST /api/v1/projects` · `DELETE /api/v1/projects/:project_uuid` | dashboard | alta (`{project_name, project_description}`) y baja (en cascada) |
| 5 | `GET /api/v1/projects/:project_uuid` | dashboard | agentes y fuentes de un proyecto |
| 6 | `POST /api/v1/projects/:project_uuid/agents` · `DELETE /api/v1/agents/:agent_uuid` | dashboard | alta (`{agent_name, heartbeat_interval_seconds}` → `api_token`, **una vez**) y baja |
| 7 | `GET` · `POST /api/v1/projects/:project_uuid/monitors` | dashboard | monitores del proyecto (con su último valor por agente) y alta |
| 8 | `PUT` · `DELETE /api/v1/monitors/:monitor_uuid` | dashboard | edición (todo menos el tipo) y baja |
| 9 | `GET /api/v1/events` | dashboard | visor: `{events, page, page_size, has_more}` |
| 10 | `GET /api/v1/events/histogram` · `GET /api/v1/events/facets` | dashboard | barras por tramo y gravedad · top de valores por campo |
| 11 | `GET /api/v1/downloads` | dashboard | binarios del agente (`public/downloads/`, tamaño y sha256) y la url de ingesta |
| 12 | `GET /downloads/<fichero>` · `GET /health` | — | descarga de un binario · vivo |

Auth = `Authorization: Bearer <token>`: "agente" = su `api_token`; "dashboard" = `APP_DASHBOARD_TOKEN` (lee y escribe).

**Filtros** de `events`, `histogram` y `facets` (los mismos en las tres): `from`, `to` (UTC `YYYY-MM-DD HH:MM:SS` sobre
`occurred_at`; sin rango, `histogram` y `facets` miran las últimas 24 h), `search` (texto en el mensaje) y, por campo
(`project_uuid`, `agent_uuid`, `severity`, `event_type`, `source_type`, `source_name`), `campo=a,b` para incluir y
`not_campo=a,b` para excluir. `events` añade `page` y `page_size` (máx. 200).

**Monitor** (`monitor_type` → `config` · valor y unidad de sus umbrales `warning_threshold` / `critical_threshold`):
`cpu`, `memory`, `gpu` → `{}` · % · `disk` → `{path}` · % usado · `network` → `{interface}` ("" = todas) · Mbit/s ·
`folder` → `{path, pattern, is_recursive}` · `log_pattern` → lo mismo + `{match (regex), severity}`. Más
`monitor_name` (único en el proyecto: es el `source_name` de sus eventos), `agent_uuid` ("" = todos los agentes),
`is_enabled` e `interval_seconds` (10-86400).

Cuerpo de un lote:

```json
{
    "batch_uuid": "b-6af9d1f5bd83a",
    "agent": {"host_name": "web-01", "host_os": "linux", "client_version": "0.1.0"},
    "events": [
        {"event_type": "log_line", "source_name": "nginx-logs", "file_path": "/var/log/nginx/error.log",
         "message": "... connect() failed ...", "occurred_at": "2026-11-15 10:00:00"},
        {"event_type": "command_result", "source_name": "disk-usage", "message": "/dev/sda1 91%",
         "exit_code": 0, "duration_ms": 12, "occurred_at": "2026-11-15 10:00:00"}
    ]
}
```

La gravedad (`debug`…`critical`) la decide el backend: palabras clave en una línea, código de salida en un comando,
los `attributes` en los eventos de Docker y los umbrales del monitor en un `host_metric` (un `severity` explícito en el
evento manda). Tipos de evento: `log_line`, `command_result`, `docker_state`, `docker_stats`, `docker_event`,
`docker_log` y `host_metric`.

## 2 - client (el agente)

### 2.1 Probarlo en local contra el backend (Docker)

1. Con el backend arrancado: `cd client && make env`, poner en `docker/.env` el `APP_SIEM_API_TOKEN` del agente.
2. `make up` — el agente vigila `client/storage/demo-logs/` y lanza los comandos de `config/siem-client.docker.json`.
3. `make demo-log m="ERROR database down"` escribe una línea; a los pocos segundos sale en el dashboard.
4. `make logs` muestra lo que hace el agente.
5. Carpeta real de Windows: `SIEM_WATCH_HOST_PATH=C:/ruta/carpeta` en `docker/.env` y `make up`. Se monta en solo
   lectura en `/watch/host` (watcher `host-folder`, todos los ficheros); lo que se añada a un fichero sale en el
   dashboard en unos segundos. Los comandos se ejecutan dentro del contenedor (Linux), no en Windows.

### 2.1.1 Editar con VS Code sin instalar Go (Dev Containers)

En `prj-go/.devcontainer/` hay `siem-backend` y `siem-client`: abrir la carpeta `prj-go` en VS Code →
*Dev Containers: Reopen in Container* → elegir uno. Se engancha al contenedor de `make up` (no crea otro) con la
extensión de Go; al cerrar VS Code el contenedor sigue. Sin `formatOnSave`: el repo va a 4 espacios, no gofmt.

### 2.2 Desplegarlo en un host

Las instrucciones paso a paso, con los comandos para copiar por sistema operativo, están en el panel: **Descargas**.

1. `make dist` (con el contenedor del cliente levantado) → `build/dist/` y una copia en `backend/public/downloads/`
   (`siem-client-{linux,windows,macos}-<versión>`), que es de donde los sirve el panel. En el host se copia como
   `siem-client` (sin versión: la unidad systemd no cambia al actualizar; la versión la dice `version`).
2. En el host, en una carpeta (p. ej. `/opt/siem-client/`): el binario y un `.env` con `APP_SIEM_ENDPOINT_URL` y
   `APP_SIEM_API_TOKEN` (salen al crear el agente en el panel). **Ya no hace falta `siem-client.json`**: lo que se
   vigila son los monitores del proyecto, que el agente descarga del backend al arrancar y cada 60 s
   (`state/remote-config.json`; si el backend no responde, sigue con la última). El fichero sigue valiendo para lo que
   no está en el panel (comandos programados, Docker) y se suma a lo remoto. `APP_SIEM_REMOTE_CONFIG=0` lo apaga.
3. `./siem-client config-check` — valida todo sin enviar nada.
4. `./siem-client agent-run` — el agente. Para dejarlo como servicio, una unidad systemd:

```ini
[Unit]
Description=siem-client
After=network-online.target

[Service]
WorkingDirectory=/opt/siem-client
ExecStart=/opt/siem-client/siem-client agent-run
Restart=always

[Install]
WantedBy=multi-user.target
```

### 2.3 Configuración (`siem-client.json`)

| # | Clave | Defecto | Qué es |
|---|---|---|---|
| 1 | `collect_interval_seconds` | 5 | cada cuánto se leen las carpetas (y se mira si toca algún comando) |
| 2 | `send_interval_seconds` | 10 | cada cuánto se envía el spool al backend |
| 3 | `heartbeat_interval_seconds` | 30 | sin nada que enviar, cada cuánto se manda el latido |
| 4 | `max_events_per_batch` | 500 | eventos por lote (máx. 2000) |
| 5 | `max_spool_batches` | 5000 | con el backend caído, lotes que se guardan antes de tirar los más viejos |
| 6 | `watchers[]` | — | `name`, `path` (absoluta), `pattern` (glob del nombre, `*`), `is_recursive`, `start_at` (`end`: solo lo nuevo · `beginning`: todo) |
| 7 | `commands[]` | — | `name`, `command` (`sh -c` / `cmd /C`), `interval_seconds` (60), `timeout_seconds` (30) |
| 8 | `docker` | apagado | `is_enabled`, `socket_path` (`/var/run/docker.sock`), `is_stats_enabled` + `stats_interval_seconds` (60), `is_events_enabled`, `is_logs_enabled` + `logs_max_lines_per_container` (200), `logs_include` / `logs_exclude` (globs del nombre del contenedor) |

Las carpetas y comandos se releen en cada pasada (no hace falta reiniciar); los intervalos, sí.

**Docker** (`docker.is_enabled: true`, el socket montado en el agente): por contenedor manda su estado de partida y cada
cambio (para, arranca, se reinicia, OOM, health), una muestra de cpu/memoria/red cada `stats_interval_seconds`, las líneas
nuevas de su log y los eventos del daemon (`create`, `die`, `oom`, `kill`, `pull`…; sin los `exec_*` de los healthchecks).
La gravedad la decide el backend (OOM o `dead` critical; parado con exit ≠ 0 o `unhealthy`, error; memoria ≥ 85 %, warning).
El agente nunca lee su propio log (bucle) y solo hace `GET`; aun así el socket da control total del daemon.

**Monitores del panel**: las métricas del host (`cpu`, `memory`, `disk`, `network`) se leen de `/proc` y `statfs`, así que
hoy solo funcionan en **Linux** (en Windows y macOS salen una vez como "no disponible"); `gpu` usa `nvidia-smi`. Dentro de
un contenedor miden el contenedor / la VM de Docker, no el host: para un servidor real, el binario suelto.

### 2.4 Comandos

| # | Comando | Qué hace |
|---|---|---|
| 1 | `config-check` | valida configuración (fichero local + monitores ya descargados) + endpoint + token |
| 2 | `config-sync` | descarga del backend los monitores de este agente |
| 3 | `agent-run` | el bucle del agente |
| 4 | `files-collect` · `commands-run --force=1` · `metrics-collect --force=1` · `docker-collect --force_stats=1` · `events-send` | los pasos del bucle, sueltos |
| 5 | `make test` · `make vet` · `make dist` · `make version` | lo habitual |

## 3 - versiones

Cada parte tiene su `APP_VERSION` (en `src/modules/shared/infrastructure/repositories/configuration/environment_reader_raw_repository.go`)
y su `changelog.md`. Subirla en cada cambio, también en los fixes.
