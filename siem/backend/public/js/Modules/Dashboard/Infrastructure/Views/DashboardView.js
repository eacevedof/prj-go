import { DateTimer } from "/js/Modules/Shared/Infrastructure/Components/DateTimer.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { NumberFormatter } from "/js/Modules/Shared/Infrastructure/Components/NumberFormatter.js";
import { TextFormatter } from "/js/Modules/Shared/Infrastructure/Components/TextFormatter.js";

import { DashboardTextEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardTextEnum.js";

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
    #textFormatter;
    #elements = {};

    constructor() {
        this.#dateTimer = DateTimer.getInstance();
        this.#htmler = Htmler.getInstance();
        this.#numberFormatter = NumberFormatter.getInstance();
        this.#textFormatter = TextFormatter.getInstance();
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
        const numberFormatter = this.#numberFormatter;
        const isWorkerOff = server.queue_worker_interval_seconds === 0;
        const stats = [
            [DashboardTextEnum.SERVER_VERSION, `${h.getEscaped(server.app_version)} · ${h.getEscaped(server.app_env)}`, false],
            [DashboardTextEnum.SERVER_NOW, h.getEscaped(this.#dateTimer.getUtcDateTime(server.now)), false],
            [DashboardTextEnum.SERVER_STARTED_AT, `${h.getEscaped(this.#dateTimer.getUtcDateTime(server.started_at))} · ${this.#dateTimer.getDuration(server.uptime_seconds)}`, false],
            [DashboardTextEnum.EVENTS_TOTAL, numberFormatter.getNumber(server.events_total), false],
            [DashboardTextEnum.DB_SIZE, numberFormatter.getBytes(server.db_size_bytes), false],
            [DashboardTextEnum.QUEUE_PENDING, `${numberFormatter.getNumber(server.queue_pending_batches)} ${this.#textFormatter.getFormatted(DashboardTextEnum.QUEUE_PENDING_EVENTS, numberFormatter.getNumber(server.queue_pending_events))}`, false],
            [DashboardTextEnum.QUEUE_FAILED, numberFormatter.getNumber(server.queue_failed_batches), server.queue_failed_batches > 0],
            [DashboardTextEnum.EVENTS_RETENTION, this.#textFormatter.getFormatted(DashboardTextEnum.RETENTION_DAYS, numberFormatter.getNumber(server.events_retention_days)), false],
            [DashboardTextEnum.QUEUE_WORKER, isWorkerOff
                ? DashboardTextEnum.WORKER_OFF
                : this.#textFormatter.getFormatted(DashboardTextEnum.WORKER_EVERY, numberFormatter.getNumber(server.queue_worker_interval_seconds)), isWorkerOff],
        ];
        return `<div class="stats">${stats.map(([label, value, isBad]) => `
            <div class="stat${isBad ? " is-bad" : ""}">
                <div class="stat-label">${label}</div>
                <div class="stat-value">${value}</div>
            </div>`).join("")}</div>`;
    }
}
