/**
 * Resultado del caso de uso: los tramos de tiempo del rango (los vacios incluidos, con ceros) con cuantos eventos hay
 * de cada gravedad. bucketSeconds = lo que dura cada tramo.
 */
export class GetEventsHistogramResultDto {

    #from;
    #to;
    #bucketSeconds;
    #buckets;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#from = String(primitives.from ?? "");
        this.#to = String(primitives.to ?? "");
        this.#bucketSeconds = Number(primitives.bucketSeconds) || 0;
        this.#buckets = Array.isArray(primitives.buckets) ? primitives.buckets : [];
    }

    /**
     * @param {Object} primitives
     * @returns {GetEventsHistogramResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new GetEventsHistogramResultDto(primitives);
    }

    /**
     * @returns {string}
     */
    getFrom() {
        return this.#from;
    }

    /**
     * @returns {string}
     */
    getTo() {
        return this.#to;
    }

    /**
     * @returns {number}
     */
    getBucketSeconds() {
        return this.#bucketSeconds;
    }

    /**
     * @returns {{bucket_at: string, total: number, debug: number, info: number, warning: number, error: number, critical: number}[]}
     */
    getBuckets() {
        return this.#buckets;
    }
}
