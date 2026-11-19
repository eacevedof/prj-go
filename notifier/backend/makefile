#!/bin/bash
# recetas a 4 espacios, sin tabuladores: el prefijo de receta es un espacio (GNU make >= 3.82)
empty :=
.RECIPEPREFIX := $(empty) $(empty)
# Stack local (API + workers en un contenedor Go, y Mailpit). Se lanza desde notifier/backend/. docker/.env es opcional.
# Cada arranque de cont-notifier-backend-app ya hace .env + go build + migraciones (docker/app/entrypoint.sh):
# los targets build y migrate quedan para relanzarlos a mano sin reiniciar.
DC := docker compose $(if $(wildcard ./docker/.env),--env-file ./docker/.env,) -f docker/docker-compose.yml
APP := cont-notifier-backend-app
HTTP_PORT ?= 9500
MAILPIT_WEB_PORT ?= 9525

# targets que no son ficheros: sin esto "make console" o "make test" no hacen nada, porque existen esas carpetas
.PHONY: help env start up down restart rebuild ps logs config ssh-app ssh-be build console migrate test vet health gitpush

help: ## Show this help message · make help
    @echo "usage: make [target]"
    @echo
    @echo "targets:"
    @awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-18s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

env: ## crea .env y docker/.env desde sus .example (no pisa los existentes) · make env
    @test -f .env || cp .env.example .env
    @test -f docker/.env || cp docker/.env.example docker/.env

start: env up ## primer arranque: env + contenedor (build y migraciones van en el arranque) · make start
    @echo "api http://localhost:$(HTTP_PORT) · mailpit http://localhost:$(MAILPIT_WEB_PORT)"

up: ## levanta el contenedor (+ go build y migraciones) · make up
    $(DC) up -d
    docker ps --filter name=notifier-

down: ## para y elimina el contenedor · make down
    $(DC) down

restart: ## reinicia el contenedor (recompila y pasa las migraciones: asi se aplican los cambios de codigo) · make restart
    $(DC) restart

rebuild: ## reconstruye la imagen y relanza · make rebuild
    $(DC) down
    $(DC) up -d --build --remove-orphans

ps: ## contenedores del stack · make ps
    docker ps --filter name=notifier-

logs: ## logs del contenedor (follow) · make logs
    $(DC) logs -f

config: ## compose resuelto · make config
    $(DC) config

ssh-app: ## shell en el contenedor de la api · make ssh-app
    docker exec -it $(APP) sh

ssh-be: ## shell en el contenedor de la api, en /app/build y mostrando el help de la consola · make ssh-be
    docker exec -it -w /app/build $(APP) sh -c "./console help; exec sh"

build: ## recompila los binarios (build/api y build/console) dentro del contenedor · make build
    docker exec $(APP) go build -o build/api ./cmd/api
    docker exec $(APP) go build -o build/console ./cmd/console

console: ## comando de consola · make console c="help"  |  make console c='sender-list'
    docker exec -it -w /app/build $(APP) ./console $(c)

migrate: ## migraciones .sql pendientes a mano (database/migrations) · make migrate
    docker exec -w /app/build $(APP) ./console run-migrations

test: ## tests · make test  |  make test f="MimeMessageBuilder"
    docker exec $(APP) go test -count=1 $(if $(f),-run $(f),) ./...

vet: ## go vet (sin gofmt: el repo va a 4 espacios) · make vet
    docker exec $(APP) go vet ./...

health: ## health-check del front controller · make health
    curl -s http://localhost:$(HTTP_PORT)/health

gitpush: ## git push · make gitpush m="mensaje"
    clear;
    git pull;
    git add .; git commit -m "$(m)"; git push;
