import { HttpRequester } from "/js/Modules/Shared/Infrastructure/Components/HttpRequester.js";

const PATH_PROJECTS = "/api/v1/projects";
const PATH_MONITORS = "/api/v1/monitors";

export class MonitorsWriterHttpRepository {

    #httpRequester;

    constructor() {
        this.#httpRequester = HttpRequester.getInstance();
    }

    /**
     * @returns {MonitorsWriterHttpRepository}
     */
    static getInstance() {
        return new MonitorsWriterHttpRepository();
    }

    /**
     * POST /api/v1/projects/:project_uuid/monitors.
     * @param {string} projectUuid
     * @param {Object} monitor claves snake_case del contrato
     * @param {string} dashboardToken
     */
    async createMonitorInProject(projectUuid, monitor, dashboardToken) {
        return this.#httpRequester.postResult(`${PATH_PROJECTS}/${encodeURIComponent(projectUuid)}/monitors`, monitor, dashboardToken);
    }

    /**
     * PUT /api/v1/monitors/:monitor_uuid (el monitor entero).
     * @param {string} monitorUuid
     * @param {Object} monitor
     * @param {string} dashboardToken
     */
    async updateMonitorByMonitorUuid(monitorUuid, monitor, dashboardToken) {
        return this.#httpRequester.putResult(`${PATH_MONITORS}/${encodeURIComponent(monitorUuid)}`, monitor, dashboardToken);
    }

    /**
     * DELETE /api/v1/monitors/:monitor_uuid.
     * @param {string} monitorUuid
     * @param {string} dashboardToken
     */
    async deleteMonitorByMonitorUuid(monitorUuid, dashboardToken) {
        return this.#httpRequester.deleteResult(`${PATH_MONITORS}/${encodeURIComponent(monitorUuid)}`, dashboardToken);
    }
}
