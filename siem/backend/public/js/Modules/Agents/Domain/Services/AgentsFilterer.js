import { AgentsKeyEnum } from "/js/Modules/Agents/Domain/Enums/AgentsKeyEnum.js";

/**
 * Filtros de la tabla de agentes de un proyecto. Se aplican en el front: el detalle ya trae todos sus agentes.
 */
export class AgentsFilterer {

    /**
     * @returns {AgentsFilterer}
     */
    static getInstance() {
        return new AgentsFilterer();
    }

    /**
     * @param {Object[]} agents tal cual los da la API
     * @param {{agentText: string, agentStatus: string}} filters vacio = sin ese filtro
     * @returns {Object[]}
     */
    getFiltered(agents, filters) {
        const agentText = String(filters.agentText ?? "").trim().toLowerCase();
        const agentStatus = String(filters.agentStatus ?? "");
        return agents.filter((agent) => {
            if (agentStatus && agent[AgentsKeyEnum.AGENT_STATUS] !== agentStatus) return false;
            if (!agentText) return true;
            return this.#getSearchableText(agent).includes(agentText);
        });
    }

    /**
     * @param {{agentText: string, agentStatus: string}} filters
     * @returns {boolean}
     */
    hasFilters(filters) {
        return Boolean(String(filters.agentText ?? "").trim() || filters.agentStatus);
    }

    /**
     * El texto busca en id, uuid, nombre y lo que el agente dice de su host (nombre, sistema, ip, version).
     */
    #getSearchableText(agent) {
        return [
            agent[AgentsKeyEnum.AGENT_ID],
            agent[AgentsKeyEnum.AGENT_UUID],
            agent[AgentsKeyEnum.AGENT_NAME],
            agent[AgentsKeyEnum.HOST_NAME],
            agent[AgentsKeyEnum.HOST_OS],
            agent[AgentsKeyEnum.HOST_IP],
            agent[AgentsKeyEnum.CLIENT_VERSION],
        ].map((value) => String(value ?? "")).join(" ").toLowerCase();
    }
}
