#docker/app/Dockerfile
# golang alpine: compila y ejecuta el siem-client (agente), sin cgo.
FROM golang:1.26.10-alpine

# apk upgrade: los parches de seguridad de alpine que la imagen aun no trae. libstdc++: lo exige el VS Code Server en alpine (Dev Containers, .devcontainer).
RUN apk upgrade --no-cache && apk add --no-cache libstdc++

ENV CGO_ENABLED=0 \
    GOFLAGS=-buildvcs=false

WORKDIR /app

# arranque: .env + build en cada up (ver entrypoint.sh).
# sed: el checkout en Windows (autocrlf) deja CRLF y el shebang no se ejecutaria.
COPY entrypoint.sh /usr/local/bin/siem-entrypoint
RUN sed -i 's/\r$//' /usr/local/bin/siem-entrypoint && chmod +x /usr/local/bin/siem-entrypoint
ENTRYPOINT ["siem-entrypoint"]
CMD ["/app/build/siem-client", "agent-run"]
