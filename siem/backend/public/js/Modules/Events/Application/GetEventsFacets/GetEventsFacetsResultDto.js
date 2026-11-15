/**
 * Resultado del caso de uso: por campo (severity, event_type, source_type, source_name, agent_name, project_name),
 * sus valores mas repetidos: [{value, filter_field, filter_value, count}], de mayor a menor. value es lo que se
 * pinta; filter_field + filter_value, con lo que se filtra por ese valor (un nombre filtra por su uuid).
 */
export class GetEventsFacetsResultDto {

    #facets;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        const facets = primitives.facets;
        this.#facets = facets && typeof facets === "object" && !Array.isArray(facets) ? facets : {};
    }

    /**
     * @param {Object} primitives
     * @returns {GetEventsFacetsResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new GetEventsFacetsResultDto(primitives);
    }

    /**
     * @returns {Object} {campo: [{value, filter_field, filter_value, count}]}
     */
    getFacets() {
        return this.#facets;
    }
}
