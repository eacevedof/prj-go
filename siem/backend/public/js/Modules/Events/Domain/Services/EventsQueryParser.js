import { RouteHashEnum } from "/js/Modules/Shared/Domain/Enums/RouteHashEnum.js";

import { EVENT_FILTER_FIELDS } from "/js/Modules/Events/Domain/Enums/EventFilterFieldEnum.js";
import { EventsKeyEnum } from "/js/Modules/Events/Domain/Enums/EventsKeyEnum.js";
import { EventsPaginationEnum } from "/js/Modules/Events/Domain/Enums/EventsPaginationEnum.js";
import { EVENTS_REFRESH_VALUES, EventsRefreshEnum } from "/js/Modules/Events/Domain/Enums/EventsRefreshEnum.js";
import { TIME_RANGE_SECONDS, TimeRangeEnum } from "/js/Modules/Events/Domain/Enums/TimeRangeEnum.js";

const UTC_DATE_TIME_PATTERN = /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/;
const VALUES_SEPARATOR = ",";

/**
 * El estado del visor de eventos vive entero en el hash (#/events?range=24h&severity=error,critical&not_source_name=x)
 * para poder guardar o compartir la url. Este servicio lo traduce en los dos sentidos y lo transforma (añadir, quitar
 * o invertir un filtro) sin tocar el original. Estado:
 *   {search, range, from, to, include: {campo: [valores]}, exclude: {campo: [valores]}, page, refresh}
 * range y from/to son excluyentes: con un atajo (range) from y to van vacios; con rango a medida, range va vacio.
 * Funciones puras: sin DOM ni reloj.
 */
export class EventsQueryParser {

    /**
     * @returns {EventsQueryParser}
     */
    static getInstance() {
        return new EventsQueryParser();
    }

    /**
     * @param {Object} query los parametros del hash
     * @returns {Object} estado completo, con los valores por defecto (ultimas 24 h, pagina 1, refresco 10 s)
     */
    getState(query = {}) {
        const state = {
            search: String(query[EventsKeyEnum.SEARCH] ?? "").trim(),
            range: TimeRangeEnum.DEFAULT,
            from: "",
            to: "",
            include: {},
            exclude: {},
            page: Math.max(EventsPaginationEnum.FIRST_PAGE, Math.trunc(Number(query[EventsKeyEnum.PAGE])) || EventsPaginationEnum.FIRST_PAGE),
            refresh: EVENTS_REFRESH_VALUES.includes(query[EventsKeyEnum.REFRESH]) ? query[EventsKeyEnum.REFRESH] : EventsRefreshEnum.DEFAULT,
        };

        const range = String(query[EventsKeyEnum.RANGE] ?? "");
        const from = String(query[EventsKeyEnum.FROM] ?? "").trim();
        const to = String(query[EventsKeyEnum.TO] ?? "").trim();
        if (TIME_RANGE_SECONDS[range]) {
            state.range = range;
        }
        else if (UTC_DATE_TIME_PATTERN.test(from) && UTC_DATE_TIME_PATTERN.test(to) && from <= to) {
            state.range = "";
            state.from = from;
            state.to = to;
        }

        EVENT_FILTER_FIELDS.forEach((field) => {
            state.include[field] = this.#getValues(query[field]);
            // un valor no puede estar incluido y excluido a la vez: gana incluir
            state.exclude[field] = this.#getValues(query[`${EventsKeyEnum.NOT_PREFIX}${field}`])
                .filter((value) => !state.include[field].includes(value));
        });
        return state;
    }

    /**
     * @param {Object} state
     * @returns {Object} los parametros del hash, sin los que valen su valor por defecto
     */
    getQuery(state) {
        const query = {};
        if (state.search) query[EventsKeyEnum.SEARCH] = state.search;
        if (state.range && state.range !== TimeRangeEnum.DEFAULT) query[EventsKeyEnum.RANGE] = state.range;
        if (!state.range) {
            query[EventsKeyEnum.FROM] = state.from;
            query[EventsKeyEnum.TO] = state.to;
        }
        EVENT_FILTER_FIELDS.forEach((field) => {
            if (state.include[field]?.length) query[field] = state.include[field].join(VALUES_SEPARATOR);
            if (state.exclude[field]?.length) query[`${EventsKeyEnum.NOT_PREFIX}${field}`] = state.exclude[field].join(VALUES_SEPARATOR);
        });
        if (state.page > EventsPaginationEnum.FIRST_PAGE) query[EventsKeyEnum.PAGE] = String(state.page);
        if (state.refresh !== EventsRefreshEnum.DEFAULT) query[EventsKeyEnum.REFRESH] = state.refresh;
        return query;
    }

    /**
     * @param {Object} state
     * @returns {string} "#/events?..." ("#/events" si todo es por defecto)
     */
    getHash(state) {
        const queryString = new URLSearchParams(this.getQuery(state)).toString();
        return queryString ? `${RouteHashEnum.EVENTS}?${queryString}` : RouteHashEnum.EVENTS;
    }

