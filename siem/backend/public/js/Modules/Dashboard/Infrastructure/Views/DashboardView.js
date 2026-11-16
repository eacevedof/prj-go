import { DateTimer } from "/js/Modules/Shared/Infrastructure/Components/DateTimer.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { NumberFormatter } from "/js/Modules/Shared/Infrastructure/Components/NumberFormatter.js";

const CSS_HIDDEN = "hidden";

/**
 * Portada: estado del servidor (tarjetas de cifras). Fechas en UTC. La tabla de proyectos es de Projects (ProjectsTableView).
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
        };
    }

    /**
     * @param {Object} server tal cual lo devuelve la API (claves snake_case)
     */
    renderServer(server) {
        this.renderErrors([]);
        this.#elements.divDashboardServer.innerHTML = this.#getServerHtml(server);
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
}
