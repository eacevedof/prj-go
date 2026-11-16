import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { ScreenEnum } from "/js/Modules/Shared/Domain/Enums/ScreenEnum.js";
import { SharedEventEnum } from "/js/Modules/Shared/Domain/Enums/SharedEventEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { EventsMessageEnum } from "/js/Modules/Events/Domain/Enums/EventsMessageEnum.js";
import { EventsPaginationEnum } from "/js/Modules/Events/Domain/Enums/EventsPaginationEnum.js";
import { EventsRefreshEnum } from "/js/Modules/Events/Domain/Enums/EventsRefreshEnum.js";
import { EventsQueryParser } from "/js/Modules/Events/Domain/Services/EventsQueryParser.js";
import { TimeRangeResolver } from "/js/Modules/Events/Domain/Services/TimeRangeResolver.js";
import { GetEventsFacetsDto } from "/js/Modules/Events/Application/GetEventsFacets/GetEventsFacetsDto.js";
import { GetEventsFacetsService } from "/js/Modules/Events/Application/GetEventsFacets/GetEventsFacetsService.js";
import { GetEventsFilterOptionsDto } from "/js/Modules/Events/Application/GetEventsFilterOptions/GetEventsFilterOptionsDto.js";
import { GetEventsFilterOptionsService } from "/js/Modules/Events/Application/GetEventsFilterOptions/GetEventsFilterOptionsService.js";
import { GetEventsHistogramDto } from "/js/Modules/Events/Application/GetEventsHistogram/GetEventsHistogramDto.js";
import { GetEventsHistogramService } from "/js/Modules/Events/Application/GetEventsHistogram/GetEventsHistogramService.js";
import { SearchEventsDto } from "/js/Modules/Events/Application/SearchEvents/SearchEventsDto.js";
import { SearchEventsService } from "/js/Modules/Events/Application/SearchEvents/SearchEventsService.js";
import { EventsView } from "/js/Modules/Events/Infrastructure/Views/EventsView.js";
import { EventsException } from "/js/Modules/Events/Domain/Exceptions/EventsException.js";

// estados de Promise.allSettled
const PROMISE_FULFILLED = "fulfilled";
const PROMISE_REJECTED = "rejected";

/**
 * Visor de eventos (#/events?...): todo su estado (texto, rango, filtros, pagina, refresco) sale del hash. Pide a la
 * vez la pagina de eventos, el histograma y las facetas con los mismos filtros; si una de las tres falla, las otras
 * dos se pintan igual y el fallo se ve en su hueco. Con un atajo de rango, from/to se recalculan en cada peticion
 * ("ultimas 24 h" avanza sola). Se repinta solo en la primera pagina y con rango relativo: mas alla, las filas se
 * moverian bajo el raton o el rango fijo no tendria nada nuevo.
 */
