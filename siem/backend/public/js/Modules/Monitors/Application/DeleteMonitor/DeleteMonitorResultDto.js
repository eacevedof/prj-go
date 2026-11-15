/**
 * Resultado del caso de uso: el uuid del monitor borrado.
 */
export class DeleteMonitorResultDto {

    #monitorUuid;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#monitorUuid = String(primitives.monitorUuid ?? "");
    }

    /**
     * @param {Object} primitives
     * @returns {DeleteMonitorResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new DeleteMonitorResultDto(primitives);
    }

    /**
     * @returns {string}
     */
    getMonitorUuid() {
        return this.#monitorUuid;
    }
}
