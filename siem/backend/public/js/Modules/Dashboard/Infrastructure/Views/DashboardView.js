import { DateTimer } from "/js/Modules/Shared/Infrastructure/Components/DateTimer.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { NumberFormatter } from "/js/Modules/Shared/Infrastructure/Components/NumberFormatter.js";

const CSS_HIDDEN = "hidden";

/**
 * Portada: estado del servidor (tarjetas de cifras) y una tarjeta por proyecto con su semaforo. Fechas en UTC.
 * Si falla un refresco se deja lo ultimo pintado y se muestra el error encima.
 * Una sola instancia: guarda las referencias al DOM.
 */
export class DashboardView {

    static #instance = null;

    #dateTimer;
    #htmler;
    #numberFormatter;
    #elements = {};

    constructor() {
        this.#dateTimer = DateTimer.getInstance();
        this.#htmler = Htmler.getInstance();
        this.#numberFormatter = NumberFormatter.getInstance();
    }

    /**
     * @returns {DashboardView}
     */
    static getInstance() {
        DashboardView.#instance ??= new DashboardView();
        return DashboardView.#instance;
    }

    init() {
        this.#elements = {
            pDashboardErrors: document.getElementById("p-dashboard-errors"),
            divDashboardServer: document.getElementById("div-dashboard-server"),
            divDashboardProjects: document.getElementById("div-dashboard-projects"),
        };
    }

    /**
     * @param {Object} server tal cual lo devuelve la API (claves snake_case)
     * @param {Object[]} projects
     */
    renderDashboard(server, projects) {
        this.renderErrors([]);
        this.#elements.divDashboardServer.innerHTML = this.#getServerHtml(server);
        this.#elements.divDashboardProjects.innerHTML = this.#getProjectsHtml(projects, server.now);
    }

    /**
     * @param {string[]} errors
     */
    renderErrors(errors) {
        this.#elements.pDashboardErrors.textContent = errors.join(" · ");
        this.#elements.pDashboardErrors.classList.toggle(CSS_HIDDEN, errors.length === 0);
    }

    #getServerHtml(server) {
        const h = this.#htmler;
        const n = this.#numberFormatter;
        const isWorkerOff = server.queue_worker_interval_seconds === 0;
        const stats = [
            ["1.1 Versión", `${h.getEscaped(server.app_version)} · ${h.getEscaped(server.app_env)}`, false],
            ["1.2 Hora del servidor (UTC)", h.getEscaped(this.#dateTimer.getUtcDateTime(server.now)), false],
            ["1.3 Arrancado (UTC)", `${h.getEscaped(this.#dateTimer.getUtcDateTime(server.started_at))} · ${this.#dateTimer.getDuration(server.uptime_seconds)}`, false],
            ["1.4 Eventos guardados", n.getNumber(server.events_total), false],
            ["1.5 Tamaño de la BD", n.getBytes(server.db_size_bytes), false],
            ["1.6 Cola: lotes pendientes", `${n.getNumber(server.queue_pending_batches)} (${n.getNumber(server.queue_pending_events)} ev.)`, false],
            ["1.7 Cola: lotes fallidos", n.getNumber(server.queue_failed_batches), server.queue_failed_batches > 0],
            ["1.8 Retención", `${n.getNumber(server.events_retention_days)} días`, false],
            ["1.9 Worker", isWorkerOff ? "apagado" : `cada ${n.getNumber(server.queue_worker_interval_seconds)} s`, isWorkerOff],
        ];
        return `<div class="stats">${stats.map(([label, value, isBad]) => `
            <div class="stat${isBad ? " is-bad" : ""}">
                <div class="stat-label">${label}</div>
                <div class="stat-value">${value}</div>
            </div>`).join("")}</div>`;
    }

    #getProjectsHtml(projects, now) {
        const h = this.#htmler;
        const n = this.#numberFormatter;
        if (projects.length === 0) {
            return `<p class="empty">Sin proyectos todavía. Crea el primero en el apartado 3.</p>`;
        }
        return `<div class="projects">${projects.map((project, index) => `
            <a class="project-card health-${h.getEscaped(project.project_health)}" href="#/projects/${encodeURIComponent(project.project_uuid)}">
                <div class="project-name">
                    <span class="project-title">2.${index + 1} ${h.getEscaped(project.project_name)}</span>
                    ${h.getBadge(project.project_health)}
                </div>
                <div class="muted">${h.getEscaped(project.project_description)}</div>
                <dl>
                    <dt>Agentes online</dt><dd>${n.getNumber(project.agents_online)} / ${n.getNumber(project.agents_total)}</dd>
                    <dt>Eventos 24 h</dt><dd>${n.getNumber(project.events_24h)}</dd>
                    <dt>Warnings 24 h</dt><dd>${n.getNumber(project.warnings_24h)}</dd>
                    <dt>Errores 24 h</dt><dd>${n.getNumber(project.errors_24h)}</dd>
                    <dt>Críticos 24 h</dt><dd>${n.getNumber(project.criticals_24h)}</dd>
                    <dt>Último evento (UTC)</dt><dd>${h.getEscaped(this.#dateTimer.getUtcDateTime(project.last_event_at))}<br>
                        <span class="muted">${this.#dateTimer.getTimeAgo(project.last_event_at, now)}</span></dd>
                </dl>
            </a>`).join("")}</div>`;
    }
}
