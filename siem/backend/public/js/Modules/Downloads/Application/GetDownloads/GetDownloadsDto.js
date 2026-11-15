/**
 * Entrada del caso de uso: solo primitivos, todo string con trim.
 */
export class GetDownloadsDto {

    #dashboardToken;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
    }

    /**
     * @param {Object} primitives
     * @returns {GetDownloadsDto}
     */
    static fromPrimitives(primitives = {}) {
        return new GetDownloadsDto(primitives);
    }

    /**
     * @returns {string}
     */
    getDashboardToken() {
        return this.#dashboardToken;
    }
}
