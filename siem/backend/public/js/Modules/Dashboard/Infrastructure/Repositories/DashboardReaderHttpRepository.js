import { HttpRequester } from "/js/Modules/Shared/Infrastructure/Components/HttpRequester.js";

const PATH_DASHBOARD = "/api/v1/dashboard";
const PATH_PROJECTS = "/api/v1/projects";

export class DashboardReaderHttpRepository {

    #httpRequester;

    constructor() {
        this.#httpRequester = HttpRequester.getInstance();
    }

    /**
     * @returns {DashboardReaderHttpRepository}
     */
    static getInstance() {
        return new DashboardReaderHttpRepository();
    }

    /**
     * GET /api/v1/dashboard: {server, projects[]}.
     * @param {string} dashboardToken
     */
    async getDashboard(dashboardToken) {
        return this.#httpRequester.getResult(PATH_DASHBOARD, {}, dashboardToken);
    }

    /**
     * GET /api/v1/projects/:project_uuid: {project, agents[]} (cada agente con sus sources[]).
     * @param {string} projectUuid
     * @param {string} dashboardToken
     */
    async getProjectDetailByProjectUuid(projectUuid, dashboardToken) {
        return this.#httpRequester.getResult(`${PATH_PROJECTS}/${encodeURIComponent(projectUuid)}`, {}, dashboardToken);
    }
}
