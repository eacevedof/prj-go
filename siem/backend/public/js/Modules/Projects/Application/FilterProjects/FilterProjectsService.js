import { ProjectsFilterer } from "/js/Modules/Projects/Domain/Services/ProjectsFilterer.js";
import { FilterProjectsResultDto } from "/js/Modules/Projects/Application/FilterProjects/FilterProjectsResultDto.js";

/**
 * Filtra la tabla de proyectos de la portada. Sin API: el dashboard ya los trae todos.
 */
export class FilterProjectsService {

    #projectsFilterer;

    constructor() {
        this.#projectsFilterer = ProjectsFilterer.getInstance();
    }

    /**
     * @returns {FilterProjectsService}
     */
    static getInstance() {
        return new FilterProjectsService();
    }

    /**
     * @param {import("./FilterProjectsDto.js").FilterProjectsDto} filterProjectsDto
     * @returns {FilterProjectsResultDto}
     */
    invoke(filterProjectsDto) {
        const filters = {
            projectText: filterProjectsDto.getProjectText(),
            projectHealth: filterProjectsDto.getProjectHealth(),
        };
        return FilterProjectsResultDto.fromPrimitives({
            filteredProjects: this.#projectsFilterer.getFiltered(filterProjectsDto.getProjects(), filters),
            projectsTotal: filterProjectsDto.getProjects().length,
            hasFilters: this.#projectsFilterer.hasFilters(filters),
        });
    }
}
