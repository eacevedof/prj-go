import { EventsPaginationEnum } from "/js/Modules/Events/Domain/Enums/EventsPaginationEnum.js";

/**
 * Resultado del caso de uso: una pagina de eventos (lo mas reciente primero, objetos planos como los da la API) y
 * si hay mas. La API no cuenta el total (seria un COUNT sobre toda la tabla): solo dice si hay pagina siguiente.
 */
export class SearchEventsResultDto {

    #events;
    #page;
    #pageSize;
    #hasMore;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        this.#events = Array.isArray(primitives.events) ? primitives.events : [];
        this.#page = Number(primitives.page) || EventsPaginationEnum.FIRST_PAGE;
        this.#pageSize = Number(primitives.pageSize) || EventsPaginationEnum.PAGE_SIZE;
        this.#hasMore = primitives.hasMore === true;
    }

    /**
     * @param {Object} primitives
     * @returns {SearchEventsResultDto}
     */
    static fromPrimitives(primitives = {}) {
        return new SearchEventsResultDto(primitives);
    }

    /**
     * @returns {Object[]}
     */
    getEvents() {
        return this.#events;
    }

    /**
     * @returns {number}
     */
    getPage() {
        return this.#page;
    }

    /**
     * @returns {number}
     */
    getPageSize() {
        return this.#pageSize;
    }

    /**
     * @returns {boolean}
     */
    hasMore() {
        return this.#hasMore;
    }
}
