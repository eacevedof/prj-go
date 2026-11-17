/**
 * Textos del camino feliz de Descargas (los ve una persona: en español).
 */
export const DownloadsTextEnum = Object.freeze({
    NO_PROJECTS: "(no hay proyectos: créalo en Proyectos)",
    TOKEN_READY: "Token listo: va en los comandos del paso 3.",
    TOKEN_MISSING: "Sin token todavía: el instalador lo pedirá al ejecutarlo.",
    // lo que hace el instalador asistido, para que se sepa antes de lanzarlo
    INSTALLER_DOES: Object.freeze([
        "pide lo que falte (dirección del panel, token)",
        "descarga el binario de este panel y comprueba su sha256",
        "escribe el .env y prueba la conexión con el panel",
        "lo deja arrancando con el sistema (systemd en Linux, tarea programada en Windows)",
    ]),
    MACOS_TROUBLESHOOTING_TITLE: "Si falla en Mac",
    MACOS_TROUBLESHOOTING_INTRO: "El binario de macOS se compila desde Linux (Go lo permite sin un Mac), pero no está firmado por Apple ni se ha probado aún en un Mac de verdad. Estos son los fallos esperables y cómo salir de ellos. <ruta> es la carpeta del agente: /opt/gosiem-client con el instalador, /usr/local/gosiem-client siguiendo el paso a paso.",
    // {problem, fix, code?}: code = comando para copiar ("" = sin comando)
    MACOS_TROUBLESHOOTING: Object.freeze([
        Object.freeze({
            problem: "«No se puede abrir porque Apple no puede comprobar si contiene software malicioso»",
            fix: "Es Gatekeeper: el binario no está firmado ni notarizado por Apple. El instalador ya le quita la marca de cuarentena; si lo descargaste a mano, quítasela tú. Desde el escritorio también vale: Ajustes del Sistema › Privacidad y seguridad › «Abrir igualmente».",
            code: "sudo xattr -d com.apple.quarantine <ruta>/siem-client",
        }),
        Object.freeze({
            problem: "«zsh: killed» o «Killed: 9» nada más ejecutarlo",
            fix: "macOS mata en Apple Silicon los binarios sin firma válida (pasa si se ha copiado o modificado). Fírmalo en local (firma ad hoc, sin cuenta de Apple) y vuelve a lanzarlo.",
            code: "sudo codesign --force --sign - <ruta>/siem-client",
        }),
        Object.freeze({
            problem: "«bad CPU type in executable»",
            fix: "El binario es solo para Mac con Apple Silicon (M1 o posterior). En un Mac con procesador Intel no arranca: hace falta otra compilación (darwin/amd64), aún no publicada.",
            code: "",
        }),
        Object.freeze({
            problem: "No arranca solo al reiniciar el Mac",
            fix: "El instalador deja el servicio con systemd en Linux y una tarea programada en Windows, pero en macOS todavía no crea el LaunchDaemon: tras instalarlo, arráncalo a mano (o crea tú el LaunchDaemon).",
            code: "cd <ruta> && sudo ./siem-client agent-run",
        }),
        Object.freeze({
            problem: "Los monitores de cpu, memoria, disco o red salen como «no disponible»",
            fix: "Esas métricas del host solo se leen en Linux (/proc). En macOS el agente sí vigila logs (patrón de log), carpetas y comandos, y envía el latido.",
            code: "",
        }),
        Object.freeze({
            problem: "El agente sale como «nunca visto» en el panel",
            fix: "Comprueba que el Mac llega al panel y que el token es bueno: config-sync baja la configuración y falla si la url o el token no valen.",
            code: "cd <ruta> && sudo ./siem-client config-sync",
        }),
    ]),
});
