import { HttpRequester } from "/js/Modules/Shared/Infrastructure/Components/HttpRequester.js";

const PATH_DOWNLOADS = "/api/v1/downloads";

export class DownloadsReaderHttpRepository {

    #httpRequester;

    constructor() {
        this.#httpRequester = HttpRequester.getInstance();
    }

    /**
     * @returns {DownloadsReaderHttpRepository}
     */
    static getInstance() {
        return new DownloadsReaderHttpRepository();
    }

    /**
     * GET /api/v1/downloads: {endpoint_url, downloads[]}.
     * @param {string} dashboardToken
     */
    async getDownloads(dashboardToken) {
        return this.#httpRequester.getResult(PATH_DOWNLOADS, {}, dashboardToken);
    }
}
