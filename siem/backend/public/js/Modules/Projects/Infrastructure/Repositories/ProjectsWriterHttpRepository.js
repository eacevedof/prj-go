import { HttpRequester } from "/js/Modules/Shared/Infrastructure/Components/HttpRequester.js";

const PATH_PROJECTS = "/api/v1/projects";

export class ProjectsWriterHttpRepository {

    #httpRequester;

    constructor() {
        this.#httpRequester = HttpRequester.getInstance();
    }

    /**
     * @returns {ProjectsWriterHttpRepository}
     */
    static getInstance() {
        return new ProjectsWriterHttpRepository();
    }

    /**
     * POST /api/v1/projects.
     * @param {{project_name: string, project_description: string}} project
     * @param {string} dashboardToken
     */
    async createProject(project, dashboardToken) {
        return this.#httpRequester.postResult(PATH_PROJECTS, project, dashboardToken);
    }

    /**
     * DELETE /api/v1/projects/:project_uuid (borra agentes, monitores y eventos).
     * @param {string} projectUuid
     * @param {string} dashboardToken
     */
    async deleteProjectByProjectUuid(projectUuid, dashboardToken) {
        return this.#httpRequester.deleteResult(`${PATH_PROJECTS}/${encodeURIComponent(projectUuid)}`, dashboardToken);
    }
}
