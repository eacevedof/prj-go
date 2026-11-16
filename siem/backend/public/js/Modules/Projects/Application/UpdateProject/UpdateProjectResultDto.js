/**
 * Resultado del caso de uso: el proyecto ya guardado.
 */
export class UpdateProjectResultDto {

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
     * @returns {UpdateProjectResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new UpdateProjectResultDto(primitives);
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
