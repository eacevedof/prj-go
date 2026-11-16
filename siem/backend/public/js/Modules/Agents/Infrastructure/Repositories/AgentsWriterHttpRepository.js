import { HttpRequester } from "/js/Modules/Shared/Infrastructure/Components/HttpRequester.js";

const PATH_PROJECTS = "/api/v1/projects";
const PATH_AGENTS = "/api/v1/agents";

export class AgentsWriterHttpRepository {

    #httpRequester;

    constructor() {
        this.#httpRequester = HttpRequester.getInstance();
    }

    /**
     * @returns {AgentsWriterHttpRepository}
     */
    static getInstance() {
        return new AgentsWriterHttpRepository();
    }

    /**
     * POST /api/v1/projects/:project_uuid/agents: la respuesta trae api_token (solo esta vez).
     * @param {string} projectUuid
     * @param {{agent_name: string, heartbeat_interval_seconds: number}} agent
     * @param {string} dashboardToken
     */
    async createAgentInProject(projectUuid, agent, dashboardToken) {
        return this.#httpRequester.postResult(`${PATH_PROJECTS}/${encodeURIComponent(projectUuid)}/agents`, agent, dashboardToken);
    }

    /**
     * PUT /api/v1/agents/:agent_uuid (nombre, latido y si esta activo; el token no cambia).
     * @param {string} agentUuid
     * @param {{agent_name: string, heartbeat_interval_seconds: number, is_enabled: boolean}} agent
     * @param {string} dashboardToken
     */
    async updateAgentByAgentUuid(agentUuid, agent, dashboardToken) {
        return this.#httpRequester.putResult(`${PATH_AGENTS}/${encodeURIComponent(agentUuid)}`, agent, dashboardToken);
    }

    /**
     * DELETE /api/v1/agents/:agent_uuid.
     * @param {string} agentUuid
     * @param {string} dashboardToken
     */
    async deleteAgentByAgentUuid(agentUuid, dashboardToken) {
        return this.#httpRequester.deleteResult(`${PATH_AGENTS}/${encodeURIComponent(agentUuid)}`, dashboardToken);
    }
}
