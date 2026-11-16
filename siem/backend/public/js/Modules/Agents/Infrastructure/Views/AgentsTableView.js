import { DateTimer } from "/js/Modules/Shared/Infrastructure/Components/DateTimer.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { NumberFormatter } from "/js/Modules/Shared/Infrastructure/Components/NumberFormatter.js";
import { TextFormatter } from "/js/Modules/Shared/Infrastructure/Components/TextFormatter.js";

import { AgentsEventEnum } from "/js/Modules/Agents/Domain/Enums/AgentsEventEnum.js";
import { AgentsKeyEnum } from "/js/Modules/Agents/Domain/Enums/AgentsKeyEnum.js";
import { AgentStatusEnum } from "/js/Modules/Agents/Domain/Enums/AgentStatusEnum.js";
import { AgentsTextEnum } from "/js/Modules/Agents/Domain/Enums/AgentsTextEnum.js";

const CSS_HIDDEN = "hidden";
const ACTION_TOGGLE_DETAIL = "toggle-detail";
const ACTION_EDIT = "edit";
const ACTION_ASK_DELETE = "ask-delete";
const ACTION_CONFIRM_DELETE = "confirm-delete";
const ACTION_CANCEL_DELETE = "cancel-delete";
// el buscador no repinta en cada tecla
const TEXT_FILTER_DEBOUNCE_MS = 300;
const COLUMNS_COUNT = 7;
const EVENTS_HASH = "#/events";
// color de cada contador de 24 h segun su gravedad (gris si es 0); las mismas clases que la tabla de proyectos
const COUNT_CSS_WARNING = "count-warning";
const COUNT_CSS_ERROR = "count-error";
const COUNT_CSS_CRITICAL = "count-critical";
const COUNT_CSS_ZERO = "count-zero";

/**
 * Tabla de agentes del detalle de un proyecto (mismo patron que la de proyectos): filtros (texto y estado), una fila
 * por agente con id, uuid, nombre (con su estado, y su host y ultimo latido debajo), warnings, errores y criticos de
 * 24 h, y acciones (ver detalle, editar, borrar). "Ver detalle" despliega bajo la fila su latido, host, cliente y sus
 * fuentes (carpeta, comando, contenedor, monitor) con el ultimo resultado. Borrar se confirma en la propia fila.
 * Filas desplegadas, confirmacion pendiente y filtros sobreviven al auto-refresco; al cambiar de proyecto se vacian.
 * Una sola instancia: guarda las referencias al DOM y ese estado.
 */
export class AgentsTableView {

    static #instance = null;

    #dateTimer;
    #dispatcher;
    #htmler;
    #numberFormatter;
    #textFormatter;
    #elements = {};
    #projectUuid = "";
    #serverNow = "";
    #agentsByUuid = new Map();
    #expandedAgentUuids = new Set();
    // agent_uuid de la fila que pregunta "¿Borrar?"; "" = ninguna
    #askingDeleteAgentUuid = "";
    #textFilterTimeoutId = null;

    constructor() {
        this.#dateTimer = DateTimer.getInstance();
        this.#dispatcher = Dispatcher.getInstance();
        this.#htmler = Htmler.getInstance();
        this.#numberFormatter = NumberFormatter.getInstance();
        this.#textFormatter = TextFormatter.getInstance();
    }

    /**
     * @returns {AgentsTableView}
     */
    static getInstance() {
        AgentsTableView.#instance ??= new AgentsTableView();
        return AgentsTableView.#instance;
    }

