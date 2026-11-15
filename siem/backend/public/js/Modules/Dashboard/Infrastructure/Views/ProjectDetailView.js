import { DateTimer } from "/js/Modules/Shared/Infrastructure/Components/DateTimer.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { NumberFormatter } from "/js/Modules/Shared/Infrastructure/Components/NumberFormatter.js";

import { AgentsEventEnum } from "/js/Modules/Agents/Domain/Enums/AgentsEventEnum.js";

const CSS_HIDDEN = "hidden";
const ACTION_ASK_DELETE = "ask-delete-agent";
const ACTION_CONFIRM_DELETE = "confirm-delete-agent";
const ACTION_CANCEL_DELETE = "cancel-delete-agent";

/**
 * Detalle de un proyecto: su cabecera y sus agentes; por agente, cada fuente (carpeta, comando, contenedor, monitor)
 * con su ultimo estado. Fechas en UTC con el "hace X" contra el now del servidor. Borrar un agente pide un segundo
 * clic dentro de la propia tarjeta (el estado de esa confirmacion sobrevive al refresco).
 * Una sola instancia: guarda las referencias al DOM y lo ultimo pintado.
 */
export class ProjectDetailView {

    static #instance = null;

    #dateTimer;
    #dispatcher;
    #htmler;
    #numberFormatter;
    #elements = {};
    #lastRender = null;
    #pendingDeleteAgentUuid = "";

    constructor() {
        this.#dateTimer = DateTimer.getInstance();
        this.#dispatcher = Dispatcher.getInstance();
        this.#htmler = Htmler.getInstance();
        this.#numberFormatter = NumberFormatter.getInstance();
    }

    /**
     * @returns {ProjectDetailView}
     */
    static getInstance() {
        ProjectDetailView.#instance ??= new ProjectDetailView();
        return ProjectDetailView.#instance;
    }

