import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";

import { ProjectsEventEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsEventEnum.js";
import { ProjectsKeyEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsKeyEnum.js";

const CSS_HIDDEN = "hidden";

/**
 * Modal "Borrar proyecto" (<dialog> nativo, patron task-js): hay que escribir el nombre exacto del proyecto para que el
 * boton se active (se lleva agentes, monitores y eventos; no se puede deshacer). Lo abren el boton rojo de la cabecera
 * del detalle (el proyecto que se esta viendo) y "Borrar" de la tabla de la portada (DELETE_PROJECT_DIALOG_REQUESTED).
 * Una sola instancia: guarda las referencias al DOM, el proyecto del detalle y el que se va a borrar.
 */
export class ProjectDeleteDialogView {

    static #instance = null;

    #dispatcher;
    #elements = {};
    // el proyecto del detalle (para el boton de su cabecera)
    #viewedProject = { projectUuid: "", projectName: "" };
    // el que pregunta el modal
    #projectUuid = "";
    #projectName = "";

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
    }

    /**
     * @returns {ProjectDeleteDialogView}
     */
    static getInstance() {
        ProjectDeleteDialogView.#instance ??= new ProjectDeleteDialogView();
        return ProjectDeleteDialogView.#instance;
    }

    init() {
        this.#elements = {
            buttonProjectDeleteOpen: document.getElementById("button-project-delete-open"),
            dialogProjectDelete: document.getElementById("dialog-project-delete"),
            formProjectDelete: document.getElementById("form-project-delete"),
            spanProjectDeleteName: document.getElementById("span-project-delete-name"),
            inputProjectDeleteName: document.getElementById("input-project-delete-name"),
            buttonProjectDelete: document.getElementById("button-project-delete"),
            buttonProjectDeleteCancel: document.getElementById("button-project-delete-cancel"),
            pProjectDeleteErrors: document.getElementById("p-project-delete-errors"),
        };
        this.#elements.buttonProjectDeleteOpen.addEventListener("click", () => this.open(this.#viewedProject));
        this.#elements.inputProjectDeleteName.addEventListener("input", () => this.#renderButton());
        this.#elements.formProjectDelete.addEventListener("submit", (event) => this.#onSubmit(event));
        this.#elements.buttonProjectDeleteCancel.addEventListener("click", () => this.close());
        // Esc tambien cierra: que no quede nada escrito para la proxima vez
        this.#elements.dialogProjectDelete.addEventListener("close", () => this.#reset());
    }

    /**
     * El proyecto que se esta viendo en el detalle: el que borra el boton rojo de su cabecera.
     * @param {Object} project tal cual lo devuelve la API
     */
    renderProject(project) {
        this.#viewedProject = {
            projectUuid: String(project[ProjectsKeyEnum.PROJECT_UUID] ?? ""),
            projectName: String(project[ProjectsKeyEnum.PROJECT_NAME] ?? ""),
        };
    }

    /**
     * @param {{projectUuid: string, projectName: string}} project
     */
    open(project) {
        if (!project.projectUuid) return;
        this.#reset();
        this.#projectUuid = project.projectUuid;
        this.#projectName = project.projectName;
        this.#elements.spanProjectDeleteName.textContent = this.#projectName;
        this.#showDialog();
        this.#elements.inputProjectDeleteName.focus();
    }

    close() {
        if (typeof this.#elements.dialogProjectDelete.close === "function") {
            this.#elements.dialogProjectDelete.close();
            return;
        }
        this.#elements.dialogProjectDelete.removeAttribute("open");
        this.#reset();
    }

    /**
     * Un error de la API: el modal sigue abierto, con el boton otra vez disponible.
     * @param {string[]} errors
     */
    renderErrors(errors) {
        this.#elements.pProjectDeleteErrors.textContent = errors.join(" · ");
        this.#elements.pProjectDeleteErrors.classList.toggle(CSS_HIDDEN, errors.length === 0);
        this.#renderButton();
    }

    #showDialog() {
        // <dialog>: modal nativo (foco atrapado, Esc cierra); sin soporte, al menos visible
        if (typeof this.#elements.dialogProjectDelete.showModal === "function") {
            if (!this.#elements.dialogProjectDelete.open) this.#elements.dialogProjectDelete.showModal();
            return;
        }
        this.#elements.dialogProjectDelete.setAttribute("open", "");
    }

    #reset() {
        this.#projectUuid = "";
        this.#projectName = "";
        this.#elements.formProjectDelete.reset();
        this.renderErrors([]);
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
        this.#elements.buttonProjectDelete.disabled = true;
        this.#dispatcher.dispatch(ProjectsEventEnum.DELETE_PROJECT_REQUESTED, { projectUuid: this.#projectUuid });
    }
}
