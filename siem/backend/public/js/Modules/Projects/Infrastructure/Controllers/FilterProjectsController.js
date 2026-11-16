import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";

import { DashboardEventEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardEventEnum.js";

import { ProjectsEventEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsEventEnum.js";
import { ProjectsMessageEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsMessageEnum.js";
import { FilterProjectsDto } from "/js/Modules/Projects/Application/FilterProjects/FilterProjectsDto.js";
import { FilterProjectsService } from "/js/Modules/Projects/Application/FilterProjects/FilterProjectsService.js";
import { ProjectsTableView } from "/js/Modules/Projects/Infrastructure/Views/ProjectsTableView.js";

/**
 * Tabla de proyectos de la portada: la pinta cada vez que se pinta el dashboard (tambien en el auto-refresco) y cada
 * vez que cambia un filtro. Projects escucha a Dashboard, no al reves. Guarda los ultimos proyectos recibidos para
 * poder filtrar sin volver a pedirlos.
 */
export class FilterProjectsController {

    #subscriber;
    #projectsTableView;
    #projects = [];

    constructor() {
        this.#subscriber = Subscriber.getInstance();
        this.#projectsTableView = ProjectsTableView.getInstance();
    }

    /**
     * @returns {FilterProjectsController}
     */
    static getInstance() {
        return new FilterProjectsController();
    }

    init() {
        this.#subscriber.subscribe(DashboardEventEnum.DASHBOARD_RENDERED, (request) => {
            this.#projects = Array.isArray(request.projects) ? request.projects : [];
            this.invoke({ filters: this.#projectsTableView.getFilters() });
        });
        this.#subscriber.subscribe(ProjectsEventEnum.PROJECTS_FILTERS_CHANGED, (request) => this.invoke(request));
    }

    /**
     * @param {{filters: {projectText: string, projectHealth: string}}} request
     */
    invoke(request) {
        try {
            const filterProjectsResultDto = FilterProjectsService.getInstance().invoke(
                FilterProjectsDto.fromPrimitives({
                    projects: this.#projects,
                    projectText: request.filters?.projectText,
                    projectHealth: request.filters?.projectHealth,
                })
            );
            this.#projectsTableView.renderProjects(
                filterProjectsResultDto.getFilteredProjects(),
                filterProjectsResultDto.getProjectsTotal(),
                filterProjectsResultDto.hasFilters(),
            );
        }
        catch (error) {
            console.error("FilterProjectsController.invoke", error);
            this.#projectsTableView.renderErrors([ProjectsMessageEnum.UNEXPECTED_ERROR]);
        }
    }
}
