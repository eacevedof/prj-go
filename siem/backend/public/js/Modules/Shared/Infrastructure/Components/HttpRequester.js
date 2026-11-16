const METHOD_GET = "GET";
const METHOD_POST = "POST";
const METHOD_PUT = "PUT";
const METHOD_DELETE = "DELETE";
const HEADER_ACCEPT = "Accept";
const HEADER_AUTHORIZATION = "Authorization";
const HEADER_CONTENT_TYPE = "Content-Type";
const CONTENT_TYPE_JSON = "application/json";

/**
 * Peticiones a la API REST del backend (la misma origin que sirve el front) con el Bearer del dashboard. Devuelve el
 * envelope de la API tal cual: {code, status, errors, data}. Una respuesta que no es JSON se normaliza a status false
 * con el codigo HTTP. No captura: un fallo de red (backend caido) propaga.
 */
export class HttpRequester {

    /**
     * @returns {HttpRequester}
     */
    static getInstance() {
        return new HttpRequester();
    }

    /**
     * @param {string} path p. ej. /api/v1/dashboard
     * @param {Object} queryParams los vacios no se envian
     * @param {string} dashboardToken
     * @returns {Promise<{code: number, status: boolean, errors: string[], data: *}>}
     */
    async getResult(path, queryParams = {}, dashboardToken = "") {
        const queryString = new URLSearchParams(
            Object.entries(queryParams).filter(([, value]) => value !== "" && value !== null && value !== undefined)
        ).toString();
        return this.#getResponse(METHOD_GET, `${path}${queryString ? `?${queryString}` : ""}`, null, dashboardToken);
    }

    /**
     * @param {string} path
     * @param {Object} body se envia como JSON
     * @param {string} dashboardToken
     */
    async postResult(path, body = {}, dashboardToken = "") {
        return this.#getResponse(METHOD_POST, path, body, dashboardToken);
    }

    /**
     * @param {string} path
     * @param {Object} body se envia como JSON
     * @param {string} dashboardToken
     */
    async putResult(path, body = {}, dashboardToken = "") {
        return this.#getResponse(METHOD_PUT, path, body, dashboardToken);
    }

    /**
     * @param {string} path
     * @param {string} dashboardToken
     */
    async deleteResult(path, dashboardToken = "") {
        return this.#getResponse(METHOD_DELETE, path, null, dashboardToken);
    }

    async #getResponse(method, url, body, dashboardToken) {
        const headers = { [HEADER_ACCEPT]: CONTENT_TYPE_JSON };
        if (dashboardToken) {
            headers[HEADER_AUTHORIZATION] = `Bearer ${dashboardToken}`;
        }
        const options = { method, headers };
        if (body !== null) {
            headers[HEADER_CONTENT_TYPE] = CONTENT_TYPE_JSON;
            options.body = JSON.stringify(body);
        }

        const response = await fetch(url, options);
        const json = await response.json().catch(() => null);
        if (!json) {
            return { code: response.status, status: false, errors: [`HTTP ${response.status}`], data: null };
        }
        return {
            code: Number(json.code ?? response.status),
            status: json.status === true,
            errors: Array.isArray(json.errors) ? json.errors : [],
            data: json.data ?? null,
        };
    }
}
