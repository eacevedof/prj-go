/**
 * Resultado del caso de uso: los agentes que cumplen los filtros y cuantos hay en total.
 */
export class FilterAgentsResultDto {

    #filteredAgents;
    #agentsTotal;
    #hasFilters;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#filteredAgents = Array.isArray(primitives.filteredAgents) ? primitives.filteredAgents : [];
        this.#agentsTotal = Number(primitives.agentsTotal) || 0;
        this.#hasFilters = primitives.hasFilters === true;
    }

    /**
     * @param {Object} primitives
     * @returns {FilterAgentsResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new FilterAgentsResultDto(primitives);
    }

    /**
     * @returns {Object[]}
     */
    getFilteredAgents() {
        return this.#filteredAgents;
    }

    /**
     * @returns {number}
     */
    getAgentsTotal() {
        return this.#agentsTotal;
    }

    /**
     * @returns {boolean}
     */
    hasFilters() {
        return this.#hasFilters;
    }
}
