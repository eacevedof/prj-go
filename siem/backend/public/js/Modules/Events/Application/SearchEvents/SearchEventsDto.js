import { EventsPaginationEnum } from "/js/Modules/Events/Domain/Enums/EventsPaginationEnum.js";
import { EventsFiltersDto } from "/js/Modules/Events/Application/EventsFiltersDto.js";

/**
 * Entrada del caso de uso: los filtros del visor (EventsFiltersDto) y la pagina.
 */
export class SearchEventsDto extends EventsFiltersDto {

    #page;

    /**
     * @param {Object} primitives
     */
    constructor(primitives) {
        super(primitives);
        this.#page = Number(primitives.page) || EventsPaginationEnum.FIRST_PAGE;
    }

    /**
     * @param {Object} primitives
     * @returns {SearchEventsDto}
     */
    static fromPrimitives(primitives = {}) {
        return new SearchEventsDto(primitives);
    }

    /**
     * @returns {number}
     */
    getPage() {
        return this.#page;
    }
}
