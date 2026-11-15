import { EVENT_FILTER_FIELDS } from "/js/Modules/Events/Domain/Enums/EventFilterFieldEnum.js";

/**
 * Lo comun de la entrada de los tres casos de uso del visor (buscar, histograma, facetas): el token y los mismos
 * filtros. Solo primitivos: strings con trim y, por campo filtrable, la lista de valores a incluir y a excluir.
 */
export class EventsFiltersDto {

    #dashboardToken;
    #from;
    #to;
    #search;
    #include;
    #exclude;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#dashboardToken = String(primitives.dashboardToken ?? "").trim();
        this.#from = String(primitives.from ?? "").trim();
        this.#to = String(primitives.to ?? "").trim();
        this.#search = String(primitives.search ?? "").trim();
        this.#include = EventsFiltersDto.#getValuesByField(primitives.include);
        this.#exclude = EventsFiltersDto.#getValuesByField(primitives.exclude);
    }

    static #getValuesByField(rawValuesByField) {
        const valuesByField = {};
        EVENT_FILTER_FIELDS.forEach((field) => {
            const rawValues = rawValuesByField?.[field];
            valuesByField[field] = Array.isArray(rawValues) ? rawValues.map((value) => String(value).trim()).filter(Boolean) : [];
        });
        return valuesByField;
    }

    /**
     * @returns {string}
     */
    getDashboardToken() {
        return this.#dashboardToken;
    }

    /**
     * @returns {{from: string, to: string, search: string, include: Object, exclude: Object}}
     */
    getFilters() {
        return {
            from: this.#from,
            to: this.#to,
            search: this.#search,
            include: this.#include,
            exclude: this.#exclude,
        };
    }
}
