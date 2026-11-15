import {
    METRIC_MONITOR_TYPES, MonitorDefaultEnum, MonitorTypeEnum,
} from "/js/Modules/Monitors/Domain/Enums/MonitorTypeEnum.js";

/**
 * Entrada del caso de uso: los campos del formulario como primitivos (strings con trim, booleanos). De ellos sale el
 * monitor con el contrato de la API: la config propia de su tipo, el intervalo (60 s por defecto) y los umbrales
 * (vacio = null = sin umbral; NaN si no era un numero, para que el validador lo rechace). Los monitores de lineas no
 * llevan umbrales.
 */
export class CreateMonitorDto {

    #dashboardToken;
    #projectUuid;
    #monitorType;
    #monitorName;
    #agentUuid;
    #isEnabled;
    #intervalSeconds;
    #warningThreshold;
    #criticalThreshold;
    #config;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
        this.#projectUuid = String(primitives.projectUuid ?? "").trim();
        this.#monitorType = String(primitives.monitorType ?? "").trim();
        this.#monitorName = String(primitives.monitorName ?? "").trim();
        this.#agentUuid = String(primitives.agentUuid ?? "").trim();
        this.#isEnabled = primitives.isEnabled !== false;

        const isMetric = METRIC_MONITOR_TYPES.includes(this.#monitorType);
        const rawInterval = String(primitives.intervalSeconds ?? "").trim();
        this.#intervalSeconds = rawInterval === "" || !isMetric ? MonitorDefaultEnum.INTERVAL_SECONDS : Number(rawInterval);
        this.#warningThreshold = isMetric ? CreateMonitorDto.#getThreshold(primitives.warningThreshold) : null;
        this.#criticalThreshold = isMetric ? CreateMonitorDto.#getThreshold(primitives.criticalThreshold) : null;
        this.#config = CreateMonitorDto.#getConfig(this.#monitorType, primitives);
    }

    /**
     * @param {Object} primitives
     * @returns {CreateMonitorDto}
     */
    static fromPrimitives(primitives = {}) {
        return new CreateMonitorDto(primitives);
    }

    static #getThreshold(rawThreshold) {
        const text = String(rawThreshold ?? "").trim().replace(",", ".");
        return text === "" ? null : Number(text);
    }

    static #getConfig(monitorType, primitives) {
        const path = String(primitives.path ?? "").trim();
        const pattern = String(primitives.pattern ?? "").trim() || MonitorDefaultEnum.FILE_PATTERN;
        const isRecursive = primitives.isRecursive === true;
        switch (monitorType) {
            case MonitorTypeEnum.DISK:
                return { path: path || MonitorDefaultEnum.DISK_PATH };
            case MonitorTypeEnum.NETWORK:
                return { interface: String(primitives.interface ?? "").trim() };
            case MonitorTypeEnum.FOLDER:
                return { path, pattern, is_recursive: isRecursive };
            case MonitorTypeEnum.LOG_PATTERN:
                return {
                    path, pattern, is_recursive: isRecursive,
                    match: String(primitives.match ?? "").trim(),
                    severity: String(primitives.severity ?? "").trim(),
                };
            default:
                return {};
        }
    }

    /**
     * @returns {string}
     */
    getDashboardToken() {
        return this.#dashboardToken;
    }

    /**
     * @returns {string}
     */
    getProjectUuid() {
        return this.#projectUuid;
    }

    /**
     * @returns {string}
     */
    getMonitorType() {
        return this.#monitorType;
    }

    /**
     * @returns {string}
     */
    getMonitorName() {
        return this.#monitorName;
    }

    /**
     * @returns {number}
     */
    getIntervalSeconds() {
        return this.#intervalSeconds;
    }

    /**
     * @returns {number|null}
     */
    getWarningThreshold() {
        return this.#warningThreshold;
    }

    /**
     * @returns {number|null}
     */
    getCriticalThreshold() {
        return this.#criticalThreshold;
    }

    /**
     * @returns {Object}
     */
    getConfig() {
        return this.#config;
    }

    /**
     * @returns {Object} el cuerpo de POST / PUT de la API (claves snake_case)
     */
    getApiMonitor() {
        return {
            monitor_type: this.#monitorType,
            monitor_name: this.#monitorName,
            agent_uuid: this.#agentUuid,
            is_enabled: this.#isEnabled,
            interval_seconds: this.#intervalSeconds,
            config: this.#config,
            warning_threshold: this.#warningThreshold,
            critical_threshold: this.#criticalThreshold,
        };
    }
}
