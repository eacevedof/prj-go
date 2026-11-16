import { BreadcrumbsTextEnum } from "/js/Modules/Shared/Domain/Enums/BreadcrumbsTextEnum.js";
import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { RouteHashEnum } from "/js/Modules/Shared/Domain/Enums/RouteHashEnum.js";
import { ScreenEnum } from "/js/Modules/Shared/Domain/Enums/ScreenEnum.js";
import { SharedEventEnum } from "/js/Modules/Shared/Domain/Enums/SharedEventEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { DashboardEventEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardEventEnum.js";
import { DashboardMessageEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardMessageEnum.js";
import { GetProjectDetailDto } from "/js/Modules/Dashboard/Application/GetProjectDetail/GetProjectDetailDto.js";
import { GetProjectDetailService } from "/js/Modules/Dashboard/Application/GetProjectDetail/GetProjectDetailService.js";
import { ProjectDetailView } from "/js/Modules/Dashboard/Infrastructure/Views/ProjectDetailView.js";
import { DashboardException } from "/js/Modules/Dashboard/Domain/Exceptions/DashboardException.js";


/**
 * Detalle de un proyecto (#/projects/<project_uuid>): su cabecera y las migas con su nombre; los agentes (con sus
 * fuentes) los pinta Agents al oir PROJECT_DETAIL_RENDERED. Se repinta sola.
 */
export class GetProjectDetailController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #navigationView;
    #projectDetailView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#navigationView = NavigationView.getInstance();
        this.#projectDetailView = ProjectDetailView.getInstance();
    }

    /**
     * @returns {GetProjectDetailController}
     */
    static getInstance() {
        return new GetProjectDetailController();
    }

    init() {
        this.#subscriber.subscribe(SharedEventEnum.SCREEN_REQUESTED, (request) => {
            if (request.screen !== ScreenEnum.PROJECT_DETAIL) return;
            this.invoke({ projectUuid: request.pathParam });
        });
    }

    /**
     * @param {{projectUuid: string}} request
     */
    async invoke(request) {
        try {
            const getProjectDetailResultDto = await GetProjectDetailService.getInstance().invoke(
                GetProjectDetailDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    projectUuid: request.projectUuid,
                })
            );
            const project = getProjectDetailResultDto.getProject();
            this.#projectDetailView.renderProject(project);
            this.#navigationView.renderBreadcrumbs([
                { label: BreadcrumbsTextEnum.PROJECTS, hash: RouteHashEnum.PROJECTS },
                { label: String(project.project_name ?? "") },
            ]);
            this.#dispatcher.dispatch(DashboardEventEnum.PROJECT_DETAIL_RENDERED, {
                project,
                agents: getProjectDetailResultDto.getAgents(),
                serverNow: getProjectDetailResultDto.getServerNow(),
            });
        }
        catch (error) {
            if (error instanceof DashboardException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof DashboardException) {
                this.#projectDetailView.renderErrors(error.getErrors());
            }
            else {
                console.error("GetProjectDetailController.invoke", error);
                this.#projectDetailView.renderErrors([DashboardMessageEnum.UNEXPECTED_ERROR]);
            }
        }
        this.#dispatcher.dispatch(SharedEventEnum.SCREEN_RENDERED, { isAutoRefreshed: true });
    }
}
