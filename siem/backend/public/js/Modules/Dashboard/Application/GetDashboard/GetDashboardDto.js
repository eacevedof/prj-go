/**
 * Entrada del caso de uso: solo primitivos, todo string con trim.
 */
export class GetDashboardDto {

    #dashboardToken;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
    }

    /**
     * @param {Object} primitives
     * @returns {GetDashboardDto}
     */
    static fromPrimitives(primitives = {}) {
        return new GetDashboardDto(primitives);
    }

    /**
     * @returns {string}
     */
    getDashboardToken() {
        return this.#dashboardToken;
    }
}
