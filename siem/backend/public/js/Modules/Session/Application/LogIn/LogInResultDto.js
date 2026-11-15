/**
 * Resultado del caso de uso: el token ya comprobado contra la API.
 */
export class LogInResultDto {

    #dashboardToken;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "");
    }

    /**
     * @param {Object} primitives
     * @returns {LogInResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new LogInResultDto(primitives);
    }

    /**
     * @returns {string}
     */
    getDashboardToken() {
        return this.#dashboardToken;
    }
}
