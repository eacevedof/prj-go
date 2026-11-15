/**
 * Resultado del caso de uso: el proyecto, sus agentes (cada uno con sus fuentes) y el now del servidor para los
 * "hace X".
 */
export class GetProjectDetailResultDto {

    #project;
    #agents;
    #serverNow;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#project = primitives.project ?? {};
        this.#agents = Array.isArray(primitives.agents) ? primitives.agents : [];
        this.#serverNow = String(primitives.serverNow ?? "");
    }

    /**
     * @param {Object} primitives
     * @returns {GetProjectDetailResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new GetProjectDetailResultDto(primitives);
    }

    /**
     * @returns {Object}
     */
    getProject() {
        return this.#project;
    }

    /**
     * @returns {Object[]} con agent_status y sources[]
     */
    getAgents() {
        return this.#agents;
    }

    /**
     * @returns {string} UTC "YYYY-MM-DD HH:MM:SS"
     */
    getServerNow() {
        return this.#serverNow;
    }
}
