/**
 * Entrada del caso de uso: solo el token del panel (el catalogo no depende de los filtros).
 */
export class GetEventsFilterOptionsDto {

    #dashboardToken;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
    }

    /**
     * @param {Object} primitives
     * @returns {GetEventsFilterOptionsDto}
     */
    static fromPrimitives(primitives = {}) {
        return new GetEventsFilterOptionsDto(primitives);
    }

    /**
     * @returns {string}
     */
    getDashboardToken() {
        return this.#dashboardToken;
    }
}
