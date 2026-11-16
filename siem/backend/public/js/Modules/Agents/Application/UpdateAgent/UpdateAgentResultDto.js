/**
 * Resultado del caso de uso: el agente ya guardado.
 */
export class UpdateAgentResultDto {

    #agentUuid;
    #agentName;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#agentUuid = String(primitives.agentUuid ?? "");
        this.#agentName = String(primitives.agentName ?? "");
    }

    /**
     * @param {Object} primitives
     * @returns {UpdateAgentResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new UpdateAgentResultDto(primitives);
    }

    /**
     * @returns {string}
     */
    getAgentUuid() {
        return this.#agentUuid;
    }

    /**
     * @returns {string}
     */
    getAgentName() {
        return this.#agentName;
    }
}
