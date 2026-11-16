import { TIME_RANGE_SECONDS } from "/js/Modules/Events/Domain/Enums/TimeRangeEnum.js";

const MILLISECONDS_PER_SECOND = 1000;
// "YYYY-MM-DDTHH:MM:SS" de toISOString(), sin milisegundos ni Z
const ISO_DATE_TIME_LENGTH = 19;

/**
 * Del rango del visor (un atajo relativo o un rango a medida) al from/to UTC que entiende la API. El reloj se le pasa:
 * funciones puras. Fechas "YYYY-MM-DD HH:MM:SS" en UTC.
 */
export class TimeRangeResolver {

    /**
     * @returns {TimeRangeResolver}
     */
    static getInstance() {
        return new TimeRangeResolver();
    }

    /**
     * @param {{range: string, from: string, to: string}} state
     * @param {Date} now
     * @returns {{from: string, to: string}} con un atajo: [now - rango, now]; a medida: tal cual
     */
    getFromTo(state, now) {
        const rangeSeconds = TIME_RANGE_SECONDS[state.range];
        if (!rangeSeconds) {
            return { from: state.from, to: state.to };
        }
        return {
            from: this.#getFormatted(new Date(now.getTime() - rangeSeconds * MILLISECONDS_PER_SECOND)),
            to: this.#getFormatted(now),
        };
    }

    /**
     * El tramo de una barra del histograma como rango a medida: [inicio, inicio + bucket - 1 s] (la API incluye los
     * dos extremos, asi que el ultimo segundo es del tramo siguiente).
     * @param {string} bucketAt inicio del tramo, UTC
     * @param {number} bucketSeconds
     * @returns {{from: string, to: string}}
     */
    getBucketFromTo(bucketAt, bucketSeconds) {
        const start = this.#getDate(bucketAt);
        const seconds = Math.max(1, Number(bucketSeconds) || 1);
        return {
            from: this.#getFormatted(start),
            to: this.#getFormatted(new Date(start.getTime() + (seconds - 1) * MILLISECONDS_PER_SECOND)),
        };
    }

    #getDate(utcDateTime) {
        return new Date(`${String(utcDateTime).trim().replace(" ", "T")}Z`);
    }

    #getFormatted(date) {
        return date.toISOString().slice(0, ISO_DATE_TIME_LENGTH).replace("T", " ");
    }
}
