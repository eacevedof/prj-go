import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

const CSS_HIDDEN = "hidden";

/**
 * Boton "Salir" de la barra superior: solo con sesion.
 */
export class SessionView {

    static #instance = null;

    #dispatcher;
    #elements = {};

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
    }

    /**
     * @returns {SessionView}
     */
    static getInstance() {
        SessionView.#instance ??= new SessionView();
        return SessionView.#instance;
    }

    init() {
        this.#elements = {
            divSession: document.getElementById("div-session"),
            buttonLogOut: document.getElementById("button-log-out"),
        };
        this.#elements.buttonLogOut.addEventListener("click", () => this.#dispatcher.dispatch(SessionEventEnum.LOG_OUT_REQUESTED));
    }

    show() {
        this.#elements.divSession.classList.remove(CSS_HIDDEN);
    }

    hide() {
        this.#elements.divSession.classList.add(CSS_HIDDEN);
    }
}
