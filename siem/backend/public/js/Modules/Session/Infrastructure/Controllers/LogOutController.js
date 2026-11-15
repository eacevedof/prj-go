import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";
import { SessionTextEnum } from "/js/Modules/Session/Domain/Enums/SessionTextEnum.js";
import { AuthView } from "/js/Modules/Session/Infrastructure/Views/AuthView.js";
import { SessionView } from "/js/Modules/Session/Infrastructure/Views/SessionView.js";

/**
 * Salir (a peticion del usuario o porque la API respondio 401). No hay caso de uso: cerrar sesion es borrar el token
 * de la pestaña; la API no guarda sesiones.
 */
export class LogOutController {

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
     * @returns {LogOutController}
     */
    static getInstance() {
        return new LogOutController();
    }

    init() {
        this.#subscriber.subscribe(SessionEventEnum.LOG_OUT_REQUESTED, () => this.invoke(SessionTextEnum.LOGGED_OUT));
        this.#subscriber.subscribe(SessionEventEnum.SESSION_EXPIRED, () => this.invoke(SessionTextEnum.SESSION_EXPIRED));
    }

    /**
     * @param {string} message aviso para la pantalla de entrada
     */
    invoke(message) {
        try {
            this.#sessionStorer.removeSession();
        }
        catch (error) {
            console.error("LogOutController.invoke", error);
        }
        this.#sessionView.hide();
        this.#navigationView.hide();
        this.#authView.show(message);
        this.#dispatcher.dispatch(SessionEventEnum.USER_LOGGED_OUT);
    }
}
