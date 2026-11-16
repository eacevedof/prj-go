/**
 * Entrada del caso de uso: solo primitivos, todo string con trim.
 */
export class UpdateProjectDto {

    #dashboardToken;
    #projectUuid;
    #projectName;
    #projectDescription;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
        this.#projectUuid = String(primitives.projectUuid ?? "").trim();
        this.#projectName = String(primitives.projectName ?? "").trim();
        this.#projectDescription = String(primitives.projectDescription ?? "").trim();
    }

    /**
     * @param {Object} primitives
     * @returns {UpdateProjectDto}
     */
    static fromPrimitives(primitives = {}) {
        return new UpdateProjectDto(primitives);
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
    getProjectUuid() {
        return this.#projectUuid;
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
