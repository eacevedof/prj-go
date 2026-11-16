/**
 * Resultado del caso de uso: el proyecto borrado y los agentes que tenia, con lo que hace falta para desinstalar su
 * siem-client a mano (el borrado no lo quita del host: solo deja de aceptarlo).
 */
export class DeleteProjectResultDto {

    #projectUuid;
    #projectName;
    #agents;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#projectUuid = String(primitives.projectUuid ?? "");
        this.#projectName = String(primitives.projectName ?? "");
        this.#agents = (Array.isArray(primitives.agents) ? primitives.agents : []).map((agent) => Object.freeze({
            agentName: String(agent.agent_name ?? ""),
            hostName: String(agent.host_name ?? ""),
            hostOs: String(agent.host_os ?? ""),
            hostIp: String(agent.host_ip ?? ""),
            clientPath: String(agent.client_path ?? ""),
        }));
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

    /**
     * @returns {string}
     */
    getProjectName() {
        return this.#projectName;
    }

    /**
     * @returns {{agentName: string, hostName: string, hostOs: string, hostIp: string, clientPath: string}[]}
     */
    getAgents() {
        return this.#agents;
    }
}
