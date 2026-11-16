/**
 * Resultado del caso de uso: los proyectos que cumplen los filtros y cuantos hay en total.
 */
export class FilterProjectsResultDto {

    #filteredProjects;
    #projectsTotal;
    #hasFilters;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#filteredProjects = Array.isArray(primitives.filteredProjects) ? primitives.filteredProjects : [];
        this.#projectsTotal = Number(primitives.projectsTotal) || 0;
        this.#hasFilters = primitives.hasFilters === true;
    }

    /**
     * @param {Object} primitives
     * @returns {FilterProjectsResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new FilterProjectsResultDto(primitives);
    }

    /**
     * @returns {Object[]}
     */
    getFilteredProjects() {
        return this.#filteredProjects;
    }

    /**
     * @returns {number}
     */
    getProjectsTotal() {
        return this.#projectsTotal;
    }

    /**
     * @returns {boolean}
     */
    hasFilters() {
        return this.#hasFilters;
    }
}
