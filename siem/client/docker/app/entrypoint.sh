#!/bin/sh
# docker/app/entrypoint.sh -- en CADA arranque del contenedor: .env si falta y binario al dia. Idempotente.
set -e

cd /app
[ -f .env ] || cp .env.example .env
go build -o build/siem-client ./cmd/console

exec "$@"
