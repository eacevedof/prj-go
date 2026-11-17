# changelog — siem-client

La versión es `APP_VERSION` en `src/modules/shared/infrastructure/repositories/configuration/environment_reader_raw_repository.go`
(la muestran `version` y `help`, viaja al backend en cada envío y da nombre a los binarios de `make dist`). Lo más nuevo
arriba.

## 0.4.2 — 2026-11-17

1. **Carpeta de instalación `gosiem-client`** (antes `siem-client`): `install.sh` instala en `/opt/gosiem-client` e
   `install.ps1` en `C:\programs\gosiem-client`. El binario sigue siendo `siem-client` y el servicio / la tarea
   programada no cambian de nombre. Un agente ya instalado no se mueve solo: si se relanza el instalador, queda una
   instalación nueva en la carpeta nueva y la vieja hay que desinstalarla a mano.

## 0.4.1 — 2026-11-16

1. **Seguridad — `install.ps1`**: la carpeta de instalación deja de heredar permisos de `C:\` (ahí "Usuarios
   autentificados" puede modificar): solo SYSTEM y Administradores (sin admin, el usuario actual) y dueño
   Administradores. Antes cualquier usuario del equipo podía cambiar el `.exe` o el `siem-client.json` que la tarea
   ejecuta como SYSTEM, o leer el token del `.env`. Relanzar el instalador corrige las instalaciones existentes.

## 0.4.0 — 2026-11-16

1. Cada envío (también el latido) lleva `agent.client_path`: la ruta absoluta del ejecutable (`os.Executable`, con los
   symlinks resueltos). El panel la usa para decir dónde desinstalar el agente cuando se borra su proyecto. Un backend
   anterior a 0.8.0 la ignora.

## 0.3.0 — 2026-11-15

1. **Config remota**: el agente baja del backend lo que tiene que vigilar (`GET .../agent-config`, los monitores que se
   crean en el panel) al arrancar y cada 60 s, y lo guarda en `state/remote-config.json` (si el backend no responde,
   sigue con la última). `APP_SIEM_REMOTE_CONFIG` (defecto 1). Comando `config-sync`.
2. `siem-client.json` pasa a ser **opcional**: con `.env` (endpoint + token) basta; si existe, se suma a lo remoto. Un
   agente sin nada que vigilar todavía arranca igual (latido) y espera a que lleguen monitores.
3. Monitores `folder` y `log_pattern` del panel como watchers (`start_at` end): `log_pattern` solo manda las líneas que
   casan con su regex `match` y puede forzar la `severity`; sus líneas llevan `attributes {monitor_uuid, monitor_type}`.
   Los watchers locales también admiten `match` y `severity`. Un monitor remoto que no vale se ignora con aviso.
4. Módulo `hostmetrics` (Linux): monitores `cpu` (`/proc/stat`), `memory` (`/proc/meminfo`), `disk` (statfs), `network`
   (`/proc/net/dev`, Mbit/s) y `gpu` (`nvidia-smi`) como eventos `host_metric` con `attributes.value`; una métrica no
   disponible se avisa una sola vez (`is_unavailable`). Comando `metrics-collect [--force=1]`.
5. `make dist` copia además los binarios a `backend/public/downloads/` (`siem-client-{linux,windows,macos}-<versión>`).

## 0.2.0 — 2026-11-15

1. Módulo `docker` (colector nativo, solo stdlib, por el socket `/var/run/docker.sock`): estados de los contenedores
   (estado de partida + cambios: para, arranca, reinicio, OOM, health, borrado), recursos (cpu/memoria/red cada
   `stats_interval_seconds`), logs de los contenedores (cursor por marca de tiempo, `logs_include`/`logs_exclude`, nunca
   el propio agente) y eventos del daemon (sin `exec_*`). Bloque `docker` en `siem-client.json`; comando `docker-collect`.
2. Eventos del daemon en ventanas contiguas con nanosegundos (con segundos enteros se perdía un `start`).
3. Compose local: socket de Docker montado en solo lectura.

## 0.1.1 — 2026-11-15

1. Prueba en local con una carpeta real de Windows: el compose monta `SIEM_WATCH_HOST_PATH` (docker/.env) en solo
   lectura en `/watch/host` y `siem-client.docker.json` la vigila (watcher `host-folder`, patrón `*`).
2. Dev Container `prj-go/.devcontainer/siem-client`: VS Code dentro del contenedor del agente (sin Go en Windows).

## 0.1.0 — 2026-11-15

1. Primera versión, esqueleto de consola de `yt-downloader`: módulos `agentconfig`, `collectors`, `shipper`, `agent`.
2. `siem-client.json`: carpetas vigiladas (glob, recursivo, `start_at`) y comandos programados (intervalo, timeout);
   endpoint y token en `.env`. `config-check` lo valida todo.
3. `agent-run`: lee las líneas nuevas (offset por fichero, detecta truncado), lanza los comandos a los que les toca y
   lo deja en un spool en disco; envía el spool al backend con reintento (backend caído = el spool espera) y latido.
4. Sin dependencias de terceros (solo la stdlib de Go): binarios para Linux, Windows y macOS.
