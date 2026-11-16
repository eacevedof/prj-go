import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { TextFormatter } from "/js/Modules/Shared/Infrastructure/Components/TextFormatter.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { ProjectsEventEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsEventEnum.js";
import { ProjectsMessageEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsMessageEnum.js";
import { ProjectsTextEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsTextEnum.js";
import { DeleteProjectDto } from "/js/Modules/Projects/Application/DeleteProject/DeleteProjectDto.js";
import { DeleteProjectService } from "/js/Modules/Projects/Application/DeleteProject/DeleteProjectService.js";
import { ProjectFormView } from "/js/Modules/Projects/Infrastructure/Views/ProjectFormView.js";
import { ProjectsTableView } from "/js/Modules/Projects/Infrastructure/Views/ProjectsTableView.js";
import { ProjectsException } from "/js/Modules/Projects/Domain/Exceptions/ProjectsException.js";

/**
 * Baja de un proyecto desde la tabla de la portada (ya confirmada en la fila): borrado, se queda en la portada con el
 * aviso y la tabla repintada. El mismo caso de uso que el borrado desde el detalle (DeleteProjectController), que en
 * cambio vuelve a la portada.
 */
export class DeleteListedProjectController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #textFormatter;
    #navigationView;
    #projectFormView;
    #projectsTableView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#textFormatter = TextFormatter.getInstance();
        this.#navigationView = NavigationView.getInstance();
        this.#projectFormView = ProjectFormView.getInstance();
        this.#projectsTableView = ProjectsTableView.getInstance();
    }

    /**
     * @returns {DeleteListedProjectController}
     */
    static getInstance() {
        return new DeleteListedProjectController();
    }

    init() {
        this.#subscriber.subscribe(ProjectsEventEnum.DELETE_LISTED_PROJECT_REQUESTED, (request) => this.invoke(request));
    }

    /**
     * @param {{projectUuid: string, projectName: string}} request
     */
    async invoke(request) {
        try {
            const deleteProjectResultDto = await DeleteProjectService.getInstance().invoke(
                DeleteProjectDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    projectUuid: request.projectUuid,
                })
            );
            this.#projectFormView.resetIfEditing(deleteProjectResultDto.getProjectUuid());
            this.#projectsTableView.renderMessage(
                this.#textFormatter.getFormatted(ProjectsTextEnum.PROJECT_DELETED, request.projectName)
            );
        }
        catch (error) {
            if (error instanceof ProjectsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof ProjectsException) {
                this.#projectsTableView.renderErrors(error.getErrors());
            }
            else {
                console.error("DeleteListedProjectController.invoke", error);
                this.#projectsTableView.renderErrors([ProjectsMessageEnum.UNEXPECTED_ERROR]);
            }
        }
        this.#navigationView.renderCurrentScreen();
    }
}
