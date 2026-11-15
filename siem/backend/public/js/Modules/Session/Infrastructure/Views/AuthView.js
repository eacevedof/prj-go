import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

const CSS_HIDDEN = "hidden";

/**
 * Pantalla de entrada: el token del dashboard (APP_DASHBOARD_TOKEN del .env del backend), un aviso (sesion cerrada,
 * token caducado) y los errores. Una sola instancia: guarda las referencias al DOM.
 */
export class AuthView {

    static #instance = null;

    #dispatcher;
    #elements = {};

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
    }

    /**
     * @returns {AuthView}
     */
    static getInstance() {
        AuthView.#instance ??= new AuthView();
        return AuthView.#instance;
    }

    init() {
        this.#elements = {
            sectionAuth: document.getElementById("section-auth"),
            formLogIn: document.getElementById("form-log-in"),
            inputDashboardToken: document.getElementById("input-dashboard-token"),
            pAuthMessage: document.getElementById("p-auth-message"),
            pAuthErrors: document.getElementById("p-auth-errors"),
        };
        this.#elements.formLogIn.addEventListener("submit", (event) => this.#onLogInSubmit(event));
    }

    /**
     * @param {string} message aviso opcional encima del formulario
     */
    show(message = "") {
        this.#elements.formLogIn.reset();
        this.#renderText(this.#elements.pAuthMessage, message);
        this.renderErrors([]);
        this.#elements.sectionAuth.classList.remove(CSS_HIDDEN);
        this.#elements.inputDashboardToken.focus();
    }

    hide() {
        this.#elements.sectionAuth.classList.add(CSS_HIDDEN);
        this.#elements.formLogIn.reset();
    }

    /**
     * @param {string[]} errors
     */
    renderErrors(errors) {
        this.#renderText(this.#elements.pAuthErrors, errors.join(" · "));
    }

    #onLogInSubmit(event) {
        event.preventDefault();
        this.renderErrors([]);
        this.#dispatcher.dispatch(SessionEventEnum.LOG_IN_REQUESTED, {
            dashboardToken: this.#elements.inputDashboardToken.value,
        });
    }

    #renderText(element, text) {
        element.textContent = text;
        element.classList.toggle(CSS_HIDDEN, text === "");
    }
}
