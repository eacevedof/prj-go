#!/bin/bash
# makefile de la raíz del monorepo prj-go: el gitpush global (sube siem, notifier, yt-downloader y lo de la raíz de una vez).
# Los targets de cada proyecto viven en su makefile: siem/backend, siem/client, notifier/backend, yt-downloader.
# recetas a 4 espacios, sin tabuladores: el prefijo de receta es un espacio (GNU make >= 3.82)
empty :=
.RECIPEPREFIX := $(empty) $(empty)

.PHONY: help gitpush

help: ## Show this help message · make help
    @echo "usage: make [target]"
    @echo
    @echo "targets:"
    @awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
    @echo
    @echo "makefiles de cada proyecto: siem/backend, siem/client, notifier/backend, yt-downloader"

gitpush: ## git push de todo el monorepo · make gitpush m="mensaje"
    clear;
    git pull;
    git add .; git commit -m "$(m)"; git push;
