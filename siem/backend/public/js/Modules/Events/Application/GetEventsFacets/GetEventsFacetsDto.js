import { EventsFiltersDto } from "/js/Modules/Events/Application/EventsFiltersDto.js";

/**
 * Entrada del caso de uso: los mismos filtros del visor (EventsFiltersDto), sin pagina.
 */
export class GetEventsFacetsDto extends EventsFiltersDto {

    /**
     * @param {Object} primitives
     * @returns {GetEventsFacetsDto}
     */
    static fromPrimitives(primitives = {}) {
        return new GetEventsFacetsDto(primitives);
    }
}
