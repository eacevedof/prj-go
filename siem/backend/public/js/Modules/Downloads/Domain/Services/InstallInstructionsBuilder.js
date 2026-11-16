import { DownloadsKeyEnum } from "/js/Modules/Downloads/Domain/Enums/DownloadsKeyEnum.js";
import { DownloadsOsEnum } from "/js/Modules/Downloads/Domain/Enums/DownloadsOsEnum.js";

const API_TOKEN_PLACEHOLDER = "<API_TOKEN>";
const LINUX_FOLDER = "/opt/siem-client";
const WINDOWS_FOLDER = "C:\\programs\\siem-client";
const MACOS_FOLDER = "/usr/local/siem-client";

/**
 * Los pasos para instalar y configurar un binario del agente, con los comandos ya montados para copiar: descarga
 * desde este mismo backend, comprobacion del sha256, el .env (endpoint + token) y como dejarlo corriendo. Ya no hace
 * falta siem-client.json: lo que se vigila se crea en el panel y el agente lo descarga solo. Funcion pura: mismo
 * binario y mismos datos = mismo texto.
 */
export class InstallInstructionsBuilder {

    /**
     * @returns {InstallInstructionsBuilder}
     */
    static getInstance() {
        return new InstallInstructionsBuilder();
    }

    /**
     * @param {{file_name: string, os: string, sha256: string, url: string}} download
     * @param {string} origin de donde se sirve el panel (https://siem.dominio.com)
     * @param {string} endpointUrl url de ingesta que va en el .env
     * @param {string} apiToken el del agente; vacio = el marcador <API_TOKEN>
     * @returns {{title: string, description: string, language: string, code: string}[]}
     */
    getSteps(download, origin, endpointUrl, apiToken = "") {
        const context = {
            fileName: String(download[DownloadsKeyEnum.FILE_NAME] ?? ""),
            sha256: String(download[DownloadsKeyEnum.SHA256] ?? ""),
            downloadUrl: `${origin}${download[DownloadsKeyEnum.URL] ?? ""}`,
            endpointUrl,
            apiToken: apiToken.trim() || API_TOKEN_PLACEHOLDER,
        };
        switch (download[DownloadsKeyEnum.OS]) {
            case DownloadsOsEnum.LINUX:
                return this.#getLinuxSteps(context);
            case DownloadsOsEnum.WINDOWS:
                return this.#getWindowsSteps(context);
            case DownloadsOsEnum.MACOS:
                return this.#getMacosSteps(context);
            default:
                return [];
        }
    }

    /**
     * La instalacion asistida en una linea: descarga el instalador del panel y lo lanza con la url y, si se conoce, el
     * token (sin token, el instalador lo pregunta: asi nunca queda un marcador sin sustituir en el servidor).
     * @param {string} os linux | macos | windows
     * @param {string} origin la direccion del backend que veran los servidores
     * @param {string} apiToken "" = que lo pregunte
     * @returns {{language: string, code: string}}
     */
    getInstallerCommand(os, origin, apiToken = "") {
        const token = apiToken.trim();
        if (os === DownloadsOsEnum.WINDOWS) {
            const tokenArgument = token ? ` -Token "${token}"` : "";
            return {
                language: "PowerShell como administrador",
                code: `Invoke-WebRequest -UseBasicParsing "${origin}/downloads/install.ps1" -OutFile "$env:TEMP\\siem-install.ps1"\n`
                    + `powershell -NoProfile -ExecutionPolicy Bypass -File "$env:TEMP\\siem-install.ps1" -Url "${origin}"${tokenArgument}`,
            };
        }
        const tokenArgument = token ? ` --token "${token}"` : "";
        return {
            language: os === DownloadsOsEnum.MACOS ? "zsh" : "bash",
            code: `curl -fsSL "${origin}/downloads/install.sh" | sudo sh -s -- --url "${origin}"${tokenArgument}`,
        };
    }

    /**
     * Todos los pasos en un solo script, para pegarlo de una vez en el servidor: cada paso precedido de un comentario
     * con su numero y su titulo (# vale en bash, zsh y PowerShell).
     * @param {{title: string, code: string}[]} steps los de getSteps
     * @returns {string}
     */
    getFullScript(steps) {
        return steps.map((step, index) => `# ${index + 1} - ${step.title}\n${step.code}`).join("\n\n");
    }

