import { AgentsFilterer } from "/js/Modules/Agents/Domain/Services/AgentsFilterer.js";
import { FilterAgentsResultDto } from "/js/Modules/Agents/Application/FilterAgents/FilterAgentsResultDto.js";

/**
 * Filtra la tabla de agentes del detalle de un proyecto. Sin API: el detalle ya los trae todos.
 */
export class FilterAgentsService {

    #agentsFilterer;

    constructor() {
        this.#agentsFilterer = AgentsFilterer.getInstance();
    }

    /**
     * @returns {FilterAgentsService}
     */
    static getInstance() {
        return new FilterAgentsService();
    }

    /**
     * @param {import("./FilterAgentsDto.js").FilterAgentsDto} filterAgentsDto
     * @returns {FilterAgentsResultDto}
     */
    invoke(filterAgentsDto) {
        const filters = {
            agentText: filterAgentsDto.getAgentText(),
            agentStatus: filterAgentsDto.getAgentStatus(),
        };
        return FilterAgentsResultDto.fromPrimitives({
            filteredAgents: this.#agentsFilterer.getFiltered(filterAgentsDto.getAgents(), filters),
            agentsTotal: filterAgentsDto.getAgents().length,
            hasFilters: this.#agentsFilterer.hasFilters(filters),
        });
    }
}
