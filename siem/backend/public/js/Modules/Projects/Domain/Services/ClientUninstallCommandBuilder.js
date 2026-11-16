import { ClientOsEnum } from "/js/Modules/Projects/Domain/Enums/ClientOsEnum.js";

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
        const installDir = this.#getParentDir(clientPath, "/");
        lines.push(this.#isInstallerDir(installDir, "/")
            ? `sudo rm -rf ${this.#getShellQuoted(installDir)}`
            : `sudo rm -f ${this.#getShellQuoted(clientPath)}`);
        return lines.join("\n");
    }

    #getWindowsCommand(clientPath) {
        const installDir = this.#getParentDir(clientPath, "\\");
        return [
            `Stop-ScheduledTask -TaskName ${ClientOsEnum.SERVICE_NAME} -ErrorAction SilentlyContinue`,
            `Unregister-ScheduledTask -TaskName ${ClientOsEnum.SERVICE_NAME} -Confirm:$false`,
            this.#isInstallerDir(installDir, "\\")
                ? `Remove-Item -Recurse -Force ${this.#getPowerShellQuoted(installDir)}`
                : `Remove-Item -Force ${this.#getPowerShellQuoted(clientPath)}`,
        ].join("\n");
    }

    #getParentDir(path, separator) {
        const lastSeparatorIndex = path.lastIndexOf(separator);
        return lastSeparatorIndex > 0 ? path.slice(0, lastSeparatorIndex) : "";
    }

    #isInstallerDir(dir, separator) {
        const dirName = dir.slice(dir.lastIndexOf(separator) + 1);
        return dirName.toLowerCase() === ClientOsEnum.INSTALL_DIR_NAME;
    }

    #getShellQuoted(value) {
        return `'${value.replaceAll("'", "'\\''")}'`;
    }

    #getPowerShellQuoted(value) {
        return `'${value.replaceAll("'", "''")}'`;
    }
}
