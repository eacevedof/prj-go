import { EventsKeyEnum } from "/js/Modules/Events/Domain/Enums/EventsKeyEnum.js";
import { EventsQueryParser } from "/js/Modules/Events/Domain/Services/EventsQueryParser.js";
import { SearchEventsResultDto } from "/js/Modules/Events/Application/SearchEvents/SearchEventsResultDto.js";
import { EventsReaderHttpRepository } from "/js/Modules/Events/Infrastructure/Repositories/EventsReaderHttpRepository.js";
import { EventsException } from "/js/Modules/Events/Domain/Exceptions/EventsException.js";

/**
 * Una pagina de eventos con los filtros del visor (rango, texto, incluir y excluir por campo). Filtrar, ordenar y
 * paginar lo hace la API (y valida los filtros: un valor no valido vuelve como 400): aqui solo se pide y se da forma
 * al resultado.
 */
export class SearchEventsService {

    #eventsQueryParser;
    #eventsReaderHttpRepository;

    constructor() {
        this.#eventsQueryParser = EventsQueryParser.getInstance();
        this.#eventsReaderHttpRepository = EventsReaderHttpRepository.getInstance();
    }

    /**
     * @returns {SearchEventsService}
     */
    static getInstance() {
        return new SearchEventsService();
    }

    /**
     * @param {import("./SearchEventsDto.js").SearchEventsDto} searchEventsDto
     * @returns {Promise<SearchEventsResultDto>}
     * @throws {EventsException}
     */
    async invoke(searchEventsDto) {
        const apiResponse = await this.#eventsReaderHttpRepository.getEventsPage(
            {
                ...this.#eventsQueryParser.getApiFilters(searchEventsDto.getFilters()),
                [EventsKeyEnum.PAGE]: searchEventsDto.getPage(),
            },
            searchEventsDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            EventsException.fromApiResponse(apiResponse);
        }

        const eventsPage = apiResponse.data ?? {};
        return SearchEventsResultDto.fromPrimitives({
            events: eventsPage[EventsKeyEnum.EVENTS],
            page: eventsPage[EventsKeyEnum.PAGE],
            pageSize: eventsPage[EventsKeyEnum.PAGE_SIZE],
            hasMore: eventsPage[EventsKeyEnum.HAS_MORE],
        });
    }
}
