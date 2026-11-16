import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { RouteHashEnum } from "/js/Modules/Shared/Domain/Enums/RouteHashEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { DashboardEventEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardEventEnum.js";

import { ProjectsEventEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsEventEnum.js";
import { ProjectsMessageEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsMessageEnum.js";
import { DeleteProjectDto } from "/js/Modules/Projects/Application/DeleteProject/DeleteProjectDto.js";
import { DeleteProjectService } from "/js/Modules/Projects/Application/DeleteProject/DeleteProjectService.js";
import { ProjectDeleteDialogView } from "/js/Modules/Projects/Infrastructure/Views/ProjectDeleteDialogView.js";
import { ProjectDeletedView } from "/js/Modules/Projects/Infrastructure/Views/ProjectDeletedView.js";
import { ProjectFormView } from "/js/Modules/Projects/Infrastructure/Views/ProjectFormView.js";
import { ProjectsTableView } from "/js/Modules/Projects/Infrastructure/Views/ProjectsTableView.js";
import { ProjectsException } from "/js/Modules/Projects/Domain/Exceptions/ProjectsException.js";

/**
 * Baja de un proyecto, venga de la tabla de la portada o de la cabecera de su detalle: el modal pide el nombre, se
 * borra y se vuelve a la portada con el aviso de los siem-client que siguen instalados (servidor, ruta y comando para
 * desinstalarlos). Tambien abre el modal y le dice que proyecto se esta viendo (Projects escucha a Dashboard, no al
 * reves).
 */
export class DeleteProjectController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #navigationView;
    #projectDeleteDialogView;
    #projectDeletedView;
    #projectFormView;
    #projectsTableView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#navigationView = NavigationView.getInstance();
        this.#projectDeleteDialogView = ProjectDeleteDialogView.getInstance();
        this.#projectDeletedView = ProjectDeletedView.getInstance();
        this.#projectFormView = ProjectFormView.getInstance();
        this.#projectsTableView = ProjectsTableView.getInstance();
    }

    /**
     * @returns {DeleteProjectController}
     */
    static getInstance() {
        return new DeleteProjectController();
    }

    init() {
        this.#subscriber.subscribe(ProjectsEventEnum.DELETE_PROJECT_DIALOG_REQUESTED, (request) => this.#projectDeleteDialogView.open(request));
        this.#subscriber.subscribe(ProjectsEventEnum.DELETE_PROJECT_REQUESTED, (request) => this.invoke(request));
        this.#subscriber.subscribe(DashboardEventEnum.PROJECT_DETAIL_RENDERED, (request) => this.#projectDeleteDialogView.renderProject(request.project ?? {}));
    }

    /**
     * @param {{projectUuid: string}} request
     */
    async invoke(request) {
        try {
            const deleteProjectResultDto = await DeleteProjectService.getInstance().invoke(
                DeleteProjectDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    projectUuid: request.projectUuid,
                })
            );
            this.#projectDeleteDialogView.close();
            this.#projectFormView.resetIfEditing(deleteProjectResultDto.getProjectUuid());
            this.#projectsTableView.renderMessage("");
            this.#projectDeletedView.renderDeletedProject(
                deleteProjectResultDto.getProjectName(), deleteProjectResultDto.getAgents()
            );
            this.#navigationView.goToHash(RouteHashEnum.PROJECTS);
        }
        catch (error) {
            if (error instanceof ProjectsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#projectDeleteDialogView.close();
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof ProjectsException) {
                this.#projectDeleteDialogView.renderErrors(error.getErrors());
                return;
            }
            console.error("DeleteProjectController.invoke", error);
            this.#projectDeleteDialogView.renderErrors([ProjectsMessageEnum.UNEXPECTED_ERROR]);
        }
    }
}
