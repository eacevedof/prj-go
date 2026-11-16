const SECONDS_PER_MINUTE = 60;
const SECONDS_PER_HOUR = 3600;
const SECONDS_PER_DAY = 86400;
const MILLISECONDS_PER_SECOND = 1000;
const EMPTY_VALUE = "—";
const NEVER = "nunca";
const TIME_AGO_PREFIX = "hace";
const UTC_DATE_TIME_PATTERN = /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/;

/**
 * Fechas siempre en UTC y en formato "YYYY-MM-DD HH:MM:SS" (24 h): es como las guarda y devuelve la API, y asi se
 * pintan, sin pasarlas a la hora local del navegador. Los "hace X" se calculan contra el "now" del servidor: el reloj
 * del navegador puede no coincidir con el del backend.
 */
export class DateTimer {

    /**
     * @returns {DateTimer}
     */
    static getInstance() {
        return new DateTimer();
    }

    /**
     * @param {string} utcDateTime de la API
     * @returns {string} "YYYY-MM-DD HH:MM:SS" UTC; "—" si viene vacia
     */
    getUtcDateTime(utcDateTime) {
        if (!utcDateTime) return EMPTY_VALUE;
        const text = String(utcDateTime).trim().replace("T", " ").replace(/Z$/, "").slice(0, 19);
        if (UTC_DATE_TIME_PATTERN.test(text)) return text;
        const date = this.#getDate(utcDateTime);
        return Number.isNaN(date.getTime()) ? String(utcDateTime) : this.#getFormatted(date);
    }

    /**
     * @returns {string} ahora en UTC, "YYYY-MM-DD HH:MM:SS"
     */
    getUtcDateTimeNow() {
        return this.#getFormatted(new Date());
    }

    /**
     * @param {Date} date
     * @returns {string} "YYYY-MM-DD HH:MM:SS" UTC
     */
    getUtcDateTimeOf(date) {
        return this.#getFormatted(date);
    }

    /**
     * @param {string} utcDateTime "YYYY-MM-DD HH:MM:SS" UTC
     * @returns {Date} Invalid Date si no es una fecha
     */
    getDateOf(utcDateTime) {
        return this.#getDate(utcDateTime);
    }

    /**
     * @param {string} utcDateTime
     * @param {string} nowUtcDateTime el now del servidor
     * @returns {string} "hace 12 s", "hace 5 min"...; "nunca" si viene vacia
     */
    getTimeAgo(utcDateTime, nowUtcDateTime) {
        if (!utcDateTime) return NEVER;
        const seconds = Math.max(0, Math.round((this.#getDate(nowUtcDateTime) - this.#getDate(utcDateTime)) / MILLISECONDS_PER_SECOND));
        return `${TIME_AGO_PREFIX} ${this.getDuration(seconds)}`;
    }

    /**
     * @param {number} seconds
     * @returns {string} "12 s", "5 min", "3 h", "2 d"
     */
    getDuration(seconds) {
        if (seconds < SECONDS_PER_MINUTE) return `${seconds} s`;
        if (seconds < SECONDS_PER_HOUR) return `${Math.floor(seconds / SECONDS_PER_MINUTE)} min`;
        if (seconds < SECONDS_PER_DAY) return `${Math.floor(seconds / SECONDS_PER_HOUR)} h`;
        return `${Math.floor(seconds / SECONDS_PER_DAY)} d`;
    }

    #getDate(utcDateTime) {
        const text = String(utcDateTime).trim().replace(" ", "T");
        return new Date(text.endsWith("Z") ? text : `${text}Z`);
    }

    // toISOString siempre es UTC: "2026-11-15T18:01:08.123Z" -> "2026-11-15 18:01:08"
    #getFormatted(date) {
        return date.toISOString().slice(0, 19).replace("T", " ");
    }
}
