import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";

import { ProjectsEventEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsEventEnum.js";

const CSS_HIDDEN = "hidden";

/**
 * "Borrar proyecto" del detalle de un proyecto. La confirmacion va dentro de la pagina: hay que escribir el nombre
 * exacto del proyecto para que el boton se active (borra agentes, monitores y eventos; no se puede deshacer).
 * Una sola instancia: guarda las referencias al DOM y el proyecto que se esta viendo.
 */
export class ProjectDangerView {

    static #instance = null;

    #dispatcher;
    #elements = {};
    #projectUuid = "";
    #projectName = "";

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
    }

    /**
     * @returns {ProjectDangerView}
     */
    static getInstance() {
        ProjectDangerView.#instance ??= new ProjectDangerView();
        return ProjectDangerView.#instance;
    }

    init() {
        this.#elements = {
            formProjectDelete: document.getElementById("form-project-delete"),
            spanProjectDeleteName: document.getElementById("span-project-delete-name"),
            inputProjectDeleteName: document.getElementById("input-project-delete-name"),
            buttonProjectDelete: document.getElementById("button-project-delete"),
            pProjectDeleteErrors: document.getElementById("p-project-delete-errors"),
        };
        this.#elements.inputProjectDeleteName.addEventListener("input", () => this.#renderButton());
        this.#elements.formProjectDelete.addEventListener("submit", (event) => this.#onSubmit(event));
    }

    /**
     * El proyecto que se esta viendo. Si es otro que el anterior, el formulario se vacia; en un refresco del mismo, no
     * (no se pierde lo escrito).
     * @param {{project_uuid: string, project_name: string}} project
     */
    renderProject(project) {
        const projectUuid = String(project.project_uuid ?? "");
        if (projectUuid !== this.#projectUuid) {
            this.#elements.formProjectDelete.reset();
            this.renderErrors([]);
        }
        this.#projectUuid = projectUuid;
        this.#projectName = String(project.project_name ?? "");
        this.#elements.spanProjectDeleteName.textContent = this.#projectName;
        this.#renderButton();
    }

    /**
     * @param {string[]} errors
     */
    renderErrors(errors) {
        this.#elements.pProjectDeleteErrors.textContent = errors.join(" · ");
        this.#elements.pProjectDeleteErrors.classList.toggle(CSS_HIDDEN, errors.length === 0);
    }

    /**
     * Tras borrar: nada del proyecto anterior queda en el formulario.
     */
    reset() {
        this.#projectUuid = "";
        this.#projectName = "";
        this.#elements.formProjectDelete.reset();
        this.#renderButton();
    }

    #isConfirmed() {
        return this.#projectName !== "" && this.#elements.inputProjectDeleteName.value.trim() === this.#projectName;
    }

    #renderButton() {
        this.#elements.buttonProjectDelete.disabled = !this.#isConfirmed();
    }

    #onSubmit(event) {
        event.preventDefault();
        if (!this.#isConfirmed()) return;
        this.renderErrors([]);
        this.#dispatcher.dispatch(ProjectsEventEnum.DELETE_PROJECT_REQUESTED, { projectUuid: this.#projectUuid });
    }
}
