import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { TextFormatter } from "/js/Modules/Shared/Infrastructure/Components/TextFormatter.js";

import { ProjectsEventEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsEventEnum.js";
import { ProjectsKeyEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsKeyEnum.js";
import { ProjectsTextEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsTextEnum.js";

const CSS_HIDDEN = "hidden";

/**
 * Formulario plegable "Nuevo proyecto" de la portada (patron task-js: el mismo para editar, con project_uuid oculto).
 * Plegado al entrar; al desplegarlo, el foco va al nombre; Ctrl+Enter (Cmd+Enter) guarda desde cualquier campo.
 * Es html fijo: el auto-refresco de la portada no lo repinta y lo que se este escribiendo no se pierde.
 * Una sola instancia: guarda las referencias al DOM y el id del proyecto que se edita.
 */
export class ProjectFormView {

    static #instance = null;

    #dispatcher;
    #textFormatter;
    #elements = {};
    // project_id del proyecto que se edita (va en el titulo); "" = alta
    #editingProjectId = "";

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#textFormatter = TextFormatter.getInstance();
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
            detailsProjectForm: document.getElementById("details-project-form"),
            summaryProjectForm: document.getElementById("summary-project-form"),
            formProject: document.getElementById("form-project"),
            hiddenProjectUuid: document.getElementById("hidden-project-uuid"),
            inputProjectName: document.getElementById("input-project-name"),
            textareaProjectDescription: document.getElementById("textarea-project-description"),
            buttonProjectSubmit: document.getElementById("button-project-submit"),
            buttonProjectCancel: document.getElementById("button-project-cancel"),
            pProjectFormErrors: document.getElementById("p-project-form-errors"),
        };
        this.#elements.detailsProjectForm.addEventListener("toggle", () => this.#onProjectFormToggle());
        this.#elements.formProject.addEventListener("submit", (event) => this.#onProjectSubmit(event));
        this.#elements.formProject.addEventListener("keydown", (event) => this.#onProjectFormKeydown(event));
        this.#elements.buttonProjectCancel.addEventListener("click", () => this.reset());
        this.#renderFormTitle();
    }

    /**
     * Vacio, plegado y en modo alta.
     */
    reset() {
        this.#elements.formProject.reset();
        this.#elements.hiddenProjectUuid.value = "";
        this.#editingProjectId = "";
        this.#renderFormTitle();
        this.#elements.buttonProjectCancel.classList.add(CSS_HIDDEN);
        this.#elements.detailsProjectForm.open = false;
        this.renderErrors([]);
        this.renderBusy(false);
    }

    /**
     * Si se estaba editando ese proyecto (lo acaban de borrar), el formulario vuelve a alta.
     * @param {string} projectUuid
     */
    resetIfEditing(projectUuid) {
        if (this.#elements.hiddenProjectUuid.value !== projectUuid) return;
        this.reset();
    }

    /**
     * "Editar" en la tabla: el proyecto entra en el formulario, que se despliega con el foco en el nombre.
     * @param {Object} project tal cual lo da el dashboard (snake_case)
     */
    renderProject(project) {
        this.#elements.hiddenProjectUuid.value = String(project[ProjectsKeyEnum.PROJECT_UUID] ?? "");
        this.#editingProjectId = String(project[ProjectsKeyEnum.PROJECT_ID] ?? "");
        this.#elements.inputProjectName.value = String(project[ProjectsKeyEnum.PROJECT_NAME] ?? "");
        this.#elements.textareaProjectDescription.value = String(project[ProjectsKeyEnum.PROJECT_DESCRIPTION] ?? "");
        this.#renderFormTitle();
        this.#elements.buttonProjectCancel.classList.remove(CSS_HIDDEN);
        this.renderErrors([]);
        this.#elements.detailsProjectForm.open = true;
        this.#elements.detailsProjectForm.scrollIntoView({ block: "nearest" });
        this.#elements.inputProjectName.focus();
    }

    /**
     * @param {string[]} errors
     */
    renderErrors(errors) {
        this.#elements.pProjectFormErrors.textContent = errors.join(" · ");
        this.#elements.pProjectFormErrors.classList.toggle(CSS_HIDDEN, errors.length === 0);
    }

    /**
     * @param {boolean} isBusy
     */
    renderBusy(isBusy) {
        this.#elements.buttonProjectSubmit.disabled = isBusy;
        this.#elements.buttonProjectCancel.disabled = isBusy;
    }

    /**
     * Al desplegar, directo a escribir el nombre (al editar ya lo pone renderProject).
     */
    #onProjectFormToggle() {
        if (!this.#elements.detailsProjectForm.open) return;
        this.#elements.inputProjectName.focus();
    }

    /**
     * requestSubmit y no submit(): pasa por #onProjectSubmit.
     */
    #onProjectFormKeydown(event) {
        if (event.key !== "Enter" || !(event.ctrlKey || event.metaKey)) return;

        event.preventDefault();
        this.#elements.formProject.requestSubmit();
    }

    #onProjectSubmit(event) {
        event.preventDefault();
        this.renderErrors([]);
        this.renderBusy(true);
        const project = {
            projectName: this.#elements.inputProjectName.value,
            projectDescription: this.#elements.textareaProjectDescription.value,
        };
        const projectUuid = this.#elements.hiddenProjectUuid.value;
        if (!projectUuid) {
            this.#dispatcher.dispatch(ProjectsEventEnum.CREATE_PROJECT_REQUESTED, project);
            return;
        }
        this.#dispatcher.dispatch(ProjectsEventEnum.UPDATE_PROJECT_REQUESTED, { projectUuid, ...project });
    }

    /**
     * "Nuevo proyecto" en el alta; "Editar proyecto #3" al editar. El boton dice lo que va a hacer.
     */
    #renderFormTitle() {
        const isEditing = this.#editingProjectId !== "";
        this.#elements.summaryProjectForm.textContent = isEditing
            ? this.#textFormatter.getFormatted(ProjectsTextEnum.FORM_TITLE_UPDATE, this.#editingProjectId)
            : ProjectsTextEnum.FORM_TITLE_CREATE;
        this.#elements.buttonProjectSubmit.textContent = isEditing ? ProjectsTextEnum.SUBMIT_UPDATE : ProjectsTextEnum.SUBMIT_CREATE;
    }
}
