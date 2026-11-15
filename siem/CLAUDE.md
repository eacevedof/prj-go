# siem — puntero

SIEM básico en **Go**: `client/` (agente que vigila carpetas y lanza comandos) → `backend/` (API de ingesta con cola,
worker, dashboard web). DDD con el patrón de `task-go`.

- Stack: Go 1.26 (compila en Docker: `golang:1.26.10-alpine`, no hay Go nativo en este PC) · SQLite en WAL
  (`modernc.org/sqlite`, sin cgo) · front HTML + JS sin build ni dependencias · el cliente solo usa la stdlib.
- Arranque/test: `README.md` de esta carpeta (`make start` en `backend/`, `make up` en `client/`).

**Documentación en Obsidian** (toda; aquí solo el puntero):
- Windows: `C:\projects\temper\ai\obsidian\projects\siem\siem-index.md`
- WSL: `/mnt/c/projects/temper/ai/obsidian/projects/siem/siem-index.md`

Al retomar: leer `siem-doing.md` (L1). Git es solo lectura: commit/push los hace Eduardo.
Config global: `claude-global\claude-global.md` del vault.
