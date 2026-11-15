/**
 * Resultado del caso de uso: el agente creado y su api_token, que la API solo devuelve esta vez.
 */
export class CreateAgentResultDto {

    #agentUuid;
    #agentName;
    #apiToken;
    #endpointUrl;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#agentUuid = String(primitives.agentUuid ?? "");
        this.#agentName = String(primitives.agentName ?? "");
        this.#apiToken = String(primitives.apiToken ?? "");
        this.#endpointUrl = String(primitives.endpointUrl ?? "");
    }

    /**
     * @param {Object} primitives
     * @returns {CreateAgentResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new CreateAgentResultDto(primitives);
    }

    /**
     * @returns {string}
     */
    getAgentUuid() {
        return this.#agentUuid;
    }

    /**
     * @returns {string}
     */
    getAgentName() {
        return this.#agentName;
    }

    /**
     * @returns {string}
     */
    getApiToken() {
        return this.#apiToken;
    }

    /**
     * @returns {string} la url de ingesta segun el backend (APP_BASE_URL): la que va en el .env del agente
     */
    getEndpointUrl() {
        return this.#endpointUrl;
    }
}
