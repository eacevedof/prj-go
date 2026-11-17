#!/bin/sh
# install.sh -- instala el siem-client (agente del SIEM) en un servidor Linux o macOS, de principio a fin:
#   1. pide lo que falte: la url del panel y el token del agente (sale al crear el agente en el panel)
#   2. descarga el binario de ese panel y comprueba su sha256 (SHA256SUMS del propio panel)
#   3. escribe el .env (a donde envia y con que token) y prueba la conexion (config-check + config-sync)
#   4. lo deja como servicio (systemd en Linux) para que arranque solo
# Lo que se vigila no se configura aqui: son los monitores del proyecto, que el agente descarga solo del panel.
#
# Uso (como root):
#   curl -fsSL https://siem.dominio.com/downloads/install.sh | sudo sh
#   curl -fsSL https://siem.dominio.com/downloads/install.sh | sudo sh -s -- --url https://siem.dominio.com --token <API_TOKEN>
#   sudo sh install.sh [--url URL] [--token TOKEN] [--dir /opt/gosiem-client] [--no-service] [--yes]
# Volver a lanzarlo actualiza el binario y conserva el token del .env si no se le pasa otro.
set -eu

SERVICE_NAME="siem-client"
INSTALL_DIR="/opt/gosiem-client"
PANEL_URL=""
API_TOKEN=""
IS_SERVICE_WANTED=1
IS_ASSUME_YES=0

say() { printf '%s\n' "$*"; }
step() { printf '\n== %s\n' "$*"; }
fail() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }

# hay terminal de verdad (no basta con que exista /dev/tty: en un contenedor o un cron existe pero no se abre)
has_tty() {
    ( : < /dev/tty ) 2>/dev/null
}

# lee de la terminal aunque el script llegue por una tuberia (curl | sh)
ask() {
    prompt="$1"
    default_value="${2:-}"
    is_secret="${3:-0}"
    has_tty || fail "falta $prompt y no hay terminal para preguntarlo: pasalo como argumento (--help)"
    if [ -n "$default_value" ]; then
        printf '%s [%s]: ' "$prompt" "$default_value" > /dev/tty
    else
        printf '%s: ' "$prompt" > /dev/tty
    fi
    if [ "$is_secret" = "1" ]; then
        stty -echo < /dev/tty 2>/dev/null || true
        IFS= read -r answer < /dev/tty || answer=""
        stty echo < /dev/tty 2>/dev/null || true
        printf '\n' > /dev/tty
    else
        IFS= read -r answer < /dev/tty || answer=""
    fi
    [ -n "$answer" ] || answer="$default_value"
    printf '%s' "$answer"
}

download() {
    if command -v curl > /dev/null 2>&1; then
        curl -fsSL -o "$2" "$1"
    elif command -v wget > /dev/null 2>&1; then
        wget -q -O "$2" "$1"
    else
        fail "hace falta curl o wget"
    fi
}

get_sha256() {
    if command -v sha256sum > /dev/null 2>&1; then
        sha256sum "$1" | cut -d' ' -f1
    else
        shasum -a 256 "$1" | cut -d' ' -f1
    fi
}

while [ $# -gt 0 ]; do
    case "$1" in
        --url) PANEL_URL="${2:-}"; shift 2 ;;
        --token) API_TOKEN="${2:-}"; shift 2 ;;
        --dir) INSTALL_DIR="${2:-}"; shift 2 ;;
        --no-service) IS_SERVICE_WANTED=0; shift ;;
        --yes|-y) IS_ASSUME_YES=1; shift ;;
        --help|-h) sed -n '2,16p' "$0" 2>/dev/null || say "uso: install.sh [--url URL] [--token TOKEN] [--dir DIR] [--no-service] [--yes]"; exit 0 ;;
        *) fail "argumento desconocido: $1 (--help)" ;;
    esac
done

case "$(uname -s)" in
    Linux) OS_NAME="linux" ;;
    Darwin) OS_NAME="macos" ;;
    *) fail "sistema no soportado: $(uname -s) (en Windows, install.bat)" ;;
esac
case "$(uname -m)" in
    x86_64|amd64) [ "$OS_NAME" = "linux" ] || fail "en macOS solo hay binario para Apple Silicon (arm64)" ;;
    arm64|aarch64) [ "$OS_NAME" = "macos" ] || fail "en Linux solo hay binario para amd64 (x86_64)" ;;
    *) fail "arquitectura no soportada: $(uname -m)" ;;
esac
[ "$(id -u)" = "0" ] || [ "$IS_SERVICE_WANTED" = "0" ] || fail "hay que lanzarlo como root (sudo) para instalar el servicio; o --no-service y un --dir con permiso"

step "1 - panel y agente"
[ -n "$PANEL_URL" ] || PANEL_URL=$(ask "url del panel (p. ej. https://siem.dominio.com)")
PANEL_URL=$(printf '%s' "$PANEL_URL" | sed 's:/*$::')
case "$PANEL_URL" in
    http://*|https://*) ;;
    *) fail "la url del panel tiene que empezar por http:// o https://" ;;
