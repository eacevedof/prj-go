const DEFAULT_HEARTBEAT_INTERVAL_SECONDS = 30;

/**
 * Entrada del caso de uso: solo primitivos, strings con trim. Sin intervalo de latido, el de por defecto (30 s).
 */
export class CreateAgentDto {

    #dashboardToken;
    #projectUuid;
    #agentName;
    #heartbeatIntervalSeconds;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
        this.#projectUuid = String(primitives.projectUuid ?? "").trim();
        this.#agentName = String(primitives.agentName ?? "").trim();
        const rawHeartbeat = String(primitives.heartbeatIntervalSeconds ?? "").trim();
        this.#heartbeatIntervalSeconds = rawHeartbeat === "" ? DEFAULT_HEARTBEAT_INTERVAL_SECONDS : Number(rawHeartbeat);
    }

    /**
     * @param {Object} primitives
     * @returns {CreateAgentDto}
     */
    static fromPrimitives(primitives = {}) {
        return new CreateAgentDto(primitives);
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
    getProjectUuid() {
        return this.#projectUuid;
    }

    /**
     * @returns {string}
     */
    getAgentName() {
        return this.#agentName;
    }

    /**
     * @returns {number} NaN si no era un numero
     */
    getHeartbeatIntervalSeconds() {
        return this.#heartbeatIntervalSeconds;
    }
}