    init() {
        this.#elements = {
            spanAgentsCount: document.getElementById("span-agents-count"),
            pAgentsMessage: document.getElementById("p-agents-message"),
            pAgentsErrors: document.getElementById("p-agents-errors"),
            formAgentsFilters: document.getElementById("form-agents-filters"),
            inputFilterAgentText: document.getElementById("input-filter-agent-text"),
            selectFilterAgentStatus: document.getElementById("select-filter-agent-status"),
            buttonClearAgentsFilters: document.getElementById("button-clear-agents-filters"),
            pAgentsEmpty: document.getElementById("p-agents-empty"),
            divAgentsTable: document.getElementById("div-agents-table"),
            tbodyAgents: document.getElementById("tbody-agents"),
        };
        this.#loadStatusOptions();
        this.#elements.formAgentsFilters.addEventListener("submit", (event) => this.#onFiltersSubmit(event));
        this.#elements.inputFilterAgentText.addEventListener("input", () => this.#onTextFilterInput());
        this.#elements.selectFilterAgentStatus.addEventListener("change", () => this.#dispatchFiltersChanged());
        this.#elements.buttonClearAgentsFilters.addEventListener("click", () => this.#onClearFiltersClick());
        this.#elements.tbodyAgents.addEventListener("click", (event) => this.#onAgentsTableClick(event));
    }

    /**
     * El proyecto que se esta viendo. Si es otro: sin filtros, filas plegadas, sin confirmacion ni avisos.
     * @param {string} projectUuid
     */
    renderProject(projectUuid) {
        if (projectUuid === this.#projectUuid) return;
        this.#projectUuid = projectUuid;
        clearTimeout(this.#textFilterTimeoutId);
        this.#elements.formAgentsFilters.reset();
        this.#expandedAgentUuids.clear();
        this.#askingDeleteAgentUuid = "";
        this.renderMessage("");
        this.renderErrors([]);
    }

    /**
     * @returns {{agentText: string, agentStatus: string}}
     */
    getFilters() {
        return {
            agentText: this.#elements.inputFilterAgentText.value,
            agentStatus: this.#elements.selectFilterAgentStatus.value,
        };
    }

    /**
     * @param {Object[]} filteredAgents los que cumplen los filtros (tal cual los da la API)
     * @param {number} agentsTotal todos, sin filtrar
     * @param {boolean} hasFilters
     * @param {string} serverNow UTC del servidor (para el "hace X")
     */
    renderAgents(filteredAgents, agentsTotal, hasFilters, serverNow) {
        this.#serverNow = serverNow;
        this.#agentsByUuid = new Map(filteredAgents.map((agent) => [agent[AgentsKeyEnum.AGENT_UUID], agent]));
        if (!this.#agentsByUuid.has(this.#askingDeleteAgentUuid)) this.#askingDeleteAgentUuid = "";

        this.#elements.spanAgentsCount.textContent = this.#getCountText(filteredAgents.length, agentsTotal);
        this.#elements.pAgentsEmpty.textContent = hasFilters ? AgentsTextEnum.LIST_NO_MATCHES : AgentsTextEnum.LIST_EMPTY;
        this.#elements.pAgentsEmpty.classList.toggle(CSS_HIDDEN, filteredAgents.length > 0);
        this.#elements.divAgentsTable.classList.toggle(CSS_HIDDEN, filteredAgents.length === 0);
        this.#renderAgentsTable();
    }

    /**
     * Aviso del camino feliz (guardado, borrado). Vacio lo oculta.
     * @param {string} message
     */
    renderMessage(message) {
        this.#elements.pAgentsMessage.textContent = message;
        this.#elements.pAgentsMessage.classList.toggle(CSS_HIDDEN, message === "");
    }

    /**
     * @param {string[]} errors
     */
    renderErrors(errors) {
        this.#elements.pAgentsErrors.textContent = errors.join(" · ");
        this.#elements.pAgentsErrors.classList.toggle(CSS_HIDDEN, errors.length === 0);
    }

    #loadStatusOptions() {
        const statusOptions = AgentStatusEnum.VALUES
            .map((agentStatus) => `<option value="${this.#htmler.getEscaped(agentStatus)}">${this.#htmler.getEscaped(AgentStatusEnum.LABELS[agentStatus])}</option>`)
            .join("");
        this.#elements.selectFilterAgentStatus.innerHTML = `<option value="">${this.#htmler.getEscaped(AgentsTextEnum.STATUS_ALL)}</option>${statusOptions}`;
    }

    #onFiltersSubmit(event) {
        event.preventDefault();
        clearTimeout(this.#textFilterTimeoutId);
        this.#dispatchFiltersChanged();
    }

    #onTextFilterInput() {
        clearTimeout(this.#textFilterTimeoutId);
        this.#textFilterTimeoutId = setTimeout(() => this.#dispatchFiltersChanged(), TEXT_FILTER_DEBOUNCE_MS);
    }

    #onClearFiltersClick() {
        clearTimeout(this.#textFilterTimeoutId);
        this.#elements.formAgentsFilters.reset();
        this.#dispatchFiltersChanged();
    }

    #dispatchFiltersChanged() {
        this.#dispatcher.dispatch(AgentsEventEnum.AGENTS_FILTERS_CHANGED, { filters: this.getFilters() });
    }

    #onAgentsTableClick(event) {
        const actionButton = event.target.closest("button[data-action]");
        if (!actionButton) return;

        const agentUuid = actionButton.dataset.agentUuid ?? "";
        const agent = this.#agentsByUuid.get(agentUuid);
        switch (actionButton.dataset.action) {
            case ACTION_TOGGLE_DETAIL:
                if (this.#expandedAgentUuids.has(agentUuid)) this.#expandedAgentUuids.delete(agentUuid);
                else this.#expandedAgentUuids.add(agentUuid);
                this.#renderAgentsTable();
                return;
            case ACTION_EDIT:
                if (!agent) return;
                this.renderMessage("");
                this.#dispatcher.dispatch(AgentsEventEnum.EDIT_AGENT_REQUESTED, { agent });
                return;
            case ACTION_ASK_DELETE:
                this.#askingDeleteAgentUuid = agentUuid;
                this.#renderAgentsTable();
                return;
            case ACTION_CANCEL_DELETE:
                this.#askingDeleteAgentUuid = "";
                this.#renderAgentsTable();
                return;
            case ACTION_CONFIRM_DELETE:
                if (!agent) return;
                this.#askingDeleteAgentUuid = "";
                actionButton.disabled = true;
                this.renderMessage("");
                this.renderErrors([]);
                this.#dispatcher.dispatch(AgentsEventEnum.DELETE_AGENT_REQUESTED, {
                    agentUuid,
                    agentName: agent[AgentsKeyEnum.AGENT_NAME],
                });
        }
    }

    #renderAgentsTable() {
        this.#elements.tbodyAgents.innerHTML = [...this.#agentsByUuid.values()]
            .map((agent) => this.#getHtmlAgentRows(agent))
            .join("");
    }

    #getCountText(filteredTotal, agentsTotal) {
        if (agentsTotal === 0) return "";
        if (filteredTotal === agentsTotal) return this.#textFormatter.getFormatted(AgentsTextEnum.COUNT_ALL, agentsTotal);
        return this.#textFormatter.getFormatted(AgentsTextEnum.COUNT_FILTERED, filteredTotal, agentsTotal);
    }

    /**
     * La fila del agente y, si esta desplegado, la de su detalle.
     */
    #getHtmlAgentRows(agent) {
        const h = this.#htmler;
        const agentUuid = String(agent[AgentsKeyEnum.AGENT_UUID] ?? "");
        const isExpanded = this.#expandedAgentUuids.has(agentUuid);
        const agentRow = `
<tr class="agent-row${isExpanded ? " is-expanded" : ""}${agent[AgentsKeyEnum.IS_ENABLED] === false ? " is-disabled" : ""}">
    <td class="num row-id">${h.getEscaped(agent[AgentsKeyEnum.AGENT_ID])}</td>
    <td class="row-uuid"><span class="copyable"><code>${h.getEscaped(agentUuid)}</code>${h.getCopyButton(agentUuid)}</span></td>
    <td class="row-main">
        <div class="row-main-head">
            <span class="row-title">${h.getEscaped(agent[AgentsKeyEnum.AGENT_NAME])}</span>
            ${h.getBadge(agent[AgentsKeyEnum.AGENT_STATUS])}
        </div>
        <div class="row-sub">${h.getEscaped(this.#getHostText(agent))} · ${h.getEscaped(this.#getLastSeenText(agent))}</div>
    </td>
    ${this.#getHtmlCountCell(agent[AgentsKeyEnum.WARNINGS_24H], COUNT_CSS_WARNING)}
    ${this.#getHtmlCountCell(agent[AgentsKeyEnum.ERRORS_24H], COUNT_CSS_ERROR)}
    ${this.#getHtmlCountCell(agent[AgentsKeyEnum.CRITICALS_24H], COUNT_CSS_CRITICAL)}
    <td class="row-actions">${this.#getHtmlRowActions(agentUuid, isExpanded)}</td>
</tr>`;
        if (!isExpanded) return agentRow;
        return `${agentRow}
<tr class="agent-detail"><td colspan="${COLUMNS_COUNT}">${this.#getHtmlAgentDetail(agent)}</td></tr>`;
    }

    #getHostText(agent) {
        const hostParts = [agent[AgentsKeyEnum.HOST_NAME], agent[AgentsKeyEnum.HOST_OS], agent[AgentsKeyEnum.HOST_IP]]
            .map((value) => String(value ?? "").trim())
            .filter((value) => value !== "");
        return hostParts.length ? hostParts.join(" ") : AgentsTextEnum.HOST_UNKNOWN;
    }

    #getLastSeenText(agent) {
        const lastSeenAt = agent[AgentsKeyEnum.LAST_SEEN_AT];
        if (!lastSeenAt) return AgentsTextEnum.NEVER_SEEN;
        return this.#dateTimer.getTimeAgo(lastSeenAt, this.#serverNow);
    }

    #getHtmlCountCell(eventsCount, countCss) {
        const count = Number(eventsCount) || 0;
        return `<td class="num ${count > 0 ? countCss : COUNT_CSS_ZERO}">${this.#numberFormatter.getNumber(count)}</td>`;
    }

    #getHtmlRowActions(agentUuid, isExpanded) {
        const h = this.#htmler;
        const escapedAgentUuid = h.getEscaped(agentUuid);
        if (agentUuid === this.#askingDeleteAgentUuid) {
            return `<span class="confirm">${h.getEscaped(AgentsTextEnum.CONFIRM_DELETE)}
                <button type="button" class="button button-small button-danger" data-action="${ACTION_CONFIRM_DELETE}" data-agent-uuid="${escapedAgentUuid}">${h.getEscaped(AgentsTextEnum.CONFIRM_DELETE_YES)}</button>
                <button type="button" class="button button-small" data-action="${ACTION_CANCEL_DELETE}">${h.getEscaped(AgentsTextEnum.CONFIRM_DELETE_NO)}</button></span>`;
        }
        const toggleText = isExpanded ? AgentsTextEnum.HIDE_DETAIL : AgentsTextEnum.SHOW_DETAIL;
        return `<button type="button" class="button button-small button-ghost" data-action="${ACTION_TOGGLE_DETAIL}" data-agent-uuid="${escapedAgentUuid}" aria-expanded="${isExpanded}">${h.getEscaped(toggleText)}</button>
            <button type="button" class="button button-small" data-action="${ACTION_EDIT}" data-agent-uuid="${escapedAgentUuid}">${h.getEscaped(AgentsTextEnum.EDIT)}</button>
            <button type="button" class="button button-small button-ghost" data-action="${ACTION_ASK_DELETE}" data-agent-uuid="${escapedAgentUuid}">${h.getEscaped(AgentsTextEnum.DELETE)}</button>`;
    }

    /**
     * Lo que antes era la tarjeta del agente: hechos y la tabla de sus fuentes con su ultimo resultado.
     */
    #getHtmlAgentDetail(agent) {
        const h = this.#htmler;
        const n = this.#numberFormatter;
        const eventsUrl = `${EVENTS_HASH}?project_uuid=${encodeURIComponent(this.#projectUuid)}&agent_uuid=${encodeURIComponent(agent[AgentsKeyEnum.AGENT_UUID] ?? "")}`;
        return `
            <dl class="facts">
                <dt>Último latido (UTC)</dt><dd>${h.getEscaped(this.#dateTimer.getUtcDateTime(agent[AgentsKeyEnum.LAST_SEEN_AT]))}
                    <span class="muted">${this.#dateTimer.getTimeAgo(agent[AgentsKeyEnum.LAST_SEEN_AT], this.#serverNow)}</span></dd>
                <dt>Latido esperado</dt><dd>cada ${n.getNumber(agent[AgentsKeyEnum.HEARTBEAT_INTERVAL_SECONDS])} s</dd>
                <dt>Host</dt><dd>${h.getEscaped(this.#getHostText(agent))}</dd>
                <dt>Cliente</dt><dd>${h.getEscaped(agent[AgentsKeyEnum.CLIENT_VERSION] || "?")}</dd>
                <dt>24 h</dt><dd>${n.getNumber(agent[AgentsKeyEnum.EVENTS_24H])} eventos · <a href="${eventsUrl}">${h.getEscaped(AgentsTextEnum.EVENTS)}</a></dd>
            </dl>
            ${this.#getHtmlSources(agent[AgentsKeyEnum.SOURCES] ?? [], eventsUrl)}`;
    }

    #getHtmlSources(sources, agentEventsUrl) {
        const h = this.#htmler;
        const n = this.#numberFormatter;
        if (sources.length === 0) {
            return `<p class="empty">${h.getEscaped(AgentsTextEnum.NO_SOURCES)}</p>`;
        }
        return `<div class="table-wrap"><table class="table">
            <thead><tr>
                <th class="num">#</th><th>Tipo</th><th>Fuente</th><th class="num">Eventos</th><th>Último (UTC)</th>
                <th>Gravedad</th><th class="num">Exit</th><th>Último mensaje</th>
            </tr></thead>
            <tbody>${sources.map((source, index) => `<tr>
                <td class="num">${index + 1}</td>
                <td class="mono">${h.getEscaped(source.source_type)}</td>
                <td><a href="${agentEventsUrl}&source_name=${encodeURIComponent(source.source_name)}">${h.getEscaped(source.source_name)}</a></td>
                <td class="num">${n.getNumber(source.events_count)}</td>
                <td class="nowrap mono">${h.getEscaped(this.#dateTimer.getUtcDateTime(source.last_event_at))}<br>
                    <span class="muted">${this.#dateTimer.getTimeAgo(source.last_event_at, this.#serverNow)}</span></td>
                <td>${source.last_severity ? h.getBadge(source.last_severity) : ""}</td>
                <td class="num">${h.getEscaped(source.last_exit_code ?? "")}</td>
                <td class="message"><div class="is-clamped">${h.getEscaped(source.last_message)}</div></td>
            </tr>`).join("")}</tbody>
        </table></div>`;
    }
}
