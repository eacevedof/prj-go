import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { TextFormatter } from "/js/Modules/Shared/Infrastructure/Components/TextFormatter.js";

import { AgentsEventEnum } from "/js/Modules/Agents/Domain/Enums/AgentsEventEnum.js";
import { AgentsKeyEnum } from "/js/Modules/Agents/Domain/Enums/AgentsKeyEnum.js";
import { AgentsTextEnum } from "/js/Modules/Agents/Domain/Enums/AgentsTextEnum.js";

const CSS_HIDDEN = "hidden";
const PATH_INGEST = "/api/v1/ingest-batches";

/**
 * Formulario plegable "Nuevo agente" del detalle de un proyecto (patron task-js, como el de proyectos): el mismo para
 * editar (agent_uuid oculto; en edicion aparece "Activo"), foco en el nombre al desplegar, Ctrl+Enter guarda.
 * Tras un alta enseña el api_token UNA sola vez: el token (con Copiar) y el .env ya montado (endpoint + token). Ese
 * bloque no se repinta con el refresco y solo desaparece al cerrarlo o al cambiar de proyecto; no se guarda en ningun
 * sitio. Una sola instancia: guarda las referencias al DOM, el proyecto que se esta viendo y el id que se edita.
 */
export class AgentFormView {

    static #instance = null;

