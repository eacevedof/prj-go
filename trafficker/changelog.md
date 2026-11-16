# changelog — trafficker

La versión es `APP_VERSION` en `src/modules/shared/infrastructure/repositories/configuration/environment_reader_raw_repository.go`
(la muestran `version` y `help`, viajará al backend del siem en cada envío y da nombre a los binarios de `make dist`).
Lo más nuevo arriba.

## 0.1.0 — 2026-11-16

1. Esqueleto copiado del `siem-client` 0.3.0 (consola, `.env`, logs, spool y estado), sin sus módulos de carpetas,
   comandos, Docker ni métricas. Solo stdlib.
2. Módulo `connections` (solo Linux): lee `/proc/net/{tcp,tcp6,udp,udp6}` y une cada socket con su proceso por el
   inodo (`/proc/<pid>/fd`); dirección `listen` / `inbound` / `outbound`; usuario de `/etc/passwd`.
3. Comando `connections-list [--filter="established,remote"] [--format=json]`: tabla en color (o JSON). Filtros del
   `.env`: loopback fuera por defecto, procesos y puertos ignorados.
4. Docker local con `network_mode: host` y `pid: host` para ver los sockets y procesos del host, más `SYS_PTRACE` y
   `apparmor:unconfined` (sin ellos todos los procesos salen como `?`).
5. **Modo local** (sin `.env` y sin token en el entorno): `./trafficker` sin argumentos = `connections-list`, como el
   `localhost-traffic` de Rust. Con `.env`, sin argumentos sale la ayuda.
6. Módulo `ipowners`: país y organización de la IP remota por **RDAP** (`rdap.org` → ARIN/RIPE/APNIC/LACNIC/AFRINIC),
   caché por rango de red en `state/ip-owners.json` (15 días; un fallo, 1 hora), máx. 40 consultas por pasada y 4 a la
   vez. IPs privadas, loopback, link-local y CGNAT se nombran sin consultar. Por defecto solo en modo local;
   `--whois=0|1` lo fuerza. Columnas `COUNTRY` y `ORGANIZATION`; en JSON, `remote_scope`, `remote_network`,
   `remote_org`, `remote_country`.
