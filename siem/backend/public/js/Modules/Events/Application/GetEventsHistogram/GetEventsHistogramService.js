import { EventsKeyEnum } from "/js/Modules/Events/Domain/Enums/EventsKeyEnum.js";
import { EventsQueryParser } from "/js/Modules/Events/Domain/Services/EventsQueryParser.js";
import { GetEventsHistogramResultDto } from "/js/Modules/Events/Application/GetEventsHistogram/GetEventsHistogramResultDto.js";
import { EventsReaderHttpRepository } from "/js/Modules/Events/Infrastructure/Repositories/EventsReaderHttpRepository.js";
import { EventsException } from "/js/Modules/Events/Domain/Exceptions/EventsException.js";

/**
 * Cuantos eventos hay por tramo de tiempo y gravedad con los filtros del visor: los datos del histograma. Los tramos
 * (unos 60 por rango) los calcula la API.
 */
export class GetEventsHistogramService {

    #eventsQueryParser;
    #eventsReaderHttpRepository;

    constructor() {
        this.#eventsQueryParser = EventsQueryParser.getInstance();
        this.#eventsReaderHttpRepository = EventsReaderHttpRepository.getInstance();
    }

    /**
     * @returns {GetEventsHistogramService}
     */
    static getInstance() {
        return new GetEventsHistogramService();
    }

    /**
     * @param {import("./GetEventsHistogramDto.js").GetEventsHistogramDto} getEventsHistogramDto
     * @returns {Promise<GetEventsHistogramResultDto>}
     * @throws {EventsException}
     */
    async invoke(getEventsHistogramDto) {
        const apiResponse = await this.#eventsReaderHttpRepository.getEventsHistogram(
            this.#eventsQueryParser.getApiFilters(getEventsHistogramDto.getFilters()),
            getEventsHistogramDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            EventsException.fromApiResponse(apiResponse);
        }

        const histogram = apiResponse.data ?? {};
        return GetEventsHistogramResultDto.fromPrimitives({
            from: histogram[EventsKeyEnum.FROM],
            to: histogram[EventsKeyEnum.TO],
            bucketSeconds: histogram[EventsKeyEnum.BUCKET_SECONDS],
            buckets: histogram[EventsKeyEnum.BUCKETS],
        });
    }
}
