import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";

import { AgentsEventEnum } from "/js/Modules/Agents/Domain/Enums/AgentsEventEnum.js";

const CSS_HIDDEN = "hidden";
const PATH_INGEST = "/api/v1/ingest-batches";

/**
 * Alta de un agente en el detalle de un proyecto y el bloque que enseña su api_token, UNA sola vez: el token (con
 * Copiar) y el .env ya montado (endpoint + token) listo para pegar en el servidor. El bloque no se repinta con el
 * refresco y solo desaparece al cerrarlo o al cambiar de proyecto. El token no se guarda en ningun sitio.
 * Una sola instancia: guarda las referencias al DOM y el proyecto que se esta viendo.
 */
export class AgentFormView {

    static #instance = null;

    #dispatcher;
    #htmler;
    #elements = {};
    #projectUuid = "";

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#htmler = Htmler.getInstance();
    }

    /**
     * @returns {AgentFormView}
     */
    static getInstance() {
        AgentFormView.#instance ??= new AgentFormView();
        return AgentFormView.#instance;
    }

    init() {
        this.#elements = {
            formAgent: document.getElementById("form-agent"),
            inputAgentName: document.getElementById("input-agent-name"),
            inputAgentHeartbeat: document.getElementById("input-agent-heartbeat"),
            buttonAgentSubmit: document.getElementById("button-agent-submit"),
            pAgentFormErrors: document.getElementById("p-agent-form-errors"),
            pAgentsErrors: document.getElementById("p-agents-errors"),
            divAgentToken: document.getElementById("div-agent-token"),
            divAgentTokenContent: document.getElementById("div-agent-token-content"),
            buttonAgentTokenHide: document.getElementById("button-agent-token-hide"),
        };
        this.#elements.formAgent.addEventListener("submit", (event) => this.#onSubmit(event));
        this.#elements.buttonAgentTokenHide.addEventListener("click", () => this.#hideToken());
    }

    /**
     * El proyecto que se esta viendo. Al cambiar de proyecto se vacia todo (el token de otro proyecto no se queda a la
     * vista); en un refresco del mismo, nada cambia.
     * @param {string} projectUuid
     */
    renderProject(projectUuid) {
        if (projectUuid === this.#projectUuid) return;
        this.#projectUuid = projectUuid;
        this.#elements.formAgent.reset();
        this.renderErrors([]);
        this.renderListErrors([]);
        this.#hideToken();
    }

    /**
     * @param {{agentName: string, apiToken: string, endpointUrl: string}} agent el recien creado; endpointUrl: la url de
     *     ingesta segun el backend (vacia = la de esta misma origin)
     */
    renderCreated(agent) {
        const h = this.#htmler;
        const endpointUrl = agent.endpointUrl || `${window.location.origin}${PATH_INGEST}`;
        const envFile = `APP_ENV=production\nAPP_SIEM_ENDPOINT_URL=${endpointUrl}\nAPP_SIEM_API_TOKEN=${agent.apiToken}`;
        this.#elements.formAgent.reset();
        this.renderErrors([]);
        this.#elements.divAgentTokenContent.innerHTML = `
            <p><strong>Agente ${h.getEscaped(agent.agentName)} creado.</strong> Este token <strong>no se vuelve a mostrar</strong>:
                cópialo ahora. Si lo pierdes, borra el agente y créalo de nuevo.</p>
            <p>${h.getCopyableValue(agent.apiToken)}</p>
            <p class="muted">El fichero <code>.env</code> que va junto al binario en el servidor (no hace falta nada más:
                lo que se vigila lo descarga el agente de este panel):</p>
            ${h.getCodeBlock(envFile, ".env")}
            <p>Binarios y pasos de instalación: <a href="#/downloads">3 - descargas</a>.</p>`;
        this.#elements.divAgentToken.classList.remove(CSS_HIDDEN);
    }

    /**
     * @param {string[]} errors del formulario de alta
     */
    renderErrors(errors) {
        this.#renderText(this.#elements.pAgentFormErrors, errors.join(" · "));
    }

    /**
     * @param {string[]} errors de la lista de agentes (una baja que falla)
     */
    renderListErrors(errors) {
        this.#renderText(this.#elements.pAgentsErrors, errors.join(" · "));
    }

    /**
     * @param {boolean} isBusy
     */
    renderBusy(isBusy) {
        this.#elements.buttonAgentSubmit.disabled = isBusy;
    }

    #hideToken() {
        this.#elements.divAgentTokenContent.innerHTML = "";
        this.#elements.divAgentToken.classList.add(CSS_HIDDEN);
    }

    #onSubmit(event) {
        event.preventDefault();
        this.renderErrors([]);
        this.#dispatcher.dispatch(AgentsEventEnum.CREATE_AGENT_REQUESTED, {
            projectUuid: this.#projectUuid,
            agentName: this.#elements.inputAgentName.value,
            heartbeatIntervalSeconds: this.#elements.inputAgentHeartbeat.value,
        });
    }

    #renderText(element, text) {
        element.textContent = text;
        element.classList.toggle(CSS_HIDDEN, text === "");
    }
}
