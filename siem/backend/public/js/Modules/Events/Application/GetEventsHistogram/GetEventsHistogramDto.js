import { EventsFiltersDto } from "/js/Modules/Events/Application/EventsFiltersDto.js";

/**
 * Entrada del caso de uso: los mismos filtros del visor (EventsFiltersDto), sin pagina.
 */
export class GetEventsHistogramDto extends EventsFiltersDto {

    /**
     * @param {Object} primitives
     * @returns {GetEventsHistogramDto}
     */
    static fromPrimitives(primitives = {}) {
        return new GetEventsHistogramDto(primitives);
    }
}
