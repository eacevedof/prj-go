const DEFAULT_HEARTBEAT_INTERVAL_SECONDS = 30;

/**
 * Entrada del caso de uso: solo primitivos, strings con trim. heartbeat vacio = el de por defecto.
 */
export class UpdateAgentDto {

    #dashboardToken;
    #agentUuid;
    #agentName;
    #heartbeatIntervalSeconds;
    #isEnabled;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
        this.#agentUuid = String(primitives.agentUuid ?? "").trim();
        this.#agentName = String(primitives.agentName ?? "").trim();
        const rawHeartbeat = String(primitives.heartbeatIntervalSeconds ?? "").trim();
        this.#heartbeatIntervalSeconds = rawHeartbeat === "" ? DEFAULT_HEARTBEAT_INTERVAL_SECONDS : Number(rawHeartbeat);
        this.#isEnabled = primitives.isEnabled !== false;
    }

    /**
     * @param {Object} primitives
     * @returns {UpdateAgentDto}
     */
    static fromPrimitives(primitives = {}) {
        return new UpdateAgentDto(primitives);
    }

    /**
     * @returns {string}
     */
    getDashboardToken() {
        return this.#dashboardToken;
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
     * @returns {number}
     */
    getHeartbeatIntervalSeconds() {
        return this.#heartbeatIntervalSeconds;
    }

    /**
     * @returns {boolean}
     */
    isEnabled() {
        return this.#isEnabled;
    }
}
