/**
 * Resultado del caso de uso: los binarios del agente que sirve el backend (objetos planos de la API: file_name, os,
 * arch, version, size_bytes, sha256, url) y la url de ingesta que debe llevar el .env del agente.
 */
export class GetDownloadsResultDto {

    #endpointUrl;
    #downloads;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#endpointUrl = String(primitives.endpointUrl ?? "");
        this.#downloads = Array.isArray(primitives.downloads) ? primitives.downloads : [];
    }

    /**
     * @param {Object} primitives
     * @returns {GetDownloadsResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new GetDownloadsResultDto(primitives);
    }

    /**
     * @returns {string}
     */
    getEndpointUrl() {
        return this.#endpointUrl;
    }

    /**
     * @returns {Object[]}
     */
    getDownloads() {
        return this.#downloads;
    }
}
