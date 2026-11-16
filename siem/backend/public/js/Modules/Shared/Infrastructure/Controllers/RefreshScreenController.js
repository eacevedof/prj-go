import { RefreshEnum } from "/js/Modules/Shared/Domain/Enums/RefreshEnum.js";
import { SharedEventEnum } from "/js/Modules/Shared/Domain/Enums/SharedEventEnum.js";
import { DateTimer } from "/js/Modules/Shared/Infrastructure/Components/DateTimer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

const MILLISECONDS_PER_SECOND = 1000;

/**
 * Repinta sola la pantalla visible: el estado de agentes y proyectos cambia sin que nadie toque nada. Cada pintado
 * (SCREEN_RENDERED) reprograma el siguiente; cambiar de pantalla o salir lo cancela. La pantalla decide si se refresca
 * y cada cuanto (el visor de eventos tiene su propio selector; el resto, RefreshEnum.INTERVAL_SECONDS).
 */
export class RefreshScreenController {

    #dateTimer;
    #subscriber;
    #navigationView;
    #timeoutId = null;

    constructor() {
        this.#dateTimer = DateTimer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#navigationView = NavigationView.getInstance();
    }

    /**
     * @returns {RefreshScreenController}
     */
    static getInstance() {
        return new RefreshScreenController();
    }

    init() {
        this.#subscriber.subscribe(SharedEventEnum.SCREEN_RENDERED, (request) => this.invoke(request));
        this.#subscriber.subscribe(SharedEventEnum.SCREEN_REQUESTED, () => this.#cancel());
        this.#subscriber.subscribe(SessionEventEnum.USER_LOGGED_OUT, () => this.#cancel());
    }

    /**
     * @param {{isAutoRefreshed: boolean, intervalSeconds?: number}} request
     */
    invoke(request) {
        this.#cancel();
        const isAutoRefreshed = request.isAutoRefreshed === true;
        const intervalSeconds = Number(request.intervalSeconds) > 0 ? Number(request.intervalSeconds) : RefreshEnum.INTERVAL_SECONDS;
        this.#navigationView.renderRefreshInfo(isAutoRefreshed, this.#dateTimer.getUtcDateTimeNow(), intervalSeconds);
        if (!isAutoRefreshed) return;

        this.#timeoutId = setTimeout(() => this.#navigationView.renderCurrentScreen(), intervalSeconds * MILLISECONDS_PER_SECOND);
    }

    #cancel() {
        clearTimeout(this.#timeoutId);
        this.#timeoutId = null;
    }
}
