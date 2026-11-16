import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { ScreenEnum } from "/js/Modules/Shared/Domain/Enums/ScreenEnum.js";
import { SharedEventEnum } from "/js/Modules/Shared/Domain/Enums/SharedEventEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { DashboardEventEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardEventEnum.js";
import { DashboardMessageEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardMessageEnum.js";
import { GetDashboardDto } from "/js/Modules/Dashboard/Application/GetDashboard/GetDashboardDto.js";
import { GetDashboardService } from "/js/Modules/Dashboard/Application/GetDashboard/GetDashboardService.js";
import { DashboardView } from "/js/Modules/Dashboard/Infrastructure/Views/DashboardView.js";
import { DashboardException } from "/js/Modules/Dashboard/Domain/Exceptions/DashboardException.js";

/**
 * Portada (#/): estado del servidor; los proyectos los pinta Projects (su tabla escucha DASHBOARD_RENDERED). Se repinta sola (SCREEN_RENDERED con auto-refresco),
 * tambien tras un error: si el backend vuelve, la pantalla se recupera sin recargar.
 */
export class GetDashboardController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #dashboardView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#dashboardView = DashboardView.getInstance();
    }

    /**
     * @returns {GetDashboardController}
     */
    static getInstance() {
        return new GetDashboardController();
    }

    init() {
        this.#subscriber.subscribe(SharedEventEnum.SCREEN_REQUESTED, (request) => {
            if (request.screen !== ScreenEnum.DASHBOARD) return;
            this.invoke();
        });
    }

    async invoke() {
        try {
            const getDashboardResultDto = await GetDashboardService.getInstance().invoke(
                GetDashboardDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                })
            );
            this.#dashboardView.renderServer(getDashboardResultDto.getServer());
            this.#dispatcher.dispatch(DashboardEventEnum.DASHBOARD_RENDERED, { projects: getDashboardResultDto.getProjects() });
        }
        catch (error) {
            if (error instanceof DashboardException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof DashboardException) {
                this.#dashboardView.renderErrors(error.getErrors());
            }
            else {
                console.error("GetDashboardController.invoke", error);
                this.#dashboardView.renderErrors([DashboardMessageEnum.UNEXPECTED_ERROR]);
            }
        }
        this.#dispatcher.dispatch(SharedEventEnum.SCREEN_RENDERED, { isAutoRefreshed: true });
    }
}
