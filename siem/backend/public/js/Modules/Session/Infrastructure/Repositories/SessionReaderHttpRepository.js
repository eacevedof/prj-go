import { HttpRequester } from "/js/Modules/Shared/Infrastructure/Components/HttpRequester.js";

// la API no tiene endpoint de login: la lectura autenticada mas ligera hace de comprobacion del token
const PATH_DASHBOARD = "/api/v1/dashboard";

export class SessionReaderHttpRepository {

    #httpRequester;

    constructor() {
        this.#httpRequester = HttpRequester.getInstance();
    }

    /**
     * @returns {SessionReaderHttpRepository}
     */
    static getInstance() {
        return new SessionReaderHttpRepository();
    }

    /**
     * GET /api/v1/dashboard con el token: status true = el token vale; 401 = no.
     * @param {string} dashboardToken
     */
    async getDashboardTokenCheck(dashboardToken) {
        return this.#httpRequester.getResult(PATH_DASHBOARD, {}, dashboardToken);
    }
}