    init() {
        this.#elements = {
            pProjectErrors: document.getElementById("p-project-errors"),
            divProjectBody: document.getElementById("div-project-body"),
            divProjectHeader: document.getElementById("div-project-header"),
            divProjectAgents: document.getElementById("div-project-agents"),
        };
        this.#elements.divProjectAgents.addEventListener("click", (event) => this.#onAgentsClick(event));
    }

    /**
     * @param {Object} project tal cual lo devuelve la API
     * @param {Object[]} agents cada uno con sus sources[]
     * @param {string} serverNow UTC
     */
    renderProjectDetail(project, agents, serverNow) {
        this.#lastRender = { project, agents, serverNow };
        // un agente que ya no existe no puede quedarse esperando confirmacion
        if (!agents.some((agent) => agent.agent_uuid === this.#pendingDeleteAgentUuid)) {
            this.#pendingDeleteAgentUuid = "";
        }
        this.renderErrors([]);
        this.#elements.divProjectBody.classList.remove(CSS_HIDDEN);
        this.#elements.divProjectHeader.innerHTML = this.#getHeaderHtml(project);
        this.#elements.divProjectAgents.innerHTML = this.#getAgentsHtml(project, agents, serverNow);
    }

    /**
     * Si el proyecto no existe o falla la API, se oculta el detalle: no se deja otro proyecto pintado.
     * @param {string[]} errors
     */
    renderErrors(errors) {
        this.#elements.pProjectErrors.textContent = errors.join(" · ");
        this.#elements.pProjectErrors.classList.toggle(CSS_HIDDEN, errors.length === 0);
        if (errors.length) {
            this.#elements.divProjectBody.classList.add(CSS_HIDDEN);
        }
    }

    #getHeaderHtml(project) {
        const h = this.#htmler;
        const projectUuid = project.project_uuid ?? "";
        return `
            <h1>1 - ${h.getEscaped(project.project_name)}</h1>
            <p class="muted">${h.getEscaped(project.project_description)}</p>
            <p>${h.getCopyableValue(projectUuid)}
                <a class="link-action" href="#/events?project_uuid=${encodeURIComponent(projectUuid)}">ver sus eventos</a></p>`;
    }

    #getAgentsHtml(project, agents, serverNow) {
        if (agents.length === 0) {
            return `<p class="empty">Sin agentes. Da de alta el primero en el apartado 3: te dará el token para instalarlo.</p>`;
        }
        return agents.map((agent, index) => this.#getAgentHtml(agent, index + 1, project.project_uuid ?? "", serverNow)).join("");
    }

    #getAgentHtml(agent, position, projectUuid, serverNow) {
        const h = this.#htmler;
        const n = this.#numberFormatter;
        const eventsUrl = `#/events?project_uuid=${encodeURIComponent(projectUuid)}&agent_uuid=${encodeURIComponent(agent.agent_uuid)}`;
        return `
            <section class="panel agent">
                <div class="panel-head">
                    <span class="agent-name">2.${position} ${h.getEscaped(agent.agent_name)}</span>
                    ${h.getBadge(agent.agent_status)}
                    <a class="link-action" href="${eventsUrl}">eventos</a>
                    <span class="panel-actions">${this.#getDeleteAgentHtml(agent)}</span>
                </div>
                <dl class="facts">
                    <dt>Último latido (UTC)</dt><dd>${h.getEscaped(this.#dateTimer.getUtcDateTime(agent.last_seen_at))}
                        <span class="muted">${this.#dateTimer.getTimeAgo(agent.last_seen_at, serverNow)}</span></dd>
                    <dt>Host</dt><dd>${h.getEscaped(agent.host_name || "sin informar")} ${h.getEscaped(agent.host_os)} ${h.getEscaped(agent.host_ip)}</dd>
                    <dt>Cliente</dt><dd>${h.getEscaped(agent.client_version || "?")}</dd>
                    <dt>24 h</dt><dd>${n.getNumber(agent.events_24h)} eventos · ${n.getNumber(agent.warnings_24h)} warnings ·
                        ${n.getNumber(agent.errors_24h)} errores · ${n.getNumber(agent.criticals_24h)} críticos</dd>
                    <dt>agent_uuid</dt><dd>${h.getCopyableValue(agent.agent_uuid)}</dd>
                </dl>
                ${this.#getSourcesHtml(agent.sources ?? [], eventsUrl, serverNow)}
            </section>`;
    }

    #getDeleteAgentHtml(agent) {
        const agentUuid = this.#htmler.getEscaped(agent.agent_uuid);
        if (agent.agent_uuid !== this.#pendingDeleteAgentUuid) {
            return `<button type="button" class="button button-small button-ghost" data-action="${ACTION_ASK_DELETE}" data-agent-uuid="${agentUuid}">Borrar agente</button>`;
        }
        return `<span class="confirm">Se borran el agente y sus eventos.
            <button type="button" class="button button-small button-danger" data-action="${ACTION_CONFIRM_DELETE}" data-agent-uuid="${agentUuid}">Sí, borrar</button>
            <button type="button" class="button button-small" data-action="${ACTION_CANCEL_DELETE}">Cancelar</button></span>`;
    }

    #getSourcesHtml(sources, agentEventsUrl, serverNow) {
        const h = this.#htmler;
        const n = this.#numberFormatter;
        if (sources.length === 0) {
            return `<p class="empty">Aún no ha enviado eventos.</p>`;
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
                    <span class="muted">${this.#dateTimer.getTimeAgo(source.last_event_at, serverNow)}</span></td>
                <td>${source.last_severity ? h.getBadge(source.last_severity) : ""}</td>
                <td class="num">${h.getEscaped(source.last_exit_code ?? "")}</td>
                <td class="message">${h.getEscaped(source.last_message)}</td>
            </tr>`).join("")}</tbody>
        </table></div>`;
    }

    #onAgentsClick(event) {
        const button = event.target.closest("button[data-action]");
        if (!button || !this.#lastRender) return;

        const action = button.dataset.action;
        if (action === ACTION_CONFIRM_DELETE) {
            this.#pendingDeleteAgentUuid = "";
            this.#dispatcher.dispatch(AgentsEventEnum.DELETE_AGENT_REQUESTED, { agentUuid: button.dataset.agentUuid });
            return;
        }
        this.#pendingDeleteAgentUuid = action === ACTION_ASK_DELETE ? button.dataset.agentUuid : "";
        this.renderProjectDetail(this.#lastRender.project, this.#lastRender.agents, this.#lastRender.serverNow);
    }
}