esac
existing_token=""
if [ -f "$INSTALL_DIR/.env" ]; then
    existing_token=$(sed -n 's/^APP_SIEM_API_TOKEN=//p' "$INSTALL_DIR/.env" | head -1)
fi
if [ -z "$API_TOKEN" ] && [ -n "$existing_token" ]; then
    say "se conserva el token del .env que ya habia en $INSTALL_DIR"
    API_TOKEN="$existing_token"
fi
[ -n "$API_TOKEN" ] || API_TOKEN=$(ask "token del agente (sale una vez al crearlo en el panel)" "" 1)
[ -n "$API_TOKEN" ] || fail "sin token no se puede enviar nada: crea el agente en el panel (Descargas o el detalle del proyecto)"
say "panel: $PANEL_URL · carpeta: $INSTALL_DIR · servicio: $([ "$IS_SERVICE_WANTED" = "1" ] && echo si || echo no)"
if [ "$IS_ASSUME_YES" = "0" ] && has_tty; then
    confirm=$(ask "seguir? (s/n)" "s")
    case "$confirm" in s|S|si|SI|y|Y|yes) ;; *) fail "cancelado" ;; esac
fi

step "2 - binario"
temp_dir=$(mktemp -d)
trap 'rm -rf "$temp_dir"' EXIT
download "$PANEL_URL/downloads/SHA256SUMS" "$temp_dir/SHA256SUMS" || fail "no se puede descargar $PANEL_URL/downloads/SHA256SUMS: comprueba la url y que el servidor llega al panel"
# la version mas alta publicada para este sistema (por numeros: la 0.10.0 va detras de la 0.9.0)
version=$(grep -o "siem-client-$OS_NAME-[0-9][0-9.]*[0-9]" "$temp_dir/SHA256SUMS" | sed "s/^siem-client-$OS_NAME-//" \
    | sort -t. -n -k1,1 -k2,2 -k3,3 | tail -1)
[ -n "$version" ] || fail "el panel no publica binario para $OS_NAME (make dist en el cliente)"
file_name="siem-client-$OS_NAME-$version"
expected_sha256=$(grep " \*\{0,1\}$file_name\$" "$temp_dir/SHA256SUMS" | cut -d' ' -f1)
say "descargando $file_name"
download "$PANEL_URL/downloads/$file_name" "$temp_dir/siem-client" || fail "no se puede descargar $file_name"
[ "$(get_sha256 "$temp_dir/siem-client")" = "$expected_sha256" ] || fail "el sha256 de $file_name no coincide: descarga corrupta o manipulada"
say "sha256 correcto"
mkdir -p "$INSTALL_DIR"
chmod +x "$temp_dir/siem-client"
if [ "$IS_SERVICE_WANTED" = "1" ] && command -v systemctl > /dev/null 2>&1; then
    systemctl stop "$SERVICE_NAME" 2>/dev/null || true
fi
mv -f "$temp_dir/siem-client" "$INSTALL_DIR/siem-client"
[ "$OS_NAME" = "macos" ] && xattr -d com.apple.quarantine "$INSTALL_DIR/siem-client" 2>/dev/null || true
say "instalado en $INSTALL_DIR/siem-client ($("$INSTALL_DIR/siem-client" version 2>/dev/null | tail -1))"

step "3 - configuracion"
umask 077
cat > "$INSTALL_DIR/.env" <<EOF
APP_ENV=production
APP_SIEM_ENDPOINT_URL=$PANEL_URL/api/v1/ingest-batches
APP_SIEM_API_TOKEN=$API_TOKEN
EOF
chmod 600 "$INSTALL_DIR/.env"
cd "$INSTALL_DIR"
"$INSTALL_DIR/siem-client" config-sync || fail "el panel no acepta el token o no responde: revisa la url y el token (config-sync)"
"$INSTALL_DIR/siem-client" config-check > /dev/null || fail "la configuracion no es valida: $INSTALL_DIR/siem-client config-check"
say "conexion y configuracion correctas"

step "4 - servicio"
if [ "$IS_SERVICE_WANTED" = "0" ]; then
    say "sin servicio (--no-service). Para arrancarlo: cd $INSTALL_DIR && ./siem-client agent-run"
elif [ "$OS_NAME" = "linux" ] && command -v systemctl > /dev/null 2>&1; then
    cat > "/etc/systemd/system/$SERVICE_NAME.service" <<EOF
[Unit]
Description=siem-client (agente del SIEM)
After=network-online.target
Wants=network-online.target

[Service]
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/siem-client agent-run
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
    systemctl daemon-reload
    systemctl enable --now "$SERVICE_NAME" > /dev/null
    say "servicio $SERVICE_NAME activo · logs: journalctl -u $SERVICE_NAME -f"
else
    say "sin systemd: arrancalo con  cd $INSTALL_DIR && ./siem-client agent-run  (o un LaunchDaemon en macOS)"
fi

step "listo"
say "en el panel, el agente pasa a online en menos de un minuto"
