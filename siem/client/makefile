#!/bin/bash
# recetas a 4 espacios, sin tabuladores: el prefijo de receta es un espacio (GNU make >= 3.82)
empty :=
.RECIPEPREFIX := $(empty) $(empty)
# siem-client: el agente. En local corre en un contenedor Go contra el backend de ../backend (make up); en un host
# real se despliega el binario de make dist + .env + siem-client.json. Se lanza desde siem/client/.
DC := docker compose $(if $(wildcard ./docker/.env),--env-file ./docker/.env,) -f docker/docker-compose.yml
APP := cont-siem-client-app
# version unica: la constante APP_VERSION del codigo (la muestra `version` y viaja al backend)
VERSION := $(shell sed -n 's/.*APP_VERSION *= *"\([^"]*\)".*/\1/p' src/modules/shared/infrastructure/repositories/configuration/environment_reader_raw_repository.go)
LDFLAGS := -ldflags "-s -w"

.PHONY: help env up down restart rebuild logs ssh-app build console test vet dist version demo-log gitpush

help: ## Show this help message · make help
    @echo "usage: make [target]"
    @echo
    @echo "targets:"
    @awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

env: ## crea .env y docker/.env desde sus .example (no pisa los existentes) · make env
    @test -f .env || cp .env.example .env
    @test -f docker/.env || cp docker/.env.example docker/.env

up: env ## levanta el agente en local (el backend ya arrancado: red siem-net) · make up
    $(DC) up -d
    docker ps --filter name=siem-client

down: ## para el agente · make down
    $(DC) down

restart: ## reinicia (recompila) · make restart
    $(DC) restart

rebuild: ## reconstruye la imagen y relanza · make rebuild
    $(DC) down
    $(DC) up -d --build --remove-orphans

logs: ## salida del agente (follow) · make logs
    $(DC) logs -f

ssh-app: ## shell en el contenedor · make ssh-app
    docker exec -it $(APP) sh

build: ## recompila build/siem-client dentro del contenedor · make build
    docker exec $(APP) go build -o build/siem-client ./cmd/console

console: ## comando de consola · make console c="config-check"  |  make console c="commands-run --force=1"
    docker exec -it -w /app/build $(APP) ./siem-client $(c)

test: ## tests · make test  |  make test f="AgentConfigValidator"
    docker exec $(APP) go test -count=1 $(if $(f),-run $(f),) ./...

vet: ## go vet (sin gofmt: el repo va a 4 espacios) · make vet
    docker exec $(APP) go vet ./...

DOWNLOADS := ../backend/public/downloads

dist: ## binarios para desplegar: build/dist/siem-client-{linux,win,macos}-<version> + copia en backend/public/downloads (los sirve el panel) · make dist
    docker exec -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 $(APP) go build $(LDFLAGS) -o build/dist/siem-client-linux-$(VERSION) ./cmd/console
    docker exec -e CGO_ENABLED=0 -e GOOS=windows -e GOARCH=amd64 $(APP) go build $(LDFLAGS) -o build/dist/siem-client-win-$(VERSION).exe ./cmd/console
    docker exec -e CGO_ENABLED=0 -e GOOS=darwin -e GOARCH=arm64 $(APP) go build $(LDFLAGS) -o build/dist/siem-client-macos-$(VERSION) ./cmd/console
    @mkdir -p $(DOWNLOADS)
    @rm -f $(DOWNLOADS)/siem-client-*
    cp build/dist/siem-client-linux-$(VERSION) $(DOWNLOADS)/siem-client-linux-$(VERSION)
    cp build/dist/siem-client-win-$(VERSION).exe $(DOWNLOADS)/siem-client-windows-$(VERSION).exe
    cp build/dist/siem-client-macos-$(VERSION) $(DOWNLOADS)/siem-client-macos-$(VERSION)
    @ls -la $(DOWNLOADS)

version: ## version actual (APP_VERSION) · make version
    @echo $(VERSION)

demo-log: ## escribe una linea en storage/demo-logs/app.log (la recoge el agente local) · make demo-log m="ERROR db down"
    @mkdir -p storage/demo-logs
    @echo "$$(date -u '+%Y-%m-%d %H:%M:%S') $(if $(m),$(m),INFO demo line)" >> storage/demo-logs/app.log

gitpush: ## git push · make gitpush m="mensaje"
    clear;
    git pull;
    git add .; git commit -m "$(m)"; git push;
