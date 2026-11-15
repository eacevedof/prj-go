/**
 * Entrada del caso de uso: solo primitivos, todo string con trim.
 */
export class LogInDto {

    #dashboardToken;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
    }

    /**
     * @param {Object} primitives
     * @returns {LogInDto}
     */
    static fromPrimitives(primitives = {}) {
        return new LogInDto(primitives);
    }

    /**
     * @returns {string}
     */
    getDashboardToken() {
        return this.#dashboardToken;
    }
}
