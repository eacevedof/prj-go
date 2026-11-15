/**
 * Resultado del caso de uso: el estado del servidor y un resumen por proyecto (objetos planos como los da la API,
 * con los estados ya calculados en el backend).
 */
export class GetDashboardResultDto {

    #server;
    #projects;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#server = primitives.server ?? {};
        this.#projects = Array.isArray(primitives.projects) ? primitives.projects : [];
    }

    /**
     * @param {Object} primitives
     * @returns {GetDashboardResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new GetDashboardResultDto(primitives);
    }

    /**
     * @returns {Object} version, uptime, cola, tamaño de la BD... y now (UTC del servidor)
     */
    getServer() {
        return this.#server;
    }

    /**
     * @returns {Object[]} uno por proyecto con su project_health
     */
    getProjects() {
        return this.#projects;
    }
}
