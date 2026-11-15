/**
 * Entrada del caso de uso: solo primitivos, todo string con trim.
 */
export class DeleteMonitorDto {

    #dashboardToken;
    #monitorUuid;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
        this.#monitorUuid = String(primitives.monitorUuid ?? "").trim();
    }

    /**
     * @param {Object} primitives
     * @returns {DeleteMonitorDto}
     */
    static fromPrimitives(primitives = {}) {
        return new DeleteMonitorDto(primitives);
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
    getMonitorUuid() {
        return this.#monitorUuid;
    }
}
