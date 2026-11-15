/**
 * Entrada del caso de uso: solo primitivos, todo string con trim.
 */
export class CreateProjectDto {

    #dashboardToken;
    #projectName;
    #projectDescription;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
        this.#projectName = String(primitives.projectName ?? "").trim();
        this.#projectDescription = String(primitives.projectDescription ?? "").trim();
    }

    /**
     * @param {Object} primitives
     * @returns {CreateProjectDto}
     */
    static fromPrimitives(primitives = {}) {
        return new CreateProjectDto(primitives);
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
    getProjectName() {
        return this.#projectName;
    }

    /**
     * @returns {string}
     */
    getProjectDescription() {
        return this.#projectDescription;
    }
}
