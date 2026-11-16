import { GetEventsFilterOptionsResultDto } from "/js/Modules/Events/Application/GetEventsFilterOptions/GetEventsFilterOptionsResultDto.js";
import { EventsReaderHttpRepository } from "/js/Modules/Events/Infrastructure/Repositories/EventsReaderHttpRepository.js";
import { EventsException } from "/js/Modules/Events/Domain/Exceptions/EventsException.js";

/**
 * Los proyectos, agentes y monitores que hay, para los desplegables del visor (Proyecto -> Agente -> Monitor). Las
 * facetas solo traen lo que tiene eventos en el rango; aqui va todo, para poder elegir antes de que haya eventos.
 */
export class GetEventsFilterOptionsService {

    #eventsReaderHttpRepository;

    constructor() {
        this.#eventsReaderHttpRepository = EventsReaderHttpRepository.getInstance();
    }

    /**
     * @returns {GetEventsFilterOptionsService}
     */
    static getInstance() {
        return new GetEventsFilterOptionsService();
    }

    /**
     * @param {import("./GetEventsFilterOptionsDto.js").GetEventsFilterOptionsDto} getEventsFilterOptionsDto
     * @returns {Promise<GetEventsFilterOptionsResultDto>}
     * @throws {EventsException}
     */
    async invoke(getEventsFilterOptionsDto) {
        const apiResponse = await this.#eventsReaderHttpRepository.getEventsFilterOptions(
            getEventsFilterOptionsDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            EventsException.fromApiResponse(apiResponse);
        }
        return GetEventsFilterOptionsResultDto.fromPrimitives(apiResponse.data ?? {});
    }
}
