#!/bin/bash
# recetas a 4 espacios, sin tabuladores: el prefijo de receta es un espacio (GNU make >= 3.82)
empty :=
.RECIPEPREFIX := $(empty) $(empty)
# trafficker: agente de red (conexiones y sockets por proceso) que manda los cambios al backend del siem. En local corre
# en un contenedor Go con la red y los procesos del host (network_mode/pid host); en un host real se despliega el
# binario de make dist + .env. Se lanza desde prj-go/trafficker/.
DC := docker compose $(if $(wildcard ./docker/.env),--env-file ./docker/.env,) -f docker/docker-compose.yml
APP := cont-trafficker-app
# version unica: la constante APP_VERSION del codigo (la muestra `version` y viaja al backend)
VERSION := $(shell sed -n 's/.*APP_VERSION *= *"\([^"]*\)".*/\1/p' src/modules/shared/infrastructure/repositories/configuration/environment_reader_raw_repository.go)
LDFLAGS := -ldflags "-s -w"

.PHONY: help env up down restart rebuild logs ssh-app build console list test vet dist version gitpush

help: ## Show this help message · make help
    @echo "usage: make [target]"
    @echo
    @echo "targets:"
    @awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

env: ## crea .env y docker/.env desde sus .example (no pisa los existentes) · make env
    @test -f .env || cp .env.example .env
    @test -f docker/.env || cp docker/.env.example docker/.env

up: env ## levanta el contenedor (red y procesos del host) · make up
    $(DC) up -d
    docker ps --filter name=$(APP)

down: ## para el contenedor · make down
    $(DC) down

restart: ## reinicia (recompila) · make restart
    $(DC) restart

rebuild: ## reconstruye la imagen y relanza · make rebuild
    $(DC) down
    $(DC) up -d --build --remove-orphans

logs: ## salida del contenedor (follow) · make logs
    $(DC) logs -f

ssh-app: ## shell en el contenedor · make ssh-app
    docker exec -it $(APP) sh

build: ## recompila build/trafficker dentro del contenedor · make build
    docker exec $(APP) go build -o build/trafficker ./cmd/console

console: ## comando de consola · make console c="help"  |  make console c="connections-list --filter=remote"
    docker exec -it -w /app/build $(APP) ./trafficker $(c)

list: ## conexiones de ahora del host · make list  |  make list f="established,remote"
    docker exec -it -w /app/build $(APP) ./trafficker connections-list $(if $(f),--filter="$(f)",)

test: ## tests · make test  |  make test f="ProcNetLineParser"
    docker exec $(APP) go test -count=1 $(if $(f),-run $(f),) ./...

vet: ## go vet (sin gofmt: el repo va a 4 espacios) · make vet
    docker exec $(APP) go vet ./...

dist: ## binario Linux en build/dist (0.1.0 solo lee /proc de Linux) · make dist
    docker exec -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 $(APP) go build $(LDFLAGS) -o build/dist/trafficker-linux-$(VERSION) ./cmd/console
    docker exec -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=arm64 $(APP) go build $(LDFLAGS) -o build/dist/trafficker-linux-arm64-$(VERSION) ./cmd/console
    @ls -la build/dist

version: ## version actual (APP_VERSION) · make version
    @echo $(VERSION)

gitpush: ## git push · make gitpush m="mensaje"
    clear;
    git pull;
    git add .; git commit -m "$(m)"; git push;
