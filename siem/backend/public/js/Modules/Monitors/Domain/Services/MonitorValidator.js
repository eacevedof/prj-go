import {
    METRIC_MONITOR_TYPES, MONITOR_SEVERITIES, MONITOR_TYPES, MonitorTypeEnum, PERCENT_MONITOR_TYPES,
} from "/js/Modules/Monitors/Domain/Enums/MonitorTypeEnum.js";
import { MonitorsMessageEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsMessageEnum.js";

const MONITOR_NAME_PATTERN = /^[A-Za-z0-9._-]{1,100}$/;
const INTERVAL_MIN_SECONDS = 10;
const INTERVAL_MAX_SECONDS = 86400;
const PERCENT_MAX = 100;

/**
 * Lo evidente de un monitor antes de llamar a la API (el resto, como el nombre duplicado, lo decide ella): nombre,
 * tipo, intervalo y umbrales de los de metricas, y la config propia de cada tipo. Devuelve la pila de errores
 * (vacia = todo bien). Un umbral null es "sin umbral".
 */
export class MonitorValidator {

    /**
     * @returns {MonitorValidator}
     */
    static getInstance() {
        return new MonitorValidator();
    }

    /**
     * @param {{monitorType: string, monitorName: string, intervalSeconds: number, warningThreshold: number|null,
     *     criticalThreshold: number|null, config: Object}} monitor
     * @returns {string[]}
     */
    getErrors(monitor) {
        const errors = [];
        if (!MONITOR_NAME_PATTERN.test(monitor.monitorName)) {
            errors.push(MonitorsMessageEnum.MONITOR_NAME_IS_NOT_VALID);
        }
        if (!MONITOR_TYPES.includes(monitor.monitorType)) {
            errors.push(MonitorsMessageEnum.MONITOR_TYPE_IS_NOT_VALID);
            return errors;
        }
        if (METRIC_MONITOR_TYPES.includes(monitor.monitorType)) {
            errors.push(...this.#getMetricErrors(monitor));
        }
        errors.push(...this.#getConfigErrors(monitor.monitorType, monitor.config ?? {}));
        return errors;
    }

    #getMetricErrors(monitor) {
        const errors = [];
        const interval = monitor.intervalSeconds;
        if (!Number.isInteger(interval) || interval < INTERVAL_MIN_SECONDS || interval > INTERVAL_MAX_SECONDS) {
            errors.push(MonitorsMessageEnum.INTERVAL_MUST_BE_10_TO_86400);
        }

        const isWarningValid = this.#isThresholdValid(monitor.warningThreshold);
        const isCriticalValid = this.#isThresholdValid(monitor.criticalThreshold);
        if (!isWarningValid) errors.push(MonitorsMessageEnum.WARNING_THRESHOLD_IS_NOT_VALID);
        if (!isCriticalValid) errors.push(MonitorsMessageEnum.CRITICAL_THRESHOLD_IS_NOT_VALID);
        if (!isWarningValid || !isCriticalValid) return errors;

        const thresholds = [monitor.warningThreshold, monitor.criticalThreshold].filter((threshold) => threshold !== null);
        if (PERCENT_MONITOR_TYPES.includes(monitor.monitorType) && thresholds.some((threshold) => threshold > PERCENT_MAX)) {
            errors.push(MonitorsMessageEnum.THRESHOLD_MAX_IS_100);
        }
        if (thresholds.length === 2 && monitor.warningThreshold >= monitor.criticalThreshold) {
            errors.push(MonitorsMessageEnum.WARNING_MUST_BE_LOWER_THAN_CRITICAL);
        }
        return errors;
    }

    #isThresholdValid(threshold) {
        return threshold === null || (typeof threshold === "number" && Number.isFinite(threshold) && threshold >= 0);
    }

    #getConfigErrors(monitorType, config) {
        const errors = [];
        const hasPath = [MonitorTypeEnum.DISK, MonitorTypeEnum.FOLDER, MonitorTypeEnum.LOG_PATTERN].includes(monitorType);
        if (hasPath && !String(config.path ?? "").trim()) {
            errors.push(MonitorsMessageEnum.CONFIG_PATH_IS_REQUIRED);
        }
        if (monitorType !== MonitorTypeEnum.LOG_PATTERN) return errors;

        // la regex NO se comprueba aqui: la compila el agente con el motor de Go (RE2), que no es el de JavaScript
        // ("(?i)error" vale en Go y no en JS; un lookahead, al reves). La valida la API con ese mismo motor (400).
        if (!String(config.match ?? "").trim()) {
            errors.push(MonitorsMessageEnum.CONFIG_MATCH_IS_REQUIRED);
        }
        const severity = String(config.severity ?? "");
        if (severity !== "" && !MONITOR_SEVERITIES.includes(severity)) {
            errors.push(MonitorsMessageEnum.CONFIG_SEVERITY_IS_NOT_VALID);
        }
        return errors;
    }
}
