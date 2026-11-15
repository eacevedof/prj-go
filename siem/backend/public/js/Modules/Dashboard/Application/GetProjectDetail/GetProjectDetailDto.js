/**
 * Entrada del caso de uso: solo primitivos, todo string con trim.
 */
export class GetProjectDetailDto {

    #dashboardToken;
    #projectUuid;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
        this.#projectUuid = String(primitives.projectUuid ?? "").trim();
    }

    /**
     * @param {Object} primitives
     * @returns {GetProjectDetailDto}
     */
    static fromPrimitives(primitives = {}) {
        return new GetProjectDetailDto(primitives);
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
}
