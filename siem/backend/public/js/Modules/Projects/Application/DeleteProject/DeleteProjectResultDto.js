/**
 * Resultado del caso de uso: el uuid del proyecto borrado.
 */
export class DeleteProjectResultDto {

    #projectUuid;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#projectUuid = String(primitives.projectUuid ?? "");
    }

    /**
     * @param {Object} primitives
     * @returns {DeleteProjectResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new DeleteProjectResultDto(primitives);
    }

    /**
     * @returns {string}
     */
    getProjectUuid() {
        return this.#projectUuid;
    }
}
