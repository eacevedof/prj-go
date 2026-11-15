import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";
import { AuthView } from "/js/Modules/Session/Infrastructure/Views/AuthView.js";
import { SessionView } from "/js/Modules/Session/Infrastructure/Views/SessionView.js";

/**
 * Token aceptado: se guarda en la pestaña, se oculta la entrada y se pinta la pantalla del hash.
 */
export class StartSessionController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #navigationView;
    #authView;
    #sessionView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#navigationView = NavigationView.getInstance();
        this.#authView = AuthView.getInstance();
        this.#sessionView = SessionView.getInstance();
    }

    /**
     * @returns {StartSessionController}
     */
    static getInstance() {
        return new StartSessionController();
    }

    init() {
        this.#subscriber.subscribe(SessionEventEnum.SESSION_STARTED, (request) => this.invoke(request));
    }

    /**
     * @param {{dashboardToken: string}} request
     */
    invoke(request) {
        try {
            this.#sessionStorer.saveSession({ dashboardToken: request.dashboardToken });
        }
        catch (error) {
            // sin sessionStorage la sesion vale hasta recargar: los repositorios leen el token de aqui y fallarian con 401
            console.error("StartSessionController.invoke", error);
        }
        this.#authView.hide();
        this.#sessionView.show();
        this.#navigationView.show();
        this.#dispatcher.dispatch(SessionEventEnum.USER_LOGGED_IN);
    }
}
