import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";

import { ProjectsEventEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsEventEnum.js";

const CSS_HIDDEN = "hidden";

/**
 * Formulario "nuevo proyecto" de la portada. Es html fijo (no se repinta con el refresco: lo que se este escribiendo
 * no se pierde). Una sola instancia: guarda las referencias al DOM.
 */
export class ProjectFormView {

    static #instance = null;

    #dispatcher;
    #elements = {};

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
    }

    /**
     * @returns {ProjectFormView}
     */
    static getInstance() {
        ProjectFormView.#instance ??= new ProjectFormView();
        return ProjectFormView.#instance;
    }

    init() {
        this.#elements = {
            formProject: document.getElementById("form-project"),
            inputProjectName: document.getElementById("input-project-name"),
            inputProjectDescription: document.getElementById("input-project-description"),
            buttonProjectSubmit: document.getElementById("button-project-submit"),
            pProjectFormErrors: document.getElementById("p-project-form-errors"),
            pProjectFormMessage: document.getElementById("p-project-form-message"),
        };
        this.#elements.formProject.addEventListener("submit", (event) => this.#onSubmit(event));
    }

    /**
     * Proyecto creado: formulario vacio y aviso.
     * @param {string} projectName
     */
    renderCreated(projectName) {
        this.#elements.formProject.reset();
        this.renderErrors([]);
        this.#renderText(this.#elements.pProjectFormMessage, `Proyecto "${projectName}" creado. Entra en él para dar de alta sus agentes y monitores.`);
    }

    /**
     * @param {string[]} errors
     */
    renderErrors(errors) {
        this.#renderText(this.#elements.pProjectFormErrors, errors.join(" · "));
    }

    /**
     * @param {boolean} isBusy
     */
    renderBusy(isBusy) {
        this.#elements.buttonProjectSubmit.disabled = isBusy;
    }

    #onSubmit(event) {
        event.preventDefault();
        this.renderErrors([]);
        this.#renderText(this.#elements.pProjectFormMessage, "");
        this.#dispatcher.dispatch(ProjectsEventEnum.CREATE_PROJECT_REQUESTED, {
            projectName: this.#elements.inputProjectName.value,
            projectDescription: this.#elements.inputProjectDescription.value,
        });
    }

    #renderText(element, text) {
        element.textContent = text;
        element.classList.toggle(CSS_HIDDEN, text === "");
    }
}
