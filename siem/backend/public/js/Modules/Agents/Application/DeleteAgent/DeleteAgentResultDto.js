/**
 * Resultado del caso de uso: el uuid del agente borrado.
 */
export class DeleteAgentResultDto {

    #agentUuid;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#agentUuid = String(primitives.agentUuid ?? "");
    }

    /**
     * @param {Object} primitives
     * @returns {DeleteAgentResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new DeleteAgentResultDto(primitives);
    }

    /**
     * @returns {string}
     */
    getAgentUuid() {
        return this.#agentUuid;
    }
}
