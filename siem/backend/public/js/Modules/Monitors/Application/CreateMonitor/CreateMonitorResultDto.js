/**
 * Resultado del caso de uso: el monitor creado, tal cual lo devuelve la API (objeto plano, claves snake_case).
 */
export class CreateMonitorResultDto {

    #monitor;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#monitor = primitives.monitor ?? {};
    }

    /**
     * @param {Object} primitives
     * @returns {CreateMonitorResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new CreateMonitorResultDto(primitives);
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