    #dispatcher;
    #htmler;
    #textFormatter;
    #elements = {};
    #projectUuid = "";
    // agent_id del agente que se edita (va en el titulo); "" = alta
    #editingAgentId = "";

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#htmler = Htmler.getInstance();
        this.#textFormatter = TextFormatter.getInstance();
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
            detailsAgentForm: document.getElementById("details-agent-form"),
            summaryAgentForm: document.getElementById("summary-agent-form"),
            formAgent: document.getElementById("form-agent"),
            hiddenAgentUuid: document.getElementById("hidden-agent-uuid"),
            inputAgentName: document.getElementById("input-agent-name"),
            inputAgentHeartbeat: document.getElementById("input-agent-heartbeat"),
            divAgentFieldEnabled: document.getElementById("div-agent-field-enabled"),
            checkboxAgentEnabled: document.getElementById("checkbox-agent-enabled"),
            buttonAgentSubmit: document.getElementById("button-agent-submit"),
            buttonAgentCancel: document.getElementById("button-agent-cancel"),
            pAgentFormErrors: document.getElementById("p-agent-form-errors"),
            divAgentToken: document.getElementById("div-agent-token"),
            divAgentTokenContent: document.getElementById("div-agent-token-content"),
            buttonAgentTokenHide: document.getElementById("button-agent-token-hide"),
        };
        this.#elements.detailsAgentForm.addEventListener("toggle", () => this.#onAgentFormToggle());
        this.#elements.formAgent.addEventListener("submit", (event) => this.#onAgentSubmit(event));
        this.#elements.formAgent.addEventListener("keydown", (event) => this.#onAgentFormKeydown(event));
        this.#elements.buttonAgentCancel.addEventListener("click", () => this.reset());
        this.#elements.buttonAgentTokenHide.addEventListener("click", () => this.#hideToken());
        this.#renderFormMode();
    }

    /**
     * El proyecto que se esta viendo. Al cambiar de proyecto se vacia todo (el token de otro proyecto no se queda a la
     * vista); en un refresco del mismo, nada cambia.
     * @param {string} projectUuid
     */
    renderProject(projectUuid) {
        if (projectUuid === this.#projectUuid) return;
        this.#projectUuid = projectUuid;
        this.reset();
        this.#hideToken();
    }

    /**
     * Vacio, plegado y en modo alta.
     */
    reset() {
        this.#elements.formAgent.reset();
        this.#elements.hiddenAgentUuid.value = "";
        this.#editingAgentId = "";
        this.#renderFormMode();
        this.#elements.detailsAgentForm.open = false;
        this.renderErrors([]);
        this.renderBusy(false);
    }

    /**
     * Si se estaba editando ese agente (lo acaban de borrar), el formulario vuelve a alta.
     * @param {string} agentUuid
     */
    resetIfEditing(agentUuid) {
        if (this.#elements.hiddenAgentUuid.value !== agentUuid) return;
        this.reset();
    }

    /**
     * "Editar" en la tabla: el agente entra en el formulario, que se despliega con el foco en el nombre.
     * @param {Object} agent tal cual lo da el detalle del proyecto (snake_case)
     */
    renderAgent(agent) {
        this.#elements.hiddenAgentUuid.value = String(agent[AgentsKeyEnum.AGENT_UUID] ?? "");
        this.#editingAgentId = String(agent[AgentsKeyEnum.AGENT_ID] ?? "");
        this.#elements.inputAgentName.value = String(agent[AgentsKeyEnum.AGENT_NAME] ?? "");
        this.#elements.inputAgentHeartbeat.value = String(agent[AgentsKeyEnum.HEARTBEAT_INTERVAL_SECONDS] ?? "");
        this.#elements.checkboxAgentEnabled.checked = agent[AgentsKeyEnum.IS_ENABLED] !== false;
        this.#renderFormMode();
        this.renderErrors([]);
        this.#elements.detailsAgentForm.open = true;
        this.#elements.detailsAgentForm.scrollIntoView({ block: "nearest" });
        this.#elements.inputAgentName.focus();
    }

    /**
     * Alta hecha: formulario vacio y plegado, y el bloque del token (una sola vez).
     * @param {{agentName: string, apiToken: string, endpointUrl: string}} agent el recien creado; endpointUrl: la url de
     *     ingesta segun el backend (vacia = la de esta misma origin)
     */
    renderCreated(agent) {
        const h = this.#htmler;
        const endpointUrl = agent.endpointUrl || `${window.location.origin}${PATH_INGEST}`;
        const envFile = `APP_ENV=production\nAPP_SIEM_ENDPOINT_URL=${endpointUrl}\nAPP_SIEM_API_TOKEN=${agent.apiToken}`;
        this.reset();
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
     * @param {string[]} errors
     */
    renderErrors(errors) {
        this.#elements.pAgentFormErrors.textContent = errors.join(" · ");
        this.#elements.pAgentFormErrors.classList.toggle(CSS_HIDDEN, errors.length === 0);
    }

    /**
     * @param {boolean} isBusy
     */
    renderBusy(isBusy) {
        this.#elements.buttonAgentSubmit.disabled = isBusy;
        this.#elements.buttonAgentCancel.disabled = isBusy;
    }

    #hideToken() {
        this.#elements.divAgentTokenContent.innerHTML = "";
        this.#elements.divAgentToken.classList.add(CSS_HIDDEN);
    }

    #onAgentFormToggle() {
        if (!this.#elements.detailsAgentForm.open) return;
        this.#elements.inputAgentName.focus();
    }

    /**
     * Ctrl+Enter (Cmd+Enter) guarda desde cualquier campo; requestSubmit pasa por #onAgentSubmit.
     */
    #onAgentFormKeydown(event) {
        if (event.key !== "Enter" || !(event.ctrlKey || event.metaKey)) return;

        event.preventDefault();
        this.#elements.formAgent.requestSubmit();
    }

    #onAgentSubmit(event) {
        event.preventDefault();
        this.renderErrors([]);
        this.renderBusy(true);
        const agentUuid = this.#elements.hiddenAgentUuid.value;
        if (!agentUuid) {
            this.#dispatcher.dispatch(AgentsEventEnum.CREATE_AGENT_REQUESTED, {
                projectUuid: this.#projectUuid,
                agentName: this.#elements.inputAgentName.value,
                heartbeatIntervalSeconds: this.#elements.inputAgentHeartbeat.value,
            });
            return;
        }
        this.#dispatcher.dispatch(AgentsEventEnum.UPDATE_AGENT_REQUESTED, {
            agentUuid,
            agentName: this.#elements.inputAgentName.value,
            heartbeatIntervalSeconds: this.#elements.inputAgentHeartbeat.value,
            isEnabled: this.#elements.checkboxAgentEnabled.checked,
        });
    }

    /**
     * Alta: "Nuevo agente", sin "Activo" ni "Cancelar edición". Edicion: "Editar agente #4" con los dos.
     */
    #renderFormMode() {
        const isEditing = this.#editingAgentId !== "";
        this.#elements.summaryAgentForm.textContent = isEditing
            ? this.#textFormatter.getFormatted(AgentsTextEnum.FORM_TITLE_UPDATE, this.#editingAgentId)
            : AgentsTextEnum.FORM_TITLE_CREATE;
        this.#elements.buttonAgentSubmit.textContent = isEditing ? AgentsTextEnum.SUBMIT_UPDATE : AgentsTextEnum.SUBMIT_CREATE;
        this.#elements.divAgentFieldEnabled.classList.toggle(CSS_HIDDEN, !isEditing);
        this.#elements.buttonAgentCancel.classList.toggle(CSS_HIDDEN, !isEditing);
    }
}
