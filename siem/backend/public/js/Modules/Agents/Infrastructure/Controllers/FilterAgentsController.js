import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";

import { DashboardEventEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardEventEnum.js";

import { AgentsEventEnum } from "/js/Modules/Agents/Domain/Enums/AgentsEventEnum.js";
import { AgentsMessageEnum } from "/js/Modules/Agents/Domain/Enums/AgentsMessageEnum.js";
import { FilterAgentsDto } from "/js/Modules/Agents/Application/FilterAgents/FilterAgentsDto.js";
import { FilterAgentsService } from "/js/Modules/Agents/Application/FilterAgents/FilterAgentsService.js";
import { AgentsTableView } from "/js/Modules/Agents/Infrastructure/Views/AgentsTableView.js";

/**
 * Tabla de agentes del detalle de un proyecto: la pinta cada vez que se pinta el detalle (tambien en el auto-refresco)
 * y cada vez que cambia un filtro. Agents escucha a Dashboard, no al reves. Guarda los ultimos agentes recibidos
 * para filtrar sin volver a pedirlos.
 */
export class FilterAgentsController {

    #subscriber;
    #agentsTableView;
    #agents = [];
    #serverNow = "";

    constructor() {
        this.#subscriber = Subscriber.getInstance();
        this.#agentsTableView = AgentsTableView.getInstance();
    }

    /**
     * @returns {FilterAgentsController}
     */
    static getInstance() {
        return new FilterAgentsController();
    }

    init() {
        this.#subscriber.subscribe(DashboardEventEnum.PROJECT_DETAIL_RENDERED, (request) => {
            this.#agents = Array.isArray(request.agents) ? request.agents : [];
            this.#serverNow = String(request.serverNow ?? "");
            this.#agentsTableView.renderProject(String(request.project?.project_uuid ?? ""));
            this.invoke({ filters: this.#agentsTableView.getFilters() });
        });
        this.#subscriber.subscribe(AgentsEventEnum.AGENTS_FILTERS_CHANGED, (request) => this.invoke(request));
    }

    /**
     * @param {{filters: {agentText: string, agentStatus: string}}} request
     */
    invoke(request) {
        try {
            const filterAgentsResultDto = FilterAgentsService.getInstance().invoke(
                FilterAgentsDto.fromPrimitives({
                    agents: this.#agents,
                    agentText: request.filters?.agentText,
                    agentStatus: request.filters?.agentStatus,
                })
            );
            this.#agentsTableView.renderAgents(
                filterAgentsResultDto.getFilteredAgents(),
                filterAgentsResultDto.getAgentsTotal(),
                filterAgentsResultDto.hasFilters(),
                this.#serverNow,
            );
        }
        catch (error) {
            console.error("FilterAgentsController.invoke", error);
            this.#agentsTableView.renderErrors([AgentsMessageEnum.UNEXPECTED_ERROR]);
        }
    }
}
