import { ClientOsEnum } from "/js/Modules/Projects/Domain/Enums/ClientOsEnum.js";

const UNIX_PATH_SEPARATOR = "/";
const WINDOWS_PATH_SEPARATOR = "\\";

/**
 * El comando para desinstalar el siem-client de un agente en su host, segun su SO y la ruta de su ejecutable
 * (client_path, la informa el cliente >= 0.4.0). Deshace lo que hace el instalador: servicio systemd o tarea programada
 * y la carpeta de instalacion (binario, .env, state, spool). Sin ruta conocida se usa la del instalador por defecto.
 */
export class ClientUninstallCommandBuilder {

    /**
     * @returns {ClientUninstallCommandBuilder}
     */
    static getInstance() {
        return new ClientUninstallCommandBuilder();
    }

    /**
     * @param {string} hostOs linux | windows | darwin ("" = aun no envio nada: se asume linux)
     * @param {string} clientPath ruta absoluta del ejecutable, o "" si no se conoce
     * @returns {string} una linea por paso
     */
    getUninstallCommand(hostOs, clientPath) {
        const path = this.getClientPathOrDefault(hostOs, clientPath);
        if (hostOs === ClientOsEnum.WINDOWS) {
            return this.#getWindowsCommand(path);
        }
        return this.#getUnixCommand(hostOs, path);
    }

    /**
     * La ruta que se muestra: la informada o, sin ella, la del instalador por defecto de su SO.
     * @param {string} hostOs
     * @param {string} clientPath
     * @returns {string}
     */
    getClientPathOrDefault(hostOs, clientPath) {
        if (clientPath) return clientPath;
        return hostOs === ClientOsEnum.WINDOWS ? ClientOsEnum.DEFAULT_PATH_WINDOWS : ClientOsEnum.DEFAULT_PATH_LINUX;
    }

    #getUnixCommand(hostOs, clientPath) {
        const lines = [];
        if (hostOs !== ClientOsEnum.MACOS) {
            lines.push(`sudo systemctl disable --now ${ClientOsEnum.SERVICE_NAME}`);
            lines.push(`sudo rm -f /etc/systemd/system/${ClientOsEnum.SERVICE_NAME}.service && sudo systemctl daemon-reload`);
        }
        const installDirectory = this.#getParentDirectory(clientPath, UNIX_PATH_SEPARATOR);
        lines.push(this.#isInstallerDirectory(installDirectory, UNIX_PATH_SEPARATOR)
            ? `sudo rm -rf ${this.#getShellQuoted(installDirectory)}`
            : `sudo rm -f ${this.#getShellQuoted(clientPath)}`);
        return lines.join("\n");
    }

    #getWindowsCommand(clientPath) {
        const installDirectory = this.#getParentDirectory(clientPath, WINDOWS_PATH_SEPARATOR);
        return [
            `Stop-ScheduledTask -TaskName ${ClientOsEnum.SERVICE_NAME} -ErrorAction SilentlyContinue`,
            `Unregister-ScheduledTask -TaskName ${ClientOsEnum.SERVICE_NAME} -Confirm:$false`,
            this.#isInstallerDirectory(installDirectory, WINDOWS_PATH_SEPARATOR)
                ? `Remove-Item -Recurse -Force ${this.#getPowerShellQuoted(installDirectory)}`
                : `Remove-Item -Force ${this.#getPowerShellQuoted(clientPath)}`,
        ].join("\n");
    }

    #getParentDirectory(path, separator) {
        const lastSeparatorIndex = path.lastIndexOf(separator);
        return lastSeparatorIndex > 0 ? path.slice(0, lastSeparatorIndex) : "";
    }

    #isInstallerDirectory(directory, separator) {
        const directoryName = directory.slice(directory.lastIndexOf(separator) + 1);
        return directoryName.toLowerCase() === ClientOsEnum.INSTALL_DIR_NAME;
    }

    #getShellQuoted(value) {
        return `'${value.replaceAll("'", "'\\''")}'`;
    }

    #getPowerShellQuoted(value) {
        return `'${value.replaceAll("'", "''")}'`;
    }
}
