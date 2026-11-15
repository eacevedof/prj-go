/**
 * Entrada del caso de uso: solo primitivos, todo string con trim.
 */
export class DeleteAgentDto {

    #dashboardToken;
    #agentUuid;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
        this.#agentUuid = String(primitives.agentUuid ?? "").trim();
    }

    /**
     * @param {Object} primitives
     * @returns {DeleteAgentDto}
     */
    static fromPrimitives(primitives = {}) {
        return new DeleteAgentDto(primitives);
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
    getAgentUuid() {
        return this.#agentUuid;
    }
}
