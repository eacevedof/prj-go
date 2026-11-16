import { EventsFilterOptionKeyEnum } from "/js/Modules/Events/Domain/Enums/EventsFilterOptionKeyEnum.js";

/**
 * Los filtros Proyecto / Agente / Monitor del visor son de seleccion multiple y van encadenados: los proyectos
 * marcados acotan los agentes que se ofrecen, y proyectos + agentes acotan los monitores (un monitor de "todos los
 * agentes" vale para cualquier agente de su proyecto). Lista vacia = sin filtro en ese nivel. Tambien dice que
 * elecciones siguen valiendo al cambiar las de arriba (las otras se quitan). Funciones puras: sin DOM.
 */
export class EventsFilterOptionsFilterer {

    /**
     * @returns {EventsFilterOptionsFilterer}
     */
    static getInstance() {
        return new EventsFilterOptionsFilterer();
    }

    /**
     * @param {Object[]} agents [{agent_uuid, agent_name, project_uuid}]
     * @param {string[]} projectUuids [] = de todos los proyectos
     * @returns {Object[]}
     */
    getAgentsOfProjects(agents, projectUuids) {
        if (projectUuids.length === 0) return agents;
        return agents.filter((agent) => projectUuids.includes(agent[EventsFilterOptionKeyEnum.PROJECT_UUID]));
    }

    /**
     * @param {Object[]} monitors [{monitor_uuid, monitor_name, monitor_type, project_uuid, agent_uuid}]
     * @param {string[]} projectUuids [] = de todos los proyectos
     * @param {string[]} agentUuids [] = de cualquier agente
     * @param {Object[]} agents para saber el proyecto de cada agente marcado
     * @returns {Object[]}
     */
    getMonitorsOfProjectsAndAgents(monitors, projectUuids, agentUuids, agents) {
        const agentProjectUuids = agents
            .filter((agent) => agentUuids.includes(agent[EventsFilterOptionKeyEnum.AGENT_UUID]))
            .map((agent) => agent[EventsFilterOptionKeyEnum.PROJECT_UUID]);
        return monitors.filter((monitor) => {
            const monitorProjectUuid = monitor[EventsFilterOptionKeyEnum.PROJECT_UUID];
            if (projectUuids.length > 0 && !projectUuids.includes(monitorProjectUuid)) return false;
            if (agentUuids.length === 0) return true;
            const monitorAgentUuid = monitor[EventsFilterOptionKeyEnum.AGENT_UUID] ?? "";
            return agentUuids.includes(monitorAgentUuid)
                || (monitorAgentUuid === "" && agentProjectUuids.includes(monitorProjectUuid));
        });
    }

    /**
     * @param {string[]} uuids los marcados
     * @param {Object[]} options las opciones que siguen ofreciendose
     * @param {string} uuidKey EventsFilterOptionKeyEnum.*_UUID
     * @returns {string[]} los marcados que siguen entre las opciones
     */
    getUuidsInOptions(uuids, options, uuidKey) {
        return uuids.filter((uuid) => options.some((option) => option[uuidKey] === uuid));
    }
}
