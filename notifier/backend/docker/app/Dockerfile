#docker/app/Dockerfile
# golang alpine: compila y ejecuta la API y la consola (sqlite en Go puro, sin cgo).
FROM golang:1.26.10-alpine

# apk upgrade: los parches de seguridad de alpine que la imagen aun no trae. libstdc++: lo exige el VS Code Server en alpine (Dev Containers, .devcontainer).
# El sendmail de busybox (/usr/sbin/sendmail) ya viene en alpine (alertas de MailTrait).
RUN apk upgrade --no-cache && apk add --no-cache libstdc++

ENV CGO_ENABLED=0 \
    GOFLAGS=-buildvcs=false

WORKDIR /app

# arranque: .env + build + migraciones en cada up (ver entrypoint.sh).
# sed: el checkout en Windows (autocrlf) deja CRLF y el shebang no se ejecutaria.
COPY entrypoint.sh /usr/local/bin/notifier-entrypoint
RUN sed -i 's/\r$//' /usr/local/bin/notifier-entrypoint && chmod +x /usr/local/bin/notifier-entrypoint
ENTRYPOINT ["notifier-entrypoint"]
CMD ["/app/build/api", "-addr=:8080"]
