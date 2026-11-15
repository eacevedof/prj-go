/**
 * Resultado del caso de uso: el monitor ya cambiado, tal cual lo devuelve la API.
 */
export class UpdateMonitorResultDto {

    #monitor;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#monitor = primitives.monitor ?? {};
    }

    /**
     * @param {Object} primitives
     * @returns {UpdateMonitorResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new UpdateMonitorResultDto(primitives);
    }

    /**
     * @returns {Object}
     */
    getMonitor() {
        return this.#monitor;
    }

    /**
     * @returns {string}
     */
    getMonitorName() {
        return String(this.#monitor.monitor_name ?? "");
    }
}
