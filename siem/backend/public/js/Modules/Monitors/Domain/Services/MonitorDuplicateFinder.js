import { MonitorDefaultEnum, MonitorTypeEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorTypeEnum.js";

const WINDOWS_PATH_PATTERN = /^[A-Za-z]:[\\/]/;
const UNIX_ROOT_LENGTH = 1;
// "C:\"
const WINDOWS_ROOT_LENGTH = 3;

export const DuplicateReasonEnum = Object.freeze({
    NAME: "name",
    RESOURCE: "resource",
});

/**
 * El monitor del proyecto con el que chocaria uno nuevo o editado, para avisar en el modal antes de enviarlo:
 * 1) mismo nombre (sin distinguir mayusculas: unico por proyecto) o 2) mismo recurso en un agente que tambien lo
 * aplicaria. Recurso: cpu, memoria y gpu, el tipo (uno por agente) · disco, la ruta · red, la interfaz · carpeta,
 * ruta + glob · log, ruta + glob + regex. Un monitor de "todos" (agent_uuid vacio) choca con cualquiera del proyecto;
 * uno de un agente, con los de ese agente y con los de todos. Gemelo del backend: MonitorResourceKeyer.go (409); si
 * cambia una regla alli, cambia aqui.
 */
export class MonitorDuplicateFinder {

    /**
     * @returns {MonitorDuplicateFinder}
     */
    static getInstance() {
        return new MonitorDuplicateFinder();
    }

    /**
     * @param {Object[]} monitors los del proyecto, tal cual los da la API (snake_case, config{})
     * @param {{monitorUuid: string, monitorType: string, monitorName: string, agentUuid: string, config: Object}} candidate
     *     monitorUuid vacio en el alta; al editar, el propio monitor no cuenta
     * @returns {{monitor: Object, reason: string}|null}
     */
    getDuplicate(monitors, candidate) {
        const otherMonitors = monitors.filter((monitor) => monitor.monitor_uuid !== candidate.monitorUuid);
        const candidateName = candidate.monitorName.trim().toLowerCase();
        const sameNameMonitor = otherMonitors.find((monitor) => String(monitor.monitor_name).toLowerCase() === candidateName);
        if (candidateName !== "" && sameNameMonitor) {
            return { monitor: sameNameMonitor, reason: DuplicateReasonEnum.NAME };
        }

        const resourceKey = this.getResourceKey(candidate.monitorType, candidate.config);
        const sameResourceMonitor = otherMonitors.find((monitor) => monitor.monitor_type === candidate.monitorType
            && this.#isSameScope(candidate.agentUuid, monitor.agent_uuid ?? "")
            && this.getResourceKey(monitor.monitor_type, monitor.config ?? {}) === resourceKey);
        return sameResourceMonitor ? { monitor: sameResourceMonitor, reason: DuplicateReasonEnum.RESOURCE } : null;
    }

    /**
     * @param {string} monitorType
     * @param {Object} config la del formulario o la de la API (con o sin valores por defecto)
     * @returns {string}
     */
    getResourceKey(monitorType, config) {
        const keyParts = [monitorType];
        switch (monitorType) {
            case MonitorTypeEnum.DISK:
                keyParts.push(this.#getComparablePath(config.path || MonitorDefaultEnum.DISK_PATH));
                break;
            case MonitorTypeEnum.NETWORK:
                keyParts.push(String(config.interface ?? "").trim().toLowerCase());
                break;
            case MonitorTypeEnum.FOLDER:
                keyParts.push(this.#getComparablePath(config.path), this.#getFilesPattern(config));
                break;
            case MonitorTypeEnum.LOG_PATTERN:
                // el regex no se recorta, igual que en el backend
                keyParts.push(this.#getComparablePath(config.path), this.#getFilesPattern(config), String(config.match ?? ""));
                break;
        }
        return keyParts.join("|");
    }

    #isSameScope(candidateAgentUuid, monitorAgentUuid) {
        return candidateAgentUuid === "" || monitorAgentUuid === "" || candidateAgentUuid === monitorAgentUuid;
    }

    #getFilesPattern(config) {
        return String(config.pattern ?? "").trim() || MonitorDefaultEnum.FILE_PATTERN;
    }

    /**
     * Sin separador final (salvo la raiz) y, en Windows, sin distinguir mayusculas: /var/log/ = /var/log y
     * C:\Logs = c:\logs.
     */
    #getComparablePath(rawPath) {
        let comparablePath = String(rawPath ?? "").trim();
        const isWindowsPath = WINDOWS_PATH_PATTERN.test(comparablePath);
        const rootLength = isWindowsPath ? WINDOWS_ROOT_LENGTH : UNIX_ROOT_LENGTH;
        while (comparablePath.length > rootLength && /[\\/]$/.test(comparablePath)) {
            comparablePath = comparablePath.slice(0, -1);
        }
        return isWindowsPath ? comparablePath.toLowerCase() : comparablePath;
    }
}
