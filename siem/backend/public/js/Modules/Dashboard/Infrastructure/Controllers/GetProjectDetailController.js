import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { ScreenEnum } from "/js/Modules/Shared/Domain/Enums/ScreenEnum.js";
import { SharedEventEnum } from "/js/Modules/Shared/Domain/Enums/SharedEventEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { DashboardEventEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardEventEnum.js";
import { DashboardMessageEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardMessageEnum.js";
import { GetProjectDetailDto } from "/js/Modules/Dashboard/Application/GetProjectDetail/GetProjectDetailDto.js";
import { GetProjectDetailService } from "/js/Modules/Dashboard/Application/GetProjectDetail/GetProjectDetailService.js";
import { ProjectDetailView } from "/js/Modules/Dashboard/Infrastructure/Views/ProjectDetailView.js";
import { DashboardException } from "/js/Modules/Dashboard/Domain/Exceptions/DashboardException.js";

/**
 * Detalle de un proyecto (#/projects/<project_uuid>): agentes y sus fuentes con el ultimo resultado. Se repinta sola.
 */
export class GetProjectDetailController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #projectDetailView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
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
            this.#projectDetailView.renderProjectDetail(
                getProjectDetailResultDto.getProject(),
                getProjectDetailResultDto.getAgents(),
                getProjectDetailResultDto.getServerNow()
            );
            this.#dispatcher.dispatch(DashboardEventEnum.PROJECT_DETAIL_RENDERED, {
                project: getProjectDetailResultDto.getProject(),
                agents: getProjectDetailResultDto.getAgents(),
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
