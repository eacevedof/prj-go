# notifier

Hub de notificaciones en Go: una API recibe un mensaje con un **token de envío**, lo encola y un worker por canal lo
entrega con la cuenta del perfil al que pertenece el token (fase 1: email por Gmail o cualquier SMTP).

- `backend/` — API + consola + workers (Go 1.26, SQLite). Versión y cambios: [`backend/changelog.md`](../backend/changelog.md).
- La documentación de diseño vive en el vault (`projects/notifier/`), no aquí.

## 1 - arrancar en local

Go se compila dentro de Docker (no hace falta Go en la máquina). Desde `backend/`:

1. `make env` y poner en `.env` una `APP_SECRETS_KEY` (`openssl rand -hex 32`).
2. `make start` → API en `http://localhost:9500` y Mailpit (los correos de prueba) en `http://localhost:9525`.
3. `make console c="help"` lista los comandos.

## 2 - primer envío

1. Usuario: `make console c='user-create --user_name=Ana --user_email=ana@example.com --password="una-clave-larga"'`.
2. Perfil de envío y su token (lo imprime una sola vez):
   - Mailpit: `make console c='sender-create --user_uuid=<uuid> --sender_name=local --provider_code=smtp_custom --from_email=avisos@example.com --smtp_host=notifier-mailpit --smtp_port=1025 --smtp_security=none'`
   - Gmail (contraseña de aplicación, exige verificación en dos pasos): `--provider_code=gmail --from_email=yo@gmail.com --secret="xxxx xxxx xxxx xxxx"`.
3. Enviar:

```bash
curl -X POST http://localhost:9500/api/v1/messages -H 'Content-Type: application/json' -d '{
  "api_token": "ntf-...",
  "to": ["ana@example.com", "luis@example.com"],
  "subject": "asunto",
  "body": "el mensaje",
  "attachments": [{"file_name": "nota.txt", "content": "texto del adjunto"}]
}'
```

4. Estado: `curl http://localhost:9500/api/v1/messages/<message_uuid> -H "Authorization: Bearer ntf-..."`.

## 3 - contrato

| # | Campo | Qué es |
|---|---|---|
| 1 | `api_token` | el token de envío (o en `Authorization: Bearer`) |
| 2 | `to` · `cc` · `bcc` | un email, varios separados por comas o una lista (máx. `APP_MESSAGE_MAX_RECIPIENTS` entre los tres); sin `to`, el `default_to` del token |
| 3 | `reply_to` | opcional |
| 4 | `subject` · `body` | obligatorio el asunto; `body_format` `text` (defecto) o `html` |
| 5 | `attachments` | lista de `{file_name, content}` o `{file_name, url}` (solo https y hosts de `APP_ATTACHMENT_ALLOWED_HOSTS`); un string suelto vale |
| 6 | `idempotency_key` | opcional: el mismo valor con el mismo token devuelve el mensaje ya creado |

Respuesta `202 {message_uuid, channel, message_status, is_duplicate}`. Reintentos: 5 intentos (1 min, 5 min, 30 min,
2 h); credenciales rechazadas o destinatario inválido (5xx) no se reintentan.