    #getEnvFile(context) {
        return `APP_ENV=production\nAPP_SIEM_ENDPOINT_URL=${context.endpointUrl}\nAPP_SIEM_API_TOKEN=${context.apiToken}`;
    }

    #getLinuxSteps(context) {
        return [
            {
                title: "Descargar, comprobar y dar permiso de ejecución",
                description: `El binario queda en ${LINUX_FOLDER}/siem-client (sin versión en el nombre: al actualizar no hay que tocar el servicio).`,
                language: "bash",
                code: `sudo mkdir -p ${LINUX_FOLDER} && cd ${LINUX_FOLDER}\n`
                    + `sudo curl -fSL -o siem-client "${context.downloadUrl}"\n`
                    + `echo "${context.sha256}  siem-client" | sha256sum -c -\n`
                    + `sudo chmod +x siem-client`,
            },
            {
                title: "Crear el .env (a dónde envía y con qué token)",
                description: "El token es el del agente: sale una sola vez al crearlo en su proyecto (apartado 2 del detalle del proyecto).",
                language: "bash",
                code: `sudo tee ${LINUX_FOLDER}/.env > /dev/null <<'EOF'\n${this.#getEnvFile(context)}\nEOF\n`
                    + `sudo chmod 600 ${LINUX_FOLDER}/.env`,
            },
            {
                title: "Comprobar la configuración y la conexión",
                description: "config-check valida endpoint y token sin enviar nada; config-sync baja del panel lo que este agente tiene que vigilar.",
                language: "bash",
                code: `cd ${LINUX_FOLDER}\nsudo ./siem-client version\nsudo ./siem-client config-check\nsudo ./siem-client config-sync`,
            },
            {
                title: "Dejarlo como servicio (systemd)",
                description: "Arranca con el sistema y se reinicia solo. En el panel, el agente pasa a online en menos de un minuto.",
                language: "bash",
                code: `sudo tee /etc/systemd/system/siem-client.service > /dev/null <<'EOF'\n`
                    + `[Unit]\nDescription=siem-client\nAfter=network-online.target\nWants=network-online.target\n\n`
                    + `[Service]\nWorkingDirectory=${LINUX_FOLDER}\nExecStart=${LINUX_FOLDER}/siem-client agent-run\nRestart=always\nRestartSec=5\n\n`
                    + `[Install]\nWantedBy=multi-user.target\nEOF\n`
                    + `sudo systemctl daemon-reload\nsudo systemctl enable --now siem-client\nsudo journalctl -u siem-client -f`,
            },
        ];
    }

    #getWindowsSteps(context) {
        return [
            {
                title: "Descargar y comprobar (PowerShell como administrador)",
                description: `El binario queda en ${WINDOWS_FOLDER}\\siem-client.exe. La última línea debe responder True.`,
                language: "PowerShell",
                code: `New-Item -ItemType Directory -Force "${WINDOWS_FOLDER}" | Out-Null\n`
                    + `Set-Location "${WINDOWS_FOLDER}"\n`
                    + `Invoke-WebRequest -Uri "${context.downloadUrl}" -OutFile siem-client.exe\n`
                    + `(Get-FileHash .\\siem-client.exe -Algorithm SHA256).Hash.ToLower() -eq "${context.sha256}"`,
            },
            {
                title: "Crear el .env (a dónde envía y con qué token)",
                description: "El token es el del agente: sale una sola vez al crearlo en su proyecto.",
                language: "PowerShell",
                code: `@"\n${this.#getEnvFile(context)}\n"@ | Set-Content -Encoding ascii "${WINDOWS_FOLDER}\\.env"`,
            },
            {
                title: "Comprobar la configuración y la conexión",
                description: "config-check valida endpoint y token sin enviar nada; config-sync baja del panel lo que este agente tiene que vigilar.",
                language: "PowerShell",
                code: `Set-Location "${WINDOWS_FOLDER}"\n.\\siem-client.exe version\n.\\siem-client.exe config-check\n.\\siem-client.exe config-sync`,
            },
            {
                title: "Dejarlo arrancando con el sistema (tarea programada)",
                description: "Corre como SYSTEM al iniciar Windows. Para probarlo antes en primer plano: .\\siem-client.exe agent-run",
                language: "PowerShell",
                code: `$action = New-ScheduledTaskAction -Execute "${WINDOWS_FOLDER}\\siem-client.exe" -Argument "agent-run" -WorkingDirectory "${WINDOWS_FOLDER}"\n`
                    + `$trigger = New-ScheduledTaskTrigger -AtStartup\n`
                    + `$settings = New-ScheduledTaskSettingsSet -RestartCount 999 -RestartInterval (New-TimeSpan -Minutes 1) -ExecutionTimeLimit ([TimeSpan]::Zero)\n`
                    + `Register-ScheduledTask -TaskName "siem-client" -Action $action -Trigger $trigger -Settings $settings -User "SYSTEM" -RunLevel Highest\n`
                    + `Start-ScheduledTask -TaskName "siem-client"`,
            },
        ];
    }

    #getMacosSteps(context) {
        return [
            {
                title: "Descargar, comprobar y quitar la cuarentena",
                description: `El binario queda en ${MACOS_FOLDER}/siem-client. macOS marca lo descargado: sin el xattr no deja ejecutarlo.`,
                language: "zsh",
                code: `sudo mkdir -p ${MACOS_FOLDER} && cd ${MACOS_FOLDER}\n`
                    + `sudo curl -fSL -o siem-client "${context.downloadUrl}"\n`
                    + `echo "${context.sha256}  siem-client" | shasum -a 256 -c -\n`
                    + `sudo chmod +x siem-client\n`
                    + `sudo xattr -d com.apple.quarantine siem-client 2>/dev/null || true`,
            },
            {
                title: "Crear el .env (a dónde envía y con qué token)",
                description: "El token es el del agente: sale una sola vez al crearlo en su proyecto.",
                language: "zsh",
                code: `sudo tee ${MACOS_FOLDER}/.env > /dev/null <<'EOF'\n${this.#getEnvFile(context)}\nEOF\n`
                    + `sudo chmod 600 ${MACOS_FOLDER}/.env`,
            },
            {
                title: "Comprobar y arrancar",
                description: "agent-run se queda en primer plano (Ctrl+C para parar). Para dejarlo fijo, un LaunchDaemon que lance ese mismo comando.",
                language: "zsh",
                code: `cd ${MACOS_FOLDER}\nsudo ./siem-client config-check\nsudo ./siem-client config-sync\nsudo ./siem-client agent-run`,
            },
        ];
    }
}
