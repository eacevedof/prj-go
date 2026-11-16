import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";

const CSS_HIDDEN = "hidden";

/**
 * Detalle de un proyecto: su cabecera (nombre, descripcion, uuid y enlace a sus eventos). Los agentes los pinta el
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
        };
    }

    /**
     * @param {Object} project tal cual lo devuelve la API
     */
    renderProject(project) {
        this.renderErrors([]);
        this.#elements.divProjectBody.classList.remove(CSS_HIDDEN);
        this.#elements.divProjectHeader.innerHTML = this.#getHeaderHtml(project);
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
}
