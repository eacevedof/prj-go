import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { NumberFormatter } from "/js/Modules/Shared/Infrastructure/Components/NumberFormatter.js";
import { TextFormatter } from "/js/Modules/Shared/Infrastructure/Components/TextFormatter.js";

import { ProjectHealthEnum } from "/js/Modules/Projects/Domain/Enums/ProjectHealthEnum.js";
import { ProjectsEventEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsEventEnum.js";
import { ProjectsKeyEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsKeyEnum.js";
import { ProjectsTextEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsTextEnum.js";

const CSS_HIDDEN = "hidden";
const ACTION_EDIT = "edit";
const ACTION_DELETE = "delete";
// el buscador no repinta en cada tecla
const TEXT_FILTER_DEBOUNCE_MS = 300;
// en la tabla, un trozo de la descripcion; la entera en el tooltip
const DESCRIPTION_EXCERPT_CHARS = 90;
const PROJECT_HASH_PREFIX = "#/projects/";
// color de cada contador de 24 h segun su gravedad (gris si es 0)
const COUNT_CSS_WARNING = "count-warning";
const COUNT_CSS_ERROR = "count-error";
const COUNT_CSS_CRITICAL = "count-critical";
const COUNT_CSS_ZERO = "count-zero";

/**
 * Tabla de proyectos de la portada (patron task-js): filtros (texto y salud), una fila por proyecto con id, uuid,
 * nombre (y su descripcion debajo), warnings, errores y criticos de las ultimas 24 h, y acciones (editar, borrar).
 * Borrar abre el modal que pide el nombre del proyecto (ProjectDeleteDialogView): se lleva agentes, monitores y eventos.
 * La tabla se repinta con cada auto-refresco; los filtros sobreviven (son estado de la vista).
 * Una sola instancia: guarda las referencias al DOM y los proyectos pintados.
 */
export class ProjectsTableView {

    static #instance = null;

    #dispatcher;
    #htmler;
    #numberFormatter;
    #textFormatter;
    #elements = {};
    #projectsByUuid = new Map();
    #textFilterTimeoutId = null;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#htmler = Htmler.getInstance();
        this.#numberFormatter = NumberFormatter.getInstance();
        this.#textFormatter = TextFormatter.getInstance();
    }

    /**
     * @returns {ProjectsTableView}
     */
    static getInstance() {
        ProjectsTableView.#instance ??= new ProjectsTableView();
        return ProjectsTableView.#instance;
    }

    init() {
        this.#elements = {
            spanProjectsCount: document.getElementById("span-projects-count"),
            pProjectsMessage: document.getElementById("p-projects-message"),
            pProjectsErrors: document.getElementById("p-projects-errors"),
            formProjectsFilters: document.getElementById("form-projects-filters"),
            inputFilterProjectText: document.getElementById("input-filter-project-text"),
            selectFilterProjectHealth: document.getElementById("select-filter-project-health"),
            buttonClearProjectsFilters: document.getElementById("button-clear-projects-filters"),
            pProjectsEmpty: document.getElementById("p-projects-empty"),
            divProjectsTable: document.getElementById("div-projects-table"),
            tbodyProjects: document.getElementById("tbody-projects"),
        };
        this.#loadHealthOptions();
        this.#elements.formProjectsFilters.addEventListener("submit", (event) => this.#onFiltersSubmit(event));
        this.#elements.inputFilterProjectText.addEventListener("input", () => this.#onTextFilterInput());
        this.#elements.selectFilterProjectHealth.addEventListener("change", () => this.#dispatchFiltersChanged());
        this.#elements.buttonClearProjectsFilters.addEventListener("click", () => this.#onClearFiltersClick());
        this.#elements.tbodyProjects.addEventListener("click", (event) => this.#onProjectsTableClick(event));
    }

    /**
     * @returns {{projectText: string, projectHealth: string}}
     */
    getFilters() {
        return {
            projectText: this.#elements.inputFilterProjectText.value,
            projectHealth: this.#elements.selectFilterProjectHealth.value,
        };
    }

    /**
     * @param {Object[]} filteredProjects los que cumplen los filtros (tal cual los da la API)
     * @param {number} projectsTotal todos, sin filtrar
     * @param {boolean} hasFilters
     */
    renderProjects(filteredProjects, projectsTotal, hasFilters) {
        this.#projectsByUuid = new Map(filteredProjects.map((project) => [project[ProjectsKeyEnum.PROJECT_UUID], project]));

        this.#elements.spanProjectsCount.textContent = this.#getCountText(filteredProjects.length, projectsTotal);
        this.#elements.pProjectsEmpty.textContent = hasFilters ? ProjectsTextEnum.LIST_NO_MATCHES : ProjectsTextEnum.LIST_EMPTY;
        this.#elements.pProjectsEmpty.classList.toggle(CSS_HIDDEN, filteredProjects.length > 0);
        this.#elements.divProjectsTable.classList.toggle(CSS_HIDDEN, filteredProjects.length === 0);
        this.#renderProjectsTable();
    }

    /**
     * Aviso del camino feliz (creado, guardado). Vacio lo oculta; el de borrado es ProjectDeletedView.
     * @param {string} message
     */
    renderMessage(message) {
        this.#elements.pProjectsMessage.textContent = message;
        this.#elements.pProjectsMessage.classList.toggle(CSS_HIDDEN, message === "");
    }

    /**
     * @param {string[]} errors
     */
    renderErrors(errors) {
        this.#elements.pProjectsErrors.textContent = errors.join(" · ");
        this.#elements.pProjectsErrors.classList.toggle(CSS_HIDDEN, errors.length === 0);
    }

    #loadHealthOptions() {
        const healthOptions = ProjectHealthEnum.VALUES
            .map((projectHealth) => `<option value="${this.#htmler.getEscaped(projectHealth)}">${this.#htmler.getEscaped(ProjectHealthEnum.LABELS[projectHealth])}</option>`)
            .join("");
        this.#elements.selectFilterProjectHealth.innerHTML = `<option value="">${this.#htmler.getEscaped(ProjectsTextEnum.HEALTH_ALL)}</option>${healthOptions}`;
    }

    /**
     * Enter en el buscador: filtra ya, sin esperar al debounce.
     */
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
        this.#elements.formProjectsFilters.reset();
        this.#dispatchFiltersChanged();
    }

    #dispatchFiltersChanged() {
        this.#dispatcher.dispatch(ProjectsEventEnum.PROJECTS_FILTERS_CHANGED, { filters: this.getFilters() });
    }

    #onProjectsTableClick(event) {
        const actionButton = event.target.closest("button[data-action]");
        if (!actionButton) return;

        const project = this.#projectsByUuid.get(actionButton.dataset.projectUuid ?? "");
        switch (actionButton.dataset.action) {
            case ACTION_EDIT:
                if (!project) return;
                this.renderMessage("");
                this.#dispatcher.dispatch(ProjectsEventEnum.EDIT_PROJECT_REQUESTED, { project });
                return;
            case ACTION_DELETE:
                if (!project) return;
                this.renderErrors([]);
                this.#dispatcher.dispatch(ProjectsEventEnum.DELETE_PROJECT_DIALOG_REQUESTED, {
                    projectUuid: project[ProjectsKeyEnum.PROJECT_UUID],
                    projectName: project[ProjectsKeyEnum.PROJECT_NAME],
                });
        }
    }

    #renderProjectsTable() {
        this.#elements.tbodyProjects.innerHTML = [...this.#projectsByUuid.values()]
            .map((project) => this.#getHtmlProjectRow(project))
            .join("");
    }

    /**
     * "(4)" sin filtrar; "(2 de 4)" si los filtros dejan fuera alguno.
     */
    #getCountText(filteredTotal, projectsTotal) {
        if (projectsTotal === 0) return "";
        if (filteredTotal === projectsTotal) return this.#textFormatter.getFormatted(ProjectsTextEnum.COUNT_ALL, projectsTotal);
        return this.#textFormatter.getFormatted(ProjectsTextEnum.COUNT_FILTERED, filteredTotal, projectsTotal);
    }

    #getHtmlProjectRow(project) {
        const h = this.#htmler;
        const projectUuid = String(project[ProjectsKeyEnum.PROJECT_UUID] ?? "");
        return `
<tr>
    <td class="num row-id">${h.getEscaped(project[ProjectsKeyEnum.PROJECT_ID])}</td>
    <td class="row-uuid"><span class="copyable"><code>${h.getEscaped(projectUuid)}</code>${h.getCopyButton(projectUuid)}</span></td>
    <td class="row-main">
        <div class="row-main-head">
            <a class="row-title project-link" href="${PROJECT_HASH_PREFIX}${encodeURIComponent(projectUuid)}">${h.getEscaped(project[ProjectsKeyEnum.PROJECT_NAME])}</a>
            ${h.getBadge(project[ProjectsKeyEnum.PROJECT_HEALTH])}
        </div>
        ${this.#getHtmlProjectDescription(project[ProjectsKeyEnum.PROJECT_DESCRIPTION])}
    </td>
    ${this.#getHtmlCountCell(project[ProjectsKeyEnum.WARNINGS_24H], COUNT_CSS_WARNING)}
    ${this.#getHtmlCountCell(project[ProjectsKeyEnum.ERRORS_24H], COUNT_CSS_ERROR)}
    ${this.#getHtmlCountCell(project[ProjectsKeyEnum.CRITICALS_24H], COUNT_CSS_CRITICAL)}
    <td class="row-actions">${this.#getHtmlRowActions(projectUuid)}</td>
</tr>`;
    }

    /**
     * Descripcion en una linea bajo el nombre; si es larga, cortada con "…" y entera en el tooltip (title).
     */
    #getHtmlProjectDescription(projectDescription) {
        const fullDescription = String(projectDescription ?? "").trim();
        if (!fullDescription) {
            return `<div class="row-sub muted">${this.#htmler.getEscaped(ProjectsTextEnum.NO_DESCRIPTION)}</div>`;
        }
        const oneLineDescription = fullDescription.replace(/\s+/g, " ");
        if (oneLineDescription.length <= DESCRIPTION_EXCERPT_CHARS) {
            return `<div class="row-sub">${this.#htmler.getEscaped(oneLineDescription)}</div>`;
        }
        const excerpt = `${oneLineDescription.slice(0, DESCRIPTION_EXCERPT_CHARS).trimEnd()}…`;
        return `<div class="row-sub is-excerpt" title="${this.#htmler.getEscaped(fullDescription)}" tabindex="0">${this.#htmler.getEscaped(excerpt)}</div>`;
    }

    /**
     * Un contador de 24 h: en gris si es 0, con el color de su gravedad si no.
     * @param {number} eventsCount
     * @param {string} countCss COUNT_CSS_* (la clase de su gravedad)
     */
    #getHtmlCountCell(eventsCount, countCss) {
        const count = Number(eventsCount) || 0;
        const countClass = count > 0 ? countCss : COUNT_CSS_ZERO;
        return `<td class="num ${countClass}">${this.#numberFormatter.getNumber(count)}</td>`;
    }

    #getHtmlRowActions(projectUuid) {
        const h = this.#htmler;
        const escapedProjectUuid = h.getEscaped(projectUuid);
        return `<button type="button" class="button button-small" data-action="${ACTION_EDIT}" data-project-uuid="${escapedProjectUuid}">${h.getEscaped(ProjectsTextEnum.EDIT)}</button>
            <button type="button" class="button button-small button-danger" data-action="${ACTION_DELETE}" data-project-uuid="${escapedProjectUuid}">${h.getEscaped(ProjectsTextEnum.DELETE)}</button>`;
    }
}
