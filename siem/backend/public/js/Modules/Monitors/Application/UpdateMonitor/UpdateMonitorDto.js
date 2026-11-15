import { CreateMonitorDto } from "/js/Modules/Monitors/Application/CreateMonitor/CreateMonitorDto.js";

/**
 * Entrada del caso de uso: los mismos campos que el alta (la API pide el monitor entero en el PUT) mas el uuid del
 * monitor que se cambia. El tipo no cambia.
 */
export class UpdateMonitorDto extends CreateMonitorDto {

    #monitorUuid;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        super(primitives);
        this.#monitorUuid = String(primitives.monitorUuid ?? "").trim();
    }

    /**
     * @param {Object} primitives
     * @returns {UpdateMonitorDto}
     */
    static fromPrimitives(primitives = {}) {
        return new UpdateMonitorDto(primitives);
    }

    /**
     * @returns {string}
     */
    getMonitorUuid() {
        return this.#monitorUuid;
    }
}
