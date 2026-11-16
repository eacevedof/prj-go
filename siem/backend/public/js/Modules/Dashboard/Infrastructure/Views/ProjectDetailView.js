import { RouteHashEnum } from "/js/Modules/Shared/Domain/Enums/RouteHashEnum.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";

import { DashboardKeyEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardKeyEnum.js";

const CSS_HIDDEN = "hidden";

/**
 * Detalle de un proyecto: su cabecera (nombre, descripcion y enlace a sus eventos; el boton rojo "Borrar proyecto" de
 * al lado es de Projects, ProjectDeleteDialogView). Los agentes los pinta el
 * modulo Agents (su tabla escucha PROJECT_DETAIL_RENDERED) y los monitores, Monitors.
 * Una sola instancia: guarda las referencias al DOM.
 */
export class ProjectDetailView {

    static #instance = null;

    #htmler;
    #elements = {};

    constructor() {
        this.#htmler = Htmler.getInstance();
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
            aProjectEvents: document.getElementById("a-project-events"),
        };
    }

    /**
     * @param {Object} project tal cual lo devuelve la API
     */
    renderProject(project) {
        this.renderErrors([]);
        this.#elements.divProjectBody.classList.remove(CSS_HIDDEN);
        this.#elements.divProjectHeader.innerHTML = this.#getHeaderHtml(project);
        this.#elements.aProjectEvents.href = `${RouteHashEnum.EVENTS}?${DashboardKeyEnum.PROJECT_UUID}=${encodeURIComponent(project[DashboardKeyEnum.PROJECT_UUID] ?? "")}`;
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
        return `
            <h1>1 - ${h.getEscaped(project[DashboardKeyEnum.PROJECT_NAME])}</h1>
            <p class="muted">${h.getEscaped(project[DashboardKeyEnum.PROJECT_DESCRIPTION])}</p>`;
    }
}