export class SearchEventsController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #eventsQueryParser;
    #timeRangeResolver;
    #eventsView;
    // si el usuario cambia de filtros antes de que llegue la respuesta, la vieja no debe pisar a la nueva
    #lastRequestNumber = 0;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#eventsQueryParser = EventsQueryParser.getInstance();
        this.#timeRangeResolver = TimeRangeResolver.getInstance();
        this.#eventsView = EventsView.getInstance();
    }

    /**
     * @returns {SearchEventsController}
     */
    static getInstance() {
        return new SearchEventsController();
    }

    init() {
        this.#subscriber.subscribe(SharedEventEnum.SCREEN_REQUESTED, (request) => {
            if (request.screen !== ScreenEnum.EVENTS) return;
            this.invoke({ query: request.query ?? {} });
        });
    }

    /**
     * @param {{query: Object}} request los parametros del hash
     */
    async invoke(request) {
        const requestNumber = ++this.#lastRequestNumber;
        const state = this.#eventsQueryParser.getState(request.query);
        const fromTo = this.#timeRangeResolver.getFromTo(state, new Date());
        this.#eventsView.renderState(state, fromTo);

        const filters = {
            dashboardToken: this.#sessionStorer.getDashboardToken(),
            from: fromTo.from,
            to: fromTo.to,
            search: state.search,
            include: state.include,
            exclude: state.exclude,
        };
        // el catalogo de los desplegables no espera a las consultas de eventos (las facetas pueden tardar segundos)
        this.#renderFilterOptions(requestNumber, filters.dashboardToken);
        const [eventsResult, histogramResult, facetsResult] = await Promise.allSettled([
            SearchEventsService.getInstance().invoke(SearchEventsDto.fromPrimitives({ ...filters, page: state.page })),
            GetEventsHistogramService.getInstance().invoke(GetEventsHistogramDto.fromPrimitives(filters)),
            GetEventsFacetsService.getInstance().invoke(GetEventsFacetsDto.fromPrimitives(filters)),
        ]);
        if (requestNumber !== this.#lastRequestNumber) return;

        if ([eventsResult, histogramResult, facetsResult].some((result) => this.#isUnauthorized(result))) {
            this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
            return;
        }

        if (eventsResult.status === PROMISE_FULFILLED) {
            this.#eventsView.renderEvents({
                events: eventsResult.value.getEvents(),
                page: eventsResult.value.getPage(),
                pageSize: eventsResult.value.getPageSize(),
                hasMore: eventsResult.value.hasMore(),
            });
        }
        else {
            this.#eventsView.renderErrors(this.#getErrors(eventsResult.reason, "events"));
        }

        if (histogramResult.status === PROMISE_FULFILLED) {
            this.#eventsView.renderHistogram({
                bucketSeconds: histogramResult.value.getBucketSeconds(),
                buckets: histogramResult.value.getBuckets(),
            });
        }
        else {
            this.#eventsView.renderHistogramErrors(this.#getErrors(histogramResult.reason, "histogram"));
        }

        if (facetsResult.status === PROMISE_FULFILLED) {
            this.#eventsView.renderFacets(facetsResult.value.getFacets());
        }
        else {
            this.#eventsView.renderFacetsErrors(this.#getErrors(facetsResult.reason, "facets"));
        }

        this.#dispatcher.dispatch(SharedEventEnum.SCREEN_RENDERED, {
            isAutoRefreshed: state.refresh !== EventsRefreshEnum.OFF && state.page === EventsPaginationEnum.FIRST_PAGE && state.range !== "",
            intervalSeconds: Number(state.refresh),
        });
    }

    /**
     * Catalogo de los desplegables Proyecto / Agente / Monitor. Si falla, se quedan con lo ultimo que hubo: los filtros
     * siguen funcionando por el hash (y un 401 lo cazan las otras tres peticiones).
     * @param {number} requestNumber
     * @param {string} dashboardToken
     */
    async #renderFilterOptions(requestNumber, dashboardToken) {
        try {
            const getEventsFilterOptionsResultDto = await GetEventsFilterOptionsService.getInstance().invoke(
                GetEventsFilterOptionsDto.fromPrimitives({ dashboardToken })
            );
            if (requestNumber !== this.#lastRequestNumber) return;
            this.#eventsView.renderFilterOptions(getEventsFilterOptionsResultDto.getFilterOptions());
        }
        catch (error) {
            if (error instanceof EventsException) return;
            console.error("SearchEventsController.invoke filter-options", error);
        }
    }

    #isUnauthorized(result) {
        return result.status === PROMISE_REJECTED && result.reason instanceof EventsException
            && result.reason.getCode() === HttpCodeEnum.UNAUTHORIZED;
    }

    #getErrors(error, part) {
        if (error instanceof EventsException) {
            return error.getErrors();
        }
        console.error(`SearchEventsController.invoke ${part}`, error);
        return [EventsMessageEnum.UNEXPECTED_ERROR];
    }
}
