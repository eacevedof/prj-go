import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
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
import { ProjectDangerView } from "/js/Modules/Projects/Infrastructure/Views/ProjectDangerView.js";
import { ProjectsException } from "/js/Modules/Projects/Domain/Exceptions/ProjectsException.js";

const DASHBOARD_HASH = "#/";

/**
 * Baja de un proyecto desde su detalle (ya confirmada en la vista): borrado, se vuelve a la portada. Tambien le dice
 * a la vista que proyecto se esta viendo (Projects escucha a Dashboard, no al reves).
 */
export class DeleteProjectController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #navigationView;
    #projectDangerView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#navigationView = NavigationView.getInstance();
        this.#projectDangerView = ProjectDangerView.getInstance();
    }

    /**
     * @returns {DeleteProjectController}
     */
    static getInstance() {
        return new DeleteProjectController();
    }

    init() {
        this.#subscriber.subscribe(ProjectsEventEnum.DELETE_PROJECT_REQUESTED, (request) => this.invoke(request));
        this.#subscriber.subscribe(DashboardEventEnum.PROJECT_DETAIL_RENDERED, (request) => this.#projectDangerView.renderProject(request.project ?? {}));
    }

    /**
     * @param {{projectUuid: string}} request
     */
    async invoke(request) {
        try {
            await DeleteProjectService.getInstance().invoke(
                DeleteProjectDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    projectUuid: request.projectUuid,
                })
            );
            this.#projectDangerView.reset();
            this.#navigationView.goToHash(DASHBOARD_HASH);
        }
        catch (error) {
            if (error instanceof ProjectsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof ProjectsException) {
                this.#projectDangerView.renderErrors(error.getErrors());
                return;
            }
            console.error("DeleteProjectController.invoke", error);
            this.#projectDangerView.renderErrors([ProjectsMessageEnum.UNEXPECTED_ERROR]);
        }
    }
}