    /**
     * Los filtros con el contrato de la API (comunes a events, histogram y facets).
     * @param {{from: string, to: string, search: string, include: Object, exclude: Object}} filters
     * @returns {Object} {from, to, search, campo: "a,b", not_campo: "a,b"}; los vacios no van
     */
    getApiFilters(filters) {
        const apiFilters = {};
        if (filters.from) apiFilters[EventsKeyEnum.FROM] = filters.from;
        if (filters.to) apiFilters[EventsKeyEnum.TO] = filters.to;
        if (filters.search) apiFilters[EventsKeyEnum.SEARCH] = filters.search;
        EVENT_FILTER_FIELDS.forEach((field) => {
            const includedValues = filters.include?.[field] ?? [];
            const excludedValues = filters.exclude?.[field] ?? [];
            if (includedValues.length) apiFilters[field] = includedValues.join(VALUES_SEPARATOR);
            if (excludedValues.length) apiFilters[`${EventsKeyEnum.NOT_PREFIX}${field}`] = excludedValues.join(VALUES_SEPARATOR);
        });
        return apiFilters;
    }

    /**
     * Las pastillas de filtro activas, en el orden de los campos: primero las de incluir, luego las de excluir.
     * @param {Object} state
     * @returns {{field: string, value: string, isExcluded: boolean}[]}
     */
    getPills(state) {
        const pills = [];
        EVENT_FILTER_FIELDS.forEach((field) => {
            (state.include[field] ?? []).forEach((value) => pills.push({ field, value, isExcluded: false }));
            (state.exclude[field] ?? []).forEach((value) => pills.push({ field, value, isExcluded: true }));
        });
        return pills;
    }

    /**
     * Incluye (o excluye) un valor. Si ya estaba en el lado contrario, cambia de lado. Vuelve a la primera pagina.
     * @param {Object} state
     * @param {string} field
     * @param {string} value
     * @param {boolean} isExcluded
     * @returns {Object} estado nuevo
     */
    getStateWithFilter(state, field, value, isExcluded) {
        const nextState = this.getStateWithoutFilter(state, field, value);
        if (!EVENT_FILTER_FIELDS.includes(field) || String(value) === "") return nextState;
        const side = isExcluded ? nextState.exclude : nextState.include;
        side[field].push(String(value));
        return nextState;
    }

    /**
     * @param {Object} state
     * @param {string} field
     * @param {string} value
     * @returns {Object} estado nuevo, sin ese valor en incluir ni en excluir, y en la primera pagina
     */
    getStateWithoutFilter(state, field, value) {
        const nextState = this.#getCopy(state);
        nextState.page = EventsPaginationEnum.FIRST_PAGE;
        if (!EVENT_FILTER_FIELDS.includes(field)) return nextState;
        nextState.include[field] = nextState.include[field].filter((included) => included !== String(value));
        nextState.exclude[field] = nextState.exclude[field].filter((excluded) => excluded !== String(value));
        return nextState;
    }

    /**
     * Lo marcado en los filtros multiples Proyecto / Agente / Monitor: cada campo pasa a incluir exactamente esos
     * valores ([] = sin filtro en ese campo) y deja de excluirlos. Los campos que no vienen no se tocan. Vuelve a la
     * primera pagina.
     * @param {Object} state
     * @param {Object} selectedValuesByField {campo: ["uuid", ...]}
     * @returns {Object} estado nuevo
     */
    getStateWithSelectedValues(state, selectedValuesByField) {
        const nextState = this.#getCopy(state);
        nextState.page = EventsPaginationEnum.FIRST_PAGE;
        Object.entries(selectedValuesByField).forEach(([field, selectedValues]) => {
            if (!EVENT_FILTER_FIELDS.includes(field)) return;
            const values = [...new Set((selectedValues ?? []).map(String).filter(Boolean))];
            nextState.include[field] = values;
            nextState.exclude[field] = nextState.exclude[field].filter((excluded) => !values.includes(excluded));
        });
        return nextState;
    }

    /**
     * @param {Object} state
     * @returns {Object} estado nuevo sin ningun filtro de campo (el texto y el rango se quedan)
     */
    getStateWithoutFilters(state) {
        const nextState = this.#getCopy(state);
        nextState.page = EventsPaginationEnum.FIRST_PAGE;
        EVENT_FILTER_FIELDS.forEach((field) => {
            nextState.include[field] = [];
            nextState.exclude[field] = [];
        });
        return nextState;
    }

    /**
     * Cambios sueltos (search, range, from/to, page, refresh). Cambiar algo que no sea la pagina vuelve a la primera.
     * @param {Object} state
     * @param {Object} changes
     * @returns {Object} estado nuevo
     */
    getStateWith(state, changes) {
        const nextState = { ...this.#getCopy(state), ...changes };
        if (!("page" in changes) && !("refresh" in changes)) {
            nextState.page = EventsPaginationEnum.FIRST_PAGE;
        }
        return nextState;
    }

    /**
     * @param {string} utcDateTime
     * @returns {boolean} "YYYY-MM-DD HH:MM:SS"
     */
    isUtcDateTime(utcDateTime) {
        return UTC_DATE_TIME_PATTERN.test(String(utcDateTime ?? "").trim());
    }

    #getValues(rawValues) {
        const values = String(rawValues ?? "").split(VALUES_SEPARATOR).map((value) => value.trim()).filter(Boolean);
        return [...new Set(values)];
    }

    #getCopy(state) {
        const copy = { ...state, include: {}, exclude: {} };
        EVENT_FILTER_FIELDS.forEach((field) => {
            copy.include[field] = [...(state.include?.[field] ?? [])];
            copy.exclude[field] = [...(state.exclude?.[field] ?? [])];
        });
        return copy;
    }
}
