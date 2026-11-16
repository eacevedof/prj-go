/**
 * Entrada del caso de uso: los agentes del detalle del proyecto (objetos planos de la API) y los filtros de la tabla.
 */
export class FilterAgentsDto {

    #agents;
    #agentText;
    #agentStatus;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#agents = Array.isArray(primitives.agents) ? primitives.agents : [];
        this.#agentText = String(primitives.agentText ?? "").trim();
        this.#agentStatus = String(primitives.agentStatus ?? "").trim();
    }

    /**
     * @param {Object} primitives
     * @returns {FilterAgentsDto}
     */
    static fromPrimitives(primitives = {}) {
        return new FilterAgentsDto(primitives);
    }

    /**
     * @returns {Object[]}
     */
    getAgents() {
        return this.#agents;
    }

    /**
     * @returns {string}
     */
    getAgentText() {
        return this.#agentText;
    }

    /**
     * @returns {string}
     */
    getAgentStatus() {
        return this.#agentStatus;
    }
}
