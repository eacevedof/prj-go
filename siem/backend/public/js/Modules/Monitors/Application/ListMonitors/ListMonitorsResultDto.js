/**
 * Resultado del caso de uso: los monitores del proyecto (objetos planos como los da la API), cada uno con states[]:
 * su ultimo valor por agente.
 */
export class ListMonitorsResultDto {

    #monitors;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#monitors = Array.isArray(primitives.monitors) ? primitives.monitors : [];
    }

    /**
     * @param {Object} primitives
     * @returns {ListMonitorsResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new ListMonitorsResultDto(primitives);
    }

    /**
     * @returns {Object[]}
     */
    getMonitors() {
        return this.#monitors;
    }
}
