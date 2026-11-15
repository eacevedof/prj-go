import { EventsKeyEnum } from "/js/Modules/Events/Domain/Enums/EventsKeyEnum.js";
import { EventsException } from "/js/Modules/Events/Domain/Exceptions/EventsException.js";
import { EventsQueryParser } from "/js/Modules/Events/Domain/Services/EventsQueryParser.js";
import { GetEventsFacetsResultDto } from "/js/Modules/Events/Application/GetEventsFacets/GetEventsFacetsResultDto.js";
import { EventsReaderHttpRepository } from "/js/Modules/Events/Infrastructure/Repositories/EventsReaderHttpRepository.js";

/**
 * Los valores mas repetidos de cada campo con los filtros del visor: el panel lateral de campos, desde el que se
 * incluye o excluye un valor con un clic.
 */
export class GetEventsFacetsService {

    #eventsQueryParser;
    #eventsReaderHttpRepository;

    constructor() {
        this.#eventsQueryParser = EventsQueryParser.getInstance();
        this.#eventsReaderHttpRepository = EventsReaderHttpRepository.getInstance();
    }

    /**
     * @returns {GetEventsFacetsService}
     */
    static getInstance() {
        return new GetEventsFacetsService();
    }

    /**
     * @param {import("./GetEventsFacetsDto.js").GetEventsFacetsDto} getEventsFacetsDto
     * @returns {Promise<GetEventsFacetsResultDto>}
     * @throws {EventsException}
     */
    async invoke(getEventsFacetsDto) {
        const apiResponse = await this.#eventsReaderHttpRepository.getEventsFacets(
            this.#eventsQueryParser.getApiFilters(getEventsFacetsDto.getFilters()),
            getEventsFacetsDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            EventsException.fromApiResponse(apiResponse);
        }
        return GetEventsFacetsResultDto.fromPrimitives({ facets: apiResponse.data?.[EventsKeyEnum.FACETS] });
    }
}
