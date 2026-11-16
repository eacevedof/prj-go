/**
 * Entrada del caso de uso: los proyectos del dashboard (objetos planos de la API) y los filtros de la tabla.
 */
export class FilterProjectsDto {

    #projects;
    #projectText;
    #projectHealth;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#projects = Array.isArray(primitives.projects) ? primitives.projects : [];
        this.#projectText = String(primitives.projectText ?? "").trim();
        this.#projectHealth = String(primitives.projectHealth ?? "").trim();
    }

    /**
     * @param {Object} primitives
     * @returns {FilterProjectsDto}
     */
    static fromPrimitives(primitives = {}) {
        return new FilterProjectsDto(primitives);
    }

    /**
     * @returns {Object[]}
     */
    getProjects() {
        return this.#projects;
    }

    /**
     * @returns {string}
     */
    getProjectText() {
        return this.#projectText;
    }

    /**
     * @returns {string}
     */
    getProjectHealth() {
        return this.#projectHealth;
    }
}
