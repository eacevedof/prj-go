#!/bin/bash
# recetas a 4 espacios, sin tabuladores: el prefijo de receta es un espacio (GNU make >= 3.82)
empty :=
.RECIPEPREFIX := $(empty) $(empty)
TODAY := $(shell date +'%Y-%m-%d')
CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

help: ## Show this help message
    @echo "usage: make [target]"
    @echo
    @echo "targets:"
    @grep -E "^(.+)\:\ ##\ (.+)" ${MAKEFILE_LIST} | column -t -c 2 -s ":#"

gitpush: ## git push m=any message
    clear;
    git pull;
    git add .; git commit -m "$(m)"; git push;

update-branch: ## update main branches
    git fetch --all;
    git checkout main; git reset --hard origin/main;
    git checkout develop; git reset --hard origin/develop;
    git checkout $(CURRENT_BRANCH);

sync-develop: ## sync main into develop
    git fetch --all;
    git checkout develop;
    git reset --hard origin/main;
    git push -f;
    git checkout $(CURRENT_BRANCH);
