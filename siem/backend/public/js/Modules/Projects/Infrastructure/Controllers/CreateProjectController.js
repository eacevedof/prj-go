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
import { CreateProjectDto } from "/js/Modules/Projects/Application/CreateProject/CreateProjectDto.js";
import { CreateProjectService } from "/js/Modules/Projects/Application/CreateProject/CreateProjectService.js";
import { ProjectFormView } from "/js/Modules/Projects/Infrastructure/Views/ProjectFormView.js";
import { ProjectsTableView } from "/js/Modules/Projects/Infrastructure/Views/ProjectsTableView.js";
import { ProjectsException } from "/js/Modules/Projects/Domain/Exceptions/ProjectsException.js";

/**
 * Alta de un proyecto desde la portada: creado, el formulario se pliega vacio y se repinta la portada (ya sale en la
 * tabla).
 */
export class CreateProjectController {

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
        try {
            const createProjectResultDto = await CreateProjectService.getInstance().invoke(
                CreateProjectDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    projectName: request.projectName,
                    projectDescription: request.projectDescription,
                })
            );
            this.#projectFormView.reset();
            this.#projectsTableView.renderMessage(
                this.#textFormatter.getFormatted(ProjectsTextEnum.PROJECT_CREATED, createProjectResultDto.getProjectName())
            );
            this.#navigationView.renderCurrentScreen();
        }
        catch (error) {
            this.#projectFormView.renderBusy(false);
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
    }
}
