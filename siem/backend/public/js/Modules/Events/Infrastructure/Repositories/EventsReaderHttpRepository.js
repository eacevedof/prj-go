import { HttpRequester } from "/js/Modules/Shared/Infrastructure/Components/HttpRequester.js";

const PATH_EVENTS = "/api/v1/events";
const PATH_EVENTS_HISTOGRAM = "/api/v1/events/histogram";
const PATH_EVENTS_FACETS = "/api/v1/events/facets";
const PATH_EVENTS_FILTER_OPTIONS = "/api/v1/events/filter-options";

/**
 * Lecturas de eventos. Las tres comparten los mismos filtros (claves de la API: from, to, search, campo=a,b,
 * not_campo=a,b); vacios = sin filtro.
 */
export class EventsReaderHttpRepository {

    #httpRequester;

    constructor() {
        this.#httpRequester = HttpRequester.getInstance();
    }

    /**
     * @returns {EventsReaderHttpRepository}
     */
    static getInstance() {
        return new EventsReaderHttpRepository();
    }

    /**
     * GET /api/v1/events: {events[], page, page_size, has_more}, lo mas reciente primero.
     * @param {Object} filters los filtros + page
     * @param {string} dashboardToken
     */
    async getEventsPage(filters, dashboardToken) {
        return this.#httpRequester.getResult(PATH_EVENTS, filters, dashboardToken);
    }

    /**
     * GET /api/v1/events/histogram: {from, to, bucket_seconds, buckets[]}.
     * @param {Object} filters
     * @param {string} dashboardToken
     */
    async getEventsHistogram(filters, dashboardToken) {
        return this.#httpRequester.getResult(PATH_EVENTS_HISTOGRAM, filters, dashboardToken);
    }

    /**
     * GET /api/v1/events/facets: {facets: {campo: [{value, filter_field, filter_value, count}]}}.
     * @param {Object} filters
     * @param {string} dashboardToken
     */
    async getEventsFacets(filters, dashboardToken) {
        return this.#httpRequester.getResult(PATH_EVENTS_FACETS, filters, dashboardToken);
    }

    /**
     * GET /api/v1/events/filter-options: {projects[], agents[], monitors[]} (el catalogo de los desplegables).
     * @param {string} dashboardToken
     */
    async getEventsFilterOptions(dashboardToken) {
        return this.#httpRequester.getResult(PATH_EVENTS_FILTER_OPTIONS, {}, dashboardToken);
    }
}
