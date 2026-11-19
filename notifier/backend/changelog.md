# notifier-backend — changelog

## 0.1.0 — 2026-11-18

Fase 1 del roadmap (backend Gmail), sin panel:

1. Esqueleto de `siem/backend` (patrón `task-go`): API REST, consola, migraciones, SQLite en WAL, tablas `app_*`.
2. `users` con roles `root` (uno, por consola) · `admin` · `owner` y políticas en `UserPolicyEnum`; `POST /api/v1/log-in`.
3. `senders`: perfiles de envío `gmail` y `smtp_custom`, secreto cifrado con AES-256-GCM (`APP_SECRETS_KEY`).
4. `apitokens`: tokens de envío `ntf-…` (solo el sha256 en BD), varios por perfil, revocables.
5. `messages`: `POST /api/v1/messages` (202, idempotencia, `to` string o lista, adjuntos de texto o por url) y
   `GET /api/v1/messages/:message_uuid`.
6. `delivery`: worker por canal (email) con `SmtpMessageSender` (starttls, ssl, none), reintentos 1 min · 5 min ·
   30 min · 2 h (5 intentos), 5xx permanente, límites por minuto y por día del perfil, huérfanos `sending` recuperados,
   descarga de adjuntos con lista blanca de hosts y sin IPs privadas.
