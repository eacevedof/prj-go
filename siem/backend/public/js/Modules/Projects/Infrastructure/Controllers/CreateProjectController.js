import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { ProjectsEventEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsEventEnum.js";
import { ProjectsMessageEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsMessageEnum.js";
import { CreateProjectDto } from "/js/Modules/Projects/Application/CreateProject/CreateProjectDto.js";
import { CreateProjectService } from "/js/Modules/Projects/Application/CreateProject/CreateProjectService.js";
import { ProjectFormView } from "/js/Modules/Projects/Infrastructure/Views/ProjectFormView.js";
import { ProjectsException } from "/js/Modules/Projects/Domain/Exceptions/ProjectsException.js";

/**
 * Alta de un proyecto desde la portada: creado, se repinta la portada (ya sale en la lista).
 */
export class CreateProjectController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #navigationView;
    #projectFormView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#navigationView = NavigationView.getInstance();
        this.#projectFormView = ProjectFormView.getInstance();
    }

    /**
     * @returns {CreateProjectController}
     */
    static getInstance() {
        return new CreateProjectController();
    }

    init() {
        this.#subscriber.subscribe(ProjectsEventEnum.CREATE_PROJECT_REQUESTED, (request) => this.invoke(request));
    }

    /**
     * @param {{projectName: string, projectDescription: string}} request
     */
    async invoke(request) {
        this.#projectFormView.renderBusy(true);
        try {
            const createProjectResultDto = await CreateProjectService.getInstance().invoke(
                CreateProjectDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    projectName: request.projectName,
                    projectDescription: request.projectDescription,
                })
            );
            this.#projectFormView.renderCreated(createProjectResultDto.getProjectName());
            this.#navigationView.renderCurrentScreen();
        }
        catch (error) {
            if (error instanceof ProjectsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof ProjectsException) {
                this.#projectFormView.renderErrors(error.getErrors());
                return;
            }
            console.error("CreateProjectController.invoke", error);
            this.#projectFormView.renderErrors([ProjectsMessageEnum.UNEXPECTED_ERROR]);
        }
        finally {
            this.#projectFormView.renderBusy(false);
        }
    }
}
