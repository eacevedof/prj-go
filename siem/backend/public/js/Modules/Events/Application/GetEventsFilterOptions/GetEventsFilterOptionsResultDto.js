import { EventsFilterOptionKeyEnum } from "/js/Modules/Events/Domain/Enums/EventsFilterOptionKeyEnum.js";

/**
 * Resultado del caso de uso: el catalogo de los desplegables del visor, tal cual lo da la API (claves snake_case).
 *   projects: [{project_uuid, project_name}]
 *   agents:   [{agent_uuid, agent_name, project_uuid}]
 *   monitors: [{monitor_uuid, monitor_name, monitor_type, project_uuid, agent_uuid}] (agent_uuid "" = todos)
 */
export class GetEventsFilterOptionsResultDto {

    #projects;
    #agents;
    #monitors;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#projects = Array.isArray(primitives[EventsFilterOptionKeyEnum.PROJECTS]) ? primitives[EventsFilterOptionKeyEnum.PROJECTS] : [];
        this.#agents = Array.isArray(primitives[EventsFilterOptionKeyEnum.AGENTS]) ? primitives[EventsFilterOptionKeyEnum.AGENTS] : [];
        this.#monitors = Array.isArray(primitives[EventsFilterOptionKeyEnum.MONITORS]) ? primitives[EventsFilterOptionKeyEnum.MONITORS] : [];
    }

    /**
     * @param {Object} primitives
     * @returns {GetEventsFilterOptionsResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new GetEventsFilterOptionsResultDto(primitives);
    }

    /**
     * @returns {{projects: Object[], agents: Object[], monitors: Object[]}}
     */
    getFilterOptions() {
        return { projects: this.#projects, agents: this.#agents, monitors: this.#monitors };
    }
}
