import { HttpRequester } from "/js/Modules/Shared/Infrastructure/Components/HttpRequester.js";

const PATH_PROJECTS = "/api/v1/projects";

export class MonitorsReaderHttpRepository {

    #httpRequester;

    constructor() {
        this.#httpRequester = HttpRequester.getInstance();
    }

    /**
     * @returns {MonitorsReaderHttpRepository}
     */
    static getInstance() {
        return new MonitorsReaderHttpRepository();
    }

    /**
     * GET /api/v1/projects/:project_uuid/monitors: {monitors[]} (cada uno con states[]).
     * @param {string} projectUuid
     * @param {string} dashboardToken
     */
    async getMonitorsByProjectUuid(projectUuid, dashboardToken) {
        return this.#httpRequester.getResult(`${PATH_PROJECTS}/${encodeURIComponent(projectUuid)}/monitors`, {}, dashboardToken);
    }
}
