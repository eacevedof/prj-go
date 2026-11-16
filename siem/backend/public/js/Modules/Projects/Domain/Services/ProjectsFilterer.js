import { ProjectsKeyEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsKeyEnum.js";

/**
 * Filtros de la tabla de proyectos. Se aplican en el front: el dashboard ya trae todos los proyectos (son pocos) y
 * asi filtrar no cuesta una peticion ni espera al siguiente refresco.
 */
export class ProjectsFilterer {

    /**
     * @returns {ProjectsFilterer}
     */
    static getInstance() {
        return new ProjectsFilterer();
    }

    /**
     * @param {Object[]} projects tal cual los da la API
     * @param {{projectText: string, projectHealth: string}} filters vacio = sin ese filtro
     * @returns {Object[]}
     */
    getFiltered(projects, filters) {
        const projectText = String(filters.projectText ?? "").trim().toLowerCase();
        const projectHealth = String(filters.projectHealth ?? "");
        return projects.filter((project) => {
            if (projectHealth && project[ProjectsKeyEnum.PROJECT_HEALTH] !== projectHealth) return false;
            if (!projectText) return true;
            return this.#getSearchableText(project).includes(projectText);
        });
    }

    /**
     * @param {{projectText: string, projectHealth: string}} filters
     * @returns {boolean}
     */
    hasFilters(filters) {
        return Boolean(String(filters.projectText ?? "").trim() || filters.projectHealth);
    }

    /**
     * El texto busca en id, uuid, nombre y descripcion.
     */
    #getSearchableText(project) {
        return [
            project[ProjectsKeyEnum.PROJECT_ID],
            project[ProjectsKeyEnum.PROJECT_UUID],
            project[ProjectsKeyEnum.PROJECT_NAME],
            project[ProjectsKeyEnum.PROJECT_DESCRIPTION],
        ].map((value) => String(value ?? "")).join(" ").toLowerCase();
    }
}
