import { DateTimer } from "/js/Modules/Shared/Infrastructure/Components/DateTimer.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { NumberFormatter } from "/js/Modules/Shared/Infrastructure/Components/NumberFormatter.js";

import {
    EVENT_FACETS, EVENT_FILTER_FIELD_LABELS, EVENT_FILTER_FIELDS, EventFilterFieldEnum,
} from "/js/Modules/Events/Domain/Enums/EventFilterFieldEnum.js";
import { EventsPaginationEnum } from "/js/Modules/Events/Domain/Enums/EventsPaginationEnum.js";
import { EVENTS_REFRESH_LABELS, EVENTS_REFRESH_VALUES } from "/js/Modules/Events/Domain/Enums/EventsRefreshEnum.js";
import { TIME_RANGE_LABELS, TIME_RANGE_SECONDS, TimeRangeEnum } from "/js/Modules/Events/Domain/Enums/TimeRangeEnum.js";
import { EventsQueryParser } from "/js/Modules/Events/Domain/Services/EventsQueryParser.js";
import { HistogramScaler } from "/js/Modules/Events/Domain/Services/HistogramScaler.js";
import { TimeRangeResolver } from "/js/Modules/Events/Domain/Services/TimeRangeResolver.js";

const CSS_HIDDEN = "hidden";
const ACTION_FILTER = "filter";
const ACTION_PILL_INVERT = "pill-invert";
const ACTION_PILL_REMOVE = "pill-remove";
const ACTION_PILLS_CLEAR = "pills-clear";
const ACTION_ROW_TOGGLE = "row-toggle";
const ACTION_BUCKET = "bucket";
const SEVERITIES_TOP_DOWN = Object.freeze(["critical", "error", "warning", "info", "debug"]);
const HISTOGRAM_BOX = Object.freeze({ width: 960, height: 170, left: 40, right: 8, top: 10, bottom: 24 });
const HISTOGRAM_HINT = "Pasa por una barra para ver su tramo; haz clic para acotar el rango a ese tramo.";
const TABLE_COLUMNS = 6;
const RANGE_ERROR = "from / to: formato YYYY-MM-DD HH:MM:SS (UTC), con from anterior o igual a to";

/**
 * Visor de eventos tipo Discover: barra de texto, rango de tiempo (atajos o a medida), auto-refresco, pastillas de
 * filtro (incluir / excluir, invertibles), histograma apilado por gravedad, panel lateral de campos con sus valores
 * mas repetidos, y tabla con filas desplegables (todos los campos del evento y sus attributes). Todo cambio es un
 * cambio de hash: el router vuelve a pedir la pantalla con el estado nuevo y la url se puede compartir.
 * Los controles son html fijo; el refresco solo repinta pastillas, histograma, facetas y tabla, y respeta las filas
 * desplegadas. Una sola instancia: guarda las referencias al DOM y lo ultimo pintado.
 */
export class EventsView {

    static #instance = null;

    #dateTimer;
    #htmler;
    #numberFormatter;
    #eventsQueryParser;
    #histogramScaler;
    #timeRangeResolver;
    #elements = {};
    #state = null;
    #eventsPage = null;
    #facets = null;
    #bucketSeconds = 0;
    #expandedEventIds = new Set();
    // uuid -> nombre de proyectos y agentes vistos (en eventos y facetas): las pastillas enseñan el nombre
    #namesByUuid = new Map();

    constructor() {
        this.#dateTimer = DateTimer.getInstance();
        this.#htmler = Htmler.getInstance();
        this.#numberFormatter = NumberFormatter.getInstance();
        this.#eventsQueryParser = EventsQueryParser.getInstance();
        this.#histogramScaler = HistogramScaler.getInstance();
        this.#timeRangeResolver = TimeRangeResolver.getInstance();
        this.#state = this.#eventsQueryParser.getState({});
    }

    /**
     * @returns {EventsView}
     */
    static getInstance() {
        EventsView.#instance ??= new EventsView();
        return EventsView.#instance;
    }

