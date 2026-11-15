/**
 * Resultado del caso de uso: el proyecto recien creado.
 */
export class CreateProjectResultDto {

    #projectUuid;
    #projectName;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#projectUuid = String(primitives.projectUuid ?? "");
        this.#projectName = String(primitives.projectName ?? "");
    }

    /**
     * @param {Object} primitives
     * @returns {CreateProjectResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new CreateProjectResultDto(primitives);
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
}