    init() {
        this.#elements = {
            sectionEvents: document.getElementById("section-events"),
            formEventsQuery: document.getElementById("form-events-query"),
            inputEventsSearch: document.getElementById("input-events-search"),
            selectEventsRange: document.getElementById("select-events-range"),
            formEventsCustomRange: document.getElementById("form-events-custom-range"),
            inputEventsFrom: document.getElementById("input-events-from"),
            inputEventsTo: document.getElementById("input-events-to"),
            pEventsRangeErrors: document.getElementById("p-events-range-errors"),
            selectEventsRefresh: document.getElementById("select-events-refresh"),
            spanEventsRangeInfo: document.getElementById("span-events-range-info"),
            divEventsPills: document.getElementById("div-events-pills"),
            divEventsHistogram: document.getElementById("div-events-histogram"),
            pEventsHistogramInfo: document.getElementById("p-events-histogram-info"),
            divEventsFacets: document.getElementById("div-events-facets"),
            pEventsErrors: document.getElementById("p-events-errors"),
            divEventsTable: document.getElementById("div-events-table"),
            buttonEventsPrevious: document.getElementById("button-events-previous"),
            buttonEventsNext: document.getElementById("button-events-next"),
            spanEventsPage: document.getElementById("span-events-page"),
        };
        this.#loadSelectOptions();
        this.#elements.formEventsQuery.addEventListener("submit", (event) => this.#onQuerySubmit(event));
        this.#elements.selectEventsRange.addEventListener("change", () => this.#onRangeChange());
        this.#elements.formEventsCustomRange.addEventListener("submit", (event) => this.#onCustomRangeSubmit(event));
        this.#elements.selectEventsRefresh.addEventListener("change", () => this.#goToState(
            this.#eventsQueryParser.getStateWith(this.#state, { refresh: this.#elements.selectEventsRefresh.value })
        ));
        this.#elements.buttonEventsPrevious.addEventListener("click", () => this.#goToPage(this.#state.page - 1));
        this.#elements.buttonEventsNext.addEventListener("click", () => this.#goToPage(this.#state.page + 1));
        this.#elements.sectionEvents.addEventListener("click", (event) => this.#onClick(event));
        this.#elements.divEventsHistogram.addEventListener("mouseover", (event) => this.#onHistogramHover(event));
        this.#elements.divEventsHistogram.addEventListener("focusin", (event) => this.#onHistogramHover(event));
        this.#elements.divEventsHistogram.addEventListener("mouseleave", () => this.#renderHistogramInfo(HISTOGRAM_HINT));
        this.#elements.divEventsHistogram.addEventListener("keydown", (event) => this.#onHistogramKeydown(event));
    }

    /**
     * El estado del hash en los controles (sin pisar lo que se este escribiendo) y las pastillas de filtro.
     * @param {Object} state ver EventsQueryParser
     * @param {{from: string, to: string}} fromTo el rango ya resuelto para esta peticion
     */
    renderState(state, fromTo) {
        // el texto del buscador solo se repone si no hay nada a medio escribir: ni con el foco dentro, ni con un texto
        // distinto del ultimo que se pinto (escrito y sin buscar todavia). Un refresco no puede pisar lo tecleado.
        const isSearchBeingEdited = document.activeElement === this.#elements.inputEventsSearch
            || this.#elements.inputEventsSearch.value !== (this.#state?.search ?? "");
        if (!isSearchBeingEdited || state.search !== (this.#state?.search ?? "")) {
            this.#elements.inputEventsSearch.value = state.search;
        }
        this.#state = state;
        this.#elements.selectEventsRange.value = state.range || TimeRangeEnum.CUSTOM;
        this.#elements.selectEventsRefresh.value = state.refresh;
        this.#elements.formEventsCustomRange.classList.toggle(CSS_HIDDEN, state.range !== "");
        if (state.range === "" && !this.#elements.formEventsCustomRange.contains(document.activeElement)) {
            this.#elements.inputEventsFrom.value = state.from;
            this.#elements.inputEventsTo.value = state.to;
        }
        this.#renderRangeErrors("");
        this.#elements.spanEventsRangeInfo.textContent = `${fromTo.from} → ${fromTo.to} UTC`;
        this.#renderPills();
    }

    /**
     * @param {{events: Object[], page: number, pageSize: number, hasMore: boolean}} eventsPage
     */
    renderEvents(eventsPage) {
        this.#eventsPage = eventsPage;
        eventsPage.events.forEach((event) => {
            this.#namesByUuid.set(event.project_uuid, event.project_name);
            this.#namesByUuid.set(event.agent_uuid, event.agent_name);
        });
        // una fila desplegada que ya no esta en la pagina no se queda en memoria
        const eventIds = new Set(eventsPage.events.map((event) => String(event.event_id)));
        this.#expandedEventIds = new Set([...this.#expandedEventIds].filter((eventId) => eventIds.has(eventId)));

        this.renderErrors([]);
        this.#renderTable();
        this.#elements.spanEventsPage.textContent = `página ${eventsPage.page}`;
        this.#elements.buttonEventsPrevious.disabled = eventsPage.page <= EventsPaginationEnum.FIRST_PAGE;
        this.#elements.buttonEventsNext.disabled = !eventsPage.hasMore;
        // los nombres recien aprendidos pueden mejorar una pastilla que salia con su uuid
        this.#renderPills();
    }

    /**
     * @param {string[]} errors de la tabla de eventos
     */
    renderErrors(errors) {
        this.#elements.pEventsErrors.textContent = errors.join(" · ");
        this.#elements.pEventsErrors.classList.toggle(CSS_HIDDEN, errors.length === 0);
        if (errors.length) {
            this.#eventsPage = null;
            this.#elements.divEventsTable.innerHTML = "";
            this.#elements.spanEventsPage.textContent = "";
            this.#elements.buttonEventsPrevious.disabled = true;
            this.#elements.buttonEventsNext.disabled = true;
        }
    }

    /**
     * @param {{bucketSeconds: number, buckets: Object[]}} histogram
     */
    renderHistogram(histogram) {
        this.#bucketSeconds = histogram.bucketSeconds;
        this.#elements.divEventsHistogram.innerHTML = this.#getHistogramHtml(histogram);
        this.#renderHistogramInfo(HISTOGRAM_HINT);
    }

    /**
     * @param {string[]} errors
     */
    renderHistogramErrors(errors) {
        this.#elements.divEventsHistogram.innerHTML = `<p class="error">Histograma: ${this.#htmler.getEscaped(errors.join(" · "))}</p>`;
        this.#renderHistogramInfo("");
    }

    /**
     * @param {Object} facets {campo: [{value, filter_field, filter_value, count}]}
     */
    renderFacets(facets) {
        this.#facets = facets;
        EVENT_FACETS.filter((eventFacet) => eventFacet.isName).forEach((eventFacet) => {
            (facets[eventFacet.facet] ?? []).forEach((item) => {
                if (item.filter_value) this.#namesByUuid.set(String(item.filter_value), String(item.value ?? ""));
            });
        });
        this.#renderFacetsList();
        this.#renderPills();
    }

    /**
     * @param {string[]} errors
     */
    renderFacetsErrors(errors) {
        this.#facets = null;
        this.#elements.divEventsFacets.innerHTML = `<p class="error">Campos: ${this.#htmler.getEscaped(errors.join(" · "))}</p>`;
    }

    #loadSelectOptions() {
        const rangeValues = [...Object.keys(TIME_RANGE_SECONDS), TimeRangeEnum.CUSTOM];
        this.#elements.selectEventsRange.innerHTML = rangeValues
            .map((range) => `<option value="${range}">${this.#htmler.getEscaped(TIME_RANGE_LABELS[range])}</option>`).join("");
        this.#elements.selectEventsRefresh.innerHTML = EVENTS_REFRESH_VALUES
            .map((refresh) => `<option value="${refresh}">${this.#htmler.getEscaped(EVENTS_REFRESH_LABELS[refresh])}</option>`).join("");
    }

    // ---- pastillas ----

    #renderPills() {
        const h = this.#htmler;
        const pills = this.#eventsQueryParser.getPills(this.#state);
        if (pills.length === 0) {
            this.#elements.divEventsPills.innerHTML = `<span class="muted">Sin filtros de campo. Añádelos con + / − en el panel de campos o al desplegar un evento.</span>`;
            return;
        }
        this.#elements.divEventsPills.innerHTML = pills.map((pill) => {
            const data = `data-field="${h.getEscaped(pill.field)}" data-value="${h.getEscaped(pill.value)}" data-excluded="${pill.isExcluded ? "1" : "0"}"`;
            return `<span class="pill${pill.isExcluded ? " is-excluded" : ""}" title="${h.getEscaped(pill.field)}: ${h.getEscaped(pill.value)}">
                <span class="pill-text">${pill.isExcluded ? `<span class="pill-not">NO</span> ` : ""}${h.getEscaped(EVENT_FILTER_FIELD_LABELS[pill.field])}:
                    <strong>${h.getEscaped(this.#getDisplayValue(pill.field, pill.value))}</strong></span>
                <button type="button" class="pill-button" data-action="${ACTION_PILL_INVERT}" ${data}
                    aria-label="Invertir el filtro" title="${pill.isExcluded ? "Pasar a incluir" : "Pasar a excluir"}">⇄</button>
                <button type="button" class="pill-button" data-action="${ACTION_PILL_REMOVE}" ${data}
                    aria-label="Quitar el filtro" title="Quitar">×</button>
            </span>`;
        }).join("") + `<button type="button" class="button button-small button-ghost" data-action="${ACTION_PILLS_CLEAR}">Quitar todos</button>`;
    }

    // un uuid de proyecto o agente se enseña con su nombre si ya se ha visto
    #getDisplayValue(field, value) {
        const isUuidField = field === EventFilterFieldEnum.PROJECT_UUID || field === EventFilterFieldEnum.AGENT_UUID;
        return (isUuidField && this.#namesByUuid.get(value)) || value;
    }

    // ---- histograma ----

    #getHistogramHtml(histogram) {
        const h = this.#htmler;
        const chart = this.#histogramScaler.getChart(histogram.buckets, histogram.bucketSeconds, HISTOGRAM_BOX);
        if (chart.bars.length === 0) {
            return `<p class="empty">Sin datos para el histograma en este rango.</p>`;
        }
        const baselineY = chart.plot.y + chart.plot.height;
        const gridHtml = chart.yTicks.map((tick) => `
            <line class="chart-grid" x1="${chart.plot.x}" x2="${chart.plot.x + chart.plot.width}" y1="${tick.y.toFixed(1)}" y2="${tick.y.toFixed(1)}"></line>
            <text class="chart-axis" x="${chart.plot.x - 6}" y="${(tick.y + 3.5).toFixed(1)}" text-anchor="end">${h.getEscaped(tick.label)}</text>`).join("");
        const barsHtml = chart.bars.map((bar) => {
            const fromTo = this.#timeRangeResolver.getBucketFromTo(bar.bucketAt, histogram.bucketSeconds);
            const info = this.#getBucketInfo(bar, fromTo);
            return `<g class="chart-bar" tabindex="0" role="button" data-action="${ACTION_BUCKET}"
                    data-from="${h.getEscaped(fromTo.from)}" data-to="${h.getEscaped(fromTo.to)}" data-info="${h.getEscaped(info)}">
                <title>${h.getEscaped(info)}</title>
                <rect class="chart-hit" x="${bar.slotX.toFixed(2)}" y="${chart.plot.y}" width="${bar.slotWidth.toFixed(2)}" height="${chart.plot.height}"></rect>
                ${bar.segments.map((segment) => `<rect class="chart-segment severity-${segment.severity}" x="${bar.x.toFixed(2)}"
                    y="${segment.y.toFixed(2)}" width="${bar.width.toFixed(2)}" height="${segment.height.toFixed(2)}"></rect>`).join("")}
            </g>`;
        }).join("");
        const xTicksHtml = chart.xTicks.map((tick) => `
            <text class="chart-axis" x="${tick.x.toFixed(1)}" y="${HISTOGRAM_BOX.height - 7}" text-anchor="middle">${h.getEscaped(tick.label)}</text>`).join("");

        return `<svg class="chart" viewBox="0 0 ${HISTOGRAM_BOX.width} ${HISTOGRAM_BOX.height}" role="img"
                aria-label="Eventos por tramo de tiempo, apilados por gravedad: ${this.#numberFormatter.getNumber(chart.total)} en el rango">
            ${gridHtml}
            <line class="chart-baseline" x1="${chart.plot.x}" x2="${chart.plot.x + chart.plot.width}" y1="${baselineY}" y2="${baselineY}"></line>
            ${barsHtml}
            ${xTicksHtml}
        </svg>
        <div class="chart-legend">
            <span class="muted">${this.#numberFormatter.getNumber(chart.total)} eventos · tramos de ${h.getEscaped(this.#dateTimer.getDuration(histogram.bucketSeconds))} · UTC</span>
            ${SEVERITIES_TOP_DOWN.map((severity) => `<span class="legend-item"><span class="legend-swatch severity-${severity}"></span>${severity}</span>`).join("")}
        </div>`;
    }

    // "2026-11-15 18:00:00 → 18:23:59 UTC · 34 eventos (2 error, 32 info)"
    #getBucketInfo(bar, fromTo) {
        const counts = SEVERITIES_TOP_DOWN
            .map((severity) => bar.segments.find((segment) => segment.severity === severity))
            .filter(Boolean)
            .map((segment) => `${segment.count} ${segment.severity}`);
        return `${fromTo.from} → ${fromTo.to} UTC · ${bar.total} eventos${counts.length ? ` (${counts.join(", ")})` : ""}`;
    }

    #renderHistogramInfo(text) {
        this.#elements.pEventsHistogramInfo.textContent = text;
    }

    #onHistogramHover(event) {
        const bar = event.target.closest?.("[data-info]");
        if (bar) this.#renderHistogramInfo(bar.dataset.info);
    }

    #onHistogramKeydown(event) {
        if (event.key !== "Enter" && event.key !== " ") return;
        const bar = event.target.closest?.(`[data-action="${ACTION_BUCKET}"]`);
        if (!bar) return;
        event.preventDefault();
        this.#goToBucket(bar);
    }

    #goToBucket(bar) {
        this.#goToState(this.#eventsQueryParser.getStateWith(this.#state, { range: "", from: bar.dataset.from, to: bar.dataset.to }));
    }

    // ---- facetas ----

    #renderFacetsList() {
        const h = this.#htmler;
        const facets = this.#facets ?? {};
        this.#elements.divEventsFacets.innerHTML = EVENT_FACETS.map((eventFacet, facetIndex) => {
            const items = Array.isArray(facets[eventFacet.facet]) ? facets[eventFacet.facet] : [];
            const maxCount = items.reduce((max, item) => Math.max(max, Number(item.count) || 0), 0) || 1;
            return `<div class="facet">
                <h3>2.${facetIndex + 1} ${h.getEscaped(eventFacet.label)}</h3>
                ${items.length === 0
                    ? `<p class="muted facet-empty">sin valores</p>`
                    : `<ul>${items.map((item) => this.#getFacetValueHtml(eventFacet, item, maxCount)).join("")}</ul>`}
            </div>`;
        }).join("");
    }

    #getFacetValueHtml(eventFacet, item, maxCount) {
        const h = this.#htmler;
        // la API dice con que campo y valor se filtra (un nombre filtra por su uuid); sin eso, el propio valor
        const filterField = EVENT_FILTER_FIELDS.includes(item.filter_field) ? item.filter_field : eventFacet.filterField;
        const filterValue = String(item.filter_value ?? (eventFacet.isName ? "" : item.value ?? ""));
        const isIncluded = (this.#state.include[filterField] ?? []).includes(filterValue);
        const isExcluded = (this.#state.exclude[filterField] ?? []).includes(filterValue);
        const widthPercent = Math.max(2, Math.round((Number(item.count) || 0) / maxCount * 100));
        return `<li class="facet-value${isIncluded ? " is-included" : ""}${isExcluded ? " is-excluded" : ""}">
            <span class="facet-bar" style="width:${widthPercent}%"></span>
            <span class="facet-label" title="${h.getEscaped(item.value)}">${h.getEscaped(item.value)}</span>
            <span class="facet-count">${this.#numberFormatter.getNumber(item.count)}</span>
            ${filterValue === "" ? "" : this.#getFilterButtonsHtml(filterField, filterValue)}
        </li>`;
    }

    #getFilterButtonsHtml(field, value) {
        const h = this.#htmler;
        const data = `data-action="${ACTION_FILTER}" data-field="${h.getEscaped(field)}" data-value="${h.getEscaped(value)}"`;
        return `<span class="filter-buttons">
            <button type="button" class="filter-button" ${data} data-excluded="0" aria-label="Incluir este valor" title="Solo este valor">+</button>
            <button type="button" class="filter-button" ${data} data-excluded="1" aria-label="Excluir este valor" title="Todo menos este valor">−</button>
        </span>`;
    }

    // ---- tabla ----

    #renderTable() {
        const h = this.#htmler;
        const eventsPage = this.#eventsPage;
        if (!eventsPage) return;
        if (eventsPage.events.length === 0) {
            this.#elements.divEventsTable.innerHTML = `<p class="empty">Sin eventos con este rango y estos filtros.</p>`;
            return;
        }
        const firstPosition = (eventsPage.page - 1) * eventsPage.pageSize;
        this.#elements.divEventsTable.innerHTML = `<div class="table-wrap table-scroll"><table class="table table-events">
            <thead><tr>
                <th class="num">#</th><th>Fecha (UTC)</th><th>Gravedad</th><th>Proyecto / agente</th><th>Fuente</th><th>Mensaje</th>
            </tr></thead>
            <tbody>${eventsPage.events.map((event, index) => {
                const eventId = String(event.event_id);
                const isExpanded = this.#expandedEventIds.has(eventId);
                return `<tr class="event-row${isExpanded ? " is-expanded" : ""}">
                    <td class="num"><button type="button" class="row-toggle" data-action="${ACTION_ROW_TOGGLE}" data-event-id="${h.getEscaped(eventId)}"
                        aria-expanded="${isExpanded}" aria-label="Desplegar el evento" title="Ver todos sus campos">${isExpanded ? "▾" : "▸"} ${firstPosition + index + 1}</button></td>
                    <td class="nowrap mono">${h.getEscaped(this.#dateTimer.getUtcDateTime(event.occurred_at))}</td>
                    <td>${h.getBadge(event.severity)}</td>
                    <td>${h.getEscaped(event.project_name)}<br><span class="muted">${h.getEscaped(event.agent_name)}</span></td>
                    <td>${h.getEscaped(event.source_name)}<br><span class="muted mono">${h.getEscaped(event.source_type)}</span></td>
                    <td class="message${isExpanded ? "" : " is-clamped"}">${h.getEscaped(event.message)}</td>
                </tr>${isExpanded ? this.#getEventDetailHtml(event) : ""}`;
            }).join("")}</tbody>
        </table></div>`;
    }

    /**
     * Todos los campos del evento y sus attributes como clave / valor; los filtrables llevan + / −.
     */
    #getEventDetailHtml(event) {
        const h = this.#htmler;
        const rows = [
            ["event_id", event.event_id],
            ["occurred_at", this.#dateTimer.getUtcDateTime(event.occurred_at)],
            ["received_at", this.#dateTimer.getUtcDateTime(event.received_at)],
            ["severity", event.severity, EventFilterFieldEnum.SEVERITY, event.severity],
            ["event_type", event.event_type, EventFilterFieldEnum.EVENT_TYPE, event.event_type],
            ["source_type", event.source_type, EventFilterFieldEnum.SOURCE_TYPE, event.source_type],
            ["source_name", event.source_name, EventFilterFieldEnum.SOURCE_NAME, event.source_name],
            ["project", `${event.project_name ?? ""} (${event.project_uuid ?? ""})`, EventFilterFieldEnum.PROJECT_UUID, event.project_uuid],
            ["agent", `${event.agent_name ?? ""} (${event.agent_uuid ?? ""})`, EventFilterFieldEnum.AGENT_UUID, event.agent_uuid],
            ["file_path", event.file_path],
            ["exit_code", event.exit_code],
            ["duration_ms", event.duration_ms],
            ["message", event.message],
        ].filter(([, value]) => value !== null && value !== undefined && value !== "");
        const attributes = event.attributes && typeof event.attributes === "object" ? event.attributes : {};
        Object.keys(attributes).sort().forEach((key) => rows.push([`attributes.${key}`, this.#getAttributeText(attributes[key])]));

        return `<tr class="event-detail"><td colspan="${TABLE_COLUMNS}">
            <table class="kv"><tbody>${rows.map(([key, value, filterField, filterValue]) => `<tr>
                <th>${h.getEscaped(key)}</th>
                <td class="kv-actions">${filterField && filterValue ? this.#getFilterButtonsHtml(filterField, String(filterValue)) : ""}</td>
                <td class="kv-value">${h.getEscaped(value)}</td>
            </tr>`).join("")}</tbody></table>
        </td></tr>`;
    }

    #getAttributeText(value) {
        return value !== null && typeof value === "object" ? JSON.stringify(value) : String(value);
    }

    // ---- eventos del DOM ----

    #onClick(event) {
        const target = event.target.closest("[data-action]");
        if (!target) return;
        const dataset = target.dataset;
        const isExcluded = dataset.excluded === "1";
        switch (dataset.action) {
            case ACTION_FILTER:
                this.#goToState(this.#eventsQueryParser.getStateWithFilter(this.#state, dataset.field, dataset.value, isExcluded));
                return;
            case ACTION_PILL_INVERT:
                this.#goToState(this.#eventsQueryParser.getStateWithFilter(this.#state, dataset.field, dataset.value, !isExcluded));
                return;
            case ACTION_PILL_REMOVE:
                this.#goToState(this.#eventsQueryParser.getStateWithoutFilter(this.#state, dataset.field, dataset.value));
                return;
            case ACTION_PILLS_CLEAR:
                this.#goToState(this.#eventsQueryParser.getStateWithoutFilters(this.#state));
                return;
            case ACTION_BUCKET:
                this.#goToBucket(target);
                return;
            case ACTION_ROW_TOGGLE:
                this.#toggleRow(dataset.eventId);
                return;
            default:
        }
    }

    #toggleRow(eventId) {
        if (this.#expandedEventIds.has(eventId)) {
            this.#expandedEventIds.delete(eventId);
        }
        else {
            this.#expandedEventIds.add(eventId);
        }
        this.#renderTable();
    }

    #onQuerySubmit(event) {
        event.preventDefault();
        this.#goToState(this.#eventsQueryParser.getStateWith(this.#state, { search: this.#elements.inputEventsSearch.value.trim() }));
    }

    #onRangeChange() {
        const range = this.#elements.selectEventsRange.value;
        if (range !== TimeRangeEnum.CUSTOM) {
            this.#goToState(this.#eventsQueryParser.getStateWith(this.#state, { range, from: "", to: "" }));
            return;
        }
        // a medida: se enseñan los dos campos con el rango que se estaba viendo, y se aplica al enviarlos
        const fromTo = this.#timeRangeResolver.getFromTo(this.#state, new Date());
        this.#elements.inputEventsFrom.value = fromTo.from;
        this.#elements.inputEventsTo.value = fromTo.to;
        this.#elements.formEventsCustomRange.classList.remove(CSS_HIDDEN);
        this.#elements.inputEventsFrom.focus();
    }

    #onCustomRangeSubmit(event) {
        event.preventDefault();
        const from = this.#elements.inputEventsFrom.value.trim();
        const to = this.#elements.inputEventsTo.value.trim();
        if (!this.#eventsQueryParser.isUtcDateTime(from) || !this.#eventsQueryParser.isUtcDateTime(to) || from > to) {
            this.#renderRangeErrors(RANGE_ERROR);
            return;
        }
        this.#goToState(this.#eventsQueryParser.getStateWith(this.#state, { range: "", from, to }));
    }

    #renderRangeErrors(text) {
        this.#elements.pEventsRangeErrors.textContent = text;
        this.#elements.pEventsRangeErrors.classList.toggle(CSS_HIDDEN, text === "");
    }

    #goToPage(page) {
        this.#goToState(this.#eventsQueryParser.getStateWith(this.#state, { page: Math.max(EventsPaginationEnum.FIRST_PAGE, page) }));
    }

    /**
     * Estado nuevo = hash nuevo. Mismo hash: no hay hashchange, se pide repintar igual.
     * @param {Object} state
     */
    #goToState(state) {
        const nextHash = this.#eventsQueryParser.getHash(state);
        if (window.location.hash === nextHash) {
            window.dispatchEvent(new HashChangeEvent("hashchange"));
            return;
        }
        window.location.hash = nextHash;
    }
}
